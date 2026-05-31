package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

const permissionAuthenticated = ""

type routePermissionResolver func(*http.Request) (string, error)

type appAPIRoute struct {
	Method             string
	Path               string
	Public             bool
	Permission         string
	PermissionResolver routePermissionResolver
	Handler            http.HandlerFunc
}

func registerAPIRoutes(mux *http.ServeMux) {
	for _, route := range appAPIRoutes() {
		mux.HandleFunc(route.Path, route.wrap())
	}
}

func appAPIRoutes() []appAPIRoute {
	return []appAPIRoute{
		{Method: http.MethodPost, Path: "/api/v1/auth/login", Public: true, Handler: handleAuthLogin},
		{Path: "/api/v1/auth/bootstrap", Public: true, Handler: handleAuthBootstrap},
		{Method: http.MethodPost, Path: "/api/v1/auth/refresh", Public: true, Handler: handleAuthRefresh},
		{Method: http.MethodPost, Path: "/api/v1/auth/logout", Public: true, Handler: handleAuthLogout},
		{Method: http.MethodGet, Path: "/api/v1/auth/session", Public: true, Handler: handleAuthSession},
		{Method: http.MethodGet, Path: "/api/v1/healthz", Public: true, Handler: handleHealthz},
		{Method: http.MethodGet, Path: "/api/v1/banks", Permission: permissionAuthenticated, Handler: handleBanks},
		{Method: http.MethodGet, Path: "/api/v1/example", Permission: "payment.create", Handler: handleExample},
		{Method: http.MethodPost, Path: "/api/v1/import/fields", PermissionResolver: importFieldsPermission, Handler: handleImportFields},
		{Method: http.MethodPost, Path: "/api/v1/import/csv", Permission: "payment.create", Handler: handleImportCSV},
		{Method: http.MethodGet, Path: "/api/v1/master-data/options", Permission: "student.view", Handler: handleMasterDataOptions},
		{Method: http.MethodGet, Path: "/api/v1/master-data/students", Permission: "student.view", Handler: handleMasterDataStudents},
		{Method: http.MethodPost, Path: "/api/v1/master-data/import/csv", Permission: "student.create", Handler: handleMasterDataImportCSV},
		{Method: http.MethodPost, Path: "/api/v1/master-data/students/save", Permission: "student.update", Handler: handleMasterDataStudentSave},
		{Method: http.MethodGet, Path: "/api/v1/school-tree", Permission: "school_tree.view", Handler: handleSchoolTree},
		{Method: http.MethodPost, Path: "/api/v1/school-tree/schools/save", Permission: "school_tree.update", Handler: handleSchoolTreeSchoolSave},
		{Method: http.MethodPost, Path: "/api/v1/school-tree/school-years/save", Permission: "school_tree.update", Handler: handleSchoolTreeSchoolYearSave},
		{Method: http.MethodPost, Path: "/api/v1/school-tree/classes/save", Permission: "school_tree.update", Handler: handleSchoolTreeClassSave},
		{Method: http.MethodGet, Path: "/api/v1/fee-schedules/options", Permission: "fee.view", Handler: handleFeeScheduleOptions},
		{Method: http.MethodGet, Path: "/api/v1/fee-schedules", Permission: "fee.view", Handler: handleFeeScheduleList},
		{Method: http.MethodPost, Path: "/api/v1/fee-schedules/preview", Permission: "fee.view", Handler: handleFeeSchedulePreview},
		{Method: http.MethodPost, Path: "/api/v1/fee-schedules/save", Permission: "fee.update", Handler: handleFeeScheduleSave},
		{Method: http.MethodGet, Path: "/api/v1/invoices/options", Permission: "invoice.view", Handler: handleInvoiceOptions},
		{Method: http.MethodGet, Path: "/api/v1/invoices", Permission: "invoice.view", Handler: handleInvoiceList},
		{Method: http.MethodPost, Path: "/api/v1/invoices/preview", Permission: "invoice.view", Handler: handleInvoicePreview},
		{Method: http.MethodPost, Path: "/api/v1/invoices/generate", Permission: "invoice.create", Handler: handleInvoiceGenerate},
		{Method: http.MethodGet, Path: "/api/v1/invoices/detail", Permission: "invoice.view", Handler: handleInvoiceDetail},
		{Method: http.MethodGet, Path: "/api/v1/invoices/payment", Permission: "invoice.view", Handler: handleInvoicePayment},
		{Method: http.MethodGet, Path: "/api/v1/invoices/pdf", Permission: "invoice.view", Handler: handleInvoicePDF},
		{Method: http.MethodGet, Path: "/api/v1/payments/providers", Permission: "payment.view", Handler: handlePaymentProviders},
		{Method: http.MethodPost, Path: "/api/v1/payments/intents", Permission: "payment.create", Handler: handlePaymentIntentCreate},
		{Method: http.MethodGet, Path: "/api/v1/payments/transactions", Permission: "payment.view", Handler: handlePaymentTransactions},
		{Method: http.MethodGet, Path: "/api/v1/payments/reconciliation", Permission: "payment.view", Handler: handlePaymentReconciliation},
		{Method: http.MethodPost, Path: "/api/v1/payments/webhooks/", Public: true, Handler: handlePaymentWebhook},
		{Method: http.MethodPost, Path: "/api/v1/payments/cash-receipts", Permission: "payment.create", Handler: handleManualCashReceipt},
		{Method: http.MethodGet, Path: "/api/v1/notifications/options", Permission: "notification.view", Handler: handleNotificationOptions},
		{Method: http.MethodGet, Path: "/api/v1/notifications/templates", Permission: "notification.view", Handler: handleNotificationTemplates},
		{Method: http.MethodGet, Path: "/api/v1/notifications/campaigns", Permission: "notification.view", Handler: handleNotificationCampaigns},
		{Method: http.MethodPost, Path: "/api/v1/notifications/campaigns/preview", Permission: "notification.send", Handler: handleNotificationCampaignPreview},
		{Method: http.MethodPost, Path: "/api/v1/notifications/campaigns/email-preview", Permission: "notification.send", Handler: handleNotificationCampaignEmailPreview},
		{Method: http.MethodPost, Path: "/api/v1/notifications/campaigns/save", Permission: "notification.create", Handler: handleNotificationCampaignSave},
		{Method: http.MethodPost, Path: "/api/v1/notifications/campaigns/send", Permission: "notification.send", Handler: handleNotificationCampaignSend},
		{Method: http.MethodGet, Path: "/api/v1/notifications/logs", Permission: "notification.view", Handler: handleNotificationLogs},
		{Method: http.MethodGet, Path: "/api/v1/admin/dashboard", Permission: "dashboard.view", Handler: handleAdminDashboard},
		{Method: http.MethodGet, Path: "/api/v1/admin/reports", Permission: "report.view", Handler: handleAdminReports},
		{Method: http.MethodGet, Path: "/api/v1/admin/reports/export", Permission: "report.export", Handler: handleAdminReportsExport},
		{Method: http.MethodGet, Path: "/api/v1/admin/audit-logs", Permission: "audit_log.view", Handler: handleAdminAuditLogs},
		{Method: http.MethodGet, Path: "/api/v1/admin/operation-logs", Permission: "operation_log.view", Handler: handleAdminOperationLogs},
		{Method: http.MethodGet, Path: "/api/v1/admin/users", Permission: "user.view", Handler: handleAdminUsers},
		{Method: http.MethodPost, Path: "/api/v1/admin/users/save", PermissionResolver: adminUserSavePermission, Handler: handleAdminUserSave},
		{Method: http.MethodPost, Path: "/api/v1/admin/users/roles", Permission: "user.assign_role", Handler: handleAdminUserRoles},
		{Method: http.MethodGet, Path: "/api/v1/admin/roles", Permission: "role.view", Handler: handleAdminRoles},
		{Method: http.MethodGet, Path: "/api/v1/qr.png", Public: true, Handler: handleQRPNG},
		{Method: http.MethodPost, Path: "/api/v1/vietqr/batch", Permission: "payment.create", Handler: handleBatch},
		{Path: "/api/v1/email/config", PermissionResolver: emailConfigPermission, Handler: handleEmailConfig},
		{Method: http.MethodPost, Path: "/api/v1/email/preview", Permission: "notification.send", Handler: handleEmailPreview},
		{Method: http.MethodPost, Path: "/api/v1/email/send", Permission: "notification.send", Handler: handleEmailSend},
		{Path: "/api/v1/email/cron", PermissionResolver: emailCronPermission, Handler: handleEmailCron},
		{Method: http.MethodPost, Path: "/api/v1/email/cron/run", Permission: "email_cron.update", Handler: handleEmailCronRun},
	}
}

