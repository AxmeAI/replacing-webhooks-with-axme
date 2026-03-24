/*
 * Payment service agent — Java example.
 *
 * Fetches an intent by ID, processes payment, and resumes with transaction result.
 *
 * Usage:
 *   export AXME_API_KEY="<agent-key>"
 *   javac -cp axme-sdk.jar Agent.java
 *   java -cp .:axme-sdk.jar Agent <intent_id>
 */

import dev.axme.sdk.AxmeClient;
import dev.axme.sdk.AxmeClientConfig;
import dev.axme.sdk.RequestOptions;
import java.time.Instant;
import java.util.Map;

public class Agent {
    public static void main(String[] args) throws Exception {
        if (args.length < 1) {
            System.err.println("Usage: java Agent <intent_id>");
            System.exit(1);
        }

        String apiKey = System.getenv("AXME_API_KEY");
        if (apiKey == null || apiKey.isEmpty()) {
            System.err.println("Error: AXME_API_KEY not set.");
            System.exit(1);
        }

        String intentId = args[0];
        var client = new AxmeClient(AxmeClientConfig.forCloud(apiKey));

        System.out.println("Processing intent: " + intentId);

        var intentData = client.getIntent(intentId, new RequestOptions());
        @SuppressWarnings("unchecked")
        var intent = (Map<String, Object>) intentData.getOrDefault("intent", intentData);
        @SuppressWarnings("unchecked")
        var payload = (Map<String, Object>) intent.getOrDefault("payload", Map.of());
        if (payload.containsKey("parent_payload")) {
            @SuppressWarnings("unchecked")
            var pp = (Map<String, Object>) payload.get("parent_payload");
            payload = pp;
        }

        String orderId = (String) payload.getOrDefault("order_id", "unknown");
        double amount = payload.containsKey("amount") ? ((Number) payload.get("amount")).doubleValue() : 0;
        String currency = (String) payload.getOrDefault("currency", "USD");
        String method = (String) payload.getOrDefault("method", "card");

        System.out.println("  Processing " + method + " payment: " + currency + " " + (int) amount + " for " + orderId + "...");
        Thread.sleep(1000);
        System.out.println("  Authorizing with payment provider...");
        Thread.sleep(1000);
        System.out.println("  Capturing funds...");
        Thread.sleep(1000);

        var result = Map.<String, Object>of(
            "action", "complete",
            "order_id", orderId,
            "transaction_id", "TXN-99001",
            "amount_charged", amount,
            "status", "captured",
            "processed_at", Instant.now().toString()
        );

        client.resumeIntent(intentId, result, new RequestOptions());
        System.out.println("  Payment captured: TXN-99001 (" + currency + " " + (int) amount + ")");
    }
}
