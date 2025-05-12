package cli

import (
	"fmt"
	"strings"
)

func Enum[T any](ptr *T, mapping map[string]T) *_enum[T] {
	e := &_enum[T]{
		mapping: mapping,
	}
	return e.Bind(ptr)
}

func EnumSlice[T any](ptr *[]T, mapping map[string]T) *_enumSlice[T] {
	e := &_enumSlice[T]{
		mapping: mapping,
	}
	return e.Bind(ptr)
}

type _enum[T any] struct {
	key     string
	mapping map[string]T
	value   *T
}

func (e *_enum[T]) Assign(str string) error {
	key, value, err := enumFind(e.mapping, str)
	if err != nil {
		return err
	}
	e.key = key
	*e.value = value
	return nil
}

func (e *_enum[T]) Bind(ptr *T) *_enum[T] {
	if ptr == nil {
		ptr = new(T)
	}
	e.value = ptr
	return e
}

func (v *_enum[T]) Consume(args *[]string) error {
	if len(*args) == 0 {
		return fmt.Errorf("enum param requires a value")
	}
	if err := v.Assign((*args)[0]); err != nil {
		return err
	}
	*args = (*args)[1:]
	return nil
}

func (v *_enum[T]) NonZero() string {
	return v.String()
}

func (e *_enum[T]) Reset() {
	e.key = ""
	*e.value = *new(T)
}

func (e _enum[T]) String() string {
	return e.key
}

func (e _enum[T]) Type() string {
	return enumType(e.mapping)
}

type _enumSlice[T any] struct {
	keys    []string
	mapping map[string]T
	values  *[]T
}

func (e *_enumSlice[T]) Assign(str string) error {
	for bit := range strings.SplitSeq(str, ",") {
		bit = strings.TrimSpace(bit)
		key, value, err := enumFind(e.mapping, bit)
		if err != nil {
			return err
		}
		e.keys = append(e.keys, key)
		*e.values = append(*e.values, value)
	}
	return nil
}

func (e *_enumSlice[T]) Bind(ptr *[]T) *_enumSlice[T] {
	if ptr == nil {
		ptr = new([]T)
	}
	e.values = ptr
	return e
}

func (v *_enumSlice[T]) Consume(args *[]string) error {
	for _, arg := range *args {
		if err := v.Assign(arg); err != nil {
			return err
		}
	}
	*args = nil
	return nil
}

func (v _enumSlice[T]) NonZero() string {
	return v.String()
}

func (e *_enumSlice[T]) Reset() {
	e.keys = nil
	*e.values = nil
}

func (e _enumSlice[T]) String() string {
	return strings.Join(e.keys, ",")
}

func (e _enumSlice[T]) Type() string {
	return enumType(e.mapping) + "..."
}

func enumFind[T any](mapping map[string]T, what string) (string, T, error) {
	what = strings.TrimSpace(what)
	for k, v := range mapping {
		if strings.EqualFold(k, what) {
			return k, v, nil
		}
	}
	return "", *new(T), fmt.Errorf("bad enum value %q", what)
}

func enumType[T any](mapping map[string]T) string {
	t := ""
	for k := range mapping {
		if t != "" {
			t += "|"
		}
		t += k
	}
	return t
}
