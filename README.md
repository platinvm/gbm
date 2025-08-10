# gbm — Generic Binary Marshaler for Go

`gbm` is a tiny, zero-dependency library for binary marshaling/unmarshaling of structs using tags.

- ✅ Int types: `int8/16/32/64`, `uint8/16/32/64` with `LE`/`BE`
- ✅ `[]byte` and `string` with `len=<FieldName>`
- ✅ UTF-8 strings by default (`enc=utf-8`)
- ✅ Length fields auto-set on marshal and enforced on unmarshal

## Install

```bash
go get github.com/platinvm/gbm@latest
```

## Usage

```go
package main

import (
  "bytes"
  "log"

  "github.com/platinvm/gbm"
)

type Packet struct {
  Id          int8   `bin:"BE"`
  PayloadSize int16  `bin:"LE"`
  Payload     []byte `bin:"len=PayloadSize"`
  MessageSize int64  `bin:"LE"`
  Message     string `bin:"enc=utf-8,len=MessageSize"`
}

func main() {
  codec, err := gbm.New[Packet]()
  if err != nil { log.Fatal(err) }

  var buf bytes.Buffer
  p := Packet{Id: 1, Payload: []byte{0xAA, 0xBB}, Message: "hi"}
  _, _ = codec.Marshal(&p, &buf)

  var q Packet
  _, _ = codec.Unmarshal(&q, &buf)
  log.Println(q.Id, q.Payload, q.Message)
}
```

### Tags

- Endianness for multi-byte integers: `LE` or `BE`
- Variable length fields:
  - `[]byte` — `len=<FieldName>`
  - `string` — `len=<FieldName>`, optional `enc=utf-8` (UTF-8 only)
- The referenced length field must appear **before** the variable-length field.

## Guarantees & Caveats

- Strings must be valid UTF-8 on marshal.
- Length fields are updated automatically if zero; otherwise must match.
- Only exported fields are processed.
- Embedded structs are supported (flattened), but variable-length references must still point to earlier fields.
