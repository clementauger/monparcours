package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

func ReadConfig(dst interface{}, filename string) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, dst)
	if err != nil {
		return err
	}

	return nil
}
