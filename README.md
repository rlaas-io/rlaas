# RLAAS

RLAAS (Rate Limiting as a Service) is a policy-driven platform for applying rate limits, quotas, and traffic shaping across APIs, services, telemetry pipelines, and business operations.

It is designed for hybrid deployment:

- embedded in Go services for low-latency local decisions
- centralized decision/API mode for shared governance
- sidecar mode for Kubernetes and polyglot environments (roadmap)

## What this application is for

Use RLAAS when you want one reusable rate-limiting platform that can enforce policies by tenant, org, app, service, endpoint, user, client, or custom dimensions.

Typical use cases:

- API throttling (HTTP/gRPC)
- tenant/org quotas
- login/abuse protection
- telemetry shaping (logs/spans)
- outbound partner API protection
- job/workflow throughput control

## Core features

- Multi-dimensional policy model with precedence and priority
- Multiple algorithms (window, bucket, quota, concurrency)
- Rich actions beyond allow/deny (delay, sample, drop, shadow)
- Shadow mode and rollout percentage for safe policy adoption
- Fail-open / fail-closed policy behavior
- Pluggable policy and counter backends
- SDK + service integration model

## Quick start

### 1) Install dependencies

```bash
go mod tidy
```

### 2) Run the server

```bash
go run ./cmd/rlaas-server
```

Optional env var:

- `RLAAS_POLICY_FILE` (default: `examples/policies.json`)

### 3) Check a decision

```powershell
$body = @{
  request_id = "r1"
  org_id = "acme"
  tenant_id = "retail"
  service = "payments"
  signal_type = "http"
  operation = "charge"
  endpoint = "/v1/charge"
  method = "POST"
  user_id = "u1"
} | ConvertTo-Json

Invoke-RestMethod -Method Post -Uri "http://localhost:8080/v1/check" -ContentType "application/json" -Body $body
```

### 4) Run tests

```bash
go test ./...
```

## Current support status

### Algorithms

- Supported now:
  - fixed window
  - token bucket
  - sliding window counter
  - concurrency limiter
  - quota limiter
- Available as base implementations (early):
  - sliding window log
  - leaky bucket

### Actions

- Supported now:
  - allow
  - deny
  - delay
  - sample
  - drop
  - shadow-only
- Planned expansion:
  - downgrade
  - drop_low_priority advanced routing policies

### Integrations

- Supported now:
  - Go SDK API (`Evaluate`, concurrency lease API)
  - HTTP check endpoint (`POST /v1/check`)
  - HTTP middleware
- In progress / scaffolded:
  - gRPC interceptor and proto contracts
  - OTEL hook abstraction
  - sidecar/agent mode

### Backend support

Policy backend status:

- Supported now:
  - file-based JSON policy store
- In progress:
  - PostgreSQL policy store
  - Oracle policy store

Counter backend status:

- Supported now:
  - in-memory counter store
  - Redis counter store
- In progress:
  - PostgreSQL counter store
  - Oracle counter store

## Language and platform support

- Backend/service language supported today: Go
- Non-Go service support:
  - possible through HTTP decision API today
  - gRPC service mode is still being completed
- Native non-Go SDKs: not yet available (planned)

## Future scope (roadmap)

- Production-grade PostgreSQL/Oracle persistence layers
- Full centralized gRPC decision service
- Sidecar mode for Kubernetes deployments
- Policy audit/versioning and rollout control-plane APIs
- Broader language SDK support (polyglot clients)
- Enhanced observability and analytics
- Advanced policy expression support

## Project maturity note

RLAAS is currently in MVP-to-early platform stage: ready for Go-first adoption and integration testing, with planned expansion toward full polyglot and enterprise control-plane capabilities.
