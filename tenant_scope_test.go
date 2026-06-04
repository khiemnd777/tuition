package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTenantScopedHandlersRejectMissingActiveTenant(t *testing.T) {
	cases := map[string]http.HandlerFunc{
		"master data options":    handleMasterDataOptions,
		"school tree":            handleSchoolTree,
		"fee schedule options":   handleFeeScheduleOptions,
		"invoice options":        handleInvoiceOptions,
		"payment reconciliation": handlePaymentReconciliation,
		"notification options":   handleNotificationOptions,
		"admin dashboard":        handleAdminDashboard,
		"admin reports export":   handleAdminReportsExport,
		"admin audit logs":       handleAdminAuditLogs,
		"admin operation logs":   handleAdminOperationLogs,
	}
	for name, handler := range cases {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
			rec := httptest.NewRecorder()

			handler(rec, req)

			if rec.Code != http.StatusForbidden {
				t.Fatalf("expected missing active tenant to be forbidden, got %d: %s", rec.Code, rec.Body.String())
			}
		})
	}
}
