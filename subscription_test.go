package main

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestNormalizeSubscriptionPlanSaveInput(t *testing.T) {
	got := normalizeSubscriptionPlanSaveInput(subscriptionPlanSaveInput{
		Code: "Plus Plan",
		Name: "  Plus  ",
		Limits: map[string]any{
			"students": 300,
		},
	})

	if got.Code != "plus_plan" {
		t.Fatalf("expected normalized code plus_plan, got %q", got.Code)
	}
	if got.Name != "Plus" {
		t.Fatalf("expected trimmed name, got %q", got.Name)
	}
	if got.Status != "active" {
		t.Fatalf("expected default active status, got %q", got.Status)
	}
	if got.DisplayOrder == nil || *got.DisplayOrder != defaultSubscriptionPlanDisplayOrder {
		t.Fatalf("expected default display order, got %+v", got.DisplayOrder)
	}
}

func TestValidateSubscriptionPlanSaveInput(t *testing.T) {
	discountPrice := 1200000
	promotionalPrice := 1500000
	displayOrder := 20
	valid := subscriptionPlanSaveInput{
		Code:                "plus",
		Name:                "Plus",
		Status:              "active",
		DisplayOrder:        &displayOrder,
		BasePriceVND:        2000000,
		PartnerPriceVND:     &discountPrice,
		PromotionalPriceVND: &promotionalPrice,
		Limits: map[string]any{
			"schools": 3,
		},
	}
	if err := validateSubscriptionPlanSaveInput(valid); err != nil {
		t.Fatalf("expected valid subscription plan, got %v", err)
	}

	archived := valid
	archived.Status = "archived"
	if err := validateSubscriptionPlanSaveInput(archived); err != nil {
		t.Fatalf("expected archived status to be valid, got %v", err)
	}

	badCode := valid
	badCode.Code = "1plus"
	if err := validateSubscriptionPlanSaveInput(badCode); err == nil {
		t.Fatal("expected invalid plan code to fail")
	}

	badStatus := valid
	badStatus.Status = "draft"
	if err := validateSubscriptionPlanSaveInput(badStatus); err == nil {
		t.Fatal("expected invalid plan status to fail")
	}

	badLimits := valid
	badLimits.Limits = nil
	if err := validateSubscriptionPlanSaveInput(badLimits); err == nil {
		t.Fatal("expected nil limits to fail validation")
	}

	badDisplayOrderValue := -1
	badDisplayOrder := valid
	badDisplayOrder.DisplayOrder = &badDisplayOrderValue
	if err := validateSubscriptionPlanSaveInput(badDisplayOrder); err == nil {
		t.Fatal("expected negative display order to fail validation")
	}

	badBasePrice := valid
	badBasePrice.BasePriceVND = -1
	if err := validateSubscriptionPlanSaveInput(badBasePrice); err == nil {
		t.Fatal("expected negative base price to fail validation")
	}

	badPartnerPriceValue := -1
	badPartnerPrice := valid
	badPartnerPrice.PartnerPriceVND = &badPartnerPriceValue
	if err := validateSubscriptionPlanSaveInput(badPartnerPrice); err == nil {
		t.Fatal("expected negative partner price to fail validation")
	}

	badPromotionalPriceValue := 2500000
	badPromotionalPrice := valid
	badPromotionalPrice.PromotionalPriceVND = &badPromotionalPriceValue
	if err := validateSubscriptionPlanSaveInput(badPromotionalPrice); err == nil {
		t.Fatal("expected promotional price above base price to fail validation")
	}
}

func TestSubscriptionPlanDisplayPricePrefersPromotion(t *testing.T) {
	promotionalPrice := 1500000
	if got := subscriptionPlanDisplayPrice(2000000, &promotionalPrice); got != 1500000 {
		t.Fatalf("expected promotional display price, got %d", got)
	}
	if got := subscriptionPlanDisplayPrice(2000000, nil); got != 2000000 {
		t.Fatalf("expected base display price, got %d", got)
	}
}

