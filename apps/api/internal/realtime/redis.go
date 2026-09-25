package realtime

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
)

type RedisBus struct {
	client *redis.Client
	prefix string
}

func NewRedisBus(client *redis.Client, prefix string) *RedisBus {
	if strings.TrimSpace(prefix) == "" {
		prefix = "salva-food"
	}
	return &RedisBus{client: client, prefix: prefix}
}

func (b *RedisBus) channel(tenantID string) string {
	return fmt.Sprintf("%s:tenant:%s:events", b.prefix, tenantID)
}

func (b *RedisBus) Publish(ctx context.Context, event Event) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return b.client.Publish(ctx, b.channel(event.TenantID), payload).Err()
}
func (b *RedisBus) Subscribe(ctx context.Context, tenantID string) (<-chan Event, error) {
	pubsub := b.client.Subscribe(ctx, b.channel(tenantID))
	if _, err := pubsub.Receive(ctx); err != nil {
		_ = pubsub.Close()
		return nil, err
	}
	out := make(chan Event, 64)
	go func() {
		defer close(out)
		defer pubsub.Close()
		ch := pubsub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-ch:
				if !ok {
					return
				}
				var event Event
				if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
					continue
				}
				select {
				case out <- event:
				default:
				}
			}
		}
	}()
	return out, nil
}
