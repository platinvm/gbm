package codec

import (
    "encoding/binary"
    "fmt"
    "io"
    "reflect"

    "github.com/platinvm/gbm/internal/rt"
)

func writeNumber(w io.Writer, order binary.ByteOrder, k reflect.Kind, fv reflect.Value) (int64, error) {
    bw := &rt.CountWriter{W: w}
    switch k {
    case reflect.Int8:
        var x = int8(fv.Int())
        return bw.N, binary.Write(bw, rt.OrderOrNative(order), x)
    case reflect.Int16:
        return bw.N, binary.Write(bw, order, int16(fv.Int()))
    case reflect.Int32:
        return bw.N, binary.Write(bw, order, int32(fv.Int()))
    case reflect.Int64:
        return bw.N, binary.Write(bw, order, fv.Int())
    case reflect.Uint8:
        var x = uint8(fv.Uint())
        return bw.N, binary.Write(bw, rt.OrderOrNative(order), x)
    case reflect.Uint16:
        return bw.N, binary.Write(bw, order, uint16(fv.Uint()))
    case reflect.Uint32:
        return bw.N, binary.Write(bw, order, uint32(fv.Uint()))
    case reflect.Uint64:
        return bw.N, binary.Write(bw, order, fv.Uint())
    default:
        return bw.N, fmt.Errorf("unsupported kind %s", k)
    }
}

func readNumber(r io.Reader, order binary.ByteOrder, k reflect.Kind, fv reflect.Value) (int64, error) {
    br := &rt.CountReader{R: r}
    switch k {
    case reflect.Int8:
        var x int8
        if err := binary.Read(br, rt.OrderOrNative(order), &x); err != nil { return br.N, err }
        fv.SetInt(int64(x))
    case reflect.Int16:
        var x int16
        if err := binary.Read(br, order, &x); err != nil { return br.N, err }
        fv.SetInt(int64(x))
    case reflect.Int32:
        var x int32
        if err := binary.Read(br, order, &x); err != nil { return br.N, err }
        fv.SetInt(int64(x))
    case reflect.Int64:
        var x int64
        if err := binary.Read(br, order, &x); err != nil { return br.N, err }
        fv.SetInt(x)
    case reflect.Uint8:
        var x uint8
        if err := binary.Read(br, rt.OrderOrNative(order), &x); err != nil { return br.N, err }
        fv.SetUint(uint64(x))
    case reflect.Uint16:
        var x uint16
        if err := binary.Read(br, order, &x); err != nil { return br.N, err }
        fv.SetUint(uint64(x))
    case reflect.Uint32:
        var x uint32
        if err := binary.Read(br, order, &x); err != nil { return br.N, err }
        fv.SetUint(uint64(x))
    case reflect.Uint64:
        var x uint64
        if err := binary.Read(br, order, &x); err != nil { return br.N, err }
        fv.SetUint(x)
    default:
        return br.N, fmt.Errorf("unsupported kind %s", k)
    }
    return br.N, nil
}
