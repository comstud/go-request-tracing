package request_tracing

import (
	"net/http"
	"testing"
)

func TestNoHeaders(t *testing.T) {
	req := &http.Request{}
	rt := NewRequestTraceFromHTTPRequest(req)

	span_id := rt.GetSpanID()
	trace_id := rt.GetTraceID()

	if len(span_id) == 0 {
		t.Error("No Span ID set")
	}

	if len(trace_id) == 0 {
		t.Error("No Trace ID set")
	}

	if span_id != trace_id {
		t.Error("Span ID != Trace ID: '%s' != '%s'", span_id, trace_id)
	}
}

func TestTraceIDHeader(t *testing.T) {
	test_trace_id := "req_1234"

	req := &http.Request{
		Header: http.Header{
			"X-Trace-Id": []string{test_trace_id},
		},
	}

	rt := NewRequestTraceFromHTTPRequest(req)

	span_id := rt.GetSpanID()
	orig_span_id := rt.GetOriginalSpanID()
	trace_id := rt.GetTraceID()

	if trace_id != test_trace_id {
		t.Errorf("Trace ID mismatch: '%s' != '%s'", trace_id, test_trace_id)
	}

	if len(span_id) == 0 {
		t.Error("No Span ID set")
	}

	if span_id == trace_id {
		t.Error("Span ID and Trace ID should not be the same but are")
	}

	if orig_span_id != span_id {
		t.Error("Original Span ID doesn't match Span ID: '%s' != '%s'",
			orig_span_id,
			span_id,
		)
	}
}

func TestTraceIDHeaderPreference(t *testing.T) {
	test_trace_id := "req_1234"
	test_req_id := "req_5678"
	test_ctreq_id := "req_9012"

	req := &http.Request{
		Header: http.Header{
			"X-Trace-Id":            []string{test_trace_id},
			"X-Request-Id":          []string{test_req_id},
			"X-Crowdtilt-Requestid": []string{test_ctreq_id},
		},
	}

	rt := NewRequestTraceFromHTTPRequest(req)

	trace_id := rt.GetTraceID()
	if trace_id != test_trace_id {
		t.Errorf("Trace ID mismatch: '%s' != '%s'", trace_id, test_trace_id)
	}

	req = &http.Request{
		Header: http.Header{
			"X-Request-Id":          []string{test_req_id},
			"X-Crowdtilt-Requestid": []string{test_ctreq_id},
		},
	}

	rt = NewRequestTraceFromHTTPRequest(req)

	trace_id = rt.GetTraceID()
	if trace_id != test_req_id {
		t.Errorf("Trace ID != Request-Id: '%s' != '%s'", trace_id, test_req_id)
	}

	req = &http.Request{
		Header: http.Header{
			"X-Crowdtilt-Requestid": []string{test_ctreq_id},
		},
	}

	rt = NewRequestTraceFromHTTPRequest(req)

	trace_id = rt.GetTraceID()
	if trace_id != test_ctreq_id {
		t.Errorf("Trace ID != Crowdtilt-RequestId: '%s' != '%s'", trace_id, test_ctreq_id)
	}
}

func TestSpanIDHeader(t *testing.T) {
	test_span_id := "req_1234"

	req := &http.Request{
		Header: http.Header{
			"X-Span-Id": []string{test_span_id},
		},
	}

	rt := NewRequestTraceFromHTTPRequest(req)

	span_id := rt.GetSpanID()
	orig_span_id := rt.GetOriginalSpanID()
	trace_id := rt.GetTraceID()

	if orig_span_id != test_span_id {
		t.Errorf("Original Span ID mismatch: '%s' != '%s'", orig_span_id, test_span_id)
	}

	if span_id == orig_span_id {
		t.Errorf("Span ID not new: '%s' == '%s''", span_id, orig_span_id)
	}

	if len(span_id) == 0 {
		t.Error("No Span ID set")
	}

	if span_id != trace_id {
		t.Error("Span ID != Trace ID: '%s' != '%s'", span_id, trace_id)
	}
}

func TestTraceAndSpanIDHeaders(t *testing.T) {
	test_span_id := "req_1234"
	test_trace_id := "req_5678"

	req := &http.Request{
		Header: http.Header{
			"X-Span-Id":  []string{test_span_id},
			"X-Trace-Id": []string{test_trace_id},
		},
	}

	rt := NewRequestTraceFromHTTPRequest(req)

	span_id := rt.GetSpanID()
	orig_span_id := rt.GetOriginalSpanID()
	trace_id := rt.GetTraceID()

	if orig_span_id != test_span_id {
		t.Errorf("Original Span ID mismatch: '%s' != '%s'", orig_span_id, test_span_id)
	}

	if span_id == orig_span_id {
		t.Errorf("Span ID not new: '%s' == '%s''", span_id, orig_span_id)
	}

	if len(span_id) == 0 {
		t.Error("No Span ID set")
	}

	if span_id == trace_id {
		t.Error("Span ID == Trace ID: '%s' == '%s'", span_id, trace_id)
	}
}
