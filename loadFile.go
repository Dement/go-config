package main

import (
	"github.com/tsuru/config"
)

func LoadFile(path string) {
	pwd := getCurrentDirectory()

	errFile := config.ReadConfigFile(pwd + path)
	if errFile != nil {
		panic(errFile)
	}
}