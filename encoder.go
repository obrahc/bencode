package bencode

import (
	"bytes"
	"fmt"
	"reflect"
)

type encoder struct {
	output bytes.Buffer
}

func (e *encoder) init() *encoder {
	return &encoder{bytes.Buffer{}}
}

func (e *encoder) encodeString(value string) {
	e.output.WriteString(fmt.Sprintf("%d", len(value)))
	e.output.WriteByte(':')
	e.output.WriteString(value)
}

func (e *encoder) encodeUint(value uint64) {
	e.output.WriteByte('i')
	e.output.WriteString(fmt.Sprintf("%d", value))
	e.output.WriteByte('e')
}
func (e *encoder) encodeList(rv reflect.Value) {
	e.output.WriteByte('l')
	indirect := reflect.Indirect(rv)
	for i := 0; i < indirect.Len(); i++ {
		switch indirect.Index(i).Kind() {
		case reflect.Uint64:
		case reflect.String:
			e.encodeString(indirect.Index(i).Interface().(string))
		case reflect.Struct:
		case reflect.Slice:
			e.encodeList(indirect.Index(i).Addr())
		}
	}
	e.output.WriteByte('e')
}

func (e *encoder) encodeDict(rv reflect.Value) error {
	e.output.WriteByte('d')
	indirect := reflect.Indirect(rv)
	for i := 0; i < indirect.NumField(); i++ {
		switch indirect.Type().Field(i).Type.Kind() {
		case reflect.Uint64:
			e.encodeString(indirect.Type().Field(i).Tag.Get("bencode"))
			e.encodeUint(indirect.Field(i).Interface().(uint64))
		case reflect.String:
			e.encodeString(indirect.Type().Field(i).Tag.Get("bencode"))
			e.encodeString(indirect.Field(i).Interface().(string))
		case reflect.Struct:
			e.encodeString(indirect.Type().Field(i).Tag.Get("bencode"))
			e.encodeDict(indirect.Field(i).Addr())
		case reflect.Slice:
			e.encodeString(indirect.Type().Field(i).Tag.Get("bencode"))
			e.encodeList(indirect.Field(i).Addr())
		}
	}
	e.output.WriteByte('e')
	return nil
}

func (e *encoder) marshal(v any) ([]byte, error) {
	e.encodeDict(reflect.ValueOf(v))
	return e.output.Bytes(), nil
}

func Encode(v any) ([]byte, error) {
	var e encoder
	e.init()
	return e.marshal(v)
}
