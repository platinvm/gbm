package codec

import "errors"

var (
    errOnlyUTF8    = errors.New("only UTF-8 strings are supported for now")
    errInvalidUTF8 = errors.New("string is not valid UTF-8")
)
