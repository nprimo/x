package csvtag

import (
	"fmt"
	"reflect"
	"strconv"
)

type InvalidUnmarshalError struct {
	Type reflect.Type
}

func (e *InvalidUnmarshalError) Error() string {
	if e.Type == nil {
		return "csv: Unmarshal(nil)"
	}
	if e.Type.Kind() != reflect.Pointer {
		return "csv: Unmarshal(non-pointer " + e.Type.String() + ")"
	}
	return "csv: Unmarshal(nil " + e.Type.String() + ")"
}

// Unmarshal parse the content of a CSV file and store the result into the
// value pointed to by v
func Unmarshal(rows [][]string, v any) error {
	if reflect.TypeOf(v).Kind() != reflect.Pointer {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}
	if reflect.TypeOf(v).Elem().Kind() != reflect.Slice {
		return &InvalidUnmarshalError{reflect.TypeOf(v)}
	}
	header := rows[0]
	if err := checkValidHeader(header, v); err != nil {
		return err
	}
	d := newDecoder()
	if err := d.init(header, v); err != nil {
		return err
	}
	return d.unmarshal(rows, v)
}

// checkValidHeader check if the header has matching tags in the struct
// passed as v.
func checkValidHeader(header []string, v any) error {
	ptrEl := reflect.TypeOf(v).Elem()
	t := ptrEl.Elem()
	for i := range t.NumField() {
		field := t.Field(i)
		tag := field.Tag.Get("csv")
		if i := index(tag, header); i == -1 {
			return fmt.Errorf("tag(%s) not present in header %+v",
				tag, header,
			)
		}
	}
	return nil
}

type decoder struct {
	// header index to field name
	headerIdToFieldName map[int]string
}

func newDecoder() decoder {
	return decoder{
		headerIdToFieldName: map[int]string{},
	}
}

func (d *decoder) init(header []string, v any) error {
	ptrEl := reflect.TypeOf(v).Elem()
	t := ptrEl.Elem()
	for i := range t.NumField() {
		field := t.Field(i)
		tag := field.Tag.Get("csv")
		id := index(tag, header)
		d.headerIdToFieldName[id] = field.Name
	}
	return nil
}

// unmarshal store the data inside the array of struct v.
func (d *decoder) unmarshal(data [][]string, v any) error {
	sliceRv := reflect.MakeSlice(
		reflect.TypeOf(v).Elem(),
		len(data)-1,
		len(data)-1,
	)
	for i, row := range data[1:] {
		rv := sliceRv.Index(i)
		for rowIndex, fieldName := range d.headerIdToFieldName {
			f := rv.FieldByName(fieldName)
			switch f.Kind() {
			case reflect.String:
				f.SetString(row[rowIndex])
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				num, err := strconv.ParseInt(row[rowIndex], 10, 64)
				if err != nil {
					return err
				}
				f.SetInt(num)
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				num, err := strconv.ParseUint(row[rowIndex], 10, 64)
				if err != nil {
					return err
				}
				f.SetUint(num)
			case reflect.Float32, reflect.Float64:
				num, err := strconv.ParseFloat(row[rowIndex], f.Type().Bits())
				if err != nil || f.OverflowFloat(num) {
					return err
				}
				f.SetFloat(num)
			default:
				// TODO: other cases:
				// - time
				// - boolean
			}
		}
	}
	reflect.ValueOf(v).Elem().Set(sliceRv)
	return nil
}

// utils function to get the index of first occurrence of a string in an array
func index(target string, ss []string) int {
	for i, s := range ss {
		if s == target {
			return i
		}
	}
	return -1
}
