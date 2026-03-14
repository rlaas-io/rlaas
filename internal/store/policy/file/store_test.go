package file

import (
	"context"
	"os"
	"path/filepath"
	"rlaas/internal/model"
	"testing"
)

func TestFileStoreCRUD(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "policies.json")
	s := New(path)
	p := model.Policy{PolicyID: "p1", Scope: model.PolicyScope{OrgID: "acme"}, Enabled: true}
	if err := s.UpsertPolicy(context.Background(), p); err != nil {
		t.Fatalf("upsert failed: %v", err)
	}
	all, err := s.LoadPolicies(context.Background(), "acme")
	if err != nil || len(all) != 1 {
		t.Fatalf("load failed")
	}
	one, err := s.GetPolicyByID(context.Background(), "p1")
	if err != nil || one.PolicyID != "p1" {
		t.Fatalf("get failed")
	}
	all2, err := s.ListPolicies(context.Background(), nil)
	if err != nil || len(all2) != 1 {
		t.Fatalf("list failed")
	}
	if err := s.DeletePolicy(context.Background(), "p1"); err != nil {
		t.Fatalf("delete failed")
	}
	if _, err := s.GetPolicyByID(context.Background(), "p1"); err == nil {
		t.Fatalf("expected missing policy")
	}
}

func TestFileStoreDirectArrayFallback(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "policies.json")
	_ = os.WriteFile(path, []byte(`[{"policy_id":"p2","enabled":true}]`), 0644)
	s := New(path)
	all, err := s.LoadPolicies(context.Background(), "")
	if err != nil || len(all) != 1 || all[0].PolicyID != "p2" {
		t.Fatalf("expected direct array fallback")
	}
}

func TestFileStoreEmptyAndInvalidFiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "policies.json")
	_ = os.WriteFile(path, []byte(""), 0644)
	s := New(path)
	if all, err := s.LoadPolicies(context.Background(), ""); err != nil || len(all) != 0 {
		t.Fatalf("expected empty policy list")
	}
	_ = os.WriteFile(path, []byte("not-json"), 0644)
	if _, err := s.LoadPolicies(context.Background(), ""); err == nil {
		t.Fatalf("expected invalid json error")
	}
}

func TestFileStoreUpsertUpdateExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "policies.json")
	s := New(path)
	p := model.Policy{PolicyID: "p1", Name: "first", Enabled: true}
	_ = s.UpsertPolicy(context.Background(), p)
	p.Name = "second"
	_ = s.UpsertPolicy(context.Background(), p)
	one, err := s.GetPolicyByID(context.Background(), "p1")
	if err != nil || one.Name != "second" {
		t.Fatalf("expected update existing policy")
	}
}
