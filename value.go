package cli

import (
	"fmt"
	"strings"
	"time"
)

type ValueBase[T any] struct {
	handler ValueHandler[T]
	value   *T
}

func NewValueBase[T any](handler ValueHandler[T]) *ValueBase[T] {
	return &ValueBase[T]{
		handler: handler,
		value:   new(T),
	}
}

func (v ValueBase[T]) Assign(str string) error {
	value, err := v.handler.Parse(str)
	if err != nil {
		return err
	}
	*v.value = value
	return nil
}

func (v *ValueBase[T]) Bind(value *T) *ValueBase[T] {
	if value != nil {
		v.value = value
	}
	return v
}

func (v ValueBase[T]) NonZero() string {
	if v.handler.IsZero(*v.value) {
		return ""
	}
	return v.String()
}

func (v ValueBase[T]) Reset() {
	*v.value = *new(T)
}

func (v ValueBase[T]) String() string {
	return v.handler.String(*v.value)
}

func (v ValueBase[T]) Type() string {
	return v.handler.Type
}

func (v ValueBase[T]) Value() T {
	return *v.value
}

type Value[T any] struct {
	ValueBase[T]
}

func NewValue[T any](handler ValueHandler[T]) *Value[T] {
	return &Value[T]{
		ValueBase[T]{
			handler: handler,
			value:   new(T),
		},
	}
}

func (v *Value[T]) Bind(value *T) *Value[T] {
	v.ValueBase.Bind(value)
	return v
}

func (v Value[T]) Consume(args *[]string) error {
	if len(*args) == 0 {
		return fmt.Errorf("%s argument required", v.handler.Type)
	}
	if err := v.Assign((*args)[0]); err != nil {
		return err
	}
	*args = (*args)[1:]
	return nil
}

type ValueFlag[T any] struct {
	ValueBase[T]
}

func NewValueFlag[T any](handler ValueHandler[T]) *ValueFlag[T] {
	return &ValueFlag[T]{
		ValueBase[T]{
			handler: handler,
			value:   new(T),
		},
	}
}

func (v *ValueFlag[T]) Bind(value *T) *ValueFlag[T] {
	v.ValueBase.Bind(value)
	return v
}

func (v ValueFlag[T]) Update() {
	*v.value = v.handler.Update(*v.value)
}

func NewValueSlice[T any](handler ValueHandler[T]) *ValueSlice[T] {
	return &ValueSlice[T]{
		handler: handler,
		value:   new([]T),
	}
}

func (v ValueSlice[T]) Assign(str string) error {
	var strs []string
	if v.handler.Split {
		strs = strings.Split(str, ",")
	} else {
		strs = []string{str}
	}
	for _, substr := range strs {
		value, err := v.handler.Parse(substr)
		if err != nil {
			return err
		}
		*v.value = append(*v.value, value)
	}
	return nil
}

func (v *ValueSlice[T]) Bind(value *[]T) *ValueSlice[T] {
	if value != nil {
		v.value = value
	}
	return v
}

func (v ValueSlice[T]) Consume(args *[]string) error {
	for _, arg := range *args {
		if err := v.Assign(arg); err != nil {
			return err
		}
	}
	*args = nil
	return nil
}

func (v ValueSlice[T]) NonZero() string {
	if len(*v.value) == 0 {
		return ""
	}
	return v.String()
}

func (v ValueSlice[T]) Reset() {
	*v.value = nil
}

func (v ValueSlice[T]) String() string {
	s := "["
	for i, value := range *v.value {
		if i > 0 {
			s += ", "
		}
		s += v.handler.String(value)
	}
	s += "]"
	return s
}

func (v ValueSlice[T]) Type() string {
	return v.handler.Type + ", ..."
}

func (v ValueSlice[T]) Value() []T {
	return *v.value
}

type ValueSlice[T any] struct {
	handler ValueHandler[T]
	value   *[]T
}

func Bool() *Value[bool] {
	return NewValue(boolHandler)
}

func BoolFlag() *ValueFlag[bool] {
	return NewValueFlag(boolHandler)
}

func BoolSlice() *ValueSlice[bool] {
	return NewValueSlice(boolHandler)
}

func Date() *Value[time.Time] {
	return NewValue(dateHandler)
}

func DateSlice() *ValueSlice[time.Time] {
	return NewValueSlice(dateHandler)
}

func Int() *Value[int] {
	return NewValue(intHandler)
}

func IntFlag() *ValueFlag[int] {
	return NewValueFlag(intHandler)
}

func IntSlice() *ValueSlice[int] {
	return NewValueSlice(intHandler)
}

func String() *Value[string] {
	return NewValue(stringHandler)
}

func StringSlice() *ValueSlice[time.Time] {
	return NewValueSlice(dateHandler)
}
