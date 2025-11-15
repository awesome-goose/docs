# Goose

## Docs

### CLI

### Usage

Once you have installed Goose CLI, you can start using it to manage your Goose applications. Here are some common commands:

```bash
# Create a new Goose CLI application
goose create app --name=my-app

# Run the Goose CLI application
goose run --app=my-app

# Deploy the Goose CLI application
goose deploy --app=my-app

# Create a new module in the Goose CLI application
goose module create --name=my-module --app=my-app

# Create a new migration for a model in the module
goose migration create --name=create_users_table --module=my-module --app=my-app

# Create a new seed file in the module
goose seed create --name=seed_users --module=my-module --app=my-app
```

See the [Reference](./reference.md) section for a complete list of commands and options available in Goose CLI.
