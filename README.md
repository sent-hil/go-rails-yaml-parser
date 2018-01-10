# go-yamlconfig

This is a small library for parsing yaml config files. The main advantage of
this over plain yaml parsing is this library falls back on `defaults` block when
getting key value, ie the way Rails defines its yaml config files.

## Usage

Initialize client:

```go
var exampleConfig = `
defaults: &defaults
 env: development

development:
 <<: *defaults // inherits from defaults

production:
 <<: *defaults   // inherits from defaults
 env: production // overrides defaults for this specific key
`

config, err := New([]byte(exampleConfig))
if err != nil {
  return err
}
```

Get default environment (development) config:

```go
val, err := config.GetString("database")
if err != nil {
  return err
}

fmt.Println(val) // development
```

Get config from another environment:

```go
// override environment
config.SetEnv(Production)
val, err = config.GetString("database")
if err != nil {
  return err
}

fmt.Println(val) // production
```

If config is not string, `Get` can be used to return an `interface{}` value:

```go
val1, err := config.Get("database")
if err != nil {
  return err
}

fmt.Println(val) // production
```

## Install

`go get github.com/sent-hil/go-yamlconfig`

## Test

`go test`
