package bencode

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

type reader struct {
	*bytes.Reader
}

func newReader(data []byte) *reader {
	return &reader{bytes.NewReader(data)}
}

func (r *reader) ReadUntil(c byte) ([]byte, error) {
	res := []byte("")
	for {
		b, err := r.ReadByte()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return []byte{}, err
		}
		if b == c {
			break
		}
		res = append(res, b)
	}
	return res, nil
}

func (r *reader) readNBytes(n uint64) ([]byte, error) {
	res := []byte("")
	var i uint64
	for i = 0; i < n; i++ {
		b, err := r.ReadByte()
		if err != nil {
			return []byte(""), err
		}
		res = append(res, b)
	}
	return res, nil
}

func (r * reader) readStringLength() (uint64, error) {
	// Read the length of string
	size, err := r.ReadUntil(':')
	if err != nil {
		return 0, err
	}
	res, err := strconv.ParseUint(string(size), 10, 64)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *reader) ReadString() (string, error) {
	size, err := r.readStringLength()
	if err != nil {
		return "", err
	}
	b, err := r.readNBytes(size)
	return string(b), err
}

func (r *reader) readInteger() (string, error) {
	r.ReadByte()
	value, err := r.ReadUntil('e')
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func (r *reader) ReadList() ([]interface{}, error) {
	res := make([]interface{}, 0)
	r.ReadByte()
	var v interface{}
	for {
		e, err := r.ReadByte() 
		if err != nil {
			return res, err
		}
		if e == 'e' {
			break
		}
		r.UnreadByte() 
		vt, err := r.getValueType()
		if err != nil {
			return res, err
		}
		switch vt {
		case reflect.Uint64:
			v, err = r.readInteger()
			if err != nil {
				return res, err
			}
		case reflect.String:
			v, err = r.ReadString()
			if err != nil {
				return res, err
			}
		case reflect.Array:
			v, err = r.ReadList()
			if err != nil {
				return res, err
			}
		case reflect.Map:
			v, err = r.ReadDictionary()
		}
		res = append(res, v)
	}
	return res, nil
}

func (r *reader) ReadDictionary() (map[string]interface{}, error) {
	res := make(map[string]interface{})
	r.ReadByte()
	var v interface{}
	for {
		s, _ := r.ReadString()
		vt, err := r.getValueType()
		if err != nil {
			return res, err
		}
		switch vt {
		case reflect.Uint64:
			v, err = r.readInteger()
			if err != nil {
				return res, err
			}
		case reflect.String:
			v, err = r.ReadString()
			if err != nil {
				return res, err
			}
		case reflect.Array:
			v, err = r.ReadList()
			if err != nil {
				return res, err
			}
		case reflect.Map:
			v, err = r.ReadDictionary()
		}
		res[s] = v
		if err != nil {
			return res, err
		}
		e, err := r.ReadByte() 
		if err != nil {
			return res, err
		}
		if e == 'e' {
			break
		}
		r.UnreadByte() 
	}
	return res, nil
}

func (r *reader) getValueType() (reflect.Kind, error) {
	b, err := r.ReadByte()
	defer r.UnreadByte()
	if err != nil {
		return reflect.Invalid, err
	}
	switch b {
	case 'i':
		return reflect.Uint64, nil
	case 'l':
		return reflect.Array, nil
	case 'd':
		return reflect.Map, nil
	default:
		return reflect.String, nil
	}
}

func (r *reader) readValue(vt reflect.Kind) (string, error) {
	var res string	
	switch vt {
	case reflect.Uint64:
		res, err := r.readInteger()
		if err != nil {
			return res, err
		}
		return res, nil
	case reflect.String:
		res, err := r.ReadString()
		if err != nil {
			return res, err
		}
		return res, nil
	}
	return res, nil
}

type decoder struct {
	reader *reader 
}

func (d *decoder) init(data []byte) {
	d.reader = newReader(data)
}

func (d *decoder) setArray(rv reflect.Value, values []interface{}) error {
	indirect := reflect.Indirect(rv)
	for k, value := range values {
		switch decodedValue := value.(type) {
		case string:
			indirect.Set(reflect.Append(indirect, reflect.ValueOf(decodedValue)))
		case []interface{}:
			n := reflect.New(indirect.Type().Elem())
			indirect.Set(reflect.Append(indirect, reflect.Indirect(n)))
			d.setArray(indirect.Index(k).Addr(), decodedValue)
		}
	}
	return nil
}

func (d *decoder) setStruct(rv reflect.Value, values map[string]interface{}) error {
	indirect := reflect.Indirect(rv)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		return fmt.Errorf("passed object is not a pointer or is nil: %v", reflect.TypeOf(rv))
	}
	if indirect.Kind() != reflect.Struct {
		return fmt.Errorf("passed object is not a struct: %v", reflect.TypeOf(rv))
	}
	fields := getStructFields(indirect.Type())
	for key, value := range values {
		if field, ok := fields[key]; ok {
			switch decodedValue := value.(type) {
			case string:
				switch field.Type.Kind() {
				case reflect.Uint64:
					insertValue, err := stringToUint(decodedValue)
					if err != nil {
						return err
					}
					indirect.FieldByName(field.Name).SetUint(insertValue)
				case reflect.String:
					indirect.FieldByName(field.Name).SetString(decodedValue)
				}
			case map[string]interface{}:
				n := reflect.New(indirect.FieldByName(field.Name).Type())
				indirect.FieldByName(field.Name).Set(reflect.Indirect(n))
				d.setStruct(indirect.FieldByName(field.Name).Addr(), decodedValue)
			case []interface{}:
				d.setArray(indirect.FieldByName(field.Name).Addr(), decodedValue)
			}
		}
	}
	return nil
}

func (d *decoder) unmarshal(v any) error{
	values, err := d.reader.ReadDictionary()
	if err != nil {
		return err
	}
	rv := reflect.ValueOf(v)
	return d.setStruct(rv, values)
}

func decode(input []byte, v any) {
	var d decoder
	d.init(input)
	d.unmarshal(v)
}

func getStructFields(t reflect.Type) map[string]reflect.StructField {
	r := make(map[string]reflect.StructField)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		r[string(f.Tag.Get("bencode"))] = f
	}
	return r
}

func stringToUint(s string) (uint64, error) {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, nil
	}
	return i, nil
}

