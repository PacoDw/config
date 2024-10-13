# Config: A Versatile Configuration Library for Go

`Config` is a Go library designed to handle configuration from environment variables, files, and defaults in a flexible and intuitive way. It supports setting default values, loading configurations from various file types (such as `.yaml`, `.json`, etc.), and reading from environment variables, making it ideal for applications that require a versatile configuration management system.

## Features
- **Default values**: Set default values that can be overridden by environment variables or file-based configurations.
- **Environment variables**: Automatically loads values from environment variables and overrides file-based or default configurations.
- **File-based configurations**: Supports loading configuration from files with different formats (`yaml`, `json`, `toml`, etc.). You can check [Viper's documentation](https://github.com/spf13/viper) for a full list of supported formats.
- **Global and internal variable precedence**: Global environment variables can be reused across multiple structs, but internal struct variables will always take precedence over global ones.

## Installation

```bash
go get github.com/PacoDw/config
```

## Usage Example

Here is an example to showcase the flexibility of the `config` library:

### Sample `.env` file
```yaml
# global environment variables
name: "MyApp1.2"
app_environment: "local"
environment: "testing3"
server_port: 3000

# server environment variables
server:
    SERVER_NAME: "MyApp"
    host: "127.0.0.1"
    environment: "dev"

# postgres environment variables
postgres:
    user: "admin2"
    host: "localhost"
    DB_PORT: 5433

# server environment variables specific to serverconfig
serverconfig:
    host: "localhost2"
    server_port: 9000
    environment: "dev3"
```

### Go Code (`main.go`)

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/PacoDw/config/config"
)

type AppConfig struct {
	// Global environment variables that can be reused in other structs
	Name        string `env:"APP_NAME" validate:"required"`
	Environment string `env:"ENVIRONMENT"`
	Port        int    `env:"PORT" default:"8080"`

	Server   ServerConfig
	Database DatabaseConfig `env:"postgres"`
}

type DatabaseConfig struct {
	User     string `env:"USER"`
	Password string `env:"DB_PASSWORD" default:"my_default_password"`
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT" default:"3052"`
}

type ServerConfig struct {
	Name        string `env:"SERVER_NAME" default:"my_default_server_name"`
	Host        string `env:"HOST" validate:"required"`
	Port        int    `env:"SERVER_PORT"`
	APIKey      string `env:"api_key" default:"my_default_api_key"`
	Environment string `env:"ENVIRONMENT" default:"user_env"`
}

func main() {
	// Initialize config with defaults and load from .env.yaml
	cfg := config.New()

	var appConfig AppConfig
	if err := cfg.Unmarshal(&appConfig); err != nil {
		log.Fatalf("Error loading AppConfig: %v", err)
	}

	var serverConfig ServerConfig
	if err := cfg.Unmarshal(&serverConfig); err != nil {
		log.Fatalf("Error loading ServerConfig: %v", err)
	}

	// Output the app configuration in JSON
	appConfigJSON, _ := json.MarshalIndent(appConfig, "", "  ")
	fmt.Printf("App Config (JSON): %s\n", appConfigJSON)

	// Output the server configuration in JSON
	serverConfigJSON, _ := json.MarshalIndent(serverConfig, "", "  ")
	fmt.Printf("Server Config (JSON): %s\n", serverConfigJSON)
}
```

### Output

```
➜ go run main.go
App Config (JSON): {
  "Name": "MyApp1.2",
  "Environment": "testing3",
  "Port": 8080,
  "Server": {
    "Name": "MyApp",
    "Host": "127.0.0.1",
    "Port": 3000,
    "APIKey": "my_default_api_key",
    "Environment": "dev"
  },
  "Database": {
    "User": "admin2",
    "Password": "my_default_password",
    "Host": "",
    "Port": 5433
  }
}
Server Config (JSON): {
  "Name": "my_default_server_name",
  "Host": "localhost2",
  "Port": 9000,
  "APIKey": "my_default_api_key",
  "Environment": "dev3"
}
```

### Explanation of Behavior
- **Default values**: If a value is not set in either the environment or the file, the default value specified in the struct tag is used. For example, the field `Server.APIKey` has a default of `"my_default_api_key"` since it is not present in the `.env` file.
  
- **Global vs internal precedence**: Global variables, like `Environment` or `Port`, will be reused across different structs (e.g., `AppConfig`, `ServerConfig`, etc.). However, if an internal struct like `ServerConfig` defines a value for the same variable (e.g., `ServerConfig.Environment`), this value will take precedence over the global one.

  In this example:
  - The global `environment` is set to `"testing3"`, but the `ServerConfig.Environment` is set to `"dev"`, so `"dev"` is used for the server.
  - `Database.Password` defaults to `"my_default_password"` because no value is provided for it in the `.env` file.

### Configuration Options
You can customize the path, file name, and file type by passing options when initializing the configuration:

```go
cfg := config.New(
    config.WithFilePath("/custom/path/"),
    config.WithFileName("custom_config"),
    config.WithFileType("json"),
)
```

This will load a configuration file located at `/custom/path/custom_config.json`.

### Validation
You can use the full suite of validation tags from the `go-playground/validator` package. The example above uses the `required` validation, but you can use many other tags like `min`, `max`, `email`, etc. Check the [official go-playground/validator documentation](https://github.com/go-playground/validator) for more examples.

### File Format Support
`Config` leverages **Viper** under the hood, which supports a wide variety of configuration file formats including `json`, `yaml`, `toml`, and more. You can refer to [Viper's documentation](https://github.com/spf13/viper) for a full list of supported formats.

## License
This project is licensed under the MIT License.