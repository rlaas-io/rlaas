package otel

import (
	"context"
	"errors"
	"testing"

	"rlaas/internal/model"
)

type processorEvalStub struct {
	decision model.Decision
	err      error
}

func (s processorEvalStub) Evaluate(_ context.Context, _ model.RequestContext) (model.Decision, error) {
	return s.decision, s.err
}

func TestProcessorProcessLogsAndSpans(t *testing.T) {
	p := NewProcessor(processorEvalStub{decision: model.Decision{Allowed: true, Action: model.ActionAllow}}, 4, true)
	logs := []LogRecord{{OrgID: "o", Service: "s", Severity: "info"}, {OrgID: "o", Service: "s", Severity: "warn"}}
	outLogs := p.ProcessLogs(context.Background(), logs)
	if len(outLogs) != 2 {
		t.Fatalf("expected all logs allowed")
	}
	spans := []SpanRecord{{OrgID: "o", Service: "s", Name: "span1"}}
	outSpans := p.ProcessSpans(context.Background(), spans)
	if len(outSpans) != 1 {
		t.Fatalf("expected span allowed")
	}
	stats := p.Stats()
	if stats.Allowed != 3 || stats.Dropped != 0 || stats.Errors != 0 {
		t.Fatalf("unexpected stats: %+v", stats)
	}
}

func TestProcessorDropsDeniedAndFailClosedErrors(t *testing.T) {
	p := NewProcessor(processorEvalStub{decision: model.Decision{Allowed: false, Action: model.ActionDeny}}, 2, false)
	logs := []LogRecord{{OrgID: "o", Service: "s", Severity: "info"}}
	if out := p.ProcessLogs(context.Background(), logs); len(out) != 0 {
		t.Fatalf("expected denied logs to drop")
	}

	p2 := NewProcessor(processorEvalStub{err: errors.New("boom")}, 1, false)
	if out := p2.ProcessSpans(context.Background(), []SpanRecord{{OrgID: "o", Service: "s", Name: "n"}}); len(out) != 0 {
		t.Fatalf("expected fail-closed error drops")
	}
	stats := p2.Stats()
	if stats.Errors != 1 || stats.Dropped != 1 {
		t.Fatalf("unexpected fail-closed stats: %+v", stats)
	}
}

func TestDecisionFilter(t *testing.T) {
	if !DecisionFilter(model.Decision{Allowed: true, Action: model.ActionAllow}) {
		t.Fatalf("allow should pass")
	}
	if DecisionFilter(model.Decision{Allowed: false, Action: model.ActionAllow}) {
		t.Fatalf("not allowed should drop")
	}
	if DecisionFilter(model.Decision{Allowed: true, Action: model.ActionDrop}) {
		t.Fatalf("drop action should drop")
	}
}
