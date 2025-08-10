package codec

import (
    "fmt"
    "io"
    "reflect"
    "unicode/utf8"

    "github.com/platinvm/gbm/internal/rt"
    "github.com/platinvm/gbm/internal/spec"
)

// prepareLengths runs BEFORE any bytes are written. It sets or validates
// all length fields referenced by []byte and string fields.
func prepareLengths(v reflect.Value, ops []spec.Op) error {
    for _, op := range ops {
        switch {
        case op.IsSlice:
            bs := v.FieldByIndex(op.Index).Bytes()
            lnField := v.FieldByIndex(op.LenIndex)
            curLen, err := rt.IntFromValue(lnField)
            if err != nil { return err }
            want := len(bs)
            if curLen == 0 {
                if err := rt.SetIntToValue(lnField, int64(want)); err != nil { return err }
            } else if curLen != want {
                return fmt.Errorf("%s: length mismatch: got=%d, %s=%d", op.Name, want, op.LenRef, curLen)
            }

        case op.IsString:
            s := v.FieldByIndex(op.Index).String()
            if !op.EncUTF8 {
                return errOnlyUTF8
            }
            if !utf8.ValidString(s) {
                return errInvalidUTF8
            }
            want := len([]byte(s)) // bytes, not runes
            lnField := v.FieldByIndex(op.LenIndex)
            curLen, err := rt.IntFromValue(lnField)
            if err != nil { return err }
            if curLen == 0 {
                if err := rt.SetIntToValue(lnField, int64(want)); err != nil { return err }
            } else if curLen != want {
                return fmt.Errorf("%s: length mismatch: got=%d, %s=%d", op.Name, want, op.LenRef, curLen)
            }
        }
    }
    return nil
}

func DoMarshal(v reflect.Value, ops []spec.Op, w io.Writer) (int64, error) {
    var total int64

    // Pre-pass: set/validate lengths before writing any bytes.
    if err := prepareLengths(v, ops); err != nil {
        return 0, err
    }

    for _, op := range ops {
        fv := v.FieldByIndex(op.Index)

        switch {
        case op.IsSlice:
            bs := fv.Bytes()
            n, err := rt.WriteFull(w, bs)
            total += n
            if err != nil { return total, err }

        case op.IsString:
            s := fv.String()
            bs := []byte(s)
            n, err := rt.WriteFull(w, bs)
            total += n
            if err != nil { return total, err }

        default:
            n, err := writeNumber(w, op.Order, op.Kind, fv)
            total += n
            if err != nil { return total, err }
        }
    }
    return total, nil
}
