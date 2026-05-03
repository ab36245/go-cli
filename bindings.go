package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Binding interface {
	Assign(string) error
	Consume(*[]string, int, int) error
	IsFlag() bool
	IsZero() bool
	Reset()
	String() string
	Type() string
	Update()
}

func Bool(ptr *bool) Binding {
	return &_simple[bool]{
		value:   ptr,
		handler: boolHandler(false),
	}
}

func BoolFlag(ptr *bool) Binding {
	return &_simple[bool]{
		value:   ptr,
		handler: boolHandler(true),
	}
}

func BoolSlice(ptr *[]bool) Binding {
	return &_slice[bool]{
		value:   ptr,
		handler: boolHandler(false),
	}
}

func Enum[T comparable](ptr *T, vals map[string]T) Binding {
	return &_simple[T]{
		value:   ptr,
		handler: enumHandler(vals),
	}
}

func EnumSlice[T comparable](ptr *[]T, vals map[string]T) Binding {
	return &_slice[T]{
		value:   ptr,
		handler: enumHandler(vals),
	}
}

func Date(ptr *time.Time) Binding {
	return &_simple[time.Time]{
		value:   ptr,
		handler: dateHandler(),
	}
}

func DateSlice(ptr *[]time.Time) Binding {
	return &_slice[time.Time]{
		value:   ptr,
		handler: dateHandler(),
	}
}

func Int(ptr *int) Binding {
	return &_simple[int]{
		value:   ptr,
		handler: intHandler(false),
	}
}

func IntFlag(ptr *int) Binding {
	return &_simple[int]{
		value:   ptr,
		handler: intHandler(true),
	}
}

func IntSlice(ptr *[]int) Binding {
	return &_slice[int]{
		value:   ptr,
		handler: intHandler(false),
	}
}

func String(ptr *string) Binding {
	return &_simple[string]{
		value:   ptr,
		handler: stringHandler(),
	}
}

func StringSlice(ptr *[]string) Binding {
	return &_slice[string]{
		value:   ptr,
		handler: stringHandler(),
	}
}

type _simple[T comparable] struct {
	value   *T
	handler _handler[T]
}

func (b *_simple[T]) Assign(str string) error {
	value, err := b.handler.parseFunc(str)
	if err != nil {
		return err
	}
	*b.value = value
	return nil
}

func (b *_simple[T]) Consume(strs *[]string, min, max int) error {
	if len(*strs) == 0 {
		if min == 0 {
			return nil
		}
		return fmt.Errorf("value required but not given")
	}
	// Do at most 1 value
	if err := b.Assign((*strs)[0]); err != nil {
		return err
	}
	*strs = (*strs)[1:]
	return nil
}

func (b _simple[T]) IsFlag() bool {
	return b.handler.updateFunc != nil
}

func (b _simple[T]) IsZero() bool {
	return *b.value == *new(T)
}

func (b *_simple[T]) Reset() {
	*b.value = *new(T)
}

func (b _simple[T]) String() string {
	return b.handler.stringFunc(*b.value)
}

func (b _simple[T]) Type() string {
	return b.handler.typeString
}

func (b *_simple[T]) Update() {
	if b.handler.updateFunc != nil {
		*b.value = b.handler.updateFunc(*b.value)
	}
}

type _slice[T comparable] struct {
	value   *[]T
	handler _handler[T]
}

func (b *_slice[T]) Assign(str string) error {
	for _, part := range b.handler.splitFunc(str) {
		value, err := b.handler.parseFunc(part)
		if err != nil {
			return err
		}
		*b.value = append(*b.value, value)
	}
	return nil
}

func (b *_slice[T]) Consume(strs *[]string, min, max int) error {
	if len(*strs) < min {
		if min == 1 {
			return fmt.Errorf("value required but not given")
		}
		return fmt.Errorf("%d values required but only %d given", min, len(*strs))
	}
	if max == 0 || max > len(*strs) {
		max = len(*strs)
	}
	for i := range max {
		if err := b.Assign((*strs)[i]); err != nil {
			return err
		}
	}
	*strs = (*strs)[max:]
	return nil
}

