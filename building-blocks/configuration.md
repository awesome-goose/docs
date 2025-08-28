<!-- config

- dotted access
- config.string
- config.bool
- config.\*\*\*


// Features
// Loads all .yaml/.yml files from a directory
// Each file becomes a top-level namespace (e.g. db.yaml → config["db"])
// Supports get/set via dotted path: db.host, auth.jwt.secret
// Expands {{ENV_VAR}} using environment values on get
// Supports ToStruct(namespace, &targetStruct) -->
<!--

Env, COnfig key and values must all be a string
string util package provide String() struct with ToBool, ToNumber etc for m=conversion  -->
