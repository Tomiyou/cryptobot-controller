package jsonLoader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// This function loads the JSON at the provided path into any struct
func LoadJSON(path string, input_type interface{}) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}

	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	err = json.Unmarshal(byteValue, input_type)
	if err != nil {
		return
	}

	fmt.Println("Successfully read", path)
	return
}

// This function loads the YAML at the provided path into any struct
func LoadYAML(path string, input_type interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)
	err = yaml.Unmarshal(byteValue, input_type)
	if err != nil {
		return err
	}

	fmt.Println("Successfully read", path)
	return nil
}
