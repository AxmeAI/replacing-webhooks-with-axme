"""Payment service agent - processes payments and resumes."""

import os, sys, time
sys.stdout.reconfigure(line_buffering=True)
from axme import AxmeClient, AxmeClientConfig

AGENT_ADDRESS = "payment-service-demo"

def handle_intent(client, intent_id):
    intent_data = client.get_intent(intent_id)
    intent = intent_data.get("intent", intent_data)
    payload = intent.get("payload", {})
    if "parent_payload" in payload:
        payload = payload["parent_payload"]

    order_id = payload.get("order_id", "unknown")
    amount = payload.get("amount", 0)
    currency = payload.get("currency", "USD")
    method = payload.get("method", "card")

    print(f"  Processing {method} payment: {currency} {amount} for {order_id}...")
    time.sleep(1)
    print(f"  Authorizing with payment provider...")
    time.sleep(1)
    print(f"  Capturing funds...")
    time.sleep(1)

    result = {
        "action": "complete",
        "order_id": order_id,
        "transaction_id": "TXN-99001",
        "amount_charged": amount,
        "status": "captured",
        "processed_at": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
    }
    client.resume_intent(intent_id, result)
    print(f"  Payment captured: TXN-99001 ({currency} {amount})")

def main():
    api_key = os.environ.get("AXME_API_KEY", "")
    if not api_key:
        print("Error: AXME_API_KEY not set."); sys.exit(1)
    client = AxmeClient(AxmeClientConfig(api_key=api_key))
    print(f"Agent listening on {AGENT_ADDRESS}...")
    print("Waiting for intents (Ctrl+C to stop)\n")
    for delivery in client.listen(AGENT_ADDRESS):
        intent_id = delivery.get("intent_id", "")
        status = delivery.get("status", "")
        if intent_id and status in ("DELIVERED", "CREATED", "IN_PROGRESS"):
            print(f"[{status}] Intent received: {intent_id}")
            try:
                handle_intent(client, intent_id)
            except Exception as e:
                print(f"  Error: {e}")

if __name__ == "__main__":
    main()
