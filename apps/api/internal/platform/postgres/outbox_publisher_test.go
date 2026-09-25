package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	db "github.com/matspectrum-ai/salva-food/apps/api/internal/platform/postgres/sqlc"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/realtime"
	"github.com/redis/go-redis/v9"
)

func TestOutboxPublisherPublishesToRedis(t *testing.T) {
	dbURL, redisURL := os.Getenv("TEST_DATABASE_URL"), os.Getenv("TEST_REDIS_URL")
	if dbURL == "" || redisURL == "" {
		t.Skip("TEST_DATABASE_URL and TEST_REDIS_URL are required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	store, err := Open(ctx, dbURL)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		t.Fatal(err)
	}
	client := redis.NewClient(opts)
	defer client.Close()
	bus := realtime.NewRedisBus(client, "salva-food-outbox-test")
	tenant := uuid.New()
	order := uuid.New()
	if _, err := store.q.CreateTenant(ctx, db.CreateTenantParams{ID: pgUUID(tenant), Name: "Outbox Test"}); err != nil {
		t.Fatal(err)
	}
	defer func() {
		_, _ = store.pool.Exec(context.Background(), `DELETE FROM tenants WHERE id = $1`, pgUUID(tenant))
	}()
	ch, err := bus.Subscribe(ctx, tenant.String())
	if err != nil {
		t.Fatal(err)
	}
	_, err = store.q.CreateOutboxEvent(ctx, db.CreateOutboxEventParams{ID: pgUUID(uuid.New()), TenantID: pgUUID(tenant), AggregateType: "order", AggregateID: pgUUID(order), EventType: realtime.EventOrderUpdated, Payload: []byte(`{"id":"fixture"}`)})
	if err != nil {
		t.Fatal(err)
	}

	publisher := NewOutboxPublisher(store, bus)
	publisher.publishBatch(ctx)
	select {
	case event := <-ch:
		if event.AggregateID != order.String() || event.Type != realtime.EventOrderUpdated {
			t.Fatalf("event=%+v", event)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("outbox event was not published")
	}
	var published int
	if err := store.pool.QueryRow(ctx, `SELECT count(*) FROM outbox_events WHERE aggregate_id = $1 AND published_at IS NOT NULL`, pgUUID(order)).Scan(&published); err != nil {
		t.Fatal(err)
	}
	if published != 1 {
		t.Fatalf("published rows=%d, want 1", published)
	}
}

func pgUUID(id uuid.UUID) pgtype.UUID { return pgtype.UUID{Bytes: id, Valid: true} }
