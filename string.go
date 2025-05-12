package cli

import (
	"fmt"
)

func String(ptr *string) *_string {
	return (&_string{}).Bind(ptr)
}

func StringSlice(ptr *[]string) *_stringSlice {
	return (&_stringSlice{}).Bind(ptr)
}

type _string struct {
	ptr *string
}

func (v *_string) Assign(str string) error {
	*v.ptr = str
	return nil
}

func (v *_string) Bind(ptr *string) *_string {
	if ptr == nil {
		ptr = new(string)
	}
	v.ptr = ptr
	return v
}

func (v *_string) Consume(args *[]string) error {
	if len(*args) == 0 {
		return fmt.Errorf("string param requires a value")
	}
	if err := v.Assign((*args)[0]); err != nil {
		return err
	}
	*args = (*args)[1:]
	return nil
}

func (v *_string) NonZero() string {
	if *v.ptr == "" {
		return ""
	}
	return v.String()
}

func (v *_string) Reset() {
	*v.ptr = ""
}

func (v _string) String() string {
	return fmt.Sprintf("%q", *v.ptr)
}

func (v _string) Type() string {
	return "string"
}

func (v _string) Value() string {
	return *v.ptr
}

type _stringSlice struct {
	ptr *[]string
}

func (v *_stringSlice) Assign(str string) error {
	*v.ptr = append(*v.ptr, str)
	return nil
}

func (v *_stringSlice) Bind(ptr *[]string) *_stringSlice {
	if ptr == nil {
		ptr = new([]string)
	}
	v.ptr = ptr
	return v
}

func (v *_stringSlice) Consume(args *[]string) error {
	for _, arg := range *args {
		if err := v.Assign(arg); err != nil {
			return err
		}
	}
	*args = nil
	return nil
}

func (v _stringSlice) NonZero() string {
	if len(*v.ptr) == 0 {
		return ""
	}
	return v.String()
}

func (v *_stringSlice) Reset() {
	*v.ptr = nil
}

func (v _stringSlice) String() string {
	return fmt.Sprintf("%v", *v.ptr)
}

func (v _stringSlice) Type() string {
	return "string..."
}

func (v _stringSlice) Value() []string {
	return *v.ptr
}
