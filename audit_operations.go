package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	adminActorHeader       = "X-ABC-Admin-Actor"
	adminActorUserIDHeader = "X-ABC-Admin-User-ID"
	requestIDHeader        = "X-Request-ID"

	defaultAuditLogLimit     = 100
	defaultOperationLogLimit = 100
	maxAuditLogLimit         = 500
	maxOperationLogLimit     = 500
)

type requestAuditContext struct {
	TenantID    string
	ActorUserID string
	ActorName   string
	RequestID   string
	IPAddress   string
	UserAgent   string
}

type auditLogInput struct {
	Context    requestAuditContext
	Action     string
	EntityType string
	EntityID   string
	Reason     string
	Metadata   map[string]any
}

type auditLogSummary struct {
	ID          string         `json:"id"`
	OccurredAt  time.Time      `json:"occurredAt"`
	TenantCode  string         `json:"tenantCode,omitempty"`
	TenantName  string         `json:"tenantName,omitempty"`
	ActorUserID string         `json:"actorUserId,omitempty"`
	ActorName   string         `json:"actorName,omitempty"`
	Action      string         `json:"action"`
	EntityType  string         `json:"entityType"`
	EntityID    string         `json:"entityId,omitempty"`
	RequestID   string         `json:"requestId,omitempty"`
	IPAddress   string         `json:"ipAddress,omitempty"`
	UserAgent   string         `json:"userAgent,omitempty"`
	Reason      string         `json:"reason,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

type auditLogCommandSummary struct {
	TotalCount       int `json:"totalCount"`
	MoneyActionCount int `json:"moneyActionCount"`
	FeeActionCount   int `json:"feeActionCount"`
	UserActionCount  int `json:"userActionCount"`
}

type auditLogFilters struct {
	TenantID   string
	Action     string
	EntityType string
	EntityID   string
	Limit      int
}

type operationLogInput struct {
	TenantID   string
	RequestID  string
	Source     string
	Level      string
	Operation  string
	Status     string
	Message    string
	EntityType string
	EntityID   string
	Metadata   map[string]any
}

type operationLogSummary struct {
	ID         string         `json:"id"`
	OccurredAt time.Time      `json:"occurredAt"`
	TenantCode string         `json:"tenantCode,omitempty"`
	TenantName string         `json:"tenantName,omitempty"`
	Source     string         `json:"source"`
	Level      string         `json:"level"`
	Operation  string         `json:"operation"`
	Status     string         `json:"status"`
	Message    string         `json:"message,omitempty"`
	EntityType string         `json:"entityType,omitempty"`
	EntityID   string         `json:"entityId,omitempty"`
	RequestID  string         `json:"requestId,omitempty"`
	Metadata   map[string]any `json:"metadata,omitempty"`
}

type operationLogCommandSummary struct {
	TotalCount              int `json:"totalCount"`
	ErrorCount              int `json:"errorCount"`
	WarnCount               int `json:"warnCount"`
	InfoCount               int `json:"infoCount"`
	WebhookErrorCount       int `json:"webhookErrorCount"`
	EmailErrorCount         int `json:"emailErrorCount"`
	CronErrorCount          int `json:"cronErrorCount"`
	BackgroundJobErrorCount int `json:"backgroundJobErrorCount"`
}

type operationLogFilters struct {
	TenantID   string
	Source     string
	Level      string
	Operation  string
	Status     string
	EntityType string
	EntityID   string
	Limit      int
}

func auditContextFromRequest(r *http.Request) requestAuditContext {
	actorUserID := strings.TrimSpace(r.Header.Get(adminActorUserIDHeader))
	actorName := strings.TrimSpace(r.Header.Get(adminActorHeader))
	if user, ok := authenticatedUserFromRequest(r); ok {
		actorUserID = firstNonEmpty(user.ID, actorUserID)
		actorName = firstNonEmpty(user.DisplayName, user.Email, actorName)
	}
	return requestAuditContext{
		TenantID:    activeTenantIDFromRequest(r),
		ActorUserID: actorUserID,
		ActorName:   actorName,
		RequestID:   strings.TrimSpace(r.Header.Get(requestIDHeader)),
		IPAddress:   requestIPAddress(r),
		UserAgent:   strings.TrimSpace(r.UserAgent()),
	}
}

func requestIPAddress(r *http.Request) string {
	for _, header := range []string{"X-Forwarded-For", "X-Real-IP"} {
		value := strings.TrimSpace(r.Header.Get(header))
		if value == "" {
			continue
		}
		first := strings.TrimSpace(strings.Split(value, ",")[0])
		if net.ParseIP(first) != nil {
			return first
		}
	}
	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil && net.ParseIP(host) != nil {
		return host
	}
	if net.ParseIP(strings.TrimSpace(r.RemoteAddr)) != nil {
		return strings.TrimSpace(r.RemoteAddr)
	}
	return ""
}

func insertAuditLog(ctx context.Context, exec masterDataExecutor, input auditLogInput) error {
	action := strings.TrimSpace(input.Action)
	entityType := strings.TrimSpace(input.EntityType)
	if action == "" || entityType == "" {
		return fmt.Errorf("audit action and entity type are required")
	}
	metadata := cloneMetadata(input.Metadata)
	reason := strings.TrimSpace(input.Reason)
	if reason != "" {
		metadata["reason"] = reason
	}
	if input.Context.ActorName != "" {
		metadata["actorName"] = input.Context.ActorName
	}
	metadataJSON, err := jsonObjectString(metadata)
	if err != nil {
		return err
	}
	_, err = exec.ExecContext(ctx, `
INSERT INTO audit_logs (tenant_id, actor_user_id, action, entity_type, entity_id, request_id, ip_address, user_agent, metadata)
VALUES (COALESCE(nullif($1, '')::uuid, (SELECT id FROM tenants WHERE code = $2)), nullif($3, '')::uuid, $4, $5, nullif($6, '')::uuid, $7, nullif($8, '')::inet, $9, $10::jsonb)`,
		strings.TrimSpace(input.Context.TenantID),
		defaultTenantCode,
		strings.TrimSpace(input.Context.ActorUserID),
		action,
		entityType,
		strings.TrimSpace(input.EntityID),
		strings.TrimSpace(input.Context.RequestID),
		strings.TrimSpace(input.Context.IPAddress),
		strings.TrimSpace(input.Context.UserAgent),
		metadataJSON,
	)
	return err
}

func recordOperationLog(ctx context.Context, exec masterDataExecutor, input operationLogInput) error {
	source := headerKey(input.Source)
	level := headerKey(input.Level)
	if level == "" {
		level = "error"
	}
	if level != "info" && level != "warn" && level != "error" {
		level = "error"
	}
	operation := strings.TrimSpace(input.Operation)
	status := headerKey(input.Status)
	if status == "" {
		status = "error"
	}
	if source == "" || operation == "" {
		return fmt.Errorf("operation log source and operation are required")
	}
	metadataJSON, err := jsonObjectString(cloneMetadata(input.Metadata))
	if err != nil {
		return err
	}
	_, err = exec.ExecContext(ctx, `
INSERT INTO operation_logs (tenant_id, source, level, operation, status, message, entity_type, entity_id, request_id, metadata)
VALUES (COALESCE(nullif($1, '')::uuid, (SELECT id FROM tenants WHERE code = $2)), $3, $4, $5, $6, $7, $8, nullif($9, '')::uuid, $10, $11::jsonb)`,
		strings.TrimSpace(input.TenantID),
		defaultTenantCode,
		source,
		level,
		operation,
		status,
		strings.TrimSpace(input.Message),
		strings.TrimSpace(input.EntityType),
		strings.TrimSpace(input.EntityID),
		strings.TrimSpace(input.RequestID),
		metadataJSON,
	)
	return err
}

func recordOperationLogBestEffort(ctx context.Context, input operationLogInput) {
	cfg, err := loadDatabaseConfig()
	if err != nil || strings.TrimSpace(cfg.URL) == "" {
		return
	}
	db, err := openConfiguredDatabase(ctx, cfg)
	if err != nil {
		return
	}
	defer db.Close()
	_ = recordOperationLog(ctx, db, input)
}

func cloneMetadata(metadata map[string]any) map[string]any {
	out := map[string]any{}
	for key, value := range metadata {
		if strings.TrimSpace(key) == "" || value == nil {
			continue
		}
		out[key] = value
	}
	return out
}

func handleAdminAuditLogs(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := resolveOperationsTenantScope(w, r, "audit_log.cross_tenant_view")
	if !ok {
		return
	}
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	logs, err := listAuditLogs(r.Context(), db, auditLogFilters{
		TenantID:   tenantID,
		Action:     strings.TrimSpace(r.URL.Query().Get("action")),
		EntityType: strings.TrimSpace(r.URL.Query().Get("entityType")),
		EntityID:   strings.TrimSpace(r.URL.Query().Get("entityId")),
		Limit:      parsePositiveInt(r.URL.Query().Get("limit"), defaultAuditLogLimit),
	})
	if err != nil {
		http.Error(w, "cannot load audit logs", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"logs": logs, "summary": buildAuditLogCommandSummary(logs)})
}

func handleAdminOperationLogs(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := resolveOperationsTenantScope(w, r, "operation_log.cross_tenant_view")
	if !ok {
		return
	}
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	logs, err := listOperationLogs(r.Context(), db, operationLogFilters{
		TenantID:   tenantID,
		Source:     headerKey(r.URL.Query().Get("source")),
		Level:      headerKey(r.URL.Query().Get("level")),
		Operation:  strings.TrimSpace(r.URL.Query().Get("operation")),
		Status:     headerKey(r.URL.Query().Get("status")),
		EntityType: strings.TrimSpace(r.URL.Query().Get("entityType")),
		EntityID:   strings.TrimSpace(r.URL.Query().Get("entityId")),
		Limit:      parsePositiveInt(r.URL.Query().Get("limit"), defaultOperationLogLimit),
	})
	if err != nil {
		http.Error(w, "cannot load operation logs", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"logs": logs, "summary": buildOperationLogCommandSummary(logs)})
}

func resolveOperationsTenantScope(w http.ResponseWriter, r *http.Request, crossTenantPermission string) (string, bool) {
	activeTenantID, ok := requireActiveTenantID(w, r)
	if !ok {
		return "", false
	}
	requestedTenantID := strings.TrimSpace(r.URL.Query().Get("tenantId"))
	if requestedTenantID == "" || requestedTenantID == activeTenantID {
		return activeTenantID, true
	}
	user, ok := authenticatedUserFromRequest(r)
	if !ok || !authenticatedUserHasPermission(user, crossTenantPermission) {
		http.Error(w, "missing required API permission: "+crossTenantPermission, http.StatusForbidden)
		return "", false
	}
	if requestedTenantID == "all" {
		return "", true
	}
	return requestedTenantID, true
}

func listAuditLogs(ctx context.Context, db *sql.DB, filters auditLogFilters) ([]auditLogSummary, error) {
	conditions := []string{"true"}
	args := []any{}
	addArg := func(value any) string {
		args = append(args, value)
		return fmt.Sprintf("$%d", len(args))
	}
	if filters.Action != "" {
		conditions = append(conditions, "action = "+addArg(filters.Action))
	}
	if filters.TenantID != "" {
		conditions = append(conditions, "tenant_id = "+addArg(filters.TenantID)+"::uuid")
	}
	if filters.EntityType != "" {
		conditions = append(conditions, "entity_type = "+addArg(filters.EntityType))
	}
	if filters.EntityID != "" {
		conditions = append(conditions, "entity_id = "+addArg(filters.EntityID)+"::uuid")
	}
	limit := filters.Limit
	if limit <= 0 {
		limit = defaultAuditLogLimit
	}
	if limit > maxAuditLogLimit {
		limit = maxAuditLogLimit
	}
	limitArg := addArg(limit)

	rows, err := db.QueryContext(ctx, `
SELECT id::text,
	occurred_at,
	COALESCE(tenant.code, ''),
	COALESCE(tenant.name, ''),
	COALESCE(actor_user_id::text, ''),
	COALESCE(metadata->>'actorName', ''),
	action,
	entity_type,
	COALESCE(entity_id::text, ''),
	request_id,
	COALESCE(ip_address::text, ''),
	user_agent,
	COALESCE(metadata->>'reason', ''),
	metadata
FROM audit_logs
LEFT JOIN tenants tenant ON tenant.id = audit_logs.tenant_id
WHERE `+strings.Join(conditions, " AND ")+`
ORDER BY occurred_at DESC
LIMIT `+limitArg, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := []auditLogSummary{}
	for rows.Next() {
		var item auditLogSummary
		var metadataBytes []byte
		if err := rows.Scan(
			&item.ID,
			&item.OccurredAt,
			&item.TenantCode,
			&item.TenantName,
			&item.ActorUserID,
			&item.ActorName,
			&item.Action,
			&item.EntityType,
			&item.EntityID,
			&item.RequestID,
			&item.IPAddress,
			&item.UserAgent,
			&item.Reason,
			&metadataBytes,
		); err != nil {
			return nil, err
		}
		item.Metadata = sanitizeLogMetadata(decodeMetadata(metadataBytes))
		logs = append(logs, item)
	}
	return logs, rows.Err()
}

func listOperationLogs(ctx context.Context, db *sql.DB, filters operationLogFilters) ([]operationLogSummary, error) {
	conditions := []string{"true"}
	args := []any{}
	addArg := func(value any) string {
		args = append(args, value)
		return fmt.Sprintf("$%d", len(args))
	}
	if filters.Source != "" {
		conditions = append(conditions, "source = "+addArg(filters.Source))
	}
	if filters.TenantID != "" {
		conditions = append(conditions, "tenant_id = "+addArg(filters.TenantID)+"::uuid")
	}
	if filters.Level != "" {
		conditions = append(conditions, "level = "+addArg(filters.Level))
	}
	if filters.Operation != "" {
		conditions = append(conditions, "operation = "+addArg(filters.Operation))
	}
	if filters.Status != "" {
		conditions = append(conditions, "status = "+addArg(filters.Status))
	}
	if filters.EntityType != "" {
		conditions = append(conditions, "entity_type = "+addArg(filters.EntityType))
	}
	if filters.EntityID != "" {
		conditions = append(conditions, "entity_id = "+addArg(filters.EntityID)+"::uuid")
	}
	limit := filters.Limit
	if limit <= 0 {
		limit = defaultOperationLogLimit
	}
	if limit > maxOperationLogLimit {
		limit = maxOperationLogLimit
	}
	limitArg := addArg(limit)

	rows, err := db.QueryContext(ctx, `
SELECT id::text,
	occurred_at,
	COALESCE(tenant.code, ''),
	COALESCE(tenant.name, ''),
	source,
	level,
	operation,
	status,
	message,
	entity_type,
	COALESCE(entity_id::text, ''),
	request_id,
	metadata
FROM operation_logs
LEFT JOIN tenants tenant ON tenant.id = operation_logs.tenant_id
WHERE `+strings.Join(conditions, " AND ")+`
ORDER BY occurred_at DESC
LIMIT `+limitArg, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	logs := []operationLogSummary{}
	for rows.Next() {
		var item operationLogSummary
		var metadataBytes []byte
		if err := rows.Scan(
			&item.ID,
			&item.OccurredAt,
			&item.TenantCode,
			&item.TenantName,
			&item.Source,
			&item.Level,
			&item.Operation,
			&item.Status,
			&item.Message,
			&item.EntityType,
			&item.EntityID,
			&item.RequestID,
			&metadataBytes,
		); err != nil {
			return nil, err
		}
		item.Metadata = sanitizeLogMetadata(decodeMetadata(metadataBytes))
		logs = append(logs, item)
	}
	return logs, rows.Err()
}

func buildOperationLogCommandSummary(logs []operationLogSummary) operationLogCommandSummary {
	var summary operationLogCommandSummary
	for _, log := range logs {
		summary.TotalCount++
		switch log.Level {
		case "error":
			summary.ErrorCount++
			switch {
			case log.Source == "webhook":
				summary.WebhookErrorCount++
			case log.Source == "email":
				summary.EmailErrorCount++
			case log.Source == "background_job" && strings.HasPrefix(log.Operation, "email.cron."):
				summary.CronErrorCount++
			case log.Source == "background_job":
				summary.BackgroundJobErrorCount++
			}
		case "warn":
			summary.WarnCount++
		case "info":
			summary.InfoCount++
		}
	}
	return summary
}

func buildAuditLogCommandSummary(logs []auditLogSummary) auditLogCommandSummary {
	var summary auditLogCommandSummary
	for _, log := range logs {
		summary.TotalCount++
		switch log.EntityType {
		case "manual_cash_receipt", "payment_transaction", "reconciliation_match":
			summary.MoneyActionCount++
		case "fee_schedule", "student_fee_adjustment":
			summary.FeeActionCount++
		case "app_user", "app_role":
			summary.UserActionCount++
		}
	}
	return summary
}

func sanitizeLogMetadata(metadata map[string]any) map[string]any {
	if len(metadata) == 0 {
		return map[string]any{}
	}
	sanitized := make(map[string]any, len(metadata))
	for key, value := range metadata {
		if isSensitiveLogMetadataKey(key) {
			sanitized[key] = "[redacted]"
			continue
		}
		sanitized[key] = sanitizeLogMetadataValue(value)
	}
	return sanitized
}

func sanitizeLogMetadataValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		return sanitizeLogMetadata(typed)
	case []any:
		items := make([]any, 0, len(typed))
		for _, item := range typed {
			items = append(items, sanitizeLogMetadataValue(item))
		}
		return items
	default:
		return value
	}
}

func isSensitiveLogMetadataKey(key string) bool {
	normalized := strings.ToLower(strings.NewReplacer("-", "", "_", "", " ", "").Replace(key))
	for _, token := range []string{"password", "secret", "token", "apikey", "authorization", "credential", "checksum", "privatekey", "rawpayload"} {
		if strings.Contains(normalized, token) {
			return true
		}
	}
	return false
}

func decodeMetadata(data []byte) map[string]any {
	if len(data) == 0 {
		return map[string]any{}
	}
	var metadata map[string]any
	if err := json.Unmarshal(data, &metadata); err != nil || metadata == nil {
		return map[string]any{}
	}
	return metadata
}
