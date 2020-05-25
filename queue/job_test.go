package queue

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

type student struct {
	Name string `json:"name"`
}

func (s *student) Handle(c *Context) {}

func TestJobToJSON(t *testing.T) {
	s := &student{"zed"}

	r := jobToJSON(s)

	if string(r) != "{\"name\":\"zed\",\"type_name\":\"student\"}" {
		t.Error("r wrong")
	}
}

func BenchmarkJobToJSON(b *testing.B) {
	s := &student{"zed"}

	for i := 0; i < b.N; i++ {
		jobToJSON(s)
	}
}

func TestJSONToJob(t *testing.T) {
	Register(&student{})
	json := "{\"name\":\"zed\",\"type_name\":\"student\"}"
	job := jsonToJob([]byte(json))
	// fmt.Printf("%+v\n", job)
	excpt := &student{"zed"}
	if !reflect.DeepEqual(excpt, job) {
		t.Error("type error")
	}
}

type hasSay interface {
	Say()
}
type dog struct {
	Name string `json:"name"`
}

func (d *dog) Say() {
	fmt.Printf("I am %s\n", d.Name)
}

func TestNormal(t *testing.T) {
	var d hasSay = &dog{"zed"}
	rt := reflect.TypeOf(reflect.ValueOf(d).Elem().Interface())
	data := "{\"name\":\"zed\"}"

	a := reflect.New(rt).Interface()
	json.Unmarshal([]byte(data), a)
	fmt.Printf("%+v\n", a)
	a.(hasSay).Say()
}
