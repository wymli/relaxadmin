package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cast"
)

type Result struct {
	v interface{}
}

func NewResult(v interface{}) Result {
	return Result{v: v}
}

func (r Result) String() string {
	switch z := r.v.(type) {
	case nil:
		return ""
	case string:
		return z
	case error:
		return z.Error()
	case fmt.Stringer:
		return z.String()
	default:
		return fmt.Sprintf("%v", r.v)
	}
}

// &string or nil
func (r Result) StringPtr() *string {
	switch r.v.(type) {
	case nil:
		return nil
	default:
		x := r.String()
		return &x
	}
}

func (r Result) Time() time.Time {
	return cast.ToTime(r.v)
}

func (r Result) Float64() float64 {
	return cast.ToFloat64(r.v)
}

func (r Result) Float32() float32 {
	return cast.ToFloat32(r.v)
}

func (r Result) Bool() bool {
	return cast.ToBool(r.v)
}

func (r Result) Int() int {
	return cast.ToInt(r.v)
}

func (r Result) Int64() int64 {
	return cast.ToInt64(r.v)
}

func (r Result) Int32() int32 {
	return cast.ToInt32(r.v)
}

func (r Result) Interface() interface{} {
	return r.v
}

func (r Result) StringSlice(sep ...string) []string {
	sep_ := ","
	if len(sep) != 0 {
		sep_ = sep[0]
	}

	switch z := r.v.(type) {
	case []string:
		return z
	case string:
		return strings.Split(z, sep_)
	default:
		return strings.Split(r.String(), sep_)
	}

}

func (r Result) IntSlice() []int {
	return cast.ToIntSlice(r.v)
}

func (r Result) Error() error {
	switch z := r.v.(type) {
	case error:
		return z
	default:
		return nil
	}
}
