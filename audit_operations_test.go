package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResolveOperationsTenantScopeAllowsPlatformOnlyAllScope(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/operation-logs?tenantId=all", nil)
	req = req.WithContext(context.WithValue(req.Context(), authenticatedUserContextKey, authenticatedUser{
		ID:              "platform-user",
		IsPlatformAdmin: true,
		PermissionSet: map[string]bool{
			"operation_log.cross_tenant_view": true,
		},
	}))
	rec := httptest.NewRecorder()

	got, ok := resolveOperationsTenantScope(rec, req, "operation_log.cross_tenant_view")
	if !ok {
		t.Fatalf("expected platform-only operation scope resolution to pass, got %d %s", rec.Code, rec.Body.String())
	}
	if got != "" {
		t.Fatalf("expected all-tenant scope to resolve as empty tenant id, got %q", got)
	}
}

func TestResolveOperationsTenantScopeUsesActiveTenantWhenPresent(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/operation-logs", nil)
	req = req.WithContext(context.WithValue(req.Context(), authenticatedUserContextKey, authenticatedUser{
		ID:           "tenant-user",
		ActiveTenant: authTenantSummary{ID: "tenant-1", Status: "active", MembershipStatus: "active"},
		PermissionSet: map[string]bool{
			"operation_log.view": true,
		},
	}))
	rec := httptest.NewRecorder()

	got, ok := resolveOperationsTenantScope(rec, req, "operation_log.cross_tenant_view")
	if !ok {
		t.Fatalf("expected active tenant scope resolution to pass, got %d %s", rec.Code, rec.Body.String())
	}
	if got != "tenant-1" {
		t.Fatalf("expected tenant-1 scope, got %q", got)
	}
}
