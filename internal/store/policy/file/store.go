package file

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"rlaas/internal/model"
	"sync"
)

type payload struct {
	Policies []model.Policy `json:"policies"`
}

// Store is a json file backed policy store for local development.
type Store struct {
	path string
	mu   sync.RWMutex
}

// New creates a file policy store for the given path.
func New(path string) *Store {
	return &Store{path: path}
}

// LoadPolicies returns policies for the tenant or org namespace.
func (s *Store) LoadPolicies(_ context.Context, tenantOrOrg string) ([]model.Policy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	all, err := s.readAll()
	if err != nil {
		return nil, err
	}
	if tenantOrOrg == "" {
		return all, nil
	}
	out := make([]model.Policy, 0)
	for _, p := range all {
		if p.Scope.TenantID == tenantOrOrg || p.Scope.OrgID == tenantOrOrg || (p.Scope.TenantID == "" && p.Scope.OrgID == "") {
			out = append(out, p)
		}
	}
	return out, nil
}

// GetPolicyByID returns one policy by id.
func (s *Store) GetPolicyByID(_ context.Context, policyID string) (*model.Policy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	all, err := s.readAll()
	if err != nil {
		return nil, err
	}
	for _, p := range all {
		if p.PolicyID == policyID {
			cp := p
			return &cp, nil
		}
	}
	return nil, errors.New("policy not found")
}

// UpsertPolicy inserts or updates one policy.
func (s *Store) UpsertPolicy(_ context.Context, p model.Policy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	all, err := s.readAll()
	if err != nil {
		return err
	}
	updated := false
	for i := range all {
		if all[i].PolicyID == p.PolicyID {
			all[i] = p
			updated = true
			break
		}
	}
	if !updated {
		all = append(all, p)
	}
	return s.writeAll(all)
}

// DeletePolicy removes one policy by id.
func (s *Store) DeletePolicy(_ context.Context, policyID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	all, err := s.readAll()
	if err != nil {
		return err
	}
	out := make([]model.Policy, 0, len(all))
	for _, p := range all {
		if p.PolicyID != policyID {
			out = append(out, p)
		}
	}
	return s.writeAll(out)
}

// ListPolicies returns all policies in this file store.
func (s *Store) ListPolicies(_ context.Context, _ map[string]string) ([]model.Policy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.readAll()
}

func (s *Store) readAll() ([]model.Policy, error) {
	if _, err := os.Stat(s.path); os.IsNotExist(err) {
		return []model.Policy{}, nil
	}
	raw, err := os.ReadFile(s.path)
	if err != nil {
		return nil, err
	}
	if len(raw) == 0 {
		return []model.Policy{}, nil
	}
	var p payload
	if err := json.Unmarshal(raw, &p); err == nil {
		return p.Policies, nil
	}
	var direct []model.Policy
	if err := json.Unmarshal(raw, &direct); err != nil {
		return nil, err
	}
	return direct, nil
}

func (s *Store) writeAll(in []model.Policy) error {
	out, err := json.MarshalIndent(payload{Policies: in}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, out, 0644)
}
