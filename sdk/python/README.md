# RLAAS Python SDK

Lightweight client for RLAAS HTTP APIs.

## Install (local)

```bash
pip install -e .
```

## Example

```python
from rlaas_sdk import RlaasClient, CheckRequest

client = RlaasClient("http://localhost:8080")

decision = client.check_limit(CheckRequest(
    request_id="req-1",
    org_id="acme",
    tenant_id="retail",
    signal_type="http",
    operation="charge",
    endpoint="/v1/charge",
    method="POST",
    user_id="u1",
))

print(decision.allowed, decision.action, decision.reason)
```
