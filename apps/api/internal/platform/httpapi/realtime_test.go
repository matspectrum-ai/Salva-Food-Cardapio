package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/matspectrum-ai/salva-food/apps/api/internal/catalog"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/orders"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/realtime"
)

func TestOrderEventsSSESnapshotAndUpdate(t *testing.T) {
	tenant := "11111111-1111-4111-8111-111111111111"
	orderStore := orders.NewMemoryStore()
	order, err := orders.NewOrder(tenant, "22222222-2222-4222-8222-222222222222", "PDV", []orders.LineItemSnapshot{{ItemID: "33333333-3333-4333-8333-333333333333", Name: "Burger", Quantity: 1, UnitPriceCents: 2500}})
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := orderStore.Create(context.Background(), order, "sse-fixture", strings.Repeat("a", 64)); err != nil {
		t.Fatal(err)
	}
	bus := realtime.NewMemoryBus()
	a := newAPIWithRealtime(catalog.NewMemoryStore(), orderStore, nil, nil, nil, false, func() string { return "fixture" }, nil, bus, true)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/orders/events", nil)
	req.Header.Set("X-Tenant-ID", tenant)
	rec := httptest.NewRecorder()
	done := make(chan struct{})
	go func() { a.Handler().ServeHTTP(rec, req); close(done) }()
	waitForBody := func(want string) {
		t.Helper()
		deadline := time.Now().Add(2 * time.Second)
		for time.Now().Before(deadline) {
			if strings.Contains(rec.Body.String(), want) {
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
		t.Fatalf("SSE body missing %q: %s", want, rec.Body.String())
	}
	waitForBody("event: orders.snapshot")
	waitForBody(order.ID)
	if err := a.publishOrderEvent(ctx, order, realtime.EventOrderUpdated); err != nil {
		t.Fatal(err)
	}
	waitForBody("event: order.updated")
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("SSE handler did not stop")
	}
	if rec.Header().Get("Content-Type") != "text/event-stream" {
		t.Fatalf("content type = %q", rec.Header().Get("Content-Type"))
	}
}
