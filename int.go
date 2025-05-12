package cli

import (
	"fmt"
	"strconv"
	"strings"
)

func Int(ptr *int) *_int {
	return (&_int{}).Bind(ptr)
}

func IntFlag(ptr *int) *_intFlag {
	return (&_intFlag{}).Bind(ptr)
}

func IntSlice(ptr *[]int) *_intSlice {
	return (&_intSlice{}).Bind(ptr)
}

type _intBase struct {
	ptr *int
}

func (v *_intBase) Assign(str string) error {
	val, err := intParse(str)
	if err != nil {
		return err
	}
	*v.ptr = val
	return nil
}

func (v *_intBase) Bind(ptr *int) *_intBase {
	if ptr == nil {
		ptr = new(int)
	}
	v.ptr = ptr
	return v
}

func (v *_intBase) Reset() {
	*v.ptr = 0
}

func (v _intBase) String() string {
	return fmt.Sprintf("%v", *v.ptr)
}

func (v _intBase) Type() string {
	return "int"
}

func (v _intBase) Value() int {
	return *v.ptr
}

type _int struct {
	_intBase
}

func (v *_int) Bind(ptr *int) *_int {
	v._intBase.Bind(ptr)
	return v
}

func (v *_int) Param() {
}

type _intFlag struct {
	_intBase
}

func (v *_intFlag) Bind(ptr *int) *_intFlag {
	v._intBase.Bind(ptr)
	return v
}

func (v *_intFlag) Update() {
	*v.ptr++
}

type _intSlice struct {
	ptr *[]int
}

func (v *_intSlice) Assign(str string) error {
	for bit := range strings.SplitSeq(str, ",") {
		val, err := intParse(bit)
		if err != nil {
			return err
		}
		*v.ptr = append(*v.ptr, val)
	}
	return nil
}

func (v *_intSlice) Bind(ptr *[]int) *_intSlice {
	if ptr == nil {
		ptr = new([]int)
	}
	v.ptr = ptr
	return v
}

func (v *_intSlice) Reset() {
	*v.ptr = nil
}

func (v *_intSlice) Param() {
}

func (v _intSlice) String() string {
	return fmt.Sprintf("%v", *v.ptr)
}

func (v _intSlice) Type() string {
	return "int..."
}

func (v _intSlice) Value() []int {
	return *v.ptr
}

func intParse(str string) (int, error) {
	str = strings.TrimSpace(str)
	val, err := strconv.ParseInt(str, 0, 0)
	if err != nil {
		return 0, fmt.Errorf("bad int value %q", str)
	}
	return int(val), nil
}
