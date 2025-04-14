package csvtag

import (
	"reflect"
	"strconv"
)

type InvalidMarshalError struct {
	Type reflect.Type
}

func (e *InvalidMarshalError) Error() string {
	if e.Type == nil {
		return "csv: Marshal(nil)"
	}
	if e.Type.Kind() != reflect.Pointer {
		return "csv: Marshal(non-pointer " + e.Type.String() + ")"
	}
	return "csv: Marshal(nil " + e.Type.String() + ")"
}

func Marshall(data any) ([][]string, error) {
	if reflect.TypeOf(data).Kind() != reflect.Slice {
		return nil, &InvalidMarshalError{reflect.TypeOf(data)}
	}
	return marshall(data)
}

func marshall(data any) ([][]string, error) {
	dataLen := reflect.ValueOf(data).Len()
	res := make([][]string, dataLen+1)

	t := reflect.TypeOf(data).Elem()
	//TODO: assume all fields have tag atm. Address this in future
	headers := make([]string, t.NumField())
	for i := range t.NumField() {
		field := t.Field(i)
		tag := field.Tag.Get("csv")
		headers[i] = tag
	}
	res[0] = headers

	for i := range dataLen {
		// TODO: len is equal to the number of entries with csv tag inside struct
		currVal := reflect.ValueOf(data).Index(i)
		row := parseRow(currVal)
		// first row is header
		res[i+1] = row
	}
	return res, nil
}

func parseRow(v reflect.Value) []string {
	row := make([]string, v.NumField())
	for i := range v.NumField() {
		currVal := v.FieldByIndex([]int{i})
		switch currVal.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			row[i] = strconv.FormatInt(currVal.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			row[i] = strconv.FormatUint(currVal.Uint(), 10)
		case reflect.String:
			row[i] = currVal.String()
		case reflect.Float32, reflect.Float64:
			row[i] = strconv.FormatFloat(currVal.Float(), 'g', 2, currVal.Type().Bits())
		}
	}
	return row
}
