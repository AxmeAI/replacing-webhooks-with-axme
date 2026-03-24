/**
 * Replacing webhooks with AXME — TypeScript example.
 *
 * Payment processing: submit a payment intent with delivery guarantees.
 * No webhook endpoint, no signature verification, no retry logic.
 *
 * Usage:
 *   npm install @axme/axme
 *   export AXME_API_KEY="your-key"
 *   npx tsx main.ts
 */

import { AxmeClient } from "@axme/axme";

async function main() {
  const client = new AxmeClient({ apiKey: process.env.AXME_API_KEY! });

  // Submit payment — platform delivers with retries, no webhook needed
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

  // Wait for completion — no webhook callback needed
  const result = await client.waitFor(intentId);
  console.log(`Final status: ${result.status}`);
}

main().catch(console.error);
