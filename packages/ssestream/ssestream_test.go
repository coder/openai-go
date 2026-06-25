// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package ssestream

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

func newTestDecoder(body string) Decoder {
	return NewDecoder(&http.Response{
		Body:   io.NopCloser(strings.NewReader(body)),
		Header: http.Header{"Content-Type": []string{"text/event-stream"}},
	})
}

func TestEventStreamDecoderSkipsCommentOnlyBlocks(t *testing.T) {
	decoder := newTestDecoder(": OPENROUTER PROCESSING\n\ndata: {\"ok\":true}\n\n")
	defer decoder.Close()

	if !decoder.Next() {
		t.Fatalf("expected event, got err: %v", decoder.Err())
	}

	event := decoder.Event()
	if event.Type != "" {
		t.Fatalf("expected empty event type, got %q", event.Type)
	}
	if string(event.Data) != "{\"ok\":true}\n" {
		t.Fatalf("expected data event, got %q", event.Data)
	}

	if decoder.Next() {
		t.Fatalf("expected no more events, got %q", decoder.Event().Data)
	}
	if decoder.Err() != nil {
		t.Fatalf("expected no error, got %v", decoder.Err())
	}
}

func TestEventStreamDecoderSkipsEventWithoutData(t *testing.T) {
	decoder := newTestDecoder("event: ping\n\nevent: message\ndata: {\"ok\":true}\n\n")
	defer decoder.Close()

	if !decoder.Next() {
		t.Fatalf("expected event, got err: %v", decoder.Err())
	}

	event := decoder.Event()
	if event.Type != "message" {
		t.Fatalf("expected event type message, got %q", event.Type)
	}
	if string(event.Data) != "{\"ok\":true}\n" {
		t.Fatalf("expected data event, got %q", event.Data)
	}

	if decoder.Next() {
		t.Fatalf("expected no more events, got %q", decoder.Event().Data)
	}
	if decoder.Err() != nil {
		t.Fatalf("expected no error, got %v", decoder.Err())
	}
}

func TestEventStreamDecoderPreservesMultilineData(t *testing.T) {
	decoder := newTestDecoder("event: message\ndata: first\ndata: second\n\n")
	defer decoder.Close()

	if !decoder.Next() {
		t.Fatalf("expected event, got err: %v", decoder.Err())
	}

	event := decoder.Event()
	if event.Type != "message" {
		t.Fatalf("expected event type message, got %q", event.Type)
	}
	if string(event.Data) != "first\nsecond\n" {
		t.Fatalf("expected multiline data, got %q", event.Data)
	}
}
