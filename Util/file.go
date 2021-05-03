package Util

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type file struct {
}

var File file

func (*file) Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func (*file) Read(path string, receiver interface{}) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, receiver)
}

func (*file) Write(path string, receiver interface{}) error {
	data, err := json.MarshalIndent(receiver, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 700)
}

func (*file) GetRootPath() string {
	t, err := os.Executable()
	if err != nil {
		ErrHandler(err)
	}
	return filepath.Dir(t)
}
