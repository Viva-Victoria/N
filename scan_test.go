package n

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
	"unsafe"
)

func Test_isSlice(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		var i int
		expectedV := reflect.ValueOf(i)
		actualV, ok := isSlice(expectedV)
		assert.Equal(t, expectedV, actualV)
		assert.False(t, ok)
	})
	t.Run("pointer", func(t *testing.T) {
		var p *int
		expectedV := reflect.ValueOf(p)
		actualV, ok := isSlice(expectedV)
		assert.Equal(t, expectedV.Elem(), actualV)
		assert.False(t, ok)
	})
	t.Run("array", func(t *testing.T) {
		var a []int
		expectedV := reflect.ValueOf(a)
		actualV, ok := isSlice(expectedV)
		assert.Equal(t, expectedV, actualV)
		assert.True(t, ok)
	})
	t.Run("array-pointer", func(t *testing.T) {
		var a []int
		expectedV := reflect.ValueOf(&a)
		actualV, ok := isSlice(expectedV)
		assert.Equal(t, expectedV.Elem(), actualV)
		assert.True(t, ok)
	})
}

type arrayTestCase struct {
	text  any
	array any
}

var (
	_intCases = map[string]any{
		"10": int(10),
		"12": int8(12),
		"14": int16(14),
		"16": int32(16),
		"18": int64(18),
	}
	_uintCases = map[string]any{
		"11": uint(11),
		"13": uint8(13),
		"15": uint16(15),
		"17": uint32(17),
		"19": uint64(19),
	}
	_floatCases = map[string]any{
		"10.11": float32(10.11),
		"12.13": float64(12.13),
	}
	_complexCases = map[string]any{
		"(2+3i)": complex64(complex(2, 3)),
		"(4+6i)": complex128(complex(4, 6)),
	}
	_boolCases = map[string]any{
		"true":  true,
		"false": false,
	}
	_stringCases = map[string]string{
		"text":   "text",
		"Russia": "Russia",
	}

	_arrayComplexCases = arrayTestCase{
		text:  []string{"(2+3i)", "(4+6i)"},
		array: []any{complex64(complex(2, 3)), complex128(complex(4, 6))},
	}
	_arrayBoolCases = arrayTestCase{
		text:  []string{"true", "false"},
		array: []bool{true, false},
	}
	_arrayStringCases = arrayTestCase{
		text:  []string{"text", "Russia"},
		array: []any{"text", "Russia"},
	}
	_arrayArrayIntCases = arrayTestCase{
		text: []string{"[1, 2]", "[-1, -2]"},
		array: [][]any{
			{
				uint(1),
				uint32(2),
			},
			{
				int(-1),
				int32(-2),
			},
		},
	}
	_arrayArrayStringCases = arrayTestCase{
		text: []string{`[text, Russia]`, `[word, verb]`},
		array: [][]any{
			{
				"text",
				"Russia",
			},
			{
				"word",
				"verb",
			},
		},
	}
)

