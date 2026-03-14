package metrics

import "testing"

func TestNewCollector(t *testing.T) {
	c := New()
	if c == nil {
		t.Fatalf("expected collector")
	}
	c.DecisionsTotal.Add(1)
	if c.DecisionsTotal.Load() != 1 {
		t.Fatalf("counter should increment")
	}
}
