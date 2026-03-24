// Payment service agent — .NET example.
//
// Fetches an intent by ID, processes payment, and resumes with transaction result.
//
// Usage:
//   export AXME_API_KEY="<agent-key>"
//   dotnet run -- <intent_id>

using Axme.Sdk;
using System.Text.Json.Nodes;

if (args.Length < 1)
{
    Console.Error.WriteLine("Usage: dotnet run -- <intent_id>");
    return 1;
}

var apiKey = Environment.GetEnvironmentVariable("AXME_API_KEY");
if (string.IsNullOrEmpty(apiKey))
{
    Console.Error.WriteLine("Error: AXME_API_KEY not set.");
    return 1;
}

var intentId = args[0];
var client = new AxmeClient(new AxmeClientConfig { ApiKey = apiKey });

Console.WriteLine($"Processing intent: {intentId}");

var intentData = await client.GetIntentAsync(intentId);
var intent = intentData["intent"]?.AsObject() ?? intentData;
var payload = intent["payload"]?.AsObject() ?? new JsonObject();
if (payload["parent_payload"] is JsonObject parentPayload)
{
    payload = parentPayload;
}

var orderId = payload["order_id"]?.ToString() ?? "unknown";
var amount = payload["amount"]?.GetValue<double>() ?? 0;
var currency = payload["currency"]?.ToString() ?? "USD";
var method = payload["method"]?.ToString() ?? "card";

Console.WriteLine($"  Processing {method} payment: {currency} {amount} for {orderId}...");
await Task.Delay(1000);
Console.WriteLine("  Authorizing with payment provider...");
await Task.Delay(1000);
Console.WriteLine("  Capturing funds...");
await Task.Delay(1000);

var result = new JsonObject
{
    ["action"] = "complete",
    ["order_id"] = orderId,
    ["transaction_id"] = "TXN-99001",
    ["amount_charged"] = amount,
    ["status"] = "captured",
    ["processed_at"] = DateTime.UtcNow.ToString("o")
};

await client.ResumeIntentAsync(intentId, result);
Console.WriteLine($"  Payment captured: TXN-99001 ({currency} {amount})");
return 0;
