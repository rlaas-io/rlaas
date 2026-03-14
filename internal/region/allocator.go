package region

import "sort"

// RegionWeight defines one region's share weight.
type RegionWeight struct {
	Region string
	Weight int64
}

// Allocation describes computed regional limit allocation.
type Allocation struct {
	Region string `json:"region"`
	Limit  int64  `json:"limit"`
}

// AllocateGlobalLimit splits a global limit across regions proportionally.
func AllocateGlobalLimit(globalLimit int64, weights []RegionWeight) []Allocation {
	if globalLimit <= 0 || len(weights) == 0 {
		return nil
	}
	valid := make([]RegionWeight, 0, len(weights))
	var totalWeight int64
	for _, w := range weights {
		if w.Region == "" || w.Weight <= 0 {
			continue
		}
		valid = append(valid, w)
		totalWeight += w.Weight
	}
	if totalWeight == 0 {
		return nil
	}
	sort.Slice(valid, func(i, j int) bool { return valid[i].Region < valid[j].Region })
	out := make([]Allocation, 0, len(valid))
	var assigned int64
	for i, w := range valid {
		limit := globalLimit * w.Weight / totalWeight
		if i == len(valid)-1 {
			limit = globalLimit - assigned
		}
		if limit < 0 {
			limit = 0
		}
		out = append(out, Allocation{Region: w.Region, Limit: limit})
		assigned += limit
	}
	return out
}

// RegionalOverflow reports how much each region exceeded its allocated limit.
func RegionalOverflow(usage map[string]int64, allocation []Allocation) map[string]int64 {
	if len(usage) == 0 || len(allocation) == 0 {
		return map[string]int64{}
	}
	limits := make(map[string]int64, len(allocation))
	for _, a := range allocation {
		limits[a.Region] = a.Limit
	}
	out := map[string]int64{}
	for region, used := range usage {
		over := used - limits[region]
		if over > 0 {
			out[region] = over
		}
	}
	return out
}
