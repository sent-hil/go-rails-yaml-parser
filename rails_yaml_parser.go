package railsyamlparser

import (
	"errors"
	"reflect"

	yaml "gopkg.in/yaml.v1"
)

var (
	Development Env = "Development"
	Test        Env = "Test"
	Staging     Env = "Staging"
	Production  Env = "Production"

	ErrKeyNotFound = errors.New("Key not defined for env.")
)

type Env string

type Client struct {
	env      Env
	Defaults struct {
		Adapter  string `yaml:"adapter"`
		Encoding string `yaml:"encoding"`
		Pool     int    `yaml:"pool"`
		Port     uint16 `yaml:"port"`
		Host     string `yaml:"host"`
	} `yaml:"defaults"`
	Test struct {
		Database string `yaml:"database"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"test"`
	Development struct {
		Database string `yaml:"database"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"development"`
}

func New(data []byte) (*Client, error) {
	var r Client
	if err := yaml.Unmarshal([]byte(data), &r); err != nil {
		return nil, err
	}

	return &r, nil
}

func (r *Client) Get(key string) (string, error) {
	// equivalent to r.<env>.<key>
	f := reflect.Indirect(reflect.ValueOf(r)).FieldByName(string(r.GetEnv()))
	k := f.FieldByName(key)

	if !k.IsValid() {
		return "", ErrKeyNotFound
	}

	return k.String(), nil
}

func (r *Client) SetEnv(env Env) {
	r.env = env
}

func (r *Client) GetEnv() Env {
	if r.env == "" {
		return Development
	}
	return r.env
}
