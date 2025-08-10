package gbm_test

import (
	"bytes"
	"testing"

	"github.com/platinvm/gbm"
)

func TestRoundTrip(t *testing.T) {
	codec, err := gbm.New[Packet]()
	if err != nil {
		t.Fatal(err)
	}

	p := Packet{
		Id:      7,
		Payload: []byte{0xAA, 0xBB, 0xCC, 0xDD},
		Message: "hello!",
	}

	var buf bytes.Buffer
	n1, err := codec.Marshal(&p, &buf)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if got := int(p.PayloadSize); got != len(p.Payload) {
		t.Fatalf("payload size not set, got %d want %d", got, len(p.Payload))
	}
	if got := int(p.MessageSize); got != len(p.Message) {
		t.Fatalf("message size not set, got %d want %d", got, len(p.Message))
	}

	var q Packet
	n2, err := codec.Unmarshal(&q, &buf)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if n1 != n2 {
		t.Fatalf("byte count mismatch: wrote %d read %d", n1, n2)
	}
	if q.Id != p.Id || q.Message != p.Message {
		t.Fatalf("roundtrip mismatch: %#v vs %#v", q, p)
	}
	if len(q.Payload) != len(p.Payload) {
		t.Fatalf("payload len mismatch: %d vs %d", len(q.Payload), len(p.Payload))
	}
	for i := range p.Payload {
		if q.Payload[i] != p.Payload[i] {
			t.Fatalf("payload mismatch at %d", i)
		}
	}
}

func TestLenMismatchOnMarshal(t *testing.T) {
	codec, _ := gbm.New[Packet]()
	p := Packet{
		Id:          1,
		PayloadSize: 99, // wrong on purpose
		Payload:     []byte{1, 2, 3},
		Message:     "ok",
	}
	var buf bytes.Buffer
	if _, err := codec.Marshal(&p, &buf); err == nil {
		t.Fatal("expected error on length mismatch")
	}
}
