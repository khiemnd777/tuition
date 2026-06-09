package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	emailCronPath         = "email_cron.local.json"
	defaultCronSendTime   = "08:00"
	defaultCronDailyLimit = 500
	emailQuotaWindow      = 24 * time.Hour
)

var emailCronMu sync.Mutex

type emailCronRequest struct {
	Rows       *[]paymentRow `json:"rows,omitempty"`
	Template   string        `json:"template,omitempty"`
	Enabled    *bool         `json:"enabled,omitempty"`
	DailyLimit int           `json:"dailyLimit,omitempty"`
	SendTime   string        `json:"sendTime,omitempty"`
}

type emailCronState struct {
	Enabled     bool              `json:"enabled"`
	DailyLimit  int               `json:"dailyLimit"`
	SendTime    string            `json:"sendTime"`
	Template    string            `json:"template"`
	Queue       []emailCronJob    `json:"queue"`
	LastRunAt   string            `json:"lastRunAt,omitempty"`
	LastRunDate string            `json:"lastRunDate,omitempty"`
	SentToday   int               `json:"sentToday,omitempty"`
	SendHistory []string          `json:"sendHistory,omitempty"`
	LastResults []emailSendResult `json:"lastResults,omitempty"`
	UpdatedAt   string            `json:"updatedAt,omitempty"`
}

type emailCronJob struct {
	ID          string     `json:"id"`
	Row         paymentRow `json:"row"`
	Status      string     `json:"status"`
	Attempts    int        `json:"attempts,omitempty"`
	ResendID    string     `json:"resendId,omitempty"`
	MessageID   string     `json:"messageId,omitempty"`
	Error       string     `json:"error,omitempty"`
	CreatedAt   string     `json:"createdAt"`
	ProcessedAt string     `json:"processedAt,omitempty"`
}

type emailCronResponse struct {
	Enabled     bool              `json:"enabled"`
	DailyLimit  int               `json:"dailyLimit"`
	SendTime    string            `json:"sendTime"`
	Template    string            `json:"template"`
	QueueTotal  int               `json:"queueTotal"`
	Queued      int               `json:"queued"`
	Sent        int               `json:"sent"`
	Errors      int               `json:"errors"`
	Skipped     int               `json:"skipped"`
	SentToday   int               `json:"sentToday"`
	SentLast24H int               `json:"sentLast24h"`
	LastRunAt   string            `json:"lastRunAt,omitempty"`
	NextRunAt   string            `json:"nextRunAt,omitempty"`
	LastResults []emailSendResult `json:"lastResults,omitempty"`
}

type emailQuotaStatus struct {
	Limit     int
	Sent      int
	Remaining int
}

func startEmailCronScheduler(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := runDueTenantEmailCronBatches(ctx); err != nil {
				fmt.Printf("tenant email cron: %v\n", err)
				recordEmailCronRunError(ctx, "tenant_scheduler", err)
			}
			state, err := runEmailCronBatch(ctx, false)
			if err != nil {
				fmt.Printf("email cron: %v\n", err)
				recordEmailCronRunError(ctx, "scheduler", err)
				continue
			}
			recordEmailCronResultErrors(ctx, state, "scheduler")
		}
	}
}

