/*
 * Replacing webhooks with AXME — Java example.
 *
 * Payment processing: submit a payment intent with delivery guarantees.
 * No webhook endpoint, no signature verification, no retry logic.
 *
 * Usage:
 *   export AXME_API_KEY="your-key"
 *   mvn compile exec:java -Dexec.mainClass="ReplacingWebhooks"
 */

import ai.axme.sdk.AxmeClient;
import ai.axme.sdk.AxmeClientConfig;
import java.util.Map;

public class ReplacingWebhooks {
    public static void main(String[] args) throws Exception {
        var client = new AxmeClient(
            AxmeClientConfig.builder()
                .apiKey(System.getenv("AXME_API_KEY"))
                .build()
        );

        // Submit payment — platform delivers with retries, no webhook needed
        String intentId = client.sendIntent(Map.of(
            "intent_type", "payment.process.v1",
            "to_agent", "agent://myorg/production/payment-service",
            "payload", Map.of(
                "order_id", "ord_12345",
                "amount_cents", 9999,
                "currency", "usd",
                "customer_email", "alice@example.com"
            )
        ));
        System.out.println("Payment submitted: " + intentId);

        // Wait for completion — no webhook callback needed
        var result = client.waitFor(intentId);
        System.out.println("Final status: " + result.getStatus());
    }
}
