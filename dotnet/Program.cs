// Replacing webhooks with AXME — .NET example.
//
// Payment processing: submit a payment intent with delivery guarantees.
// No webhook endpoint, no signature verification, no retry logic.
//
// Usage:
//   export AXME_API_KEY="your-key"
//   dotnet run

using Axme.Sdk;

var client = new AxmeClient(new AxmeClientConfig
{
    ApiKey = Environment.GetEnvironmentVariable("AXME_API_KEY")!
});

// Submit payment — platform delivers with retries, no webhook needed
var intentId = await client.SendIntentAsync(new
{
    intent_type = "payment.process.v1",
    to_agent = "agent://myorg/production/payment-service",
    payload = new
    {
        order_id = "ord_12345",
        amount_cents = 9999,
        currency = "usd",
        customer_email = "alice@example.com"
    }
});
Console.WriteLine($"Payment submitted: {intentId}");

// Wait for completion — no webhook callback needed
var result = await client.WaitForAsync(intentId);
Console.WriteLine($"Final status: {result.Status}");