func handleEmailCron(w http.ResponseWriter, r *http.Request) {
	tenantID := activeTenantIDFromRequest(r)
	switch r.Method {
	case http.MethodGet:
		state, err := loadEmailCronStateForTenant(r.Context(), tenantID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, emailCronPublic(state, time.Now()))
	case http.MethodPost:
		r.Body = http.MaxBytesReader(w, r.Body, 4<<20)
		var req emailCronRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}
		if req.Rows != nil && len(*req.Rows) > maxRows {
			http.Error(w, fmt.Sprintf("too many rows, max is %d", maxRows), http.StatusBadRequest)
			return
		}
		state, err := loadEmailCronStateForTenant(r.Context(), tenantID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		state = applyEmailCronRequest(state, req, time.Now())
		state = normalizeEmailCronState(state)
		if state.Enabled {
			cfg, _ := loadEmailConfigForTenant(r.Context(), tenantID)
			if err := validateEmailConfigForSend(cfg); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		if err := saveEmailCronStateForTenant(r.Context(), tenantID, state); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, emailCronPublic(state, time.Now()))
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleEmailCronRun(w http.ResponseWriter, r *http.Request) {
	state, err := runEmailCronBatchForTenant(r.Context(), activeTenantIDFromRequest(r), true)
	if err != nil {
		recordEmailCronRunError(r.Context(), "manual", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	recordEmailCronResultErrors(r.Context(), state, "manual")
	writeJSON(w, http.StatusOK, emailCronPublic(state, time.Now()))
}

func applyEmailCronRequest(state emailCronState, req emailCronRequest, now time.Time) emailCronState {
	if req.Enabled != nil {
		state.Enabled = *req.Enabled
	}
	if req.DailyLimit > 0 {
		state.DailyLimit = req.DailyLimit
	}
	if strings.TrimSpace(req.SendTime) != "" {
		state.SendTime = strings.TrimSpace(req.SendTime)
	}
	if strings.TrimSpace(req.Template) != "" {
		state.Template = strings.TrimSpace(req.Template)
	}
	if req.Rows != nil {
		state.Queue = makeEmailCronQueue(*req.Rows, now)
		state.LastResults = nil
		state.LastRunAt = ""
		state.LastRunDate = ""
		state.SentToday = 0
		state.SendHistory = nil
	}
	return state
}

func recordEmailCronRunError(ctx context.Context, trigger string, err error) {
	if err == nil {
		return
	}
	recordOperationLogBestEffort(ctx, operationLogInput{
		Source:    "background_job",
		Level:     "error",
		Operation: "email.cron.run",
		Status:    "error",
		Message:   err.Error(),
		Metadata: map[string]any{
			"trigger": trigger,
		},
	})
}

func recordEmailCronResultErrors(ctx context.Context, state emailCronState, trigger string) {
	for _, result := range state.LastResults {
		if result.Status != "error" {
			continue
		}
		recordOperationLogBestEffort(ctx, operationLogInput{
			Source:    "background_job",
			Level:     "error",
			Operation: "email.cron.send",
			Status:    "error",
			Message:   result.Error,
			Metadata: map[string]any{
				"trigger":   trigger,
				"rowId":     result.ID,
				"recipient": result.Email,
				"transient": result.Transient,
			},
		})
	}
}

func runEmailCronBatch(ctx context.Context, force bool) (emailCronState, error) {
	emailCronMu.Lock()
	defer emailCronMu.Unlock()

	state, err := loadEmailCronStateUnlocked()
	if err != nil {
		return state, err
	}
	state = normalizeEmailCronState(state)
	now := time.Now()
	if !force && !emailCronDue(state, now) {
		return state, nil
	}
	if !state.Enabled && !force {
		return state, nil
	}
	queued := queuedEmailCronJobs(state)
	if len(queued) == 0 {
		state.LastRunAt = now.Format(time.RFC3339)
		state.LastRunDate = localDate(now)
		state.LastResults = nil
		_ = saveEmailCronStateUnlocked(state)
		return state, nil
	}

	cfg, err := loadEmailConfig()
	if err != nil {
		return state, err
	}
	if err := validateEmailConfigForSend(cfg); err != nil {
		return state, err
	}

	state.SendHistory = seedLegacySendHistory(state, now)
	state.SendHistory = pruneEmailSendHistory(state.SendHistory, now)
	sentBefore := len(state.SendHistory)
	remaining := state.DailyLimit - sentBefore
	if remaining <= 0 {
		return state, nil
	}

	results := make([]emailSendResult, 0, remaining)
	sent := 0
	baseURL := schedulerBaseURL(cfg)
	for _, idx := range queued {
		if sent >= remaining {
			break
		}
		job := &state.Queue[idx]
		result := sendPaymentEmailRow(ctx, cfg, job.Row, state.Template, baseURL, false)
		job.Attempts++
		job.ProcessedAt = now.Format(time.RFC3339)
		job.Status = result.Status
		job.ResendID = result.ResendID
		job.MessageID = result.MessageID
		job.Error = result.Error
		results = append(results, result)
		if result.Transient {
			job.Status = "queued"
			break
		}
		if result.Status == "sent" {
			sent++
			sleepEmailSendPace(ctx, cfg)
		}
	}

	state.LastRunAt = now.Format(time.RFC3339)
	state.LastRunDate = localDate(now)
	addEmailSentToState(&state, sent, now)
	state.LastResults = results
	state.UpdatedAt = now.Format(time.RFC3339)
	if err := saveEmailCronStateUnlocked(state); err != nil {
		return state, err
	}
	return state, nil
}

func runEmailCronBatchForTenant(ctx context.Context, tenantID string, force bool) (emailCronState, error) {
	tenantID = strings.TrimSpace(tenantID)
	if tenantID == "" {
		return runEmailCronBatch(ctx, force)
	}
	emailCronMu.Lock()
	defer emailCronMu.Unlock()

	db, err := openMasterDataDatabase(ctx)
	if err != nil {
		return emailCronState{}, err
	}
	defer db.Close()

	state, err := loadTenantEmailCronStateUnlocked(ctx, db, tenantID)
	if err != nil {
		return state, err
	}
	state = normalizeEmailCronState(state)
	now := time.Now()
	if !force && !emailCronDue(state, now) {
		return state, nil
	}
	if !state.Enabled && !force {
		return state, nil
	}
	queued := queuedEmailCronJobs(state)
	if len(queued) == 0 {
		state.LastRunAt = now.Format(time.RFC3339)
		state.LastRunDate = localDate(now)
		state.LastResults = nil
		_ = saveTenantEmailCronStateUnlocked(ctx, db, tenantID, state, "")
		return state, nil
	}

	cfg, err := loadTenantEmailConfig(ctx, db, tenantID)
	if err != nil {
		return state, err
	}
	if err := validateEmailConfigForSend(cfg); err != nil {
		return state, err
	}

	state.SendHistory = seedLegacySendHistory(state, now)
	state.SendHistory = pruneEmailSendHistory(state.SendHistory, now)
	sentBefore := len(state.SendHistory)
	remaining := state.DailyLimit - sentBefore
	if remaining <= 0 {
		return state, nil
	}

	results := make([]emailSendResult, 0, remaining)
	sent := 0
	baseURL := schedulerBaseURL(cfg)
	for _, idx := range queued {
		if sent >= remaining {
			break
		}
		job := &state.Queue[idx]
		result := sendPaymentEmailRow(ctx, cfg, job.Row, state.Template, baseURL, false)
		job.Attempts++
		job.ProcessedAt = now.Format(time.RFC3339)
		job.Status = result.Status
		job.ResendID = result.ResendID
		job.MessageID = result.MessageID
		job.Error = result.Error
		results = append(results, result)
		if result.Transient {
			job.Status = "queued"
			break
		}
		if result.Status == "sent" {
			sent++
			sleepEmailSendPace(ctx, cfg)
		}
	}

	state.LastRunAt = now.Format(time.RFC3339)
	state.LastRunDate = localDate(now)
	addEmailSentToState(&state, sent, now)
	state.LastResults = results
	state.UpdatedAt = now.Format(time.RFC3339)
	if err := saveTenantEmailCronStateUnlocked(ctx, db, tenantID, state, ""); err != nil {
		return state, err
	}
	return state, nil
}

func runDueTenantEmailCronBatches(ctx context.Context) error {
	db, err := openMasterDataDatabase(ctx)
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.QueryContext(ctx, `
SELECT tenant_id::text
FROM tenant_email_cron_states
WHERE COALESCE((state->>'enabled')::boolean, false) = true
ORDER BY updated_at ASC`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var tenantID string
		if err := rows.Scan(&tenantID); err != nil {
			return err
		}
		state, err := runEmailCronBatchForTenant(ctx, tenantID, false)
		if err != nil {
			recordEmailCronRunError(ctx, "tenant_scheduler:"+tenantID, err)
			continue
		}
		recordEmailCronResultErrors(ctx, state, "tenant_scheduler:"+tenantID)
	}
	return rows.Err()
}

func loadEmailCronState() (emailCronState, error) {
	emailCronMu.Lock()
	defer emailCronMu.Unlock()
	return loadEmailCronStateUnlocked()
}

func loadEmailCronStateForTenant(ctx context.Context, tenantID string) (emailCronState, error) {
	tenantID = strings.TrimSpace(tenantID)
	if tenantID == "" {
		return loadEmailCronState()
	}
	emailCronMu.Lock()
	defer emailCronMu.Unlock()
	db, err := openMasterDataDatabase(ctx)
	if err != nil {
		return emailCronState{}, err
	}
	defer db.Close()
	return loadTenantEmailCronStateUnlocked(ctx, db, tenantID)
}

func loadTenantEmailCronStateUnlocked(ctx context.Context, db *sql.DB, tenantID string) (emailCronState, error) {
	state := defaultEmailCronState()
	var data []byte
	err := db.QueryRowContext(ctx, `
SELECT state
FROM tenant_email_cron_states
WHERE tenant_id = $1::uuid`, strings.TrimSpace(tenantID)).Scan(&data)
	if err == sql.ErrNoRows {
		return state, nil
	}
	if err != nil {
		return state, err
	}
	if err := json.Unmarshal(data, &state); err != nil {
		return state, err
	}
	return normalizeEmailCronState(state), nil
}

func loadEmailCronStateUnlocked() (emailCronState, error) {
	state := defaultEmailCronState()
	data, err := os.ReadFile(emailCronPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return state, nil
		}
		return state, err
	}
	if err := json.Unmarshal(data, &state); err != nil {
		return state, err
	}
	return normalizeEmailCronState(state), nil
}

func saveEmailCronState(state emailCronState) error {
	emailCronMu.Lock()
	defer emailCronMu.Unlock()
	return saveEmailCronStateUnlocked(state)
}

func saveEmailCronStateForTenant(ctx context.Context, tenantID string, state emailCronState) error {
	return saveEmailCronStateForTenantByUser(ctx, tenantID, state, "")
}

func saveEmailCronStateForTenantByUser(ctx context.Context, tenantID string, state emailCronState, userID string) error {
	tenantID = strings.TrimSpace(tenantID)
	if tenantID == "" {
		return saveEmailCronState(state)
	}
	emailCronMu.Lock()
	defer emailCronMu.Unlock()
	db, err := openMasterDataDatabase(ctx)
	if err != nil {
		return err
	}
	defer db.Close()
	return saveTenantEmailCronStateUnlocked(ctx, db, tenantID, state, userID)
}

func saveTenantEmailCronStateUnlocked(ctx context.Context, db *sql.DB, tenantID string, state emailCronState, userID string) error {
	state = normalizeEmailCronState(state)
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, `
INSERT INTO tenant_email_cron_states (tenant_id, state, updated_by_user_id)
VALUES ($1::uuid, $2::jsonb, nullif($3, '')::uuid)
ON CONFLICT (tenant_id) DO UPDATE
SET state = EXCLUDED.state,
	updated_by_user_id = EXCLUDED.updated_by_user_id,
	updated_at = now()`, strings.TrimSpace(tenantID), string(data), strings.TrimSpace(userID))
	return err
}

func saveEmailCronStateUnlocked(state emailCronState) error {
	state = normalizeEmailCronState(state)
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(emailCronPath, data, 0600)
}

func defaultEmailCronState() emailCronState {
	return emailCronState{
		Enabled:    false,
		DailyLimit: defaultCronDailyLimit,
		SendTime:   defaultCronSendTime,
		Template:   "payment_due",
		Queue:      []emailCronJob{},
	}
}

func normalizeEmailCronState(state emailCronState) emailCronState {
	if state.DailyLimit <= 0 {
		state.DailyLimit = defaultCronDailyLimit
	}
	if state.DailyLimit > defaultCronDailyLimit {
		state.DailyLimit = defaultCronDailyLimit
	}
	if _, err := parseCronSendTime(state.SendTime); err != nil {
		state.SendTime = defaultCronSendTime
	}
	state.Template = strings.TrimSpace(state.Template)
	if state.Template == "" {
		state.Template = "payment_due"
	}
	for idx := range state.Queue {
		if state.Queue[idx].Status == "" {
			state.Queue[idx].Status = "queued"
		}
	}
	return state
}

func makeEmailCronQueue(rows []paymentRow, now time.Time) []emailCronJob {
	jobs := make([]emailCronJob, 0, len(rows))
	for idx, row := range rows {
		if strings.TrimSpace(row.ID) == "" {
			row.ID = fmt.Sprintf("row-%03d", idx+1)
		}
		row = cleanRow(row)
		jobs = append(jobs, emailCronJob{
			ID:        row.ID,
			Row:       row,
			Status:    "queued",
			CreatedAt: now.Format(time.RFC3339),
		})
	}
	return jobs
}

func emailCronPublic(state emailCronState, now time.Time) emailCronResponse {
	state = normalizeEmailCronState(state)
	resp := emailCronResponse{
		Enabled:     state.Enabled,
		DailyLimit:  state.DailyLimit,
		SendTime:    state.SendTime,
		Template:    state.Template,
		QueueTotal:  len(state.Queue),
		SentToday:   sentLast24hForState(state, now),
		SentLast24H: sentLast24hForState(state, now),
		LastRunAt:   state.LastRunAt,
		NextRunAt:   nextEmailCronRunAt(state, now).Format(time.RFC3339),
		LastResults: state.LastResults,
	}
	for _, job := range state.Queue {
		switch job.Status {
		case "sent":
			resp.Sent++
		case "error":
			resp.Errors++
		case "skipped":
			resp.Skipped++
		default:
			resp.Queued++
		}
	}
	return resp
}

func queuedEmailCronJobs(state emailCronState) []int {
	indexes := make([]int, 0, len(state.Queue))
	for idx, job := range state.Queue {
		if job.Status == "" || job.Status == "queued" {
			indexes = append(indexes, idx)
		}
	}
	return indexes
}

func emailCronDue(state emailCronState, now time.Time) bool {
	if !state.Enabled {
		return false
	}
	if len(queuedEmailCronJobs(state)) == 0 {
		return false
	}
	sendAt := emailCronTimeOnDate(state.SendTime, now)
	if now.Before(sendAt) {
		return false
	}
	return state.LastRunDate != localDate(now)
}

func nextEmailCronRunAt(state emailCronState, now time.Time) time.Time {
	sendAt := emailCronTimeOnDate(state.SendTime, now)
	if !now.Before(sendAt) || state.LastRunDate == localDate(now) {
		sendAt = sendAt.Add(24 * time.Hour)
	}
	return sendAt
}

func emailCronTimeOnDate(sendTime string, now time.Time) time.Time {
	parsed, err := parseCronSendTime(sendTime)
	if err != nil {
		parsed, _ = parseCronSendTime(defaultCronSendTime)
	}
	return time.Date(now.Year(), now.Month(), now.Day(), parsed.Hour(), parsed.Minute(), 0, 0, now.Location())
}

func parseCronSendTime(value string) (time.Time, error) {
	return time.Parse("15:04", strings.TrimSpace(value))
}

func sentTodayForState(state emailCronState, now time.Time) int {
	return sentLast24hForState(state, now)
}

func recordEmailCronSent(count int, now time.Time) {
	recordEmailCronSentForTenant(context.Background(), "", count, now)
}

func recordEmailCronSentForTenant(ctx context.Context, tenantID string, count int, now time.Time) {
	if count <= 0 {
		return
	}
	tenantID = strings.TrimSpace(tenantID)
	emailCronMu.Lock()
	defer emailCronMu.Unlock()
	if tenantID != "" {
		db, err := openMasterDataDatabase(ctx)
		if err != nil {
			return
		}
		defer db.Close()
		state, err := loadTenantEmailCronStateUnlocked(ctx, db, tenantID)
		if err != nil {
			return
		}
		addEmailSentToState(&state, count, now)
		state.LastRunDate = localDate(now)
		state.UpdatedAt = now.Format(time.RFC3339)
		_ = saveTenantEmailCronStateUnlocked(ctx, db, tenantID, state, "")
		return
	}
	state, err := loadEmailCronStateUnlocked()
	if err != nil {
		return
	}
	addEmailSentToState(&state, count, now)
	state.LastRunDate = localDate(now)
	state.UpdatedAt = now.Format(time.RFC3339)
	_ = saveEmailCronStateUnlocked(state)
}

func remainingEmailSendQuota(now time.Time) (int, error) {
	status, err := emailSendQuotaStatus(now)
	if err != nil {
		return 0, err
	}
	return status.Remaining, nil
}

func emailSendQuotaStatus(now time.Time) (emailQuotaStatus, error) {
	return emailSendQuotaStatusForTenant(context.Background(), "", now)
}

func emailSendQuotaStatusForTenant(ctx context.Context, tenantID string, now time.Time) (emailQuotaStatus, error) {
	tenantID = strings.TrimSpace(tenantID)
	if tenantID != "" {
		state, err := loadEmailCronStateForTenant(ctx, tenantID)
		if err != nil {
			return emailQuotaStatus{}, err
		}
		state = normalizeEmailCronState(state)
		sent := sentLast24hForState(state, now)
		remaining := state.DailyLimit - sent
		if remaining < 0 {
			remaining = 0
		}
		return emailQuotaStatus{Limit: state.DailyLimit, Sent: sent, Remaining: remaining}, nil
	}
	state, err := loadEmailCronState()
	if err != nil {
		return emailQuotaStatus{}, err
	}
	state = normalizeEmailCronState(state)
	sent := sentLast24hForState(state, now)
	remaining := state.DailyLimit - sent
	if remaining < 0 {
		remaining = 0
	}
	return emailQuotaStatus{Limit: state.DailyLimit, Sent: sent, Remaining: remaining}, nil
}

func addEmailSentToState(state *emailCronState, count int, now time.Time) {
	state.SendHistory = seedLegacySendHistory(*state, now)
	state.SendHistory = pruneEmailSendHistory(state.SendHistory, now)
	for i := 0; i < count; i++ {
		state.SendHistory = append(state.SendHistory, now.Format(time.RFC3339))
	}
	state.SendHistory = pruneEmailSendHistory(state.SendHistory, now)
	state.SentToday = len(state.SendHistory)
}

func sentLast24hForState(state emailCronState, now time.Time) int {
	history := pruneEmailSendHistory(state.SendHistory, now)
	if len(history) > 0 {
		return len(history)
	}
	if state.LastRunDate == localDate(now) {
		return state.SentToday
	}
	return 0
}

func seedLegacySendHistory(state emailCronState, now time.Time) []string {
	if len(state.SendHistory) > 0 || state.LastRunDate != localDate(now) || state.SentToday <= 0 {
		return state.SendHistory
	}
	history := make([]string, 0, state.SentToday)
	for i := 0; i < state.SentToday; i++ {
		history = append(history, now.Format(time.RFC3339))
	}
	return history
}

func pruneEmailSendHistory(history []string, now time.Time) []string {
	if len(history) == 0 {
		return nil
	}
	cutoff := now.Add(-emailQuotaWindow)
	out := make([]string, 0, len(history))
	for _, raw := range history {
		t, err := time.Parse(time.RFC3339, raw)
		if err != nil {
			continue
		}
		if !t.Before(cutoff) {
			out = append(out, t.Format(time.RFC3339))
		}
	}
	return out
}

func localDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func schedulerBaseURL(cfg emailConfig) string {
	if cfg.PublicBaseURL != "" {
		return cfg.PublicBaseURL
	}
	return "http://localhost:" + defaultHTTPPort
}
