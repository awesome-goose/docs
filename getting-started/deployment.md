<!--
    Packing
    Cross-compiling
    Systemd Service
    Using a Proxy
    Cloud Providers -->

# Goose

## Docs

### Getting Started

### Deployment

#### Introduction

Goose can be deployed in various ways to suit your needs. Here are some common deployment methods:

1. **CLI**: Deploy Goose applications using the command line interface.
2. **Docker**: Use Docker to containerize and deploy Goose applications.
3. **Goose Cloud**: Deploy your applications to Goose Cloud, a managed service for Goose applications.

#### Run Your App Directly

To run your Goose application on your local machine, private server, use the following command:

```bash
$ goose run
```

This command will package your application and prepare it for deployment. You can then run the packaged application on any server that has Goose installed.

#### Docker Deployment

To deploy your Goose application using Docker, a dockerfile is provided by default. You can build and run your application in a Docker container with the following commands:

```bash
$ docker build -t my-app .
$ docker run -d -p 8080:8080 my-app
```

This will build your application into a Docker image and run it in a container, exposing port 8080.

#### Goose Cloud Deployment

To deploy your Goose application to Goose Cloud, you can use the following command:

```bash
$ goose deploy --platform=cloud --name=my-app
```
