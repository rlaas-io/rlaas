package region

import "testing"

func TestAllocateGlobalLimit(t *testing.T) {
	alloc := AllocateGlobalLimit(10000, []RegionWeight{{Region: "US", Weight: 5}, {Region: "EU", Weight: 3}, {Region: "APAC", Weight: 2}})
	if len(alloc) != 3 {
		t.Fatalf("expected 3 allocations")
	}
	var total int64
	for _, a := range alloc {
		total += a.Limit
	}
	if total != 10000 {
		t.Fatalf("expected exact total allocation, got %d", total)
	}
}

func TestAllocateGlobalLimitEdgeCases(t *testing.T) {
	if got := AllocateGlobalLimit(0, []RegionWeight{{Region: "US", Weight: 1}}); len(got) != 0 {
		t.Fatalf("expected empty allocation for zero limit")
	}
	if got := AllocateGlobalLimit(100, []RegionWeight{{Region: "", Weight: 1}, {Region: "US", Weight: 0}}); len(got) != 0 {
		t.Fatalf("expected empty allocation for invalid weights")
	}
}

func TestRegionalOverflow(t *testing.T) {
	alloc := []Allocation{{Region: "US", Limit: 5000}, {Region: "EU", Limit: 3000}, {Region: "APAC", Limit: 2000}}
	over := RegionalOverflow(map[string]int64{"US": 5200, "EU": 2900, "APAC": 2200}, alloc)
	if over["US"] != 200 || over["APAC"] != 200 {
		t.Fatalf("unexpected overflow map: %+v", over)
	}
	if _, ok := over["EU"]; ok {
		t.Fatalf("eu should not overflow")
	}
}
