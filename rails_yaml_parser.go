package railsyamlparser

import (
	"errors"
	"reflect"

	yaml "gopkg.in/yaml.v1"
)

var (
	// preset environments
	Development Env = "Development"
	Test        Env = "Test"
	Staging     Env = "Staging"
	Production  Env = "Production"

	// name of block that'll be used to look for key, if key is not in env.
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

// New initializes Client after unmarshalling given data arg into yaml.
func New(data []byte) (*Client, error) {
	var r Client
	if err := yaml.Unmarshal([]byte(data), &r); err != nil {
		return nil, err
	}

	return &r, nil
}

// Get gets given key from set env block. If key is not found in env, it looks
// in `defaults` blocks. If key doesn't exist there either, it returns
// ErrKeyNotFound.
func (r *Client) Get(key string) (string, error) {
	if fromEnv := r.getKeyFromBlock(string(r.GetEnv()), key); fromEnv.IsValid() {
		return fromEnv.String(), nil
	}

	if fromDefault := r.getKeyFromBlock(Default, key); fromDefault.IsValid() {
		return fromDefault.String(), nil
	}

	return "", ErrKeyNotFound
}

// SetEnv sets given env. This should be used to set env based on env vars or
// user specified env.
func (r *Client) SetEnv(env Env) {
	r.env = env
}

// GetEnv returns current env. If not set, it returns Development.
func (r *Client) GetEnv() Env {
	if r.env == "" {
		return Development
	}
	return r.env
}

// getKeyFromBlock attempts to get given key from give block. Equivalent to
// `r.<block>.<key>`.
func (r *Client) getKeyFromBlock(block, key string) reflect.Value {
	f := reflect.Indirect(reflect.ValueOf(r)).FieldByName(block)
	k := f.FieldByName(key)

	return k
}
