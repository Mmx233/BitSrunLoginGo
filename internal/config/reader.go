package config

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"path"
)

func newReaderFromPath(pathname string) Reader {
	switch path.Ext(pathname) {
	case ".json":
		return Json{}
	case ".yaml":
		return Yaml{}
	default:
		log.Warnf("未知配置类型，使用 yaml 进行解析")
		return Yaml{}
	}
}

type Reader interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

type Json struct {
}

func (Json) Marshal(v any) ([]byte, error) {
	return json.MarshalIndent(v, "", " ")
}
func (Json) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

type Yaml struct {
}

func (Yaml) Marshal(v any) ([]byte, error) {
	return yaml.Marshal(v)
}
func (Yaml) Unmarshal(data []byte, v any) error {
	return yaml.Unmarshal(data, v)
}
