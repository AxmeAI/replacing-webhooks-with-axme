# Replacing Webhooks with AXME

Webhooks fail silently. Retries are unreliable. Signature verification is error-prone. Dead letter queues need their own monitoring. You build all this infrastructure for every service-to-service callback.

**There is a better way.** Replace webhook callbacks with intent lifecycle — built-in delivery guarantees, retries, and real-time observability.

> **Alpha** · Built with [AXME](https://github.com/AxmeAI/axme) (AXP Intent Protocol).
> [cloud.axme.ai](https://cloud.axme.ai) · [contact@axme.ai](mailto:contact@axme.ai)

---

## The Problem

Every payment integration ships with the same webhook headaches:

```
Payment provider → POST /webhooks/payment-complete
  → Verify HMAC signature (wrong? silent failure)
  → Parse payload (schema changed? silent failure)
  → Process event (DB down? retry later... maybe)
  → Return 200 within 5 seconds (timeout? provider retries... duplicate events)
```

What you end up building:
- **Signature verification** — HMAC, RSA, or custom schemes per provider
- **Idempotency layer** — deduplicate retried webhooks by event ID
- **Retry infrastructure** — exponential backoff, max retries, dead letter queue
- **Monitoring** — alert when webhooks stop arriving (how do you know?)
- **Public endpoint** — expose an HTTPS endpoint, manage certificates, firewall rules

---

## The Solution: Intent Delivery

```
Client → send_intent("process payment") → intent_id
Platform → delivers to payment service → retries on failure
Client → observe(intent_id) ← real-time lifecycle events
```

No public endpoint. No signature verification. No retry logic. No dead letter queue. The platform guarantees delivery.

---

## Quick Start

### Python

```bash
pip install axme
export AXME_API_KEY="your-key"   # Get one: axme login
```

```python
from axme import AxmeClient, AxmeClientConfig
import os

client = AxmeClient(AxmeClientConfig(api_key=os.environ["AXME_API_KEY"]))

# Submit payment — platform delivers to payment service with retries
intent_id = client.send_intent({
    "intent_type": "payment.process.v1",
    "to_agent": "agent://myorg/production/payment-service",
    "payload": {
        "order_id": "ord_12345",
        "amount_cents": 9999,
        "currency": "usd",
        "customer_email": "alice@example.com",
    },
})

print(f"Payment submitted: {intent_id}")

# Observe delivery and processing — no webhook endpoint needed
for event in client.observe(intent_id):
    print(f"  [{event['status']}] {event['event_type']}")
    if event["status"] in ("COMPLETED", "FAILED", "TIMED_OUT"):
        break
```

### TypeScript

```bash
npm install @axme/axme
```

```typescript
import { AxmeClient } from "@axme/axme";

const client = new AxmeClient({ apiKey: process.env.AXME_API_KEY! });

const intentId = await client.sendIntent({
  intentType: "payment.process.v1",
  toAgent: "agent://myorg/production/payment-service",
  payload: {
    orderId: "ord_12345",
    amountCents: 9999,
    currency: "usd",
    customerEmail: "alice@example.com",
  },
});

console.log(`Payment submitted: ${intentId}`);

const result = await client.waitFor(intentId);
console.log(`Done: ${result.status}`);
```

---

## More Languages

Full implementations in all 5 languages:

| Language | Directory | Install |
|----------|-----------|---------|
| [Python](python/) | `python/` | `pip install axme` |
| [TypeScript](typescript/) | `typescript/` | `npm install @axme/axme` |
| [Go](go/) | `go/` | `go get github.com/AxmeAI/axme-sdk-go` |
| [Java](java/) | `java/` | Maven Central: `ai.axme:axme-sdk` |
| [.NET](dotnet/) | `dotnet/` | `dotnet add package Axme.Sdk` |

---

## Before / After

### Before: Webhook Infrastructure (250+ lines)

```python
@app.post("/webhooks/payment-complete")
async def payment_webhook(request: Request):
    # Step 1: Verify signature
    signature = request.headers.get("x-webhook-signature")
    body = await request.body()
    expected = hmac.new(WEBHOOK_SECRET, body, hashlib.sha256).hexdigest()
    if not hmac.compare_digest(signature, expected):
        raise HTTPException(401, "Invalid signature")

    # Step 2: Parse and deduplicate
    payload = json.loads(body)
    event_id = payload["event_id"]
    if redis.exists(f"processed:{event_id}"):
        return {"status": "duplicate"}  # Already processed
    redis.setex(f"processed:{event_id}", 86400, "1")

    # Step 3: Process (must return 200 within 5 seconds or provider retries)
    try:
        await process_payment(payload)
    except Exception:
        # Queue for retry... but what if the queue is down?
        dead_letter.put(payload)
        raise

    return {"status": "ok"}

# Plus: HTTPS endpoint, TLS certificates, firewall rules,
# dead letter queue consumer, monitoring for missing webhooks...
```

### After: AXME Intent Delivery (15 lines)

```python
from axme import AxmeClient, AxmeClientConfig

client = AxmeClient(AxmeClientConfig(api_key=os.environ["AXME_API_KEY"]))

intent_id = client.send_intent({
    "intent_type": "payment.process.v1",
    "to_agent": "agent://myorg/production/payment-service",
    "payload": {
        "order_id": "ord_12345",
        "amount_cents": 9999,
        "currency": "usd",
    },
})

result = client.wait_for(intent_id)
print(result["status"])  # COMPLETED, FAILED, or TIMED_OUT
```

No webhook endpoint. No signature verification. No idempotency layer. No dead letter queue. No retry logic.

---

## How It Works

```
┌────────────┐  send_intent()   ┌────────────────┐   deliver    ┌──────────────┐
│            │ ───────────────> │                │ (guaranteed) │              │
│   Order    │                  │   AXME Cloud   │ ──────────>  │   Payment    │
│   Service  │ <─ observe(SSE)  │   (platform)   │              │   Service    │
│            │                  │                │ <─ resume()  │   (agent)    │
└────────────┘                  │   retries,     │  with result │              │
                                │   delivery     │              │  processes   │
                                │   guarantees   │              │  payment     │
                                └────────────────┘              └──────────────┘

Before:                         After:
  Provider -> webhook -> you      You -> intent -> platform -> service
  (fire & pray)                   (guaranteed delivery + observability)
```

1. Order service submits a payment **intent** via AXME SDK
2. Platform **delivers** it to the payment service agent — with retries and delivery guarantees
3. Payment service processes and **resumes** with result (success/failure)
4. Order service **observes** the full lifecycle — no webhook endpoint required
5. If delivery fails, platform retries automatically — no dead letter queue needed

---

## Run the Full Example

### Prerequisites

```bash
# Install CLI (one-time)
curl -fsSL https://raw.githubusercontent.com/AxmeAI/axme-cli/main/install.sh | sh
# Open a new terminal, or run the "source" command shown by the installer

# Log in
axme login

# Install Python SDK
pip install axme
```

### Terminal 1 - submit the intent

```bash
axme scenarios apply scenario.json
# Note the intent_id in the output
```

### Terminal 2 - start the agent

Get the agent key after scenario apply:

```bash
# macOS
cat ~/Library/Application\ Support/axme/scenario-agents.json | grep -A2 payment-service-demo

# Linux
cat ~/.config/axme/scenario-agents.json | grep -A2 payment-service-demo
```

Then run the agent in your language of choice:

```bash
# Python (SSE stream listener)
AXME_API_KEY=<agent-key> python agent.py

# TypeScript (SSE stream listener, requires Node 20+)
cd typescript && npm install
AXME_API_KEY=<agent-key> npx tsx agent.ts

# Go (SSE stream listener)
cd go && go run ./cmd/agent/

# Java (processes a single intent by ID)
cd java/agent && mvn compile
AXME_API_KEY=<agent-key> mvn -q exec:java -Dexec.mainClass="Agent" -Dexec.args="<step-intent-id>"

# .NET (processes a single intent by ID)
cd dotnet/agent && dotnet run -- <step-intent-id>
```

### Verify

```bash
axme intents get <intent_id>
# lifecycle_status: COMPLETED
```

---

## Related

- [AXME](https://github.com/AxmeAI/axme) — project overview
- [AXP Spec](https://github.com/AxmeAI/axp-spec) — open Intent Protocol specification
- [AXME Examples](https://github.com/AxmeAI/axme-examples) — 20+ runnable examples across 5 languages
- [AXME CLI](https://github.com/AxmeAI/axme-cli) — manage intents, agents, scenarios from the terminal

---

Built with [AXME](https://github.com/AxmeAI/axme) (AXP Intent Protocol).
