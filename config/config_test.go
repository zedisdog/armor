package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"testing"
)

func TestConfig(t *testing.T) {
	var data = `
a: Easy!
b:
  c: 2
  d: 5
  e:
    f: '123'
    g: true
`
	c := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(data), &c)
	if err != nil {
		fmt.Println(err)
	}

	config := &yamlConfig{config: c}

	a := config.String("a")
	if a != "Easy!" {
		t.Errorf("except <Easy!> got %s", a)
	}

	d := config.Int("b.d")
	if d != 5 {
		t.Errorf("except <5> got %d", d)
	}

	f := config.String("b.e.f")
	if f != "123" {
		t.Errorf("except <'123'> got %s", f)
	}

	g := config.Bool("b.e.g")
	if g != true {
		t.Errorf("except <true> got %t", g)
	}
}
