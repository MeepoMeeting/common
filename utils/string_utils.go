package utils

import (
	"encoding/json"
	"io/ioutil"
)

func WriteToJSON(val interface{}) string {
	str, err := json.Marshal(val)
	if err != nil {
		return ""
	}
	return string(str)
}

func ReadFromFile(path string) (string, error) {
	b, err := ioutil.ReadFile(path) // just pass the file name
	if err != nil {
		return "", err
	}
	return string(b), nil
}
