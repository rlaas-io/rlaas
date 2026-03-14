package matcher

import (
	"rlaas/internal/model"
	"testing"
)

func TestSelectWinnerSpecificity(t *testing.T) {
	m := New()
	req := model.RequestContext{OrgID: "acme", TenantID: "retail", Service: "payments", UserID: "u1", SignalType: "http"}
	policies := []model.Policy{
		{PolicyID: "org", Enabled: true, Priority: 10, Scope: model.PolicyScope{OrgID: "acme", SignalType: "http"}},
		{PolicyID: "user", Enabled: true, Priority: 1, Scope: model.PolicyScope{OrgID: "acme", UserID: "u1", SignalType: "http"}},
	}
	matched, err := m.Match(req, policies)
	if err != nil {
		t.Fatalf("match failed: %v", err)
	}
	winner, err := m.SelectWinner(req, matched)
	if err != nil {
		t.Fatalf("select failed: %v", err)
	}
	if winner.PolicyID != "user" {
		t.Fatalf("expected user policy, got %s", winner.PolicyID)
	}
}

func TestSelectWinnerErrorAndTieBreak(t *testing.T) {
	m := New()
	if _, err := m.SelectWinner(model.RequestContext{}, nil); err == nil {
		t.Fatalf("expected error when no policies")
	}
	policies := []model.Policy{
		{PolicyID: "a", Enabled: true, Priority: 1, Scope: model.PolicyScope{OrgID: "acme"}},
		{PolicyID: "b", Enabled: true, Priority: 1, Scope: model.PolicyScope{OrgID: "acme"}},
	}
	w, err := m.SelectWinner(model.RequestContext{}, policies)
	if err != nil || w.PolicyID != "b" {
		t.Fatalf("expected deterministic policy id tie-break")
	}
}

func TestMatchTagMismatch(t *testing.T) {
	m := New()
	req := model.RequestContext{OrgID: "acme", Tags: map[string]string{"env": "dev"}}
	policies := []model.Policy{{PolicyID: "p", Scope: model.PolicyScope{OrgID: "acme", Tags: map[string]string{"env": "prod"}}}}
	matched, err := m.Match(req, policies)
	if err != nil {
		t.Fatal(err)
	}
	if len(matched) != 0 {
		t.Fatalf("expected no match due to tag mismatch")
	}
}
