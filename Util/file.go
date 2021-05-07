package Util

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type file struct{}

var File file

func (a *file) Exists(path string) bool {
	_, err := os.Stat(a.GetRootPath() + path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func (a *file) Read(path string) ([]byte, error) {
	return ioutil.ReadFile(a.GetRootPath() + path)
}

func (a *file) ReadJson(path string, receiver interface{}) error {
	data, err := a.Read(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, receiver)
}

func (a *file) Write(path string, data []byte) error {
	return ioutil.WriteFile(a.GetRootPath()+path, data, 700)
}

func (a *file) WriteJson(path string, receiver interface{}) error {
	data, err := json.MarshalIndent(receiver, "", " ")
	if err != nil {
		return err
	}
	return a.Write(path, data)
}

func (*file) GetRootPath() string {
	t, err := os.Executable()
	if err != nil {
		ErrHandler(err)
	}
	return filepath.Dir(t) + "/"
}

func (a *file) Add(path string, c string) error {
	file, err := os.OpenFile(a.GetRootPath()+path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 700)
	defer file.Close()
	if err != nil {
		return err
	}
	w := bufio.NewWriter(file)
	if _, err = w.WriteString(c + "\n"); err != nil {
		return err
	}
	return w.Flush()
}
