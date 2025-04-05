package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Webhook interface {
	Send(ctx context.Context, ev Event) error
}

type NopWebhook struct{}

func (NopWebhook) Send(_ context.Context, _ Event) error {
	return nil
}

type PostWebhook struct {
	Url     string
	Timeout time.Duration
	Client  *http.Client
}

func (wh PostWebhook) Send(ev Event) error {
	data, err := json.Marshal(ev)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), wh.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", wh.Url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := wh.Client.Do(req)
	if err != nil {
		return err
	}
	_ = resp.Body.Close()

	if resp.StatusCode > 299 {
		return fmt.Errorf("webhook response status code: %d", resp.StatusCode)
	}

	return nil
}
