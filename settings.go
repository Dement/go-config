package main

import (
	"github.com/tsuru/config"
)

func LoadConfig() {
	LoadFile("/config/parameters.yml")
}

func GetParameters(key string) (interface {})  {
	val, errConf := config.Get(key)

	if errConf != nil {
		panic(errConf)
	}

	return val
}

func GetListParameters(key string) ([]string) {
	val, errConf := config.GetList(key)

	if errConf != nil {
		panic(errConf)
	}

	return val
}

func FindListParameters(key, value string) (bool) {
	maps := GetListParameters(key)

	for _, val := range maps {
		if val == value {
			return true
		}
	}

	return false
}