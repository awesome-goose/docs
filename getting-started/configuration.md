# Goose

## Docs

### Getting Started

### Configuration

#### Introduction

Goose is a highly configurable platform. You can set various options to customize its behavior. The configuration can be set in the following ways:

- through environment variables
- through a configuration file
- dynamically in the code

Environment variables are useful for sensitive information such as API keys, database credentials, or when you want to override specific settings without changing the configuration file. Environment variables feeds into the configuration file, so you can use them to override specific settings in the configuration file.

The configuration file is the preferred way to set the configuration, as it allows you to manage all settings in one place.

Dynamic configuration is useful for settings that need to be changed at runtime, such as logging levels or feature flags. However, it is generally recommended to use environment variables or a configuration file for most settings, as they provide a more consistent and manageable way to configure your application.

#### Environment Variables

For apps generated using Goose CLI, the environment variables are set in the `.env` file in the root of your application. You can also set them in your shell or in a `.env` file in the root of your project.

There are several environment variables that you can set to configure Goose. Here is a list of the most commonly used environment variables, these are provided as defaults in the `.env` file:

```bash
# Set the environment variables in the .env file
ENV=development
PORT=8080
LOG_LEVEL=info

```

Goose also provides a way to set and get environment variables in the code. This is useful for settings that need to be changed at runtime.

```go
package main

import (
    "github.com/awesome-goose/framework/env"
)

func main() {
    // Set an environment variable
    env.Set("MY_ENV_VAR", "my_value")

    // Get an environment variable
    value := goose.Get("MY_ENV_VAR")
    fmt.Println("MY_ENV_VAR:", value)
}
```

#### Configuration File

The configuration file is a YAML file that contains all the settings for your application. It is located in the `config` directory of your application. The default configuration file is `config/app.yaml`.

However, you can end up with multiple configuration files depending on the additional goose modules imported by your application. Each module can have its own configuration file, which will be merged with the main configuration file.

Goose is a module first framework, so you can create and import your own modules, each with its own configuration file. This allows you to organize your configuration settings in a modular way.

```yaml
# config/app.yaml
name: my-app
environment: { { ENV } }
port: { { PORT } }
log_level: { { LOG_LEVEL } }
```

The above example shows how easy it is to incorporate environment variables into your configuration file. The `{{ ENV }}` syntax is used to reference environment variables, allowing you to set the values dynamically based on the environment.

#### Dynamic Configuration

You can also set configuration options dynamically in your code. This is useful for settings that need to be changed at runtime, such as logging levels or feature flags.

```go
package main

import (
    "github.com/awesome-goose/framework/config"
)

func main() {
    // Set a configuration option dynamically
    config.Set("app.name", "my-dynamic-app")

    // Get a configuration option
    appName := config.Get("app.name")
    fmt.Println("App Name:", appName)
}
```

#### Configuration Hierarchy

Goose follows a configuration hierarchy where the order of precedence is as follows:

1. **Dynamic Configuration**: Settings set dynamically in the code take precedence over other settings.
2. **Environment Variables**: Environment variables override settings in the configuration file.
3. **Configuration File**: Settings in the configuration file are used as defaults.
   This means that if you set a configuration option dynamically, it will override any settings in the environment variables or configuration file. If you set an environment variable, it will override the corresponding setting in the configuration file. If a setting is not set in any of these, the default value from the configuration file will be used.

#### Conclusion

Goose provides a flexible and powerful configuration system that allows you to customize your application easily. You can set configuration options through environment variables, a configuration file, or dynamically in the code. This flexibility allows you to manage your application's settings in a way that best suits your needs, whether you're working in development, testing, or production environments.
