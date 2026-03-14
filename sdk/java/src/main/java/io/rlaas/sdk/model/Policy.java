package io.rlaas.sdk.model;

import java.util.Map;

public class Policy {
    private String policyId;
    private String name;
    private boolean enabled;
    private int priority;
    private Map<String, Object> scope;
    private Map<String, Object> algorithm;
    private String action;
    private String failureMode = "fail_open";
    private String enforcementMode = "enforce";
    private int rolloutPercent = 100;

    public String getPolicyId() { return policyId; }
    public void setPolicyId(String policyId) { this.policyId = policyId; }
    public String getName() { return name; }
    public void setName(String name) { this.name = name; }
    public boolean isEnabled() { return enabled; }
    public void setEnabled(boolean enabled) { this.enabled = enabled; }
    public int getPriority() { return priority; }
    public void setPriority(int priority) { this.priority = priority; }
    public Map<String, Object> getScope() { return scope; }
    public void setScope(Map<String, Object> scope) { this.scope = scope; }
    public Map<String, Object> getAlgorithm() { return algorithm; }
    public void setAlgorithm(Map<String, Object> algorithm) { this.algorithm = algorithm; }
    public String getAction() { return action; }
    public void setAction(String action) { this.action = action; }
    public String getFailureMode() { return failureMode; }
    public void setFailureMode(String failureMode) { this.failureMode = failureMode; }
    public String getEnforcementMode() { return enforcementMode; }
    public void setEnforcementMode(String enforcementMode) { this.enforcementMode = enforcementMode; }
    public int getRolloutPercent() { return rolloutPercent; }
    public void setRolloutPercent(int rolloutPercent) { this.rolloutPercent = rolloutPercent; }
}
