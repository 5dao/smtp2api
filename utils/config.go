package utils

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

//LoadConfig toml
func LoadConfig(configFileName string, cfg interface{}) (err error) {
	configFileName, err = filepath.Abs(configFileName)
	if err != nil {
		return
	}

	var configFile *os.File
	configFile, err = os.Open(configFileName)
	if err != nil {
		return
	}
	defer configFile.Close()

	if _, err = toml.DecodeReader(configFile, cfg); err != nil {
		return
	}
	return
}