func (b _slice[T]) IsFlag() bool {
	return false
}

func (b _slice[T]) IsZero() bool {
	return len(*b.value) == 0
}

func (b *_slice[T]) Reset() {
	*b.value = nil
}

func (b _slice[T]) String() string {
	var s strings.Builder
	for i, v := range *b.value {
		if i > 0 {
			s.WriteString(", ")
		}
		s.WriteString(b.handler.stringFunc(v))
	}
	return fmt.Sprintf("[%s]", s.String())
}

func (b _slice[T]) Type() string {
	return b.handler.typeString + "..."
}

func (b *_slice[T]) Update() {
	// Do nothing
}

type _handler[T comparable] struct {
	parseFunc  func(string) (T, error)
	splitFunc  func(string) []string
	stringFunc func(T) string
	typeString string
	updateFunc func(T) T
}

func boolHandler(isFlag bool) _handler[bool] {
	handler := _handler[bool]{
		parseFunc: func(s string) (bool, error) {
			s = strings.TrimSpace(s)
			value, err := strconv.ParseBool(s)
			if err != nil {
				err = fmt.Errorf("bad int value %q", s)
			}
			return value, err
		},
		splitFunc: func(s string) []string {
			return strings.Split(s, ",")
		},
		stringFunc: func(v bool) string {
			return fmt.Sprintf("%v", v)
		},
		typeString: "bool",
	}
	if isFlag {
		handler.updateFunc = func(_ bool) bool {
			return true
		}
	}
	return handler
}

func dateHandler() _handler[time.Time] {
	return _handler[time.Time]{
		parseFunc: func(s string) (time.Time, error) {
			s = strings.TrimSpace(s)
			value, err := time.Parse("2006-01-02", s)
			if err != nil {
				err = fmt.Errorf("bad date value %q", s)
			}
			return value, err
		},
		splitFunc: func(s string) []string {
			return strings.Split(s, ",")
		},
		stringFunc: func(v time.Time) string {
			return v.Format("2006-01-02")
		},
		typeString: "date",
	}
}

func enumHandler[T comparable](vals map[string]T) _handler[T] {
	handler := _handler[T]{
		parseFunc: func(s string) (T, error) {
			if v, ok := vals[s]; ok {
				return v, nil
			}
			return *new(T), fmt.Errorf("unknown enum value %q", s)
		},
		splitFunc: func(s string) []string {
			return strings.Split(s, ",")
		},
		stringFunc: func(val T) string {
			for k, v := range vals {
				if v == val {
					return k
				}
			}
			return "(none)"
		},
	}
	var s strings.Builder
	for k := range vals {
		if s.String() != "" {
			s.WriteString("|")
		}
		s.WriteString(k)
	}
	handler.typeString = s.String()
	return handler
}

func intHandler(isFlag bool) _handler[int] {
	handler := _handler[int]{
		parseFunc: func(s string) (int, error) {
			s = strings.TrimSpace(s)
			value, err := strconv.ParseInt(s, 10, 0)
			if err != nil {
				err = fmt.Errorf("bad int value %q", s)
			}
			return int(value), err
		},
		splitFunc: func(s string) []string {
			return strings.Split(s, ",")
		},
		stringFunc: func(v int) string {
			return fmt.Sprintf("%v", v)
		},
		typeString: "int",
	}
	if isFlag {
		handler.updateFunc = func(v int) int {
			return v + 1
		}
	}
	return handler
}

func stringHandler() _handler[string] {
	return _handler[string]{
		parseFunc: func(s string) (string, error) {
			return s, nil
		},
		splitFunc: func(s string) []string {
			return []string{s}
		},
		stringFunc: func(v string) string {
			return fmt.Sprintf("%q", v)
		},
		typeString: "string",
	}
}
