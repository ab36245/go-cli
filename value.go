package cli

import (
	"fmt"
	"strconv"
	"time"
)

type value struct {
	kind  string
	array bool
	dflt  string

	doInc func()
	doSet func(string) error
}

func newValue(binding any, flag bool) (*value, error) {
	if binding == nil {
		return nil, fmt.Errorf("no binding")
	}

	switch b := binding.(type) {
	case *bool:
		v := &value{
			kind: "bool",
			dflt: fmt.Sprintf("%t", *b),
			doInc: func() {
				*b = true
			},
			doSet: func(s string) error {
				return setValue(b, toBool, s)
			},
		}
		return v, nil

	case *[]bool:
		v := &value{
			kind:  "bool...",
			array: true,
			doSet: func(s string) error {
				return setSlice(b, toBool, s)
			},
		}
		return v, nil

	case *float64:
		v := &value{
			kind: "float",
			dflt: fmt.Sprintf("%.1f", *b),
			doSet: func(s string) error {
				return setValue(b, toFloat, s)
			},
		}
		return v, nil

	case *[]float64:
		v := &value{
			kind:  "float...",
			array: true,
			doSet: func(s string) error {
				return setSlice(b, toFloat, s)
			},
		}
		return v, nil

	case *int:
		v := &value{
			kind: "int",
			dflt: fmt.Sprintf("%d", *b),
			doSet: func(s string) error {
				return setValue(b, toInt, s)
			},
		}
		if flag {
			v.doInc = func() {
				*b++
			}
		}
		return v, nil

	case *[]int:
		v := &value{
			kind:  "int...",
			array: true,
			doSet: func(s string) error {
				return setSlice(b, toInt, s)
			},
		}
		return v, nil

	case *string:
		v := &value{
			kind: "string",
			dflt: *b,
			doSet: func(s string) error {
				return setValue(b, toString, s)
			},
		}
		if *b != "" {
			v.dflt = fmt.Sprintf("\"%s\"", *b)
		}
		return v, nil

	case *[]string:
		v := &value{
			kind:  "string...",
			array: true,
			doSet: func(s string) error {
				return setSlice(b, toString, s)
			},
		}
		return v, nil

	case *Date:
		v := &value{
			kind: "date",
			doSet: func(s string) error {
				return setValue(b, toDate, s)
			},
		}
		return v, nil

	case *[]Date:
		v := &value{
			kind:  "date...",
			array: true,
			doSet: func(s string) error {
				return setSlice(b, toDate, s)
			},
		}
		return v, nil

	default:
		return nil, fmt.Errorf("invalid binding %T", b)
	}
}

func (v *value) inc() bool {
	if v.doInc == nil {
		return false
	}
	v.doInc()
	return true
}

func (v *value) set(s string) error {
	return v.doSet(s)
}

type converter[T any] func(string) (T, error)

func setSlice[T any](binding *[]T, converter converter[T], s string) error {
	val, err := converter(s)
	if err != nil {
		return err
	}
	if *binding == nil {
		*binding = []T{val}
	} else {
		*binding = append(*binding, val)
	}
	return nil
}

func setValue[T any](binding *T, converter converter[T], s string) error {
	val, err := converter(s)
	if err != nil {
		return err
	}
	*binding = val
	return nil
}

func toBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

func toFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func toInt(s string) (int, error) {
	val, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		return 0, err
	}
	return int(val), nil
}

func toString(s string) (string, error) {
	return s, nil
}

func toDate(s string) (Date, error) {
	val, err := time.Parse("2006-01-02", s)
	if err != nil {
		return Date{}, err
	}
	return Date(val), nil
}
