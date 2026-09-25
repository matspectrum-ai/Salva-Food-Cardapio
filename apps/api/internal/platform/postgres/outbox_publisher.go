package postgres

import (
	"context"
	"time"

	db "github.com/matspectrum-ai/salva-food/apps/api/internal/platform/postgres/sqlc"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/realtime"
)

type OutboxPublisher struct {
	store     *Store
	bus       realtime.Bus
	batchSize int32
	interval  time.Duration
}

func NewOutboxPublisher(store *Store, bus realtime.Bus) *OutboxPublisher {
	return &OutboxPublisher{store: store, bus: bus, batchSize: 100, interval: 250 * time.Millisecond}
}

func (p *OutboxPublisher) Run(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()
	for {
		p.publishBatch(ctx)
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func (p *OutboxPublisher) publishBatch(ctx context.Context) {
	events, err := p.store.q.ListPendingOutboxEvents(ctx, p.batchSize)
	if err != nil {
		return
	}

	for _, row := range events {
		event := mapOutboxEvent(row)
		if err := p.bus.Publish(ctx, event); err != nil {
			_ = p.store.q.IncrementOutboxEventAttempts(ctx, row.ID)
			continue
		}
		_ = p.store.q.MarkOutboxEventPublished(ctx, row.ID)
	}
}

func mapOutboxEvent(row db.OutboxEvent) realtime.Event {
	return realtime.Event{
		ID:            formatUUID(row.ID),
		TenantID:      formatUUID(row.TenantID),
		AggregateType: row.AggregateType,
		AggregateID:   formatUUID(row.AggregateID),
		Type:          row.EventType,
		Payload:       row.Payload,
		OccurredAt:    row.OccurredAt.Time,
	}
}
