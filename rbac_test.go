package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequirePermissionRejectsHeaderSpoofing(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users/save", nil)
	req.Header.Set("X-ABC-Admin-Permission", "system.users.write")
	req = req.WithContext(context.WithValue(req.Context(), authenticatedUserContextKey, authenticatedUser{
		ID:            "user-1",
		Email:         "viewer@example.edu.vn",
		ActiveTenant:  authTenantSummary{ID: "tenant-1", Status: "active", MembershipStatus: "active"},
		PermissionSet: map[string]bool{"system.users.read": true},
	}))
	rec := httptest.NewRecorder()

	handler := requirePermissionForAuthenticated("system.users.write", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	handler(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected header spoof to be forbidden, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestRequirePermissionAllowsAuthenticatedUserPermission(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/reports", nil)
	req = req.WithContext(context.WithValue(req.Context(), authenticatedUserContextKey, authenticatedUser{
		ID:            "user-1",
		Email:         "billing@example.edu.vn",
		ActiveTenant:  authTenantSummary{ID: "tenant-1", Status: "active", MembershipStatus: "active"},
		PermissionSet: map[string]bool{"admin.reports.read": true},
	}))
	rec := httptest.NewRecorder()

	handler := requirePermissionForAuthenticated("admin.reports.read", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	handler(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected request to pass, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestRequirePermissionRejectsMissingActiveTenant(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/reports", nil)
	req = req.WithContext(context.WithValue(req.Context(), authenticatedUserContextKey, authenticatedUser{
		ID:            "user-1",
		Email:         "billing@example.edu.vn",
		PermissionSet: map[string]bool{"admin.reports.read": true},
	}))
	rec := httptest.NewRecorder()

	handler := requirePermissionForAuthenticated("admin.reports.read", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	handler(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected missing active tenant to be forbidden, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestAppAPIRoutesHaveExplicitRBACClassification(t *testing.T) {
	authOnly := map[string]bool{}
	for _, route := range appAPIRoutes() {
		if route.Handler == nil {
			t.Fatalf("route %s %s has no handler", route.Method, route.Path)
		}
		if !strings.HasPrefix(route.Path, "/api/v1/") {
			t.Fatalf("route %s %s is outside the API namespace", route.Method, route.Path)
		}
		if route.Public {
			continue
		}
		if route.Permission == permissionAuthenticated && route.PermissionResolver == nil {
			authOnly[route.Path] = true
			continue
		}
		if route.Permission == "" && route.PermissionResolver == nil {
			t.Fatalf("protected route %s %s has no permission or resolver", route.Method, route.Path)
		}
	}
	if len(authOnly) != 1 || !authOnly["/api/v1/banks"] {
		t.Fatalf("expected only banks to be authenticated without a named permission, got %+v", authOnly)
	}
}

func TestAppAPIRoutePermissionMapCoversSensitiveRoutes(t *testing.T) {
	routes := map[string]appAPIRoute{}
	for _, route := range appAPIRoutes() {
		routes[route.Method+" "+route.Path] = route
	}
	cases := map[string]string{
		"POST /api/v1/auth/tenant/switch":                   "tenant.switch",
		"GET /api/v1/tenants":                               "tenant.view",
		"GET /api/v1/subscriptions/plans":                   "subscription.view",
		"GET /api/v1/subscriptions/invoices":                "subscription.view",
		"POST /api/v1/tenants/subscription/save":            "subscription.update",
		"POST /api/v1/subscriptions/invoices/generate":      "subscription.update",
		"POST /api/v1/subscriptions/invoices/mark-paid":     "subscription.update",
		"POST /api/v1/subscriptions/dunning/run":            "subscription.update",
		"POST /api/v1/master-data/import/csv":               "student.create",
		"POST /api/v1/master-data/students/save":            "student.update",
		"GET /api/v1/school-tree":                           "school_tree.view",
		"POST /api/v1/school-tree/classes/save":             "school_tree.update",
		"POST /api/v1/fee-schedules/save":                   "fee.update",
		"POST /api/v1/invoices/generate":                    "invoice.create",
		"POST /api/v1/payments/cash-receipts":               "payment.create",
		"POST /api/v1/notifications/campaigns/send":         "notification.send",
		"POST /api/v1/notifications/paid-confirmation/send": "notification.send",
		"GET /api/v1/admin/reports/export":                  "report.export",
		"GET /api/v1/admin/operation-logs":                  "operation_log.view",
		"POST /api/v1/admin/users/roles":                    "user.assign_role",
		"POST /api/v1/vietqr/batch":                         "payment.create",
		"POST /api/v1/email/preview":                        "notification.send",
		"POST /api/v1/email/cron/run":                       "email_cron.update",
	}
	for key, want := range cases {
		route, ok := routes[key]
		if !ok {
			t.Fatalf("missing route %s", key)
		}
		if route.Permission != want {
			t.Fatalf("route %s permission = %q, want %q", key, route.Permission, want)
		}
	}
	if routes["POST /api/v1/admin/users/save"].PermissionResolver == nil {
		t.Fatal("admin user save must resolve create/update permission from request body")
	}
	if routes["POST /api/v1/tenants/save"].PermissionResolver == nil {
		t.Fatal("tenant save must resolve create/update permission from request body")
	}
	if !routes["POST /api/v1/payments/webhooks/"].Public {
		t.Fatal("payment webhook endpoint must remain public for providers")
	}
	if !routes["GET /api/v1/qr.png"].Public {
		t.Fatal("QR PNG endpoint must remain public for scan links")
	}
	if !routes["GET /api/v1/healthz"].Public {
		t.Fatal("health endpoint must remain public for Docker healthchecks")
	}
}

func TestDynamicPermissionResolvers(t *testing.T) {
	masterReq := httptest.NewRequest(http.MethodPost, "/api/v1/import/fields?target=master_data", nil)
	if got, err := importFieldsPermission(masterReq); err != nil || got != "student.create" {
		t.Fatalf("expected master import fields permission, got %q, %v", got, err)
	}
	paymentReq := httptest.NewRequest(http.MethodPost, "/api/v1/import/fields?target=payments", nil)
	if got, err := importFieldsPermission(paymentReq); err != nil || got != "payment.create" {
		t.Fatalf("expected payment import fields permission, got %q, %v", got, err)
	}
	getConfigReq := httptest.NewRequest(http.MethodGet, "/api/v1/email/config", nil)
	if got, err := emailConfigPermission(getConfigReq); err != nil || got != "email_config.view" {
		t.Fatalf("expected email config read permission, got %q, %v", got, err)
	}
	postConfigReq := httptest.NewRequest(http.MethodPost, "/api/v1/email/config", nil)
	if got, err := emailConfigPermission(postConfigReq); err != nil || got != "email_config.update" {
		t.Fatalf("expected email config write permission, got %q, %v", got, err)
	}
	deleteConfigReq := httptest.NewRequest(http.MethodDelete, "/api/v1/email/config", nil)
	if got, err := emailConfigPermission(deleteConfigReq); err != nil || got != permissionAuthenticated {
		t.Fatalf("expected unsupported email config method to reach handler auth-only, got %q, %v", got, err)
	}
	createUserReq := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users/save", strings.NewReader(`{"email":"a@example.edu.vn"}`))
	if got, err := adminUserSavePermission(createUserReq); err != nil || got != "user.create" {
		t.Fatalf("expected user create permission, got %q, %v", got, err)
	}
	updateUserReq := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users/save", strings.NewReader(`{"id":"user-id","email":"a@example.edu.vn"}`))
	if got, err := adminUserSavePermission(updateUserReq); err != nil || got != "user.update" {
		t.Fatalf("expected user update permission, got %q, %v", got, err)
	}
	createTenantReq := httptest.NewRequest(http.MethodPost, "/api/v1/tenants/save", strings.NewReader(`{"code":"SCHOOL_B","name":"School B"}`))
	if got, err := tenantSavePermission(createTenantReq); err != nil || got != "tenant.create" {
		t.Fatalf("expected tenant create permission, got %q, %v", got, err)
	}
	updateTenantReq := httptest.NewRequest(http.MethodPost, "/api/v1/tenants/save", strings.NewReader(`{"id":"tenant-id","code":"SCHOOL_B","name":"School B"}`))
	if got, err := tenantSavePermission(updateTenantReq); err != nil || got != "tenant.update" {
		t.Fatalf("expected tenant update permission, got %q, %v", got, err)
	}
	getCronReq := httptest.NewRequest(http.MethodGet, "/api/v1/email/cron", nil)
	if got, err := emailCronPermission(getCronReq); err != nil || got != "email_cron.view" {
		t.Fatalf("expected email cron view permission, got %q, %v", got, err)
	}
	postCronReq := httptest.NewRequest(http.MethodPost, "/api/v1/email/cron", nil)
	if got, err := emailCronPermission(postCronReq); err != nil || got != "email_cron.update" {
		t.Fatalf("expected email cron update permission, got %q, %v", got, err)
	}
}

func TestRequirePermissionAllowsCrossTenantOperationAlias(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/operation-logs?tenantId=all", nil)
	req = req.WithContext(context.WithValue(req.Context(), authenticatedUserContextKey, authenticatedUser{
		ID:            "user-1",
		Email:         "ops@example.edu.vn",
		ActiveTenant:  authTenantSummary{ID: "tenant-1", Status: "active", MembershipStatus: "active"},
		PermissionSet: map[string]bool{"operations.cross_tenant.read": true},
	}))
	rec := httptest.NewRecorder()

	handler := requirePermissionForAuthenticated("operation_log.cross_tenant_view", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	handler(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected cross-tenant alias permission to pass, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestRequirePermissionBlocksWriteWhenTenantSubscriptionIsSuspended(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/invoices/generate", nil)
	req = req.WithContext(context.WithValue(req.Context(), authenticatedUserContextKey, authenticatedUser{
		ID:    "user-1",
		Email: "billing@example.edu.vn",
		ActiveTenant: authTenantSummary{
			ID:                 "tenant-1",
			Status:             "active",
			MembershipStatus:   "active",
			SubscriptionStatus: "suspended",
		},
		PermissionSet: map[string]bool{"invoice.create": true},
	}))
	rec := httptest.NewRecorder()

	handler := requirePermissionForAuthenticated("invoice.create", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	handler(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected suspended subscription write to be forbidden, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestRequirePermissionAllowsReadWhenTenantSubscriptionIsSuspended(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/invoices", nil)
	req = req.WithContext(context.WithValue(req.Context(), authenticatedUserContextKey, authenticatedUser{
		ID:    "user-1",
		Email: "billing@example.edu.vn",
		ActiveTenant: authTenantSummary{
			ID:                 "tenant-1",
			Status:             "active",
			MembershipStatus:   "active",
			SubscriptionStatus: "suspended",
		},
		PermissionSet: map[string]bool{"invoice.view": true},
	}))
	rec := httptest.NewRecorder()

	handler := requirePermissionForAuthenticated("invoice.view", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	handler(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected suspended subscription read to pass, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestRequirePermissionAllowsSubscriptionUpdateWhenTenantSubscriptionIsSuspended(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tenants/subscription/save", nil)
	req = req.WithContext(context.WithValue(req.Context(), authenticatedUserContextKey, authenticatedUser{
		ID:    "user-1",
		Email: "owner@example.edu.vn",
		ActiveTenant: authTenantSummary{
			ID:                 "tenant-1",
			Status:             "active",
			MembershipStatus:   "active",
			SubscriptionStatus: "suspended",
		},
		PermissionSet: map[string]bool{"subscription.update": true},
	}))
	rec := httptest.NewRecorder()

	handler := requirePermissionForAuthenticated("subscription.update", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})
	handler(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected subscription update to pass on suspended tenant, got %d: %s", rec.Code, rec.Body.String())
	}
}
