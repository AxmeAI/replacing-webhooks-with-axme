// Replacing webhooks with AXME — Go example.
//
// Payment processing: submit a payment intent with delivery guarantees.
// No webhook endpoint, no signature verification, no retry logic.
//
// Usage:
//
//	export AXME_API_KEY="your-key"
//	go run main.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/AxmeAI/axme-sdk-go/axme"
)

func main() {
	client, err := axme.NewClient(axme.ClientConfig{
		APIKey: os.Getenv("AXME_API_KEY"),
	})
	if err != nil {
		log.Fatalf("create client: %v", err)
	}

	ctx := context.Background()

	// Submit payment — platform delivers with retries, no webhook needed
	intentID, err := client.SendIntent(ctx, map[string]any{
		"intent_type":    "payment.process.v1",
		"to_agent":       "agent://myorg/production/payment-service",
		"order_id":       "ord_12345",
		"amount_cents":   9999,
		"currency":       "usd",
		"customer_email": "alice@example.com",
	}, axme.RequestOptions{})
	if err != nil {
		log.Fatalf("send intent: %v", err)
	}
	fmt.Printf("Payment submitted: %s\n", intentID)

	// Wait for completion — no webhook callback needed
	result, err := client.WaitFor(ctx, intentID, axme.ObserveOptions{})
	if err != nil {
		log.Fatalf("wait: %v", err)
	}
	fmt.Printf("Final status: %v\n", result["status"])
}
