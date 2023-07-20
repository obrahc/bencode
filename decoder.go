package bencode

import (
	"fmt"
	"reflect"
	"strconv"
)

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

func (d *decoder) readStringLength() (uint64, error) {
	size, err := d.reader.readUntil(':')
	if err != nil {
		return 0, err
	}
	res, err := strconv.ParseUint(string(size), 10, 64)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (d *decoder) readString() (string, error) {
	size, err := d.readStringLength()
	if err != nil {
		return "", err
	}
	b, err := d.reader.readNBytes(size)
	return string(b), err
}

func (d *decoder) readInteger() (string, error) {
	d.reader.ReadByte()
	value, err := d.reader.readUntil('e')
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func (d *decoder) getValueType() (reflect.Kind, error) {
	b, err := d.reader.ReadByte()
	defer d.reader.UnreadByte()
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

func (d *decoder) readList() ([]interface{}, error) {
	res := make([]interface{}, 0)
	d.reader.ReadByte()
	var v interface{}
	for {
		e, err := d.reader.ReadByte() 
		if err != nil {
			return res, err
		}
		if e == 'e' {
			break
		}
		d.reader.UnreadByte() 
		vt, err := d.getValueType()
		if err != nil {
			return res, err
		}
		switch vt {
		case reflect.Uint64:
			v, err = d.readInteger()
			if err != nil {
				return res, err
			}
		case reflect.String:
			v, err = d.readString()
			if err != nil {
				return res, err
			}
		case reflect.Array:
			v, err = d.readList()
			if err != nil {
				return res, err
			}
		case reflect.Map:
			v, err = d.readDictionary()
		}
		res = append(res, v)
	}
	return res, nil
}

func (d *decoder) readDictionary() (map[string]interface{}, error) {
	res := make(map[string]interface{})
	d.reader.ReadByte()
	var v interface{}
	for {
		s, _ := d.readString()
		vt, err := d.getValueType()
		if err != nil {
			return res, err
		}
		switch vt {
		case reflect.Uint64:
			v, err = d.readInteger()
			if err != nil {
				return res, err
			}
		case reflect.String:
			v, err = d.readString()
			if err != nil {
				return res, err
			}
		case reflect.Array:
			v, err = d.readList()
			if err != nil {
				return res, err
			}
		case reflect.Map:
			v, err = d.readDictionary()
		}
		res[s] = v
		if err != nil {
			return res, err
		}
		e, err := d.reader.ReadByte() 
		if err != nil {
			return res, err
		}
		if e == 'e' {
			break
		}
		d.reader.UnreadByte() 
	}
	return res, nil
}

func (d *decoder) unmarshal(v any) error{
	values, err := d.readDictionary()
	if err != nil {
		return err
	}
	rv := reflect.ValueOf(v)
	return d.setStruct(rv, values)
}

func Decode(input []byte, v any) error {
	var d decoder
	d.init(input)
	return d.unmarshal(v)
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

