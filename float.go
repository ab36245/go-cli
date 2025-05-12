package cli

import (
	"fmt"
	"strconv"
	"strings"
)

func Float(ptr *float64) *_float {
	return (&_float{}).Bind(ptr)
}

func FloatSlice(ptr *[]float64) *_floatSlice {
	return (&_floatSlice{}).Bind(ptr)
}

type _float struct {
	ptr *float64
}

func (v *_float) Assign(str string) error {
	val, err := floatParse(str)
	if err != nil {
		return err
	}
	*v.ptr = val
	return nil
}

func (v *_float) Bind(ptr *float64) *_float {
	if ptr == nil {
		ptr = new(float64)
	}
	v.ptr = ptr
	return v
}

func (v *_float) Param() {
}

func (v *_float) Reset() {
	*v.ptr = 0
}

func (v _float) String() string {
	return fmt.Sprintf("%v", *v.ptr)
}

func (v _float) Type() string {
	return "float"
}

func (v _float) Value() float64 {
	return *v.ptr
}

type _floatSlice struct {
	ptr *[]float64
}

func (v *_floatSlice) Assign(str string) error {
	for bit := range strings.SplitSeq(str, ",") {
		val, err := floatParse(bit)
		if err != nil {
			return err
		}
		*v.ptr = append(*v.ptr, val)
	}
	return nil
}

func (v *_floatSlice) Bind(ptr *[]float64) *_floatSlice {
	if ptr == nil {
		ptr = new([]float64)
	}
	v.ptr = ptr
	return v
}

func (v *_floatSlice) Param() {
}

func (v *_floatSlice) Reset() {
	*v.ptr = nil
}

func (v _floatSlice) String() string {
	return fmt.Sprintf("%v", *v.ptr)
}

func (v _floatSlice) Type() string {
	return "float..."
}

func (v _floatSlice) Value() []float64 {
	return *v.ptr
}

func floatParse(str string) (float64, error) {
	str = strings.TrimSpace(str)
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, fmt.Errorf("bad float value %q", str)
	}
	return val, nil
}
