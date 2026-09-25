package realtime

import (
	"context"
	"testing"
	"time"
)

func TestMemoryBusPublishesOnlyToTenantSubscribers(t *testing.T) {
	bus := NewMemoryBus()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	first, err := bus.Subscribe(ctx, "tenant-a")
	if err != nil {
		t.Fatal(err)
	}
	second, err := bus.Subscribe(ctx, "tenant-b")
	if err != nil {
		t.Fatal(err)
	}
	event := Event{ID: "event-1", TenantID: "tenant-a", Type: EventOrderUpdated, OccurredAt: time.Now()}
	if err := bus.Publish(ctx, event); err != nil {
		t.Fatal(err)
	}
	select {
	case got := <-first:
		if got.ID != event.ID {
			t.Fatalf("id = %s, want %s", got.ID, event.ID)
		}
	case <-time.After(time.Second):
		t.Fatal("tenant-a subscriber did not receive event")
	}
	select {
	case got := <-second:
		t.Fatalf("tenant-b received event: %+v", got)
	case <-time.After(50 * time.Millisecond):
	}
}
