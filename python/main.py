"""
Replacing webhooks with AXME — Python example.

Payment processing: submit a payment intent with delivery guarantees.
No webhook endpoint, no signature verification, no retry logic.

Usage:
    pip install axme
    export AXME_API_KEY="your-key"
    python main.py
"""

import os
from axme import AxmeClient, AxmeClientConfig


def main():
    client = AxmeClient(
        AxmeClientConfig(api_key=os.environ["AXME_API_KEY"])
    )

    # Submit payment — platform delivers with retries, no webhook needed
    intent_id = client.send_intent(
        {
            "intent_type": "payment.process.v1",
            "to_agent": "agent://myorg/production/payment-service",
            "payload": {
                "order_id": "ord_12345",
                "amount_cents": 9999,
                "currency": "usd",
                "customer_email": "alice@example.com",
            },
        }
    )
    print(f"Payment submitted: {intent_id}")

    # Observe delivery and processing — no webhook endpoint needed
    print("Watching delivery...")
    for event in client.observe(intent_id):
        status = event.get("status", "")
        print(f"  [{status}] {event.get('event_type', '')}")
        if status in ("COMPLETED", "FAILED", "TIMED_OUT", "CANCELLED"):
            break

    # Fetch final state
    intent = client.get_intent(intent_id)
    print(f"\nFinal status: {intent['intent']['lifecycle_status']}")


if __name__ == "__main__":
    main()
