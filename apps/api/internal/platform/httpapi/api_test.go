package httpapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/matspectrum-ai/salva-food/apps/api/internal/catalog"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/orders"
)

const (
	tenantA = "11111111-1111-4111-8111-111111111111"
	tenantB = "22222222-2222-4222-8222-222222222222"
)

type testHarness struct {
	handler http.Handler
}

func newTestHarness() testHarness {
	catalogStore := catalog.NewMemoryStore()
	orderStore := orders.NewMemoryStore()
	var sequence atomic.Int64
	newID := func() string { return fmt.Sprintf("00000000-0000-4000-8000-%012d", sequence.Add(1)) }
	return testHarness{handler: New(catalogStore, orderStore, newID).Handler()}
}
func (h testHarness) do(t *testing.T, method, path, tenantID string, body any, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()
	var payload []byte
	var err error
	if body != nil {
		payload, err = json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
	}
	req := httptest.NewRequest(method, path, bytes.NewReader(payload))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if tenantID != "" {
		req.Header.Set("X-Tenant-ID", tenantID)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	rec := httptest.NewRecorder()
	h.handler.ServeHTTP(rec, req)
	return rec
}

func decodeResponse[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()
	var value T
	if err := json.Unmarshal(rec.Body.Bytes(), &value); err != nil {
		t.Fatalf("decode response: %v\nbody=%s", err, rec.Body.String())
	}
	return value
}
func TestCatalogOrderLifecycleHTTP(t *testing.T) {
	h := newTestHarness()
	categoryRec := h.do(t, http.MethodPost, "/api/v1/catalog/categories", tenantA, map[string]any{
		"name": "Lanches", "sort_order": 0,
	}, nil)
	if categoryRec.Code != http.StatusCreated {
		t.Fatalf("create category: got %d body=%s", categoryRec.Code, categoryRec.Body.String())
	}
	category := decodeResponse[catalog.Category](t, categoryRec)

	itemRec := h.do(t, http.MethodPost, "/api/v1/catalog/items", tenantA, map[string]any{
		"category_id":    category.ID,
		"name":           "X-Burger",
		"price_cents":    2590,
		"sort_order":     0,
		"sold_by_weight": false,
	}, nil)
	if itemRec.Code != http.StatusCreated {
		t.Fatalf("create item: got %d body=%s", itemRec.Code, itemRec.Body.String())
	}
	item := decodeResponse[catalog.Item](t, itemRec)

	orderBody := map[string]any{
		"source": "PDV",
		"items":  []map[string]any{{"item_id": item.ID, "quantity": 2}},
	}
	orderRec := h.do(t, http.MethodPost, "/api/v1/orders", tenantA, orderBody, map[string]string{"Idempotency-Key": "order-001"})
	if orderRec.Code != http.StatusCreated {
		t.Fatalf("create order: got %d body=%s", orderRec.Code, orderRec.Body.String())
	}
	order := decodeResponse[orders.Order](t, orderRec)
	if order.Status != orders.StatusAnalysis || order.TotalCents != 5180 {
		t.Fatalf("unexpected order: %+v", order)
	}

	for _, next := range []orders.Status{orders.StatusProduction, orders.StatusReady, orders.StatusFinalized} {
		rec := h.do(t, http.MethodPost, fmt.Sprintf("/api/v1/orders/%s/transitions", order.ID), tenantA, map[string]any{
			"status": next,
		}, nil)
		if rec.Code != http.StatusOK {
			t.Fatalf("transition %s: got %d body=%s", next, rec.Code, rec.Body.String())
		}
		order = decodeResponse[orders.Order](t, rec)
		if order.Status != next {
			t.Fatalf("transition status: got %s want %s", order.Status, next)
		}
	}

	invalid := h.do(t, http.MethodPost, fmt.Sprintf("/api/v1/orders/%s/transitions", order.ID), tenantA, map[string]any{
		"status": orders.StatusCancelled,
	}, nil)
	if invalid.Code != http.StatusBadRequest {
		t.Fatalf("terminal transition: got %d body=%s", invalid.Code, invalid.Body.String())
	}
}

func TestOrderIdempotencyHTTP(t *testing.T) {
	h := newTestHarness()
	categoryRec := h.do(t, http.MethodPost, "/api/v1/catalog/categories", tenantA, map[string]any{"name": "Pizzas", "sort_order": 0}, nil)
	category := decodeResponse[catalog.Category](t, categoryRec)
	itemRec := h.do(t, http.MethodPost, "/api/v1/catalog/items", tenantA, map[string]any{
		"category_id": category.ID, "name": "Calabresa", "price_cents": 4000, "sort_order": 0,
	}, nil)
	item := decodeResponse[catalog.Item](t, itemRec)
	body := map[string]any{"source": "STOREFRONT", "items": []map[string]any{{"item_id": item.ID, "quantity": 1}}}

	first := h.do(t, http.MethodPost, "/api/v1/orders", tenantA, body, map[string]string{"Idempotency-Key": "checkout-001"})
	if first.Code != http.StatusCreated {
		t.Fatalf("first create: got %d body=%s", first.Code, first.Body.String())
	}
	firstOrder := decodeResponse[orders.Order](t, first)

	replay := h.do(t, http.MethodPost, "/api/v1/orders", tenantA, body, map[string]string{"Idempotency-Key": "checkout-001"})
	if replay.Code != http.StatusCreated || replay.Header().Get("Idempotent-Replay") != "true" {
		t.Fatalf("replay: got %d replay=%q body=%s", replay.Code, replay.Header().Get("Idempotent-Replay"), replay.Body.String())
	}
	replayedOrder := decodeResponse[orders.Order](t, replay)
	if replayedOrder.ID != firstOrder.ID {
		t.Fatalf("replay returned different order: %s != %s", replayedOrder.ID, firstOrder.ID)
	}

	conflictBody := map[string]any{"source": "STOREFRONT", "items": []map[string]any{{"item_id": item.ID, "quantity": 2}}}
	conflict := h.do(t, http.MethodPost, "/api/v1/orders", tenantA, conflictBody, map[string]string{"Idempotency-Key": "checkout-001"})
	if conflict.Code != http.StatusConflict {
		t.Fatalf("idempotency conflict: got %d body=%s", conflict.Code, conflict.Body.String())
	}
}

func TestTenantIsolationHTTP(t *testing.T) {
	h := newTestHarness()
	created := h.do(t, http.MethodPost, "/api/v1/catalog/categories", tenantA, map[string]any{"name": "Bebidas", "sort_order": 0}, nil)
	if created.Code != http.StatusCreated {
		t.Fatalf("create category: got %d", created.Code)
	}

	otherTenant := h.do(t, http.MethodGet, "/api/v1/catalog/categories", tenantB, nil, nil)
	if otherTenant.Code != http.StatusOK {
		t.Fatalf("list tenant-b: got %d", otherTenant.Code)
	}
	var payload struct {
		Data []catalog.Category `json:"data"`
	}
	payload = decodeResponse[struct {
		Data []catalog.Category `json:"data"`
	}](t, otherTenant)
	if len(payload.Data) != 0 {
		t.Fatalf("tenant leak: %+v", payload.Data)
	}

	missingTenant := h.do(t, http.MethodGet, "/api/v1/orders", "", nil, nil)
	if missingTenant.Code != http.StatusBadRequest {
		t.Fatalf("missing tenant header: got %d", missingTenant.Code)
	}

	invalidTenant := h.do(t, http.MethodGet, "/api/v1/orders", "not-a-uuid", nil, nil)
	if invalidTenant.Code != http.StatusBadRequest {
		t.Fatalf("invalid tenant UUID: got %d", invalidTenant.Code)
	}
}
