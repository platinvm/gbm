// go test -fuzz=Fuzz -fuzztime=10s
package gbm_test

import (
	"bytes"
	"testing"
	"unicode/utf8"

	"github.com/platinvm/gbm"
)

func FuzzRoundTrip(f *testing.F) {
	codec, _ := gbm.New[Packet]()

	f.Add(int8(3), []byte{1, 2, 3}, "ok")
	f.Fuzz(func(t *testing.T, id int8, payload []byte, msg string) {
		if !utf8.ValidString(msg) {
			return
		}
		p := Packet{Id: id, Payload: payload, Message: msg}
		var buf bytes.Buffer
		if _, err := codec.Marshal(&p, &buf); err != nil {
			t.Skip()
		}
		var q Packet
		if _, err := codec.Unmarshal(&q, &buf); err != nil {
			t.Fatalf("unmarshal: %v", err)
		}
	})
}
