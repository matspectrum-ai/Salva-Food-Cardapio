package realtime

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestRedisBusRoundTrip(t *testing.T) {
	url := os.Getenv("TEST_REDIS_URL")
	if url == "" {
		t.Skip("TEST_REDIS_URL is not set")
	}
	opts, err := redis.ParseURL(url)
	if err != nil {
		t.Fatal(err)
	}
	client := redis.NewClient(opts)
	defer client.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		t.Fatal(err)
	}
	bus := NewRedisBus(client, "salva-food-test")
	ch, err := bus.Subscribe(ctx, "tenant-test")
	if err != nil {
		t.Fatal(err)
	}
	want := Event{ID: "redis-1", TenantID: "tenant-test", Type: EventOrderCreated, Payload: []byte(`{"id":"order-1"}`), OccurredAt: time.Now()}
	if err := bus.Publish(ctx, want); err != nil {
		t.Fatal(err)
	}
	select {
	case got := <-ch:
		if got.ID != want.ID || got.TenantID != want.TenantID || got.Type != want.Type || string(got.Payload) != string(want.Payload) {
			t.Fatalf("got %+v, want %+v", got, want)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("did not receive redis event")
	}
}