func Test_anyToString(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		for s, i := range _intCases {
			assert.EqualValues(t, []string{s}, anyToStrings(i))
		}
	})
	t.Run("uint", func(t *testing.T) {
		for s, i := range _uintCases {
			assert.EqualValues(t, []string{s}, anyToStrings(i))
		}
	})
	t.Run("float", func(t *testing.T) {
		for s, f := range _floatCases {
			assert.EqualValues(t, []string{s}, anyToStrings(f))
		}
	})
	t.Run("complex", func(t *testing.T) {
		for s, c := range _complexCases {
			assert.EqualValues(t, []string{s}, anyToStrings(c))
		}
	})
	t.Run("string", func(t *testing.T) {
		for s, b := range _stringCases {
			assert.EqualValues(t, []string{s}, anyToStrings(b))
		}
	})
	t.Run("bool", func(t *testing.T) {
		for s, b := range _boolCases {
			assert.EqualValues(t, []string{s}, anyToStrings(b))
		}
	})
	t.Run("array", func(t *testing.T) {
		t.Run("int", func(t *testing.T) {
			assert.EqualValues(t,
				[]string{"10", "20", "30", "15", "30", "45"},
				anyToStrings([]any{int(10), int8(20), int32(30), uint(15), uint16(30), uint64(45)}),
			)
		})
		t.Run("float", func(t *testing.T) {
			assert.EqualValues(t,
				[]string{"10.11", "12.13"},
				anyToStrings([]any{float32(10.11), float64(12.13)}),
			)
		})
		t.Run("complex", func(t *testing.T) {
			assert.EqualValues(t,
				_arrayComplexCases.text,
				anyToStrings(_arrayComplexCases.array),
			)
		})
		t.Run("bool", func(t *testing.T) {
			assert.EqualValues(t,
				_arrayBoolCases.text,
				anyToStrings(_arrayBoolCases.array),
			)
		})
		t.Run("string", func(t *testing.T) {
			assert.EqualValues(t,
				_arrayStringCases.text,
				anyToStrings(_arrayStringCases.array),
			)
		})
		t.Run("array", func(t *testing.T) {
			t.Run("int", func(t *testing.T) {
				assert.EqualValues(t,
					_arrayArrayIntCases.text,
					anyToStrings(_arrayArrayIntCases.array),
				)
			})
			t.Run("string", func(t *testing.T) {
				assert.EqualValues(t,
					_arrayArrayStringCases.text,
					anyToStrings(_arrayArrayStringCases.array),
				)
			})
		})
	})
}

