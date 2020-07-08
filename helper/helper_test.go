package helper

import (
	"testing"
)

func TestStruct2Map(t *testing.T) {
	s := struct {
		Number int    `json:"number"`
		Str    string `json:"str"`
	}{1, "123"}

	m := Struct2Map(&s)

	if m["number"].(float64) != 1 || m["str"].(string) != "123" {
		t.Fatal("covert failed")
	}
}
