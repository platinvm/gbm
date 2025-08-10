package gbm_test

import (
	"bytes"
	"fmt"

	"github.com/platinvm/gbm"
)

func Example() {
	codec, _ := gbm.New[Packet]()
	p := Packet{Id: 1, Payload: []byte{1, 2, 3}, Message: "hi"}
	var buf bytes.Buffer
	_, _ = codec.Marshal(&p, &buf)
	var q Packet
	_, _ = codec.Unmarshal(&q, &buf)
	fmt.Println(q.Id, q.Payload, q.Message)
	// Output: 1 [1 2 3] hi
}
