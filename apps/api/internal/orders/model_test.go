package orders

import "testing"

func validItems() []LineItemSnapshot {
	return []LineItemSnapshot{
		{ItemID: "item-1", Name: "Pizza", Quantity: 2, UnitPriceCents: 3000},
		{ItemID: "item-2", Name: "Suco", Quantity: 1, UnitPriceCents: 900},
	}
}

func TestNewOrderStartsInAnalysisAndSnapshotsTotal(t *testing.T) {
	order, err := NewOrder("tenant-1", "order-1", "PDV", validItems())
	if err != nil {
		t.Fatal(err)
	}
	if order.Status != StatusAnalysis {
		t.Fatalf("status = %s, want %s", order.Status, StatusAnalysis)
	}
	if order.TotalCents != 6900 {
		t.Fatalf("total = %d, want 6900", order.TotalCents)
	}
}

func TestNewOrderRejectsInvalidBoundaryData(t *testing.T) {
	cases := []struct {
		name, tenantID, id, source string
		items                      []LineItemSnapshot
	}{
		{"missing tenant", "", "o1", "PDV", validItems()},
		{"missing id", "t1", "", "PDV", validItems()},
		{"missing source", "t1", "o1", "", validItems()},
		{"empty items", "t1", "o1", "PDV", nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := NewOrder(tc.tenantID, tc.id, tc.source, tc.items); err == nil {
				t.Fatal("expected validation error")
			}
		})
	}
}
func TestOrderLifecycle(t *testing.T) {
	order, _ := NewOrder("tenant-1", "order-1", "PDV", validItems())
	for _, next := range []Status{StatusProduction, StatusReady, StatusFinalized} {
		if err := order.Transition(next); err != nil {
			t.Fatalf("transition to %s: %v", next, err)
		}
	}
	if order.Status != StatusFinalized {
		t.Fatalf("status = %s", order.Status)
	}
	if err := order.Transition(StatusReady); err == nil {
		t.Fatal("expected finalized order to be terminal")
	}
}

func TestOrderCannotSkipLifecycleStages(t *testing.T) {
	order, _ := NewOrder("tenant-1", "order-1", "PDV", validItems())
	if err := order.Transition(StatusReady); err == nil {
		t.Fatal("expected invalid transition")
	}
}

func TestOrderCanBeCancelledBeforeFinalization(t *testing.T) {
	for _, start := range []Status{StatusAnalysis, StatusProduction, StatusReady} {
		t.Run(string(start), func(t *testing.T) {
			order, _ := NewOrder("tenant-1", "order-1", "PDV", validItems())
			if start == StatusProduction {
				_ = order.Transition(StatusProduction)
			}
			if start == StatusReady {
				_ = order.Transition(StatusProduction)
				_ = order.Transition(StatusReady)
			}
			if err := order.Transition(StatusCancelled); err != nil {
				t.Fatalf("cancel: %v", err)
			}
			if err := order.Transition(StatusProduction); err == nil {
				t.Fatal("cancelled order must be terminal")
			}
		})
	}
}
func TestOrderRejectsInvalidLineItems(t *testing.T) {
	cases := []LineItemSnapshot{
		{ItemID: "", Name: "Pizza", Quantity: 1, UnitPriceCents: 1000},
		{ItemID: "item-1", Name: "", Quantity: 1, UnitPriceCents: 1000},
		{ItemID: "item-1", Name: "Pizza", Quantity: 0, UnitPriceCents: 1000},
		{ItemID: "item-1", Name: "Pizza", Quantity: 1, UnitPriceCents: -1},
	}
	for i, item := range cases {
		if _, err := NewOrder("tenant-1", "order-1", "PDV", []LineItemSnapshot{item}); err == nil {
			t.Fatalf("case %d: expected invalid item error", i)
		}
	}
}

func TestTransitionToCurrentStateIsIdempotent(t *testing.T) {
	order, _ := NewOrder("tenant-1", "order-1", "PDV", validItems())
	if err := order.Transition(StatusAnalysis); err != nil {
		t.Fatalf("same-state transition should be idempotent: %v", err)
	}
	if order.Status != StatusAnalysis {
		t.Fatalf("unexpected state %s", order.Status)
	}
}
