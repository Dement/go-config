package goconfig

import (
	"github.com/tsuru/config"
	"os"
)

func LoadFile(path string) {
	pwd := getCurrentDirectory()

	errFile := config.ReadConfigFile(pwd + path)
	if errFile != nil {
		panic(errFile)
	}
}

func getCurrentDirectory() string {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return pwd
}