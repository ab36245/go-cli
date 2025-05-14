package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ValueHandler[T any] struct {
	IsZero func(T) bool
	Parse  func(string) (T, error)
	Split  bool
	String func(T) string
	Type   string
	Update func(T) T
}

var boolHandler = ValueHandler[bool]{
	IsZero: func(value bool) bool {
		return !value
	},

	Parse: func(str string) (bool, error) {
		str = strings.TrimSpace(str)
		value, err := strconv.ParseBool(str)
		if err != nil {
			err = fmt.Errorf("bad bool value %q", str)
		}
		return value, err
	},

	Split: true,

	String: func(value bool) string {
		return fmt.Sprintf("%v", value)
	},

	Type: "bool",

	Update: func(value bool) bool {
		return true
	},
}

var dateHandler = ValueHandler[time.Time]{
	IsZero: func(value time.Time) bool {
		return value.IsZero()
	},

	Parse: func(str string) (time.Time, error) {
		str = strings.TrimSpace(str)
		value, err := time.Parse("2006-01-02", str)
		if err != nil {
			err = fmt.Errorf("bad date value %q", str)
		}
		return value, err
	},

	Split: true,

	String: func(value time.Time) string {
		return fmt.Sprintf("%v", value)
	},

	Type: "date",

	Update: nil,
}

var floatHandler = ValueHandler[float64]{
	IsZero: func(value float64) bool {
		return value == 0
	},

	Parse: func(str string) (float64, error) {
		str = strings.TrimSpace(str)
		value, err := strconv.ParseFloat(str, 64)
		if err != nil {
			err = fmt.Errorf("bad float value %q", str)
		}
		return value, err
	},

	Split: true,

	String: func(value float64) string {
		return fmt.Sprintf("%v", value)
	},

	Type: "float",

	Update: nil,
}

var intHandler = ValueHandler[int]{
	IsZero: func(value int) bool {
		return value == 0
	},

	Parse: func(str string) (int, error) {
		str = strings.TrimSpace(str)
		value, err := strconv.ParseInt(str, 0, 0)
		if err != nil {
			err = fmt.Errorf("bad int value %q", str)
		}
		return int(value), err
	},

	Split: true,

	String: func(value int) string {
		return fmt.Sprintf("%v", value)
	},

	Type: "int",

	Update: func(value int) int {
		return value + 1
	},
}

var stringHandler = ValueHandler[string]{
	IsZero: func(value string) bool {
		return value == ""
	},

	Parse: func(str string) (string, error) {
		return str, nil
	},

	Split: false,

	String: func(value string) string {
		return fmt.Sprintf("%q", value)
	},

	Type: "string",

	Update: nil,
}
