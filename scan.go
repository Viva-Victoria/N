package n

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

func anyToStrings(a any) []string {
	val := reflect.ValueOf(a)
	if val, slice := isSlice(val); slice {
		return stringSlice(val)
	}

	return []string{fmt.Sprintf("%v", a)}
}

func stringSlice(val reflect.Value) []string {
	items := make([]string, 0, val.Len())

	for i := 0; i < val.Len(); i++ {
		var (
			itemVal = val.Index(i)
			item    string
		)

		if itemVal, slice := isSlice(itemVal); slice {
			item = fmt.Sprintf(`[%s]`, strings.Join(stringSlice(itemVal), ", "))
		} else {
			item = fmt.Sprintf("%v", itemVal.Interface())
		}

		items = append(items, item)
	}

	return items
}

func isSlice(val reflect.Value) (reflect.Value, bool) {
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	return val, val.Kind() == reflect.Slice
}

func stringToAny(values []string, a any) error {
	ptr := reflect.ValueOf(a)

	if ptr.Kind() != reflect.Pointer {
		return fmt.Errorf(`type "%s" not a pointer`, ptr.Type())
	}

	val := ptr.Elem()
	if val.Kind() == reflect.Slice {
		return scanSlice(val, val, values, 0)
	}

	return scanOne(val, val, values[0])
}

func scanSlice(pointer, typer reflect.Value, values []string, deep int) error {
	var (
		typ   = typer.Type()
		count = len(values)
		err   error
	)

	array := reflect.MakeSlice(typ, count, count)
	defer func() {
		pointer.Set(array)
	}()

	for i, stringValue := range values {
		reflectValue := array.Index(i)
		typer = reflectValue

		if reflectValue.Kind() == reflect.Interface {
			typer, err = detectType(0, stringValue)
			if err != nil {
				return err
			}
		}

		if typer.Kind() == reflect.Slice {
			log.Println(stringValue, deep)
			if deep > MaxDeep {
				return ErrTooDeep
			}

			err = scanSlice(reflectValue, typer, splitByComma(stringValue), deep+1)
			if err != nil {
				return err
			}

			continue
		}

		err = scanOne(reflectValue, typer, stringValue)
		if err != nil {
			return err
		}
	}

	return nil
}

func scanOne(pointer, typer reflect.Value, s string) error {
	s = strings.Trim(s, " \t\n\r")

	var (
		directSet = pointer.Kind() != reflect.Interface
		err       error
	)

	var resultValue reflect.Value
	switch typer.Kind() {
	case reflect.Bool:
		var b bool
		b, err = strconv.ParseBool(s)
		if directSet {
			pointer.SetBool(b)
			return err
		}

		resultValue = reflect.ValueOf(b)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var i64 int64
		i64, err = strconv.ParseInt(s, 10, 64)
		if directSet {
			pointer.SetInt(i64)
			return err
		}

		resultValue = reflect.ValueOf(i64)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if s[0] == 'u' {
			s = s[1:]
		}

		var u64 uint64
		u64, err = strconv.ParseUint(s, 10, 64)
		if directSet {
			pointer.SetUint(u64)
			return err
		}

		resultValue = reflect.ValueOf(u64)

	case reflect.Float32, reflect.Float64:
		var f64 float64
		f64, err = strconv.ParseFloat(s, 64)
		if directSet {
			pointer.SetFloat(f64)
			return err
		}

		resultValue = reflect.ValueOf(f64)

	case reflect.Complex64, reflect.Complex128:
		var c128 complex128
		c128, err = strconv.ParseComplex(s, 128)
		if directSet {
			pointer.SetComplex(c128)
			return err
		}

		resultValue = reflect.ValueOf(c128)

	case reflect.String:
		if directSet {
			pointer.SetString(s)
			return err
		}

		resultValue = reflect.ValueOf(s)

	default:
		return fmt.Errorf(`unknown type: "%s"`, typer.Kind())
	}

	pointer.Set(resultValue)
	return nil
}
