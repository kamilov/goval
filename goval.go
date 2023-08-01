package goval

import (
	"encoding"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
)

type Settable interface {
	Set(value string) error
}

func Val(ptr interface{}, value string) error {
	return setValue(reflect.ValueOf(ptr), value)
}

func indirect(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}

		v = v.Elem()
	}

	return v
}

func setValue(v reflect.Value, value string) error {
	v = indirect(v)
	t := v.Type()

	if !v.CanAddr() {
		return errors.New("The value is unaddressable")
	}

	ival := v.Addr().Interface()

	if ival, ok := ival.(Settable); ok {
		return ival.Set(value)
	}

	if ival, ok := ival.(encoding.TextUnmarshaler); ok {
		return ival.UnmarshalText([]byte(value))
	}

	if ival, ok := ival.(encoding.BinaryUnmarshaler); ok {
		return ival.UnmarshalBinary([]byte(value))
	}

	switch t.Kind() {
	case reflect.String:
		v.SetString(value)
		break

	case reflect.Bool:
		val, err := strconv.ParseBool(value)

		if err != nil {
			return err
		}

		v.SetBool(val)
		break

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := strconv.ParseInt(value, 0, t.Bits())

		if err != nil {
			return err
		}

		v.SetInt(val)
		break

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(value, 0, t.Bits())

		if err != nil {
			return err
		}

		v.SetUint(val)
		break

	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(value, t.Bits())

		if err != nil {
			return err
		}

		v.SetFloat(val)
		break

	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			sl := reflect.ValueOf([]byte(value))

			v.Set(sl)
			return nil
		}
		fallthrough

	default:
		return json.Unmarshal([]byte(value), v.Addr().Interface())
	}

	return nil
}

func is(v reflect.Value, t reflect.Kind) bool {
	v = indirect(v)

	return t == v.Kind()
}

func IsString(ptr interface{}) bool {
	return is(reflect.ValueOf(ptr), reflect.String)
}

func IsBool(ptr interface{}) bool {
	return is(reflect.ValueOf(ptr), reflect.Bool)
}

func IsInt(ptr interface{}) bool {
	return is(reflect.ValueOf(ptr), reflect.Int)
}

func IsInt8(ptr interface{}) bool {
	return is(reflect.ValueOf(ptr), reflect.Int8)
}

func IsInt16(ptr interface{}) bool {
	return is(reflect.ValueOf(ptr), reflect.Int16)
}

func IsInt32(ptr interface{}) bool {
	return is(reflect.ValueOf(ptr), reflect.Int32)
}

func IsInt64(ptr interface{}) bool {
	return is(reflect.ValueOf(ptr), reflect.Int64)
}

func IsUint(ptr interface{}) bool {
	return is(reflect.ValueOf(ptr), reflect.Uint)
}

func IsUint8(ptr interface{}) bool {
	return is(reflect.ValueOf(ptr), reflect.Uint8)
}

func IsUint16(ptr interface{}) bool {
	return is(reflect.ValueOf(ptr), reflect.Uint16)
}

func IsUint32(ptr interface{}) bool {
	return is(reflect.ValueOf(ptr), reflect.Uint32)
}

func IsUint64(ptr interface{}) bool {
	return is(reflect.ValueOf(ptr), reflect.Uint64)
}

func IsFloat32(ptr interface{}) bool {
	return is(reflect.ValueOf(ptr), reflect.Float32)
}

func IsFloat64(ptr interface{}) bool {
	return is(reflect.ValueOf(ptr), reflect.Float64)
}

func IsSlice(ptr interface{}) bool {
	return is(reflect.ValueOf(ptr), reflect.Slice)
}

func IsMap(ptr interface{}) bool {
	return is(reflect.ValueOf(ptr), reflect.Map)
}

func IsStruct(ptr interface{}) bool {
	return is(reflect.ValueOf(ptr), reflect.Struct)
}
