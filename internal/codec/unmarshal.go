package codec

import (
    "fmt"
    "io"
    "reflect"

    "github.com/platinvm/gbm/internal/rt"
    "github.com/platinvm/gbm/internal/spec"
)

func DoUnmarshal(v reflect.Value, ops []spec.Op, r io.Reader) (int64, error) {
    var total int64
    for _, op := range ops {
        fv := v.FieldByIndex(op.Index)

        switch {
        case op.IsSlice:
            lnField := v.FieldByIndex(op.LenIndex)
            ln, err := rt.IntFromValue(lnField)
            if err != nil { return total, err }
            if ln < 0 { return total, fmt.Errorf("%s: negative length %d", op.Name, ln) }
            buf := make([]byte, ln)
            n, err := io.ReadFull(r, buf)
            total += int64(n)
            if err != nil { return total, err }
            fv.SetBytes(buf)

        case op.IsString:
            lnField := v.FieldByIndex(op.LenIndex)
            ln, err := rt.IntFromValue(lnField)
            if err != nil { return total, err }
            if ln < 0 { return total, fmt.Errorf("%s: negative length %d", op.Name, ln) }
            buf := make([]byte, ln)
            n, err := io.ReadFull(r, buf)
            total += int64(n)
            if err != nil { return total, err }
            fv.SetString(string(buf))

        default:
            n, err := readNumber(r, op.Order, op.Kind, fv)
            total += n
            if err != nil { return total, err }
        }
    }
    return total, nil
}
