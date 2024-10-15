// config.go
package main

import (
	"io/ioutil"

	"github.com/pelletier/go-toml"
)

type ColorConfig struct {
	Title          string `toml:"title"`
	NormalText     string `toml:"normal_text"`
	Cursor         string `toml:"cursor"`
	Selected       string `toml:"selected"`
	Border         string `toml:"border"`
	Instruction    string `toml:"instruction"`
	ActiveColumnBg string `toml:"active_column_bg"`
}

type Config struct {
	Colors ColorConfig `toml:"colors"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}