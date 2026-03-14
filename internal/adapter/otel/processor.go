package otel

import (
	"context"
	"sync"
	"sync/atomic"

	"rlaas/internal/model"
)

// LogRecord is a lightweight telemetry log model for processor-style filtering.
type LogRecord struct {
	OrgID    string
	Service  string
	Severity string
	Tags     map[string]string
	Body     string
}

// SpanRecord is a lightweight telemetry span model for processor-style filtering.
type SpanRecord struct {
	OrgID   string
	Service string
	Name    string
	Tags    map[string]string
}

// ProcessorStats tracks allow/drop/error counters for OTEL pipeline processing.
type ProcessorStats struct {
	Allowed int64 `json:"allowed"`
	Dropped int64 `json:"dropped"`
	Errors  int64 `json:"errors"`
}

// Processor applies RLAAS decisions to telemetry batches.
type Processor struct {
	hook     Hook
	workers  int
	failOpen bool
	allowCnt atomic.Int64
	dropCnt  atomic.Int64
	errorCnt atomic.Int64
}

// NewProcessor creates an OTEL processor with bounded workers.
func NewProcessor(eval Evaluator, workers int, failOpen bool) *Processor {
	if workers <= 0 {
		workers = 1
	}
	return &Processor{hook: Hook{Eval: eval}, workers: workers, failOpen: failOpen}
}

// ProcessLogs applies rate limiting on a log batch and returns kept records.
func (p *Processor) ProcessLogs(ctx context.Context, logs []LogRecord) []LogRecord {
	if len(logs) == 0 {
		return nil
	}
	jobs := make(chan int, len(logs))
	for i := range logs {
		jobs <- i
	}
	close(jobs)

	keep := make([]bool, len(logs))
	for i := range keep {
		keep[i] = true
	}
	var wg sync.WaitGroup
	for i := 0; i < p.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				logRec := logs[idx]
				allowed, _, err := p.hook.AllowLog(ctx, logRec.OrgID, logRec.Service, logRec.Severity, logRec.Tags)
				if err != nil {
					p.errorCnt.Add(1)
					if !p.failOpen {
						keep[idx] = false
						p.dropCnt.Add(1)
						continue
					}
				}
				if !allowed {
					keep[idx] = false
					p.dropCnt.Add(1)
					continue
				}
				p.allowCnt.Add(1)
			}
		}()
	}
	wg.Wait()

	out := make([]LogRecord, 0, len(logs))
	for i, ok := range keep {
		if ok {
			out = append(out, logs[i])
		}
	}
	return out
}

// ProcessSpans applies rate limiting on a span batch and returns kept records.
func (p *Processor) ProcessSpans(ctx context.Context, spans []SpanRecord) []SpanRecord {
	if len(spans) == 0 {
		return nil
	}
	jobs := make(chan int, len(spans))
	for i := range spans {
		jobs <- i
	}
	close(jobs)

	keep := make([]bool, len(spans))
	for i := range keep {
		keep[i] = true
	}
	var wg sync.WaitGroup
	for i := 0; i < p.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for idx := range jobs {
				s := spans[idx]
				allowed, _, err := p.hook.AllowSpan(ctx, s.OrgID, s.Service, s.Name, s.Tags)
				if err != nil {
					p.errorCnt.Add(1)
					if !p.failOpen {
						keep[idx] = false
						p.dropCnt.Add(1)
						continue
					}
				}
				if !allowed {
					keep[idx] = false
					p.dropCnt.Add(1)
					continue
				}
				p.allowCnt.Add(1)
			}
		}()
	}
	wg.Wait()

	out := make([]SpanRecord, 0, len(spans))
	for i, ok := range keep {
		if ok {
			out = append(out, spans[i])
		}
	}
	return out
}

// Stats returns aggregate processor counters.
func (p *Processor) Stats() ProcessorStats {
	return ProcessorStats{Allowed: p.allowCnt.Load(), Dropped: p.dropCnt.Load(), Errors: p.errorCnt.Load()}
}

// DecisionFilter reports whether to keep one decision based on action.
func DecisionFilter(d model.Decision) bool {
	if !d.Allowed {
		return false
	}
	switch d.Action {
	case model.ActionDrop, model.ActionDropLowPriority, model.ActionDeny:
		return false
	default:
		return true
	}
}
