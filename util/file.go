package util

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
	root, err := a.GetRootPath()
	if err != nil {
		return false
	}
	_, err = os.Stat(root + path)
	return err == nil || os.IsExist(err)
}

func (a *file) Read(path string) ([]byte, error) {
	root, err := a.GetRootPath()
	if err != nil {
		return nil, err
	}
	return ioutil.ReadFile(root + path)
}

func (a *file) ReadJson(path string, receiver interface{}) error {
	data, err := a.Read(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, receiver)
}

func (a *file) Write(path string, data []byte) error {
	root, err := a.GetRootPath()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(root+path, data, 700)
}

func (a *file) WriteJson(path string, receiver interface{}) error {
	data, err := json.MarshalIndent(receiver, "", " ")
	if err != nil {
		return err
	}
	return a.Write(path, data)
}

func (*file) GetRootPath() (string, error) {
	t, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(t) + "/", nil
}

func (a *file) Add(path string, c string) error {
	root, err := a.GetRootPath()
	if err != nil {
		return err
	}
	file, err := os.OpenFile(root+path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 700)
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
