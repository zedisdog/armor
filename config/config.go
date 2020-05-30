package config

import (
	"github.com/joho/godotenv"
	"github.com/zedisdog/armor/file"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

func init() {
	if file.FileExists("./.env") {
		if err := godotenv.Load("./.env"); err != nil {
			panic(err)
		}
	}
	if _, err := LoadYaml(os.Getenv("ARMOR_CONFIG_FILE")); err != nil {
		panic(err)
	}
}

var Conf *Config

type Configure interface {
	String(key string) string
	Int(key string) int
	Interface(key string) interface{}
	Bool(key string) bool
	Bytes(key string) []byte
}

type Config struct {
	config map[string]interface{}
}

func (y *Config) String(key string) string {
	return y.getValue(key).(string)
}

func (y *Config) Int(key string) int {
	return y.getValue(key).(int)
}

func (y *Config) Interface(key string) interface{} {
	return y.getValue(key)
}

func (y *Config) Bool(key string) bool {
	return y.getValue(key).(bool)
}

func (y *Config) Bytes(key string) []byte {
	return []byte(y.getValue(key).(string))
}

func (y *Config) getValue(key string) interface{} {
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

// LoadYaml load yaml file in config, only once
func LoadYaml(file string) (Configure, error) {
	if Conf == nil {
		data, err := getFileContent(file)
		if err != nil {
			return nil, err
		}
		c := make(map[string]interface{})
		err = yaml.Unmarshal(data, &c)
		if err != nil {
			return nil, err
		}
		Conf = &Config{config: c}
	}

	return Conf, nil
}

func getFileContent(file string) ([]byte, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}
