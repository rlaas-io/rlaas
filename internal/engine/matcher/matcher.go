package matcher

import (
	"errors"
	"rlaas/internal/model"
	"sort"
)

// Matcher finds policies that apply to a request and picks the final winner.
type Matcher interface {
	Match(req model.RequestContext, policies []model.Policy) ([]model.Policy, error)
	SelectWinner(req model.RequestContext, policies []model.Policy) (*model.Policy, error)
}

// DefaultMatcher uses strict field matching and deterministic tie breaking.
type DefaultMatcher struct{}

// New creates the default policy matcher.
func New() *DefaultMatcher {
	return &DefaultMatcher{}
}

// Match returns all policies whose scope matches the request.
func (m *DefaultMatcher) Match(req model.RequestContext, policies []model.Policy) ([]model.Policy, error) {
	matched := make([]model.Policy, 0, len(policies))
	for _, policy := range policies {
		if matchesScope(req, policy.Scope) {
			matched = append(matched, policy)
		}
	}
	return matched, nil
}

// SelectWinner chooses one policy using specificity, priority, and policy id.
func (m *DefaultMatcher) SelectWinner(_ model.RequestContext, policies []model.Policy) (*model.Policy, error) {
	if len(policies) == 0 {
		return nil, errors.New("no matching policy")
	}
	sort.SliceStable(policies, func(i, j int) bool {
		si, sj := specificityScore(policies[i].Scope), specificityScore(policies[j].Scope)
		if si != sj {
			return si > sj
		}
		if policies[i].Priority != policies[j].Priority {
			return policies[i].Priority > policies[j].Priority
		}
		if policies[i].PolicyID != policies[j].PolicyID {
			return policies[i].PolicyID > policies[j].PolicyID
		}
		return false
	})
	winner := policies[0]
	return &winner, nil
}

// matchesScope checks every configured scope field and required tags.
func matchesScope(req model.RequestContext, s model.PolicyScope) bool {
	if !matchString(s.OrgID, req.OrgID) || !matchString(s.TenantID, req.TenantID) || !matchString(s.Application, req.Application) || !matchString(s.Service, req.Service) || !matchString(s.Environment, req.Environment) || !matchString(s.SignalType, req.SignalType) || !matchString(s.Operation, req.Operation) || !matchString(s.Endpoint, req.Endpoint) || !matchString(s.Method, req.Method) || !matchString(s.UserID, req.UserID) || !matchString(s.APIKey, req.APIKey) || !matchString(s.ClientID, req.ClientID) || !matchString(s.SourceIP, req.SourceIP) || !matchString(s.Region, req.Region) || !matchString(s.Resource, req.Resource) || !matchString(s.Severity, req.Severity) || !matchString(s.SpanName, req.SpanName) || !matchString(s.Topic, req.Topic) || !matchString(s.ConsumerGroup, req.ConsumerGroup) || !matchString(s.JobType, req.JobType) {
		return false
	}
	for k, v := range s.Tags {
		if req.Tags[k] != v {
			return false
		}
	}
	return true
}

// matchString treats empty scope values as wildcards.
func matchString(scope, val string) bool {
	return scope == "" || scope == val
}

// specificityScore ranks policies by how many concrete scope fields they set.
func specificityScore(s model.PolicyScope) int {
	score := 0
	weights := []string{s.UserID, s.APIKey, s.ClientID, s.Endpoint, s.Method, s.Operation, s.Service, s.Application, s.TenantID, s.OrgID, s.SignalType, s.Environment, s.SourceIP, s.Region, s.Resource, s.Severity, s.SpanName, s.Topic, s.ConsumerGroup, s.JobType}
	for _, field := range weights {
		if field != "" {
			score += 10
		}
	}
	score += len(s.Tags)
	return score
}
