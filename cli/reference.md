# Goose

## Docs

### CLI

### Reference

The Goose CLI provides a set of commands to help you manage your Goose applications. Below is a reference of the available commands and their options.

#### App Commands

- `app --name=<app-name> init`: Initializes a new Goose CLI application with the specified name.
- `app --name=<app-name> run [--watch]`: Runs the specified Goose CLI application. Use the `--watch` flag to enable live reloading.
- `app --name=<app-name> deploy`: Deploys the specified Goose CLI application.

#### Module Commands

- `module --name=<module-name> create`: Creates a new module within the Goose CLI application.
- `module --name=<module-name> model create`: Creates a new model within the specified module.
- `module --name=<module-name> service create`: Creates a new service within the specified module.
- `module --name=<module-name> controller create`: Creates a new controller within the specified module.
- `module --name=<module-name> route create`: Creates a new route within the specified module.

#### Migration Commands

- `migration --name=<migration-name> create [for=model-name]`: Creates a new migration file. Use the `for` option to associate the migration with a specific model.
- `migration --name=<migration-name> up`: Applies the specified migration.
- `migration --name=<migration-name> down`: Rolls back the specified migration.

#### Seed Commands

- `seed --name=<seed-name> create`: Creates a new seed file within the specified module.
- `seed --name=<seed-name> up`: Applies the specified seed.
- `seed --name=<seed-name> down`: Rolls back the specified seed.

#### Test Commands

- `test run`: Runs the tests for the Goose CLI application.