func Test_stringToAny(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		testValue(t, []string{"10"}, int(10))
		testValue(t, []string{"12"}, int8(12))
		testValue(t, []string{"14"}, int16(14))
		testValue(t, []string{"16"}, int32(16))
		testValue(t, []string{"18"}, int64(18))
	})
	t.Run("uint", func(t *testing.T) {
		testValue(t, []string{"11"}, uint(11))
		testValue(t, []string{"13"}, uint8(13))
		testValue(t, []string{"15"}, uint16(15))
		testValue(t, []string{"17"}, uint32(17))
		testValue(t, []string{"19"}, uint64(19))
	})
	t.Run("float", func(t *testing.T) {
		testValue(t, []string{"10.11"}, float32(10.11))
		testValue(t, []string{"12.13"}, float64(12.13))
	})
	t.Run("complex", func(t *testing.T) {
		testValue(t, []string{"(2+3i)"}, complex64(complex(2, 3)))
		testValue(t, []string{"(4+6i)"}, complex128(complex(4, 6)))
	})
	t.Run("string", func(t *testing.T) {
		testValue(t, []string{"text"}, "text")
		testValue(t, []string{"Russia"}, "Russia")
	})
	t.Run("bool", func(t *testing.T) {
		testValue(t, []string{"true"}, true)
		testValue(t, []string{"false"}, false)
	})
	t.Run("array", func(t *testing.T) {
		t.Run("bad", func(t *testing.T) {
			t.Run("parse-error", func(t *testing.T) {
				expected := []int64{10}
				actual := expected

				require.Error(t, stringToAny([]string{"A10"}, &actual))
			})
			t.Run("bad-value", func(t *testing.T) {
				var a []any
				require.Error(t, stringToAny([]string{"    "}, &a))
			})
			t.Run("too-deep", func(t *testing.T) {
				var a [][][][][]any
				require.Error(t, stringToAny([]string{"[[[[[1]]]]]"}, &a))
			})
		})
		t.Run("int", func(t *testing.T) {
			t.Run("auto", func(t *testing.T) {
				expected := []any{int64(10), int64(20), int64(30), int64(40), int64(50)}
				actual := expected

				require.NoError(t, stringToAny([]string{"10", "20", "30", "40", "50"}, &actual))
				assert.EqualValues(t, expected, actual)
			})
			t.Run("manual", func(t *testing.T) {
				expected := []int{10, 20, 30, 40, 50}
				actual := expected

				require.NoError(t, stringToAny([]string{"10", "20", "30", "40", "50"}, &actual))
				assert.EqualValues(t, expected, actual)
			})
		})
		t.Run("float", func(t *testing.T) {
			t.Run("auto", func(t *testing.T) {
				expected := []any{float64(1.1), float64(1.2), float64(1.3)}
				actual := expected

				require.NoError(t, stringToAny([]string{"1.1", "1.2", "1.3"}, &actual))
				assert.EqualValues(t, expected, actual)
			})
			t.Run("manual", func(t *testing.T) {
				expected := []float64{1.1, 1.2, 1.3}
				actual := expected

				require.NoError(t, stringToAny([]string{"1.1", "1.2", "1.3"}, &actual))
				assert.EqualValues(t, expected, actual)
			})
		})
		t.Run("complex", func(t *testing.T) {
			t.Run("auto", func(t *testing.T) {
				expected := []any{complex(2, 3), complex(4, 6)}
				actual := expected

				require.NoError(t, stringToAny([]string{"(2+3i)", "(4+6i)"}, &actual))
				assert.EqualValues(t, expected, actual)
			})
			t.Run("manual", func(t *testing.T) {
				expected := []complex128{complex(2, 3), complex(4, 6)}
				actual := expected

				require.NoError(t, stringToAny([]string{"(2+3i)", "(4+6i)"}, &actual))
				assert.EqualValues(t, expected, actual)
			})
		})
		t.Run("bool", func(t *testing.T) {
			t.Run("auto", func(t *testing.T) {
				expected := []any{true, false}
				actual := expected

				require.NoError(t, stringToAny([]string{"true", "false"}, &actual))
				assert.EqualValues(t, expected, actual)
			})
			t.Run("manual", func(t *testing.T) {
				expected := []bool{true, false}
				actual := expected

				require.NoError(t, stringToAny([]string{"true", "false"}, &actual))
				assert.EqualValues(t, expected, actual)
			})
		})
		t.Run("string", func(t *testing.T) {
			t.Run("auto", func(t *testing.T) {
				expected := []any{"text", "Russia"}
				actual := expected

				require.NoError(t, stringToAny([]string{"text", "Russia"}, &actual))
				assert.EqualValues(t, expected, actual)
			})
			t.Run("manual", func(t *testing.T) {
				expected := []string{"text", "Russia"}
				actual := expected

				require.NoError(t, stringToAny([]string{"text", "Russia"}, &actual))
				assert.EqualValues(t, expected, actual)
			})
		})
		t.Run("array", func(t *testing.T) {
			t.Run("int", func(t *testing.T) {
				expected := []any{[]string{"a", "b"}, []int64{1, 2}}
				actual := expected

				require.NoError(t, stringToAny([]string{"[a, b]", "[1, 2]"}, &actual))
				assert.EqualValues(t, expected, actual)
			})
			t.Run("uint", func(t *testing.T) {
				expected := []any{[]string{"a", "b"}, []uint64{1, 2}}
				actual := expected

				require.NoError(t, stringToAny([]string{"[a, b]", "[u1, u2]"}, &actual))
				assert.EqualValues(t, expected, actual)
			})
			t.Run("string", func(t *testing.T) {
				expected := [][]string{{"a", "b"}, {"1", "2"}}
				actual := expected

				require.NoError(t, stringToAny([]string{"[a, b]", "[1, 2]"}, &actual))
				assert.EqualValues(t, expected, actual)
			})
		})
	})
	t.Run("non-pointer", func(t *testing.T) {
		var a int
		err := stringToAny([]string{"10"}, a)
		assert.Error(t, err)
	})
	t.Run("unknown-type", func(t *testing.T) {
		var a unsafe.Pointer
		err := stringToAny([]string{"a"}, &a)
		assert.Error(t, err)
	})
}

func testValue[T any](t *testing.T, input []string, v T) {
	var temp T
	require.NoError(t, stringToAny(input, &temp))
	assert.Equal(t, v, temp)
}
