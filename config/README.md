
# Config
Flex Config is a lightweight yet powerful configuration manager for Go projects.
It takes advantage of Dot-env (.env) files and OS environment variables alongside config files (JSON, YAML) to meet all of your requirements.

## Documentation
### Required Go Version
It requires Go `v1.16` or newer versions.

### Installation
To install this package run the following command in the root of your project.

```bash
go get github.com/20326/flexbox/config
```

### Quick Start
The following example demonstrates how to use a JSON configuration file.

```go
// The configuration struct
type MyConfig struct {
    Addr       string
    Debug      bool
    Rand       float64
}

// Create an instance of the configuration struct
myConfig := MyConfig{}

// Create a Config instance and load it from a yaml file
errs := config.New("etc").LoadFile("app.yaml", &myConfig).End()

if len(errs) > 0 {
	log.Fatalf("errors: %v", errs)
}

// Use `myConfig`...
```

### Types

* `Json`: It loads using a JSON file.
* `Yaml`: It feeds using a YAML file.
* `DotEnv`: It feeds using a dot env (.env) file.

The `.env` file:
```env
ENV=production
````

### Load Priority
the priority of loading configuration data:

filename.${env}.ext (if exists) > filename.ext

PS: Env: default `local`

#### DotEnv
The `DotEnv` package to load `.env` files.
The example below shows how to use the `DotEnv` feeder.

The `.env` file: 

## License
Flex Config is released under the [MIT License](http://opensource.org/licenses/mit-license.php).