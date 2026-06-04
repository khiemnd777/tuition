package main

import (
	"net/http"
	"strings"
)

func requireActiveTenantID(w http.ResponseWriter, r *http.Request) (string, bool) {
	tenantID := activeTenantIDFromRequest(r)
	if tenantID == "" {
		http.Error(w, "active tenant required", http.StatusForbidden)
		return "", false
	}
	return tenantID, true
}

func normalizeTenantID(value string) string {
	return strings.TrimSpace(value)
}
