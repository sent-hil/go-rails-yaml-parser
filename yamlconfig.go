package yamlconfig

import (
	"errors"

	yaml "gopkg.in/yaml.v1"
)

var (
	// preset environments
	Development Env = "development"
	Test        Env = "test"
	Staging     Env = "staging"
	Production  Env = "production"

	// name of block that'll be used to look for key, if key is not in env.
	Default = "defaults"

	ErrKeyNotFound = errors.New("Key not defined for env.")
)

type Env string

type Client struct {
	env        Env
	YamlStruct map[string]interface{}
}

// New initializes Client after unmarshalling given data arg into yaml.
func New(data []byte) (*Client, error) {
	client := &Client{YamlStruct: map[string]interface{}{}}

	if err := yaml.Unmarshal([]byte(data), &client.YamlStruct); err != nil {
		return nil, err
	}

	return client, nil
}

// Get gets given key from set env block. If key is not found in env, it looks
// in `defaults` blocks. If key doesn't exist there either, it returns
// ErrKeyNotFound.
func (r *Client) Get(key string) (interface{}, error) {
	if v, ok := r.getKeyFromBlock(string(r.GetEnv()), key); ok {
		return v, nil
	}

	if v, ok := r.getKeyFromBlock(Default, key); ok {
		return v, nil
	}

	return "", ErrKeyNotFound
}

// GetString gets string value from given key from set env block. Use `Get` for interface return value.
func (r *Client) GetString(key string) (string, error) {
	value, err := r.Get(key)
	if err != nil {
		return "", err
	}

	return value.(string), nil
}

func (r *Client) MustGet(key string) interface{} {
	val, err := r.Get(key)
	if err != nil {
		panic(ErrKeyNotFound)
	}

	return val
}

func (r *Client) MustGetString(key string) string {
	return r.MustGet(key).(string)
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
func (r *Client) getKeyFromBlock(block, key string) (interface{}, bool) {
	if v, ok := r.YamlStruct[block]; ok {
		k, ok := v.(map[interface{}]interface{})[key]
		return k, ok
	}

	return nil, false
}
