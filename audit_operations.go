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

type auditLogFilters struct {
	Action     string
	EntityType string
	Limit      int
}

type operationLogInput struct {
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

type operationLogFilters struct {
	Source string
	Level  string
	Limit  int
}

func auditContextFromRequest(r *http.Request) requestAuditContext {
	actorUserID := strings.TrimSpace(r.Header.Get(adminActorUserIDHeader))
	actorName := strings.TrimSpace(r.Header.Get(adminActorHeader))
	if user, ok := authenticatedUserFromRequest(r); ok {
		actorUserID = firstNonEmpty(user.ID, actorUserID)
		actorName = firstNonEmpty(user.DisplayName, user.Email, actorName)
	}
	return requestAuditContext{
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
INSERT INTO audit_logs (actor_user_id, action, entity_type, entity_id, request_id, ip_address, user_agent, metadata)
VALUES (nullif($1, '')::uuid, $2, $3, nullif($4, '')::uuid, $5, nullif($6, '')::inet, $7, $8::jsonb)`,
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
INSERT INTO operation_logs (source, level, operation, status, message, entity_type, entity_id, request_id, metadata)
VALUES ($1, $2, $3, $4, $5, $6, nullif($7, '')::uuid, $8, $9::jsonb)`,
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
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	logs, err := listAuditLogs(r.Context(), db, auditLogFilters{
		Action:     strings.TrimSpace(r.URL.Query().Get("action")),
		EntityType: strings.TrimSpace(r.URL.Query().Get("entityType")),
		Limit:      parsePositiveInt(r.URL.Query().Get("limit"), defaultAuditLogLimit),
	})
	if err != nil {
		http.Error(w, "cannot load audit logs", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"logs": logs})
}

func handleAdminOperationLogs(w http.ResponseWriter, r *http.Request) {
	db, err := openMasterDataDatabase(r.Context())
	if err != nil {
		writeMasterDataDBError(w, err)
		return
	}
	defer db.Close()

	logs, err := listOperationLogs(r.Context(), db, operationLogFilters{
		Source: headerKey(r.URL.Query().Get("source")),
		Level:  headerKey(r.URL.Query().Get("level")),
		Limit:  parsePositiveInt(r.URL.Query().Get("limit"), defaultOperationLogLimit),
	})
	if err != nil {
		http.Error(w, "cannot load operation logs", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"logs": logs})
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
	if filters.EntityType != "" {
		conditions = append(conditions, "entity_type = "+addArg(filters.EntityType))
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
		item.Metadata = decodeMetadata(metadataBytes)
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
	if filters.Level != "" {
		conditions = append(conditions, "level = "+addArg(filters.Level))
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
		item.Metadata = decodeMetadata(metadataBytes)
		logs = append(logs, item)
	}
	return logs, rows.Err()
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