func (route appAPIRoute) wrap() http.HandlerFunc {
	next := route.Handler
	if !route.Public {
		if route.PermissionResolver != nil {
			next = requireResolvedPermission(route.PermissionResolver, next)
		} else {
			next = requirePermission(route.Permission, next)
		}
	}
	if route.Method != "" {
		next = method(route.Method, next)
	}
	return next
}

func requirePermission(permission string, next http.HandlerFunc) http.HandlerFunc {
	return requireAuthenticated(requirePermissionForAuthenticated(permission, next))
}

func requireResolvedPermission(resolve routePermissionResolver, next http.HandlerFunc) http.HandlerFunc {
	return requireAuthenticated(func(w http.ResponseWriter, r *http.Request) {
		permission, err := resolve(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		requirePermissionForAuthenticated(permission, next)(w, r)
	})
}

func requirePermissionForAuthenticated(permission string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if permission == permissionAuthenticated {
			next(w, r)
			return
		}
		user, ok := authenticatedUserFromRequest(r)
		if !ok {
			http.Error(w, "not authenticated", http.StatusUnauthorized)
			return
		}
		if !authenticatedUserHasPermission(user, permission) {
			http.Error(w, "missing required API permission: "+permission, http.StatusForbidden)
			return
		}
		next(w, r)
	}
}

