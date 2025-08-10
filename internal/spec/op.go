package spec

import (
    "encoding/binary"
    "reflect"
)

// Op describes how to handle a field during (un)marshal.
// Exported to allow cross-internal package use; still private to the module.
type Op struct {
    Name     string
    Index    []int
    Kind     reflect.Kind
    Order    binary.ByteOrder
    IsSlice  bool
    IsString bool
    LenRef   string
    LenIndex []int
    EncUTF8  bool
}
