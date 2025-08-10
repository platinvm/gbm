package gbm_test

import (
	"bytes"
	"testing"

	"github.com/platinvm/gbm"
)

func BenchmarkMarshalUnmarshal(b *testing.B) {
	codec, _ := gbm.New[Packet]()
	p := Packet{
		Id:      1,
		Payload: make([]byte, 1024),
		Message: "hello world",
	}
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		_, _ = codec.Marshal(&p, &buf)
		var q Packet
		_, _ = codec.Unmarshal(&q, &buf)
	}
}
