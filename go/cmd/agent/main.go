// Payment service agent — Go example.
//
// Processes payments and resumes with transaction result.
//
// Usage:
//
//	export AXME_API_KEY="<agent-key>"
//	go run agent.go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AxmeAI/axme-sdk-go/axme"
)

const agentAddress = "payment-service-demo"

func handleIntent(ctx context.Context, client *axme.Client, intentID string) error {
	intentData, err := client.GetIntent(ctx, intentID, axme.RequestOptions{})
	if err != nil {
		return fmt.Errorf("get intent: %w", err)
	}

	intent, _ := intentData["intent"].(map[string]any)
	if intent == nil {
		intent = intentData
	}
	payload, _ := intent["payload"].(map[string]any)
	if payload == nil {
		payload = map[string]any{}
	}
	if pp, ok := payload["parent_payload"].(map[string]any); ok {
		payload = pp
	}

	orderID, _ := payload["order_id"].(string)
	if orderID == "" {
		orderID = "unknown"
	}
	amount, _ := payload["amount"].(float64)
	currency, _ := payload["currency"].(string)
	if currency == "" {
		currency = "USD"
	}
	method, _ := payload["method"].(string)
	if method == "" {
		method = "card"
	}

	fmt.Printf("  Processing %s payment: %s %.0f for %s...\n", method, currency, amount, orderID)
	time.Sleep(1 * time.Second)
	fmt.Println("  Authorizing with payment provider...")
	time.Sleep(1 * time.Second)
	fmt.Println("  Capturing funds...")
	time.Sleep(1 * time.Second)

	result := map[string]any{
		"action":         "complete",
		"order_id":       orderID,
		"transaction_id": "TXN-99001",
		"amount_charged": amount,
		"status":         "captured",
		"processed_at":   time.Now().UTC().Format(time.RFC3339),
	}

	_, err = client.ResumeIntent(ctx, intentID, result, axme.RequestOptions{})
	if err != nil {
		return fmt.Errorf("resume intent: %w", err)
	}
	fmt.Printf("  Payment captured: TXN-99001 (%s %.0f)\n", currency, amount)
	return nil
}

func main() {
	apiKey := os.Getenv("AXME_API_KEY")
	if apiKey == "" {
		log.Fatal("Error: AXME_API_KEY not set.")
	}

	client, err := axme.NewClient(axme.ClientConfig{APIKey: apiKey})
	if err != nil {
		log.Fatalf("create client: %v", err)
	}

	ctx := context.Background()

	fmt.Printf("Agent listening on %s...\n", agentAddress)
	fmt.Println("Waiting for intents (Ctrl+C to stop)")

	intents, errCh := client.Listen(ctx, agentAddress, axme.ListenOptions{})

	go func() {
		for err := range errCh {
			log.Printf("Listen error: %v", err)
		}
	}()

	for delivery := range intents {
		intentID, _ := delivery["intent_id"].(string)
		status, _ := delivery["status"].(string)
		if intentID == "" {
			continue
		}
		if status == "DELIVERED" || status == "CREATED" || status == "IN_PROGRESS" {
			fmt.Printf("[%s] Intent received: %s\n", status, intentID)
			if err := handleIntent(ctx, client, intentID); err != nil {
				fmt.Printf("  Error: %v\n", err)
			}
		}
	}
}
