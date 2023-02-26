package n

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	MaxDeep = 3
)

var (
	ErrTooDeep = errors.New("can't parse more arrays deeper than [][][][]")
)

func detectType(deep int, a string) (reflect.Value, error) {
	if deep > MaxDeep+1 {
		return reflect.Value{}, ErrTooDeep
	}

	if strings.EqualFold(a, "true") || strings.EqualFold(a, "false") {
		var b bool
		return reflect.ValueOf(b), nil
	}

	l := len(a)
	if l > 1 && a[0] == '[' && a[l-1] == ']' {
		return detectArrayType(deep+1, a)
	}

	var (
		complexStage uint8
		integer      bool
		float        bool
		unsignedInt  bool
	)
	for _, r := range []rune(a) {
		if r == ' ' || r == '\n' || r == '\t' || r == '\r' {
			continue
		}

		if r == 'u' {
			unsignedInt = true
			continue
		}

		if r >= '0' && r <= '9' && complexStage == 0 {
			integer = !float
			continue
		}

		if r == '.' && !float && complexStage == 0 {
			integer = false
			float = true
			continue
		}

		if r == '(' && complexStage == 0 {
			complexStage++
			continue
		}

		if complexStage > 0 {
			if r == '+' && complexStage == 1 {
				complexStage++
				continue
			}

			if r == 'i' && complexStage == 2 {
				complexStage++
				continue
			}

			if r == ')' && complexStage == 3 {
				var c complex128
				return reflect.ValueOf(c), nil
			}

			if r >= '0' && r <= '9' {
				continue
			}
		}

		return reflect.ValueOf(a), nil
	}

	if float {
		var f float64
		return reflect.ValueOf(f), nil
	}
	if integer {
		if unsignedInt {
			var i uint64
			return reflect.ValueOf(i), nil
		}

		var i int64
		return reflect.ValueOf(i), nil
	}

	return reflect.Value{}, fmt.Errorf(`can't detect type for "%s"`, a)
}

func detectArrayType(deep int, a string) (reflect.Value, error) {
	var (
		items      = splitByComma(a)
		values     = make([]reflect.Value, 0, len(items))
		consistent = true
	)

	for i, item := range items {
		itemValue, err := detectType(deep, item)
		if err != nil {
			return reflect.Value{}, err
		}

		if i > 0 && values[i-1].Type() != itemValue.Type() {
			consistent = false
		}

		values = append(values, itemValue)
	}

	if consistent {
		return reflect.New(reflect.SliceOf(values[0].Type())).Elem(), nil
	}

	anyArray := make([]any, 0, len(values))
	sliceValue := reflect.ValueOf(anyArray)
	sliceValue = reflect.Append(sliceValue, values...)
	return sliceValue, nil
}
