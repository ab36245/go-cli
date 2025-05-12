package cli

import (
	"fmt"
	"strings"
	"time"
)

func Date(ptr *time.Time) *_date {
	return (&_date{}).Bind(ptr)
}

func DateSlice(ptr *[]time.Time) *_dateSlice {
	return (&_dateSlice{}).Bind(ptr)
}

type _date struct {
	ptr *time.Time
}

func (v *_date) Assign(str string) error {
	val, err := dateParse(str)
	if err != nil {
		return err
	}
	*v.ptr = val
	return nil
}

func (v *_date) Bind(ptr *time.Time) *_date {
	if ptr == nil {
		ptr = new(time.Time)
	}
	v.ptr = ptr
	return v
}

func (v *_date) Consume(args *[]string) error {
	if len(*args) == 0 {
		return fmt.Errorf("date param requires a value")
	}
	if err := v.Assign((*args)[0]); err != nil {
		return err
	}
	*args = (*args)[1:]
	return nil
}

func (v *_date) NonZero() string {
	if v.ptr.IsZero() {
		return ""
	}
	return v.String()
}

func (v *_date) Reset() {
	*v.ptr = time.Time{}
}

func (v _date) String() string {
	return fmt.Sprintf("%v", *v.ptr)
}

func (v _date) Type() string {
	return "date"
}

func (v _date) Value() time.Time {
	return *v.ptr
}

type _dateSlice struct {
	ptr *[]time.Time
}

func (v *_dateSlice) Assign(str string) error {
	for bit := range strings.SplitSeq(str, ",") {
		val, err := dateParse(bit)
		if err != nil {
			return err
		}
		*v.ptr = append(*v.ptr, val)
	}
	return nil
}

func (v *_dateSlice) Bind(ptr *[]time.Time) *_dateSlice {
	if ptr == nil {
		ptr = new([]time.Time)
	}
	v.ptr = ptr
	return v
}

func (v *_dateSlice) Consume(args *[]string) error {
	for _, arg := range *args {
		if err := v.Assign(arg); err != nil {
			return err
		}
	}
	*args = nil
	return nil
}

func (v _dateSlice) NonZero() string {
	if len(*v.ptr) == 0 {
		return ""
	}
	return v.String()
}

func (v *_dateSlice) Reset() {
	*v.ptr = nil
}

func (v _dateSlice) String() string {
	return fmt.Sprintf("%v", *v.ptr)
}

func (v _dateSlice) Type() string {
	return "date..."
}

func (v _dateSlice) Value() []time.Time {
	return *v.ptr
}

func dateParse(str string) (time.Time, error) {
	str = strings.TrimSpace(str)
	val, err := time.Parse("2006-01-02", str)
	if err != nil {
		return time.Time{}, fmt.Errorf("bad date value %q", str)
	}
	return val, nil
}