func TestSaveSubscriptionPlanUpdatesCodeByID(t *testing.T) {
	db, state := openFakeSubscriptionPlanDB(t, false)
	partnerPrice := 1200000
	promotionalPrice := 1500000
	displayOrder := 10

	got, err := saveSubscriptionPlan(context.Background(), db, subscriptionPlanSaveInput{
		ID:                  "11111111-1111-1111-1111-111111111111",
		Code:                "plus_2026",
		Name:                "Plus 2026",
		Status:              "active",
		Description:         "Updated plan code",
		ContactPrice:        true,
		DisplayOrder:        &displayOrder,
		BasePriceVND:        2000000,
		PartnerPriceVND:     &partnerPrice,
		PromotionalPriceVND: &promotionalPrice,
		Limits: map[string]any{
			"schools": 3,
		},
	}, authenticatedUser{}, requestAuditContext{})
	if err != nil {
		t.Fatalf("expected update by id to succeed, got %v", err)
	}
	if got.ID != "11111111-1111-1111-1111-111111111111" || got.Code != "plus_2026" {
		t.Fatalf("expected updated plan identity/code, got %+v", got)
	}
	if got.DisplayPriceVND != promotionalPrice {
		t.Fatalf("expected promotional display price, got %d", got.DisplayPriceVND)
	}
	if !got.ContactPrice || got.DisplayOrder != displayOrder {
		t.Fatalf("expected contact price/display order to round-trip, got %+v", got)
	}
	if state.updateCount() != 1 {
		t.Fatalf("expected one update query, got %d", state.updateCount())
	}
	if !strings.Contains(state.lastUpdateQuery(), "WHERE id = $1::uuid") {
		t.Fatalf("expected update query to target id, got %q", state.lastUpdateQuery())
	}
}

func TestSaveSubscriptionPlanRejectsDuplicateCodeOnIDUpdate(t *testing.T) {
	db, state := openFakeSubscriptionPlanDB(t, true)
	_, err := saveSubscriptionPlan(context.Background(), db, subscriptionPlanSaveInput{
		ID:           "11111111-1111-1111-1111-111111111111",
		Code:         "plus",
		Name:         "Plus",
		Status:       "active",
		DisplayOrder: intPtr(20),
		BasePriceVND: 2000000,
		Limits: map[string]any{
			"schools": 3,
		},
	}, authenticatedUser{}, requestAuditContext{})
	if err == nil || !strings.Contains(err.Error(), "plan code already exists") {
		t.Fatalf("expected duplicate code error, got %v", err)
	}
	if state.updateCount() != 0 {
		t.Fatalf("expected duplicate guard to stop update, got %d update queries", state.updateCount())
	}
}

var registerFakeSubscriptionPlanDriver sync.Once
var fakeSubscriptionPlanStates sync.Map

type fakeSubscriptionPlanState struct {
	mu               sync.Mutex
	name             string
	conflict         bool
	updates          int
	lastUpdateSQL    string
	conflictCheckSQL string
}

func (state *fakeSubscriptionPlanState) updateCount() int {
	state.mu.Lock()
	defer state.mu.Unlock()
	return state.updates
}

func (state *fakeSubscriptionPlanState) lastUpdateQuery() string {
	state.mu.Lock()
	defer state.mu.Unlock()
	return state.lastUpdateSQL
}

func openFakeSubscriptionPlanDB(t *testing.T, conflict bool) (*sql.DB, *fakeSubscriptionPlanState) {
	t.Helper()
	registerFakeSubscriptionPlanDriver.Do(func() {
		sql.Register("fake_subscription_plans", fakeSubscriptionPlanDriver{})
	})
	name := strings.ReplaceAll(t.Name(), "/", "_")
	state := &fakeSubscriptionPlanState{name: name, conflict: conflict}
	fakeSubscriptionPlanStates.Store(name, state)
	db, err := sql.Open("fake_subscription_plans", name)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = db.Close()
		fakeSubscriptionPlanStates.Delete(name)
	})
	return db, state
}

type fakeSubscriptionPlanDriver struct{}

