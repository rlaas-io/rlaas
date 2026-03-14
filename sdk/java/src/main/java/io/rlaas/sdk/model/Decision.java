package io.rlaas.sdk.model;

public class Decision {
    private boolean allowed;
    private String action;
    private String reason;
    private long remaining;
    private String retryAfter;

    public boolean isAllowed() { return allowed; }
    public void setAllowed(boolean allowed) { this.allowed = allowed; }
    public String getAction() { return action; }
    public void setAction(String action) { this.action = action; }
    public String getReason() { return reason; }
    public void setReason(String reason) { this.reason = reason; }
    public long getRemaining() { return remaining; }
    public void setRemaining(long remaining) { this.remaining = remaining; }
    public String getRetryAfter() { return retryAfter; }
    public void setRetryAfter(String retryAfter) { this.retryAfter = retryAfter; }
}
