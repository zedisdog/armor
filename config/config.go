package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type Config interface {
	String(key string) string
	Int(key string) int
	Interface(key string) interface{}
	Bool(key string) bool
}

type yamlConfig struct {
	config map[string]interface{}
}

func (y *yamlConfig) String(key string) string {
	return y.getValue(key).(string)
}

func (y *yamlConfig) Int(key string) int {
	return y.getValue(key).(int)
}

func (y *yamlConfig) Interface(key string) interface{} {
	return y.getValue(key)
}

func (y *yamlConfig) Bool(key string) bool {
	return y.getValue(key).(bool)
}

func (y *yamlConfig) getValue(key string) interface{} {
	keys := strings.Split(key, ".")
	config := y.config
	for index, key := range keys {
		result, ok := config[key]

		if !ok {
			return nil
		}

		if index+1 < len(keys) {
			config = make(map[string]interface{})
			for key, item := range result.(map[interface{}]interface{}) {
				config[key.(string)] = item
			}
		} else {
			return result
		}
	}

	return nil
}

func LoadYaml(file string) (Config, error) {
	data, err := getFileContent(file)
	if err != nil {
		return nil, err
	}
	c := make(map[string]interface{})
	err = yaml.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}

	return &yamlConfig{config: c}, nil
}

func getFileContent(file string) ([]byte, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}