func (fakeSubscriptionPlanDriver) Open(name string) (driver.Conn, error) {
	value, ok := fakeSubscriptionPlanStates.Load(name)
	if !ok {
		return nil, errors.New("unknown fake subscription plan database")
	}
	return &fakeSubscriptionPlanConn{state: value.(*fakeSubscriptionPlanState)}, nil
}

type fakeSubscriptionPlanConn struct {
	state *fakeSubscriptionPlanState
}

func (conn *fakeSubscriptionPlanConn) Prepare(query string) (driver.Stmt, error) {
	return nil, errors.New("prepared statements are not implemented")
}

func (conn *fakeSubscriptionPlanConn) Close() error {
	return nil
}

func (conn *fakeSubscriptionPlanConn) Begin() (driver.Tx, error) {
	return fakeSubscriptionPlanTx{}, nil
}

func (conn *fakeSubscriptionPlanConn) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	return fakeSubscriptionPlanTx{}, nil
}

func (conn *fakeSubscriptionPlanConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	return driver.RowsAffected(1), nil
}

func (conn *fakeSubscriptionPlanConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	normalized := strings.Join(strings.Fields(query), " ")
	switch {
	case strings.Contains(normalized, "SELECT id::text FROM subscription_plans WHERE code = $1 AND id <> $2::uuid"):
		conn.state.mu.Lock()
		conn.state.conflictCheckSQL = normalized
		conflict := conn.state.conflict
		conn.state.mu.Unlock()
		if conflict {
			return &fakeSubscriptionPlanRows{
				columns: []string{"id"},
				values:  [][]driver.Value{{"22222222-2222-2222-2222-222222222222"}},
			}, nil
		}
		return &fakeSubscriptionPlanRows{columns: []string{"id"}}, nil
	case strings.HasPrefix(normalized, "UPDATE subscription_plans SET"):
		conn.state.mu.Lock()
		conn.state.updates++
		conn.state.lastUpdateSQL = normalized
		conn.state.mu.Unlock()
		return &fakeSubscriptionPlanRows{
			columns: []string{"id", "code", "name", "status", "description", "contact_price", "display_order", "base_price_vnd", "partner_price_vnd", "promotional_price_vnd", "limits", "updated_at"},
			values: [][]driver.Value{{
				stringDriverValue(args[0].Value),
				stringDriverValue(args[1].Value),
				stringDriverValue(args[2].Value),
				stringDriverValue(args[3].Value),
				stringDriverValue(args[4].Value),
				boolDriverValue(args[5].Value),
				int64DriverValue(args[6].Value),
				int64DriverValue(args[7].Value),
				nullableInt64DriverValue(args[8].Value),
				nullableInt64DriverValue(args[9].Value),
				[]byte(stringDriverValue(args[10].Value)),
				time.Date(2026, time.June, 10, 0, 0, 0, 0, time.UTC),
			}},
		}, nil
	default:
		return nil, errors.New("unexpected query: " + query)
	}
}

type fakeSubscriptionPlanTx struct{}

func (fakeSubscriptionPlanTx) Commit() error {
	return nil
}

func (fakeSubscriptionPlanTx) Rollback() error {
	return nil
}

type fakeSubscriptionPlanRows struct {
	columns []string
	values  [][]driver.Value
	index   int
}

func (rows *fakeSubscriptionPlanRows) Columns() []string {
	return rows.columns
}

func (rows *fakeSubscriptionPlanRows) Close() error {
	return nil
}

func (rows *fakeSubscriptionPlanRows) Next(dest []driver.Value) error {
	if rows.index >= len(rows.values) {
		return io.EOF
	}
	copy(dest, rows.values[rows.index])
	rows.index++
	return nil
}

func stringDriverValue(value any) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(value.(string))
}

func int64DriverValue(value any) int64 {
	switch typed := value.(type) {
	case int64:
		return typed
	case int:
		return int64(typed)
	default:
		return 0
	}
}

func boolDriverValue(value any) bool {
	typed, _ := value.(bool)
	return typed
}

func nullableInt64DriverValue(value any) driver.Value {
	if value == nil {
		return nil
	}
	return int64DriverValue(value)
}
