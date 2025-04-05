package webhook

import (
	"context"
)

type Webhook interface {
	Send(ctx context.Context, ev Event) error
}

type NopWebhook struct{}

func (NopWebhook) Send(_ context.Context, _ Event) error {
	return nil
}
