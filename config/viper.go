package config

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
)

func New(path string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return v, nil
}

var ProviderSet = wire.NewSet(New)
