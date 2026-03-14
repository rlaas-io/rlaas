# RLAAS .NET SDK

## Build

```bash
dotnet build ./Rlaas.Sdk/Rlaas.Sdk.csproj
```

## Example

```csharp
using Rlaas.Sdk;
using Rlaas.Sdk.Models;

var client = new RlaasClient("http://localhost:8080");

var decision = await client.CheckLimitAsync(new CheckRequest
{
    RequestId = "req-1",
    OrgId = "acme",
    TenantId = "retail",
    SignalType = "http",
    Operation = "charge",
    Endpoint = "/v1/charge",
    Method = "POST",
    UserId = "u1"
});

Console.WriteLine($"{decision.Allowed} {decision.Action} {decision.Reason}");
```
