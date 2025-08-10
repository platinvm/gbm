package gbm

import (
	"fmt"
	"io"
	"reflect"

	"github.com/platinvm/gbm/internal/build"
	"github.com/platinvm/gbm/internal/codec"
)

// Codec is a generic binary (un)marshaler for T configured via struct tags.
type Codec[T any] struct {
	marshal   func(*T, io.Writer) (int64, error)
	unmarshal func(*T, io.Reader) (int64, error)
}

// New analyzes T and returns a ready-to-use Codec.
func New[T any]() (*Codec[T], error) {
	var zero *T
	rt := reflect.TypeOf(zero).Elem()
	if rt.Kind() != reflect.Struct {
		return nil, fmt.Errorf("gbm: T must be a struct, got %s", rt.Kind())
	}

	// build field index map for len= references
	nameToIndex := map[string][]int{}
	var collect func(t reflect.Type, prefix []int)
	collect = func(t reflect.Type, prefix []int) {
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if f.PkgPath != "" { // unexported
				continue
			}
			idx := append(append([]int{}, prefix...), i)
			if f.Anonymous && f.Type.Kind() == reflect.Struct {
				collect(f.Type, idx)
				continue
			}
			nameToIndex[f.Name] = idx
		}
	}

	collect(rt, nil)

	ops, err := build.BuildOps(rt, nameToIndex, nil)
	if err != nil {
		return nil, err
	}

	m := func(sp *T, w io.Writer) (int64, error) {
		return codec.DoMarshal(reflect.ValueOf(sp).Elem(), ops, w)
	}
	u := func(sp *T, r io.Reader) (int64, error) {
		return codec.DoUnmarshal(reflect.ValueOf(sp).Elem(), ops, r)
	}
	return &Codec[T]{marshal: m, unmarshal: u}, nil
}

// Marshal writes the binary representation of v to w.
func (c *Codec[T]) Marshal(v *T, w io.Writer) (int64, error) { return c.marshal(v, w) }

// Unmarshal reads from r into v.
func (c *Codec[T]) Unmarshal(v *T, r io.Reader) (int64, error) { return c.unmarshal(v, r) }

// To use with gbm.New(), if error != nil panics
func Must[T any](c *Codec[T], err error) *Codec[T] {
	if err != nil {
		panic(err)
	}

	return c
}
