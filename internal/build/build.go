package build

import (
    "fmt"
    "reflect"

    "github.com/platinvm/gbm/internal/rt"
    "github.com/platinvm/gbm/internal/spec"
    "github.com/platinvm/gbm/internal/tags"
)

// BuildOps walks the struct type and generates the operation plan.
func BuildOps(t reflect.Type, nameToIndex map[string][]int, prefix []int) ([]spec.Op, error) {
    var ops []spec.Op
    for i := 0; i < t.NumField(); i++ {
        f := t.Field(i)
        if f.PkgPath != "" { // unexported
            continue
        }
        idx := append(append([]int{}, prefix...), i)
        if f.Anonymous && f.Type.Kind() == reflect.Struct {
            sub, err := BuildOps(f.Type, nameToIndex, idx)
            if err != nil { return nil, err }
            ops = append(ops, sub...)
            continue
        }

        tag := f.Tag.Get("bin")
        order, _ := tags.ParseOrder(tag)
        lenRef, hasLen := tags.ParseLen(tag)
        isUTF8 := tags.ParseEnc(tag)

        k := f.Type.Kind()
        switch k {
        case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
            reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
            if k != reflect.Int8 && k != reflect.Uint8 && order == nil {
                return nil, fmt.Errorf("field %s: must specify LE/BE for %s", f.Name, k)
            }
            ops = append(ops, spec.Op{
                Name:  f.Name,
                Index: idx,
                Kind:  k,
                Order: order,
            })

        case reflect.Slice:
            if f.Type.Elem().Kind() != reflect.Uint8 {
                return nil, fmt.Errorf("field %s: only []byte is supported", f.Name)
            }
            if !hasLen {
                return nil, fmt.Errorf("field %s: []byte requires len=<FieldName>", f.Name)
            }
            lenIdx, ok := nameToIndex[lenRef]
            if !ok {
                return nil, fmt.Errorf("field %s: len ref %q not found", f.Name, lenRef)
            }
            if !rt.IsBefore(lenIdx, idx) {
                return nil, fmt.Errorf("field %s: len ref %q must appear before the field", f.Name, lenRef)
            }
            ops = append(ops, spec.Op{
                Name:     f.Name,
                Index:    idx,
                Kind:     k,
                IsSlice:  true,
                LenRef:   lenRef,
                LenIndex: lenIdx,
            })

        case reflect.String:
            if !hasLen {
                return nil, fmt.Errorf("field %s: string requires len=<FieldName>", f.Name)
            }
            lenIdx, ok := nameToIndex[lenRef]
            if !ok {
                return nil, fmt.Errorf("field %s: len ref %q not found", f.Name, lenRef)
            }
            if !rt.IsBefore(lenIdx, idx) {
                return nil, fmt.Errorf("field %s: len ref %q must appear before the field", f.Name, lenRef)
            }
            ops = append(ops, spec.Op{
                Name:     f.Name,
                Index:    idx,
                Kind:     k,
                IsString: true,
                LenRef:   lenRef,
                LenIndex: lenIdx,
                EncUTF8:  isUTF8,
            })

        default:
            return nil, fmt.Errorf("field %s: unsupported kind %s", f.Name, k)
        }
    }
    return ops, nil
}
