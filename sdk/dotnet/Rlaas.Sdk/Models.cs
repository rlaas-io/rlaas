namespace Rlaas.Sdk.Models;

public sealed class CheckRequest
{
    public string RequestId { get; set; } = string.Empty;
    public string OrgId { get; set; } = string.Empty;
    public string TenantId { get; set; } = string.Empty;
    public string SignalType { get; set; } = string.Empty;
    public string Operation { get; set; } = string.Empty;
    public string Endpoint { get; set; } = string.Empty;
    public string Method { get; set; } = string.Empty;
    public string? UserId { get; set; }
    public string? ApiKey { get; set; }
    public Dictionary<string, string>? Tags { get; set; }
}

public sealed class Decision
{
    public bool Allowed { get; set; }
    public string Action { get; set; } = string.Empty;
    public string Reason { get; set; } = string.Empty;
    public long Remaining { get; set; }
    public string RetryAfter { get; set; } = string.Empty;
}

public sealed class Policy
{
    public string PolicyId { get; set; } = string.Empty;
    public string Name { get; set; } = string.Empty;
    public bool Enabled { get; set; }
    public int Priority { get; set; }
    public Dictionary<string, object> Scope { get; set; } = new();
    public Dictionary<string, object> Algorithm { get; set; } = new();
    public string Action { get; set; } = string.Empty;
    public string FailureMode { get; set; } = "fail_open";
    public string EnforcementMode { get; set; } = "enforce";
    public int RolloutPercent { get; set; } = 100;
}
