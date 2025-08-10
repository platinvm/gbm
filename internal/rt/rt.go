package rt

import (
    "encoding/binary"
    "errors"
    "fmt"
    "io"
    "reflect"
)

func IsBefore(a, b []int) bool {
    for i := 0; i < len(a) && i < len(b); i++ {
        if a[i] < b[i] { return true }
        if a[i] > b[i] { return false }
    }
    return len(a) < len(b)
}

type CountWriter struct {
    W io.Writer
    N int64
}

func (cw *CountWriter) Write(p []byte) (int, error) {
    n, err := cw.W.Write(p)
    cw.N += int64(n)
    return n, err
}

type CountReader struct {
    R io.Reader
    N int64
}

func (cr *CountReader) Read(p []byte) (int, error) {
    n, err := cr.R.Read(p)
    cr.N += int64(n)
    return n, err
}

func OrderOrNative(o binary.ByteOrder) binary.ByteOrder {
    if o == nil { return binary.BigEndian }
    return o
}

func WriteFull(w io.Writer, p []byte) (int64, error) {
    var total int64
    for len(p) > 0 {
        n, err := w.Write(p)
        total += int64(n)
        if err != nil { return total, err }
        p = p[n:]
    }
    return total, nil
}

func IntFromValue(v reflect.Value) (int, error) {
    switch v.Kind() {
    case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
        return int(v.Int()), nil
    case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
        u := v.Uint()
        if u > uint64(^uint(0)>>1) {
            return 0, errors.New("length overflows int")
        }
        return int(u), nil
    default:
        return 0, fmt.Errorf("length field must be int/uint, got %s", v.Kind())
    }
}

func SetIntToValue(v reflect.Value, x int64) error {
    switch v.Kind() {
    case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
        v.SetInt(x); return nil
    case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
        if x < 0 { return errors.New("negative length for unsigned field") }
        v.SetUint(uint64(x)); return nil
    default:
        return fmt.Errorf("length field must be int/uint, got %s", v.Kind())
    }
}
