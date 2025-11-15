<!-- gooselet
Goose modules are affectionately called "Gooselet". They are the building blocks of your Goose application, allowing you to organize your code into reusable components. Each module can contain its own models, controllers, routes, and other components.
[]: # Gooselets can be created using the Goose CLI, and they can be easily shared and reused across different applications. This modular approach allows you to build complex applications by composing smaller, focused modules, making your codebase more maintainable and scalable.
[]: #
[]: # For more information on how to create and use Gooselets, see the [Gooselets documentation](./modules/index.md).



type Module struct {
Imports []Importable
Exports []Exportable
Providers []Injectable
Declarations []Declarable
} -->
<!--
Goose app are a composition of Gooselets (modules). Each Gooselet is a self-contained unit that encapsulates a specific feature or functionality of the application. Gooselets can contain modules, models, controllers, services, and other components that work together to provide a cohesive feature set. -->
