/**
 * Payment service agent — TypeScript example.
 *
 * Processes payments and resumes with transaction result.
 *
 * Usage:
 *   export AXME_API_KEY="<agent-key>"
 *   npx tsx agent.ts
 */

import { AxmeClient } from "@axme/axme";

const AGENT_ADDRESS = "payment-service-demo";

async function handleIntent(client: AxmeClient, intentId: string) {
  const intentData = await client.getIntent(intentId);
  const intent = intentData.intent ?? intentData;
  let payload = intent.payload ?? {};
  if (payload.parent_payload) {
    payload = payload.parent_payload;
  }

  const orderId = payload.order_id ?? "unknown";
  const amount = payload.amount ?? 0;
  const currency = payload.currency ?? "USD";
  const method = payload.method ?? "card";

  console.log(`  Processing ${method} payment: ${currency} ${amount} for ${orderId}...`);
  await new Promise((r) => setTimeout(r, 1000));
  console.log(`  Authorizing with payment provider...`);
  await new Promise((r) => setTimeout(r, 1000));
  console.log(`  Capturing funds...`);
  await new Promise((r) => setTimeout(r, 1000));

  const result = {
    action: "complete",
    order_id: orderId,
    transaction_id: "TXN-99001",
    amount_charged: amount,
    status: "captured",
    processed_at: new Date().toISOString(),
  };

  await client.resumeIntent(intentId, result, { ownerAgent: "payment-service-demo" });
  console.log(`  Payment captured: TXN-99001 (${currency} ${amount})`);
}

async function main() {
  const apiKey = process.env.AXME_API_KEY;
  if (!apiKey) {
    console.error("Error: AXME_API_KEY not set.");
    process.exit(1);
  }

  const client = new AxmeClient({ apiKey });

  console.log(`Agent listening on ${AGENT_ADDRESS}...`);
  console.log("Waiting for intents (Ctrl+C to stop)\n");

  for await (const delivery of client.listen(AGENT_ADDRESS)) {
    const intentId = delivery.intent_id;
    const status = delivery.status;
    if (intentId && ["DELIVERED", "CREATED", "IN_PROGRESS"].includes(status)) {
      console.log(`[${status}] Intent received: ${intentId}`);
      try {
        await handleIntent(client, intentId);
      } catch (e) {
        console.error(`  Error: ${e}`);
      }
    }
  }
}

main().catch(console.error);
