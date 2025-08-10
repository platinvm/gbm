package gbm_test

type Packet struct {
	Id          int8   `bin:"BE"`
	PayloadSize int16  `bin:"LE"`
	Payload     []byte `bin:"len=PayloadSize"`
	MessageSize int64  `bin:"LE"`
	Message     string `bin:"enc=utf-8,len=MessageSize"`
}
