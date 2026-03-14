# RLAAS TypeScript SDK

## Build

```bash
npm install
npm run build
```

## Example

```ts
import { RlaasClient } from "./dist";

const client = new RlaasClient("http://localhost:8080");

async function run() {
  const decision = await client.checkLimit({
    request_id: "req-1",
    org_id: "acme",
    tenant_id: "retail",
    signal_type: "http",
    operation: "charge",
    endpoint: "/v1/charge",
    method: "POST",
    user_id: "u1",
  });
  console.log(decision.allowed, decision.action, decision.reason);
}

run();
```
