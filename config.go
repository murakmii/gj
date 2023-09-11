package gj

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"text/template"
)

type Config struct {
	ClassPath []string `json:"class_path"`
}

func ReadConfig(r io.Reader) (*Config, error) {
	confTpl, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	tpl, err := template.New("config").Parse(string(confTpl))
	if err != nil {
		return nil, err
	}

	envs := make(map[string]string)
	for _, pair := range os.Environ() {
		kv := strings.Split(pair, "=")
		envs[kv[0]] = kv[1]
	}

	confJson := bytes.NewBuffer(nil)
	if err := tpl.Execute(confJson, envs); err != nil {
		return nil, err
	}

	conf := &Config{}
	if err := json.NewDecoder(confJson).Decode(conf); err != nil {
		return nil, err
	}

	return conf, nil
}
