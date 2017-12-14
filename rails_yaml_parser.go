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

	Default = "Defaults"

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

func (r *Client) getKeyFromBlock(block, key string) reflect.Value {
	// equivalent to r.<block>.<key>
	f := reflect.Indirect(reflect.ValueOf(r)).FieldByName(block)
	k := f.FieldByName(key)

	return k
}

func (r *Client) Get(key string) (string, error) {
	if fromEnv := r.getKeyFromBlock(string(r.GetEnv()), key); fromEnv.IsValid() {
		return fromEnv.String(), nil
	}

	if fromDefault := r.getKeyFromBlock(Default, key); fromDefault.IsValid() {
		return fromDefault.String(), nil
	}

	return "", ErrKeyNotFound
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
