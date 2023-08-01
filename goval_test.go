package goval

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

type (
	set             bool
	unmarshalText   uint64
	unmarshalBinary string
)

func (s *set) Set(value string) error {
	v, err := strconv.ParseBool(strings.ToLower(value))

	if err != nil {
		return err
	}

	*s = set(v)

	return nil
}

func (u *unmarshalText) UnmarshalText(data []byte) error {
	var v uint64

	v, err := strconv.ParseUint(string(data), 10, 0)

	if err != nil {
		return err
	}

	*u = unmarshalText(v)

	return nil
}

func (u *unmarshalBinary) UnmarshalBinary(data []byte) error {
	*u = unmarshalBinary(string(data) + "ok")
	return nil
}

func TestAll(t *testing.T) {
	var (
		tString          string
		tStringPoint     *string
		tBool            bool
		tInt             int
		tUint64          uint64
		tFloat32         float32
		tSliceByte       []byte
		tSliceInt        []int
		tSliceString     []string
		tMap             map[string]int
		tSet             set
		tUnmarshalText   unmarshalText
		tUnmarshalBibary unmarshalBinary
		tests            = []struct {
			name      string
			point     reflect.Value
			checkFunc func(ptr interface{}) bool
			value     string
			expected  interface{}
			isEqual   bool
			isError   bool
		}{
			{"test_string", reflect.ValueOf(&tString), IsString, "abc", "abc", true, false},
			{"test_point_string", reflect.ValueOf(&tStringPoint), IsString, "abc", "abc", true, false},
			{"test_bool", reflect.ValueOf(&tBool), IsBool, "TRUE", true, true, false},
			{"test_bool_failure", reflect.ValueOf(&tBool), IsBool, "TRuE", true, true, true},
			{"test_int", reflect.ValueOf(&tInt), IsInt, "1", int(1), true, false},
			{"test_int_failure", reflect.ValueOf(&tStringPoint), IsInt, "1", int64(1), false, false},
			{"test_uint", reflect.ValueOf(&tUint64), IsUint64, "1", uint64(1), true, false},
			{"test_uint_failure_type", reflect.ValueOf(&tUint64), IsUint32, "1", uint32(1), false, false},
			{"test_uint_failure_value", reflect.ValueOf(&tUint64), IsUint64, "a1", uint64(1), false, true},
			{"test_float", reflect.ValueOf(&tFloat32), IsFloat32, "1", float32(1), true, false},
			{"test_float_failure_type", reflect.ValueOf(&tFloat32), IsFloat64, "1.1", float64(1.1), false, false},
			{"test_float_failure_value", reflect.ValueOf(&tFloat32), IsFloat32, "a1.1", float32(1.1), false, true},
			{"test_slice_byte", reflect.ValueOf(&tSliceByte), IsSlice, "abc", []byte("abc"), true, false},
			{"test_slice_int", reflect.ValueOf(&tSliceInt), IsSlice, "[1, 2]", []int{1, 2}, true, false},
			{"test_slice_string", reflect.ValueOf(&tSliceString), IsSlice, `["a", "b", "c"]`, []string{"a", "b", "c"}, true, false},
			{"test_map", reflect.ValueOf(&tMap), IsMap, `{"a": 1, "b": 2}`, map[string]int{"a": 1, "b": 2}, true, false},
			{"test_map_unparsed", reflect.ValueOf(&tMap), IsMap, `"a": 1, "b": 2`, "", true, true},
			{"test_settable_interface_1", reflect.ValueOf(&tSet), IsBool, "1", set(true), true, false},
			{"test_settable_interface_True", reflect.ValueOf(&tSet), IsBool, "True", set(true), true, false},
			{"test_settable_interface_TrUe", reflect.ValueOf(&tSet), IsBool, "TrUe", set(true), true, false},
			{"test_unmarshal_text", reflect.ValueOf(&tUnmarshalText), IsUint64, "1", unmarshalText(1), true, false},
			{"test_unmarshal_binary", reflect.ValueOf(&tUnmarshalBibary), IsString, "1", unmarshalBinary("1ok"), true, false},
			{"test_unaddressable", reflect.ValueOf("test"), IsString, "test", "test", true, true},
		}
	)

	for _, test := range tests {
		err := Val(test.point.Interface(), test.value)

		if test.isError {
			assert.NotNil(t, err, test.name)
		} else if assert.Nil(t, err, test.name) {
			actual := indirect(test.point)

			if test.isEqual {
				assert.True(t, reflect.DeepEqual(test.expected, actual.Interface()), test.name)
				assert.True(t, test.checkFunc(actual.Interface()), test.name)
			} else {
				assert.False(t, reflect.DeepEqual(test.expected, actual.Interface()), test.name)
				assert.False(t, test.checkFunc(actual.Interface()), test.name)
			}
		}
	}
}