func authenticatedUserHasPermission(user authenticatedUser, permission string) bool {
	if permission == permissionAuthenticated {
		return true
	}
	if user.PermissionSet[permission] {
		return true
	}
	for _, alias := range permissionAliases[permission] {
		if user.PermissionSet[alias] {
			return true
		}
	}
	for _, item := range user.Permissions {
		if item.Code == permission {
			return true
		}
		for _, alias := range permissionAliases[permission] {
			if item.Code == alias {
				return true
			}
		}
	}
	return false
}

var permissionAliases = map[string][]string{
	"user.view":           {"system.users.read"},
	"user.create":         {"system.users.write"},
	"user.update":         {"system.users.write"},
	"user.assign_role":    {"system.users.assign_roles", "system.users.write"},
	"role.view":           {"system.roles.read"},
	"role.update":         {"system.roles.write"},
	"student.view":        {"master_data.read"},
	"student.create":      {"master_data.write"},
	"student.update":      {"master_data.write"},
	"school_tree.view":    {"master_data.read"},
	"school_tree.update":  {"master_data.write"},
	"fee.view":            {"fee_schedules.read"},
	"fee.create":          {"fee_schedules.write"},
	"fee.update":          {"fee_schedules.write"},
	"invoice.view":        {"invoices.read"},
	"invoice.create":      {"invoices.write"},
	"invoice.update":      {"invoices.write"},
	"payment.view":        {"payments.read"},
	"payment.create":      {"payments.write"},
	"payment.reconcile":   {"payments.reconcile"},
	"notification.view":   {"notifications.read"},
	"notification.create": {"notifications.write"},
	"notification.send":   {"notifications.send", "email.send"},
	"email_config.view":   {"email.config.read"},
	"email_config.update": {"email.config.write"},
	"email_cron.view":     {"email.cron.manage"},
	"email_cron.update":   {"email.cron.manage"},
	"report.view":         {"admin.reports.read"},
	"report.export":       {"admin.reports.export"},
	"dashboard.view":      {"admin.dashboard.read"},
	"operation_log.view":  {"operations.read"},
	"audit_log.view":      {"audit.read"},

	"system.users.read":         {"user.view"},
	"system.users.write":        {"user.create", "user.update"},
	"system.users.assign_roles": {"user.assign_role"},
	"system.roles.read":         {"role.view"},
	"system.roles.write":        {"role.update"},
	"master_data.read":          {"student.view", "school_tree.view"},
	"master_data.write":         {"student.create", "student.update", "school_tree.update"},
	"fee_schedules.read":        {"fee.view"},
	"fee_schedules.write":       {"fee.create", "fee.update"},
	"invoices.read":             {"invoice.view"},
	"invoices.write":            {"invoice.create", "invoice.update"},
	"payments.read":             {"payment.view"},
	"payments.write":            {"payment.create"},
	"payments.reconcile":        {"payment.reconcile"},
	"notifications.read":        {"notification.view"},
	"notifications.write":       {"notification.create"},
	"notifications.send":        {"notification.send"},
	"email.config.read":         {"email_config.view"},
	"email.config.write":        {"email_config.update"},
	"email.send":                {"notification.send"},
	"email.cron.manage":         {"email_cron.view", "email_cron.update"},
	"admin.reports.read":        {"report.view"},
	"admin.reports.export":      {"report.export"},
	"admin.dashboard.read":      {"dashboard.view"},
	"operations.read":           {"operation_log.view"},
	"audit.read":                {"audit_log.view"},
}

func importFieldsPermission(r *http.Request) (string, error) {
	switch r.URL.Query().Get("target") {
	case "master_data":
		return "student.create", nil
	case "", "payments":
		return "payment.create", nil
	default:
		return "payment.create", nil
	}
}

func adminUserSavePermission(r *http.Request) (string, error) {
	if r.Method != http.MethodPost {
		return permissionAuthenticated, nil
	}
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		return "", err
	}
	r.Body = io.NopCloser(bytes.NewReader(body))
	var input adminUserSaveInput
	if strings.TrimSpace(string(body)) != "" {
		_ = json.Unmarshal(body, &input)
	}
	if strings.TrimSpace(input.ID) == "" {
		return "user.create", nil
	}
	return "user.update", nil
}

func emailConfigPermission(r *http.Request) (string, error) {
	switch r.Method {
	case http.MethodGet:
		return "email_config.view", nil
	case http.MethodPost:
		return "email_config.update", nil
	default:
		return permissionAuthenticated, nil
	}
}

func emailCronPermission(r *http.Request) (string, error) {
	switch r.Method {
	case http.MethodGet:
		return "email_cron.view", nil
	case http.MethodPost:
		return "email_cron.update", nil
	default:
		return permissionAuthenticated, nil
	}
}
