# Goose

<!-- - app lifecycle (=platform)
- - modules (register, boot, shutdown)
- - service container & providers (register, boot, shutdown)
- request lifecycle (=app)
- kernel (http/cli/etc)
- - receive request
- - return response
- - in between
- - - parsing (pre & post hooks)
- - - request (pre & post hooks)
- - - validation (pre & post hooks)
- - - request scope services
- - - routing, middlewares, pre & post hooks
- - - serialisation (pre & post hooks)
- - - response (pre & post hooks)
      nsure


core is required
only one of angular/react/vue/svelte/web/api/cli is required
use package init to auto register platform
also define platfom under config/app.go (config.app.platform)

platform/core
on boot:

- load the user defined platform type
- load the user defined packages
- register services into container for the platform and all defined packages
- register boot hooks for the platform and all defined packages
- register shutdown hooks for the platform and all defined packages
- boot app (also invokes book hooks for packages & services)

on request:

- Parsing, Routing, Middleware, Controller, Serialisation, Response

platform/core interface:

- request() - return an instance of request( + context). How the request is formed depend on the platform, cli would use args and flags, while api, web would use http.
- response(request) - returns an instance of response for the given request
- Handle(request) response - the only public method, takes a raw request, turn it to goose request, propage it thru the router, middleware pipe, controller, serializer and finally return a response
  - request represents the raw input from the corresponding platform, this is input to request()
  - response response a platform compatible response, returned by response()

      -->

<!-- - platform
- - any of web, api, cli, angular, react, vue, (flutter, mobile, desktop)
- app
- - request & response handling machine
- request
- - raw input from client, abstracted into context for userland
- response
- - raw output to client, abstracted into reply for userland
- context
- - userland container for request information/actions
- - clean, highlevel request
- reply
- - userland container for response information/actions
- - clean, highlevel request -->

<!--
- platform
- - init
- app
- - run
- router
- hook
- middleware
- request
- response
- middleware
- context
- validator -->

<!--
// - make provision for handling request
// - handle request
// - - receive request
// - - return response
// - - in between
// - - - parsing (pre & post hooks)
// - - - request (pre & post hooks)
// - - - validation (pre & post hooks)
// - - - request scope services
// - - - routing, middlewares, pre & post hooks
// - - - serialisation (pre & post hooks)
// - - - response (pre & post hooks)
// - ensure

// - load env
// - load config, merge env if necessary
// - load module
// - register services
// - boot services
// - create & return app instance for the platform

// system/platform hooks
// register default services
// override with custom services
// etc -->

## Lifecycle

Two separate processes work together to ensure goose perform as expected: the application lifecycle and the request lifecycle. The application lifecycle is responsible for the overall state of the application, while the request lifecycle is responsible for handling individual requests.

### App Lifecycle

The application lifecycle is the process that governs the initialization, execution, and termination of the application. It consists of three main phases: **Register**, **Boot**, and **Shutdown**. Each phase plays a crucial role in ensuring that the application is set up correctly, runs efficiently, and cleans up resources properly.

1. **Register**: This phase is responsible for setting up the application, including registering modules, services, and configurations.
   During this phase: - Modules and services are registered with the application. - The service container is populated with bindings and singletons. - Configuration files and environment variables are loaded.

2. **Boot**: In this phase, the application is fully initialized and ready to handle requests. This includes booting modules, initializing services, and setting up event listeners.

   - Modules and services are initialized and made ready for use.
   - Event listeners, middleware, and routes are registered.
   - Any deferred services are resolved and bootstrapped.

3. **Shutdown**: The final phase where the application cleans up resources, closes connections, and performs any necessary shutdown tasks.
   During this phase: - Resources such as database connections, file handles, and caches are released. - Any cleanup tasks are performed. - Logs and final metrics are flushed.

### Request Lifecycle

The request lifecycle is the process that handles incoming requests and generates responses. It consists of several stages, including **Parsing**, **Routing**, **Middleware**, **Controller**, **Serialisation** **Response**. Each stage plays a crucial role in processing the request and generating the appropriate response.

1. **Parsing**: This stage is responsible for parsing the incoming request, extracting relevant information, and preparing it for further processing.

2. **Routing**: In this stage, the application determines which route or endpoint should handle the request based on the request's URL and HTTP method. The router matches the request to a specific controller and action.

3. **Middleware**: Middleware is a series of functions that are executed in sequence before the request reaches the controller. Middleware can perform tasks such as authentication, logging, and modifying the request or response.

4. **Controller**: The controller is responsible for handling the request and generating a response. It contains the business logic and interacts with models, services, and other components to process the request.

5. **Serialisation**: In this stage, the response data is serialized into a format suitable for transmission over the network, such as JSON, HTML or XML. This ensures that the response can be easily consumed by the client.

6. **Response**: The final stage where the application sends the generated response back to the client. This includes setting the appropriate HTTP status code, headers, and body content.

### Summary

The application lifecycle and request lifecycle work together to ensure that the application is initialized, handles requests efficiently, and cleans up resources properly. By understanding these lifecycles, you can create robust and maintainable applications that perform well under various conditions.
