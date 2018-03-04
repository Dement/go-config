package goconfig

import (
	"fmt"
	"os"
	"bufio"
	"io/ioutil"
	"strings"
	"github.com/tsuru/config"
	"gopkg.in/yaml.v2"
)

var result map[interface{}]interface{}

func main() {
	checkFile("/config")

	LoadFile("/config/parameters.yml.dist")

	parameters, err := config.Get("parameters")
	if err != nil {
		panic(err)
	}

	if parameters == nil {
		os.Exit(0)
	}

	md, ok := parameters.(map[interface{}]interface{})
	if ok != true {
		panic(ok)
	}

	loadFile("/config/parameters.yml")

	getParams(md, "")

	ymlTempl, err := yaml.Marshal(&result)
	if err != nil {
		panic(err)
	}

	recordFile(ymlTempl, getCurrentDirectory() + "/config/parameters.yml", 0664)
}

func getParams(maps map[interface{}]interface{}, key string) {
	var prefix, example string
	var typeVal string

	for k, v := range maps {
		if len(key) == 0 {
			prefix = k.(string)
		} else {
			prefix = key + ":" + k.(string)
		}

		switch t := v.(type) {
		case string, []interface{}:
			value, error := getParameter(prefix)

			if error == true || value == nil {
				reader := bufio.NewReader(os.Stdin)

				switch tp := v.(type) {
				case string:
					example = v.(string)
					typeVal = "string"
				case []interface{}:
					var count int = 0
					example = "["
					for _, valEx := range tp {
						if count > 0 {
							example = example + " "
						}

						example = example + valEx.(string)
						count++
					}
					example = example + "]"
					typeVal = "maps"
				default:
					panic("unsupported type")
				}

				fmt.Print(prefix + " (" + example + ") : ")
				text, _ := reader.ReadString('\n')

				if text == "\n" {
					value = v
				} else {
					if typeVal == "string" {
						value = strings.TrimSpace(text)
					} else if typeVal == "maps" {
						var tempVal string

						tempVal = strings.Replace(text, "[", "", 10)
						tempVal = strings.Replace(tempVal, "]", "", 10)
						tempVal = strings.Replace(tempVal, "\n", "", 10)

						value = strings.Split(tempVal, " ")
					}
				}
			}

			prefixKeys := strings.Split(prefix, ":")

			last := map[interface{}]interface{}{
				prefixKeys[len(prefixKeys)-1]: value,
			}
			for i := len(prefixKeys) - 2; i >= 0; i-- {
				last = map[interface{}]interface{}{
					prefixKeys[i]: last,
				}
			}

			result = (mergeMaps(result, last))

		case map[interface{}]interface{}:
			getParams(t, prefix)

		default:
			panic("unsupported type")
		}
	}
}

func loadFile(path string) {
	pwd := getCurrentDirectory()

	errFile := config.ReadConfigFile(pwd + path)
	if errFile != nil {
		panic(errFile)
	}
}

func getParameter(key string) (interface {}, bool)  {
	var error bool
	val, errConf := config.Get(key)

	if errConf != nil {
		error = true
	} else {
		error = false
	}

	return val, error
}

func mergeMaps(map1, map2 map[interface{}]interface{}) map[interface{}]interface{} {
	result := make(map[interface{}]interface{})
	for k, v2 := range map2 {
		if v1, ok := map1[k]; !ok {
			result[k] = v2
		} else {
			map1Inner, ok1 := v1.(map[interface{}]interface{})
			map2Inner, ok2 := v2.(map[interface{}]interface{})
			if ok1 && ok2 {
				result[k] = mergeMaps(map1Inner, map2Inner)
			} else {
				result[k] = v2
			}
		}
	}
	for k, v := range map1 {
		if v2, ok := map2[k]; !ok {
			result[k] = v
		} else {
			map1Inner, ok1 := v.(map[interface{}]interface{})
			map2Inner, ok2 := v2.(map[interface{}]interface{})
			if ok1 && ok2 {
				result[k] = mergeMaps(map1Inner, map2Inner)
			}
		}
	}
	return result
}

func checkFile(path string)  {
	pwd := getCurrentDirectory()

	if dirOrFileIsExists(pwd + path) == false {
		os.Mkdir(pwd + path, 0755)
	}

	if dirOrFileIsExists(pwd + path + "/parameters.yml.dist") == false {
		saveFile(generatorDefaultTemplate(), pwd + path + "/parameters.yml.dist", 0664)
	}

	if dirOrFileIsExists(pwd + path + "/parameters.yml") == false {
		saveFile("", pwd + path + "/parameters.yml", 0664)
	}

	if dirOrFileIsExists(pwd + path + "/.gitignore") == false {
		saveFile(generatorGitignoreTemplate(), pwd + path + "/.gitignore", 0664)
	}
}

func dirOrFileIsExists(f string) bool {
	if _, err := os.Stat(f); os.IsNotExist(err) {
		return false
	}

	return true
}

func getCurrentDirectory() string {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return pwd
}

func generatorDefaultTemplate() string  {
	var temp string = "parameters:"
	return temp
}

func generatorGitignoreTemplate() string  {
	var temp string = `*
!parameters.yml.dist`
	return temp
}

func saveFile(template, path string, perm os.FileMode)  {
	d := []byte(template)
	recordFile(d, path, perm)
}

func recordFile(d []byte, path string, perm os.FileMode)  {
	err := ioutil.WriteFile(path, d, perm)

	if err != nil {
		panic(err)
	}
}
