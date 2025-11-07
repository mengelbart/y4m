package y4m

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var errInvalidRatio = errors.New("invalid ratio")

type Ratio struct {
	Numerator   int
	Denominator int
}

func (r *Ratio) parse(s string) (err error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return errInvalidRatio
	}
	r.Numerator, err = strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("%w: invalid numerator: %w", errInvalidRatio, err)
	}
	r.Denominator, err = strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("%w: invalid denominator: %w", errInvalidRatio, err)
	}
	return nil
}
