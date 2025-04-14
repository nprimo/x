package csvtag

import (
	"reflect"
	"testing"
)

type User struct {
	Name  string `csv:"name"`
	Email string `csv:"email"`
}

func TestUnmarshallUser(t *testing.T) {
	data := [][]string{
		{"name", "email"},
		{"mario", "m@c.com"},
		{"luigi", "l@c.com"},
	}
	want := []User{
		{"mario", "m@c.com"},
		{"luigi", "l@c.com"},
	}

	var vv []User
	err := Unmarshal(data, &vv)
	if err != nil {
		t.Fatalf("did not expect to fail: %s\n", err)
	}

	if !reflect.DeepEqual(want, vv) {
		t.Fatalf("got %+v, want %+v\n", vv, want)
	}
}

type A struct {
	A string `csv:"a"`
	B int    `csv:"b"`
}

func TestUnmarshallWithNum(t *testing.T) {
	data := [][]string{
		{"a", "b"},
		{"a", "1"},
		{"b", "0"},
		{"c", "-1"},
	}
	want := []A{
		{"a", 1},
		{"b", 0},
		{"c", -1},
	}

	var vv []A
	err := Unmarshal(data, &vv)
	if err != nil {
		t.Fatalf("did not expect to fail: %s\n", err)
	}

	if !reflect.DeepEqual(want, vv) {
		t.Fatalf("got %+v, want %+v\n", vv, want)
	}
}

type B struct {
	A string  `csv:"a"`
	B float32 `csv:"b"`
	C float64 `csv:"c"`
}

func TestUnmarshallWitFloat(t *testing.T) {
	data := [][]string{
		{"a", "b", "c"},
		{"a", "1", "1.2"},
		{"b", "0", "0.3"},
		{"c", "-1", "-0.3"},
	}
	want := []B{
		{"a", 1, 1.2},
		{"b", 0, 0.3},
		{"c", -1, -0.3},
	}

	var vv []B
	err := Unmarshal(data, &vv)
	if err != nil {
		t.Fatalf("did not expect to fail: %s\n", err)
	}

	if !reflect.DeepEqual(want, vv) {
		t.Fatalf("got %+v, want %+v\n", vv, want)
	}
}

func TestCheckValidHeader(t *testing.T) {
	header := []string{
		"name", "email",
	}
	v := []User{}

	err := checkValidHeader(header, &v)
	if err != nil {
		t.Fatalf("expect no error with %+v and %+v: %s\n", header, v, err)
	}
}

func TestDecoderInit(t *testing.T) {
	d := newDecoder()
	header := []string{"name", "email"}
	vv := []User{}
	if err := d.init(header, &vv); err != nil {
		t.Fatalf("should not fail: %q\n", err)
	}
	if d.headerIdToFieldName[0] != "Name" {
		t.Fatalf("expected mapper 0 to be name: %s\n", d.headerIdToFieldName[0])
	}
	if d.headerIdToFieldName[1] != "Email" {
		t.Fatalf("expected mapper 0 to be name: %s\n", d.headerIdToFieldName[0])
	}
}
