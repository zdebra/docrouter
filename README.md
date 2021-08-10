# docrouter

Golang router for a web server with builtin OpenAPI documentation

## Router vs Server

docrouter is a router not a server has following implications:

- no `ListenAndServe` method - serving exposed http.Handler is not this package responsibility
- should I remove middleware functionality completely? Or add router-level middleware?
