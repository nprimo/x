package csvtag

import (
	"reflect"
	"testing"
)

func TestMarshallStringInt(t *testing.T) {
	type A struct {
		A string `csv:"a"`
		B int    `csv:"b"`
	}

	data := []A{
		{"b", 1},
		{"c", 2},
	}

	want := [][]string{
		{"a", "b"},
		{"b", "1"},
		{"c", "2"},
	}

	got, err := Marshall(data)
	if err != nil {
		t.Fatalf("didn't expect to fail Marshall with data %v: %s\n", data, err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %+v, got %+v\n", want, got)
	}
}

func TestMarshallStringFloat32(t *testing.T) {
	type A struct {
		A string  `csv:"a"`
		B float32 `csv:"b"`
	}

	data := []A{
		{"b", 1.2},
		{"c", 2},
	}

	want := [][]string{
		{"a", "b"},
		{"b", "1.2"},
		{"c", "2"},
	}

	got, err := Marshall(data)
	if err != nil {
		t.Fatalf("didn't expect to fail Marshall with data %v: %s\n", data, err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Fatalf("want %+v, got %+v\n", want, got)
	}
}
