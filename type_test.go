package n

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func Test_detectType(t *testing.T) {
	for s, typ := range map[string]reflect.Value{
		"true":          reflect.ValueOf(true),
		"false":         reflect.ValueOf(false),
		"text":          reflect.ValueOf(""),
		"[notarray!":    reflect.ValueOf(""),
		"samenotarray]": reflect.ValueOf(""),
		"semi;colon":    reflect.ValueOf(""),
		"10":            reflect.ValueOf(int64(0)),
		"u10":           reflect.ValueOf(uint64(0)),
		"1.1":           reflect.ValueOf(float64(1.1)),
		"(2+3i)":        reflect.ValueOf(complex128(complex(2, 3))),
		"[one]":         reflect.ValueOf([]string{}),
		"[a, b]":        reflect.ValueOf([]string{}),
		"[1, 2]":        reflect.ValueOf([]int64{}),
		"[1, a]":        reflect.ValueOf([]any{}),
		"[[a, b], [1, 2]]": reflect.ValueOf([]any{
			[]string{},
			[]int64{},
		}),
		"[[[[are],[you]],[crazy],[[?],[!]]],[1,2],[3,4]]": reflect.ValueOf([]any{
			[]any{
				[]any{
					[]string{},
					[]string{},
				},
				[]string{},
				[]any{
					[]string{},
					[]string{},
				},
			},
			[]int64{},
			[]int64{},
		}),
	} {
		t.Run(s, func(t *testing.T) {
			actualTyp, err := detectType(0, s)
			require.NoError(t, err)

			assert.Equal(t, typ.Kind().String(), actualTyp.Kind().String())
			assert.Equal(t, typ.Type().String(), actualTyp.Type().String())
		})
	}
	t.Run("too-deep", func(t *testing.T) {
		_, err := detectType(4, `[1,2]`)
		assert.Error(t, ErrTooDeep, err)
	})
	t.Run("empty", func(t *testing.T) {
		_, err := detectType(0, `  `)
		assert.Error(t, err)
	})
}
