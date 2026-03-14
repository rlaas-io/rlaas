export type CheckRequest = {
  request_id: string;
  org_id: string;
  tenant_id: string;
  signal_type: string;
  operation: string;
  endpoint: string;
  method: string;
  user_id?: string;
  api_key?: string;
  tags?: Record<string, string>;
};

export type Decision = {
  allowed: boolean;
  action: string;
  reason: string;
  remaining?: number;
  retry_after?: string;
};

export type Policy = {
  policy_id: string;
  name: string;
  enabled: boolean;
  priority: number;
  scope: Record<string, unknown>;
  algorithm: Record<string, unknown>;
  action: string;
  failure_mode?: string;
  enforcement_mode?: string;
  rollout_percent?: number;
};
