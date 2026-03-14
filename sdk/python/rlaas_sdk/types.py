from dataclasses import dataclass, asdict
from typing import Any, Dict, Optional


@dataclass
class CheckRequest:
    request_id: str
    org_id: str
    tenant_id: str
    signal_type: str
    operation: str
    endpoint: str
    method: str
    user_id: Optional[str] = None
    api_key: Optional[str] = None
    tags: Optional[Dict[str, str]] = None

    def to_dict(self) -> Dict[str, Any]:
        payload = asdict(self)
        return {k: v for k, v in payload.items() if v is not None}


@dataclass
class Decision:
    allowed: bool
    action: str
    reason: str
    remaining: int = 0
    retry_after: str = ""


@dataclass
class Policy:
    policy_id: str
    name: str
    enabled: bool
    priority: int
    scope: Dict[str, Any]
    algorithm: Dict[str, Any]
    action: str
    failure_mode: str = "fail_open"
    enforcement_mode: str = "enforce"
    rollout_percent: int = 100

    def to_dict(self) -> Dict[str, Any]:
        return asdict(self)
