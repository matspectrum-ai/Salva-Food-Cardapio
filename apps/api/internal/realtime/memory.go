package realtime

import (
	"context"
	"sync"
)

type MemoryBus struct {
	mu   sync.RWMutex
	next uint64
	subs map[string]map[uint64]chan Event
}

func NewMemoryBus() *MemoryBus {
	return &MemoryBus{subs: make(map[string]map[uint64]chan Event)}
}

func (b *MemoryBus) Publish(_ context.Context, event Event) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for _, ch := range b.subs[event.TenantID] {
		select {
		case ch <- event:
		default:
		}
	}
	return nil
}
func (b *MemoryBus) Subscribe(ctx context.Context, tenantID string) (<-chan Event, error) {
	b.mu.Lock()
	b.next++
	id := b.next
	if b.subs[tenantID] == nil {
		b.subs[tenantID] = make(map[uint64]chan Event)
	}
	ch := make(chan Event, 64)
	b.subs[tenantID][id] = ch
	b.mu.Unlock()

	go func() {
		<-ctx.Done()
		b.mu.Lock()
		delete(b.subs[tenantID], id)
		if len(b.subs[tenantID]) == 0 {
			delete(b.subs, tenantID)
		}
		close(ch)
		b.mu.Unlock()
	}()

	return ch, nil
}
