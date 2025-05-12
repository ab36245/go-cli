package cli

import (
	"fmt"
	"strconv"
	"strings"
)

func Bool(ptr *bool) *_bool {
	return (&_bool{}).Bind(ptr)
}

func BoolFlag(ptr *bool) *_boolFlag {
	return (&_boolFlag{}).Bind(ptr)
}

func BoolSlice(ptr *[]bool) *_boolSlice {
	return (&_boolSlice{}).Bind(ptr)
}

type _boolBase struct {
	ptr *bool
}

func (v *_boolBase) Assign(str string) error {
	val, err := boolParse(str)
	if err != nil {
		return err
	}
	*v.ptr = val
	return nil
}

func (v *_boolBase) Bind(ptr *bool) *_boolBase {
	if ptr == nil {
		ptr = new(bool)
	}
	v.ptr = ptr
	return v
}

func (v *_boolBase) NonZero() string {
	if !*v.ptr {
		return ""
	}
	return v.String()
}

func (v *_boolBase) Reset() {
	*v.ptr = false
}

func (v _boolBase) String() string {
	return fmt.Sprintf("%v", *v.ptr)
}

func (v _boolBase) Type() string {
	return "bool"
}

func (v _boolBase) Value() bool {
	return *v.ptr
}

type _bool struct {
	_boolBase
}

func (v *_bool) Bind(ptr *bool) *_bool {
	v._boolBase.Bind(ptr)
	return v
}

func (v *_bool) Consume(args *[]string) error {
	if len(*args) == 0 {
		return fmt.Errorf("bool param requires a value")
	}
	if err := v.Assign((*args)[0]); err != nil {
		return err
	}
	*args = (*args)[1:]
	return nil
}

type _boolFlag struct {
	_boolBase
}

func (v *_boolFlag) Bind(ptr *bool) *_boolFlag {
	v._boolBase.Bind(ptr)
	return v
}

func (v *_boolFlag) Update() {
	*v.ptr = true
}

type _boolSlice struct {
	ptr *[]bool
}

func (v *_boolSlice) Assign(str string) error {
	for bit := range strings.SplitSeq(str, ",") {
		val, err := boolParse(bit)
		if err != nil {
			return err
		}
		*v.ptr = append(*v.ptr, val)
	}
	return nil
}

func (v *_boolSlice) Bind(ptr *[]bool) *_boolSlice {
	if ptr == nil {
		ptr = new([]bool)
	}
	v.ptr = ptr
	return v
}

func (v *_boolSlice) Consume(args *[]string) error {
	for _, arg := range *args {
		if err := v.Assign(arg); err != nil {
			return err
		}
	}
	*args = nil
	return nil
}

func (v _boolSlice) NonZero() string {
	if len(*v.ptr) == 0 {
		return ""
	}
	return v.String()
}

func (v *_boolSlice) Reset() {
	*v.ptr = nil
}

func (v _boolSlice) String() string {
	return fmt.Sprintf("%v", *v.ptr)
}

func (v _boolSlice) Type() string {
	return "bool..."
}

func (v _boolSlice) Value() []bool {
	return *v.ptr
}

func boolParse(str string) (bool, error) {
	str = strings.TrimSpace(str)
	val, err := strconv.ParseBool(str)
	if err != nil {
		return false, fmt.Errorf("bad bool value %q", str)
	}
	return val, nil
}
