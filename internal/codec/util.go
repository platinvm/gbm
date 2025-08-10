package codec

import (
    "fmt"

    "github.com/platinvm/gbm/internal/spec"
)

func lengthMismatch(op spec.Op, got, expect int) error {
    return fmt.Errorf("%s: length mismatch: got=%d, %s=%d", op.Name, got, op.LenRef, expect)
}
