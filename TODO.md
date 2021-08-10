# docrouter

Golang router with OpenAPI specification built-in.

### Project Bootstrap

- [x] create TODO.md
- [ ] create changelog https://keepachangelog.com/en/1.0.0/
- [ ] first release

### Create Router instance

- [ ] Support Source of truth is code use-case
  - [x] router is created with handlers along with required metadata for OpenAPI specification
  - [ ] Input Params
    - [x] generate doc for QueryParams with code reflection
    - [x] generate doc for PathParams with code reflection
    - [x] generate doc for CookieParams with code reflection
    - [x] generate doc for HeadersParams with code reflection
    - [ ] support all OpenAPI types
    - [ ] validate parameters in AddRoute func
  - [ ] generate doc for Request with code reflection
  - [ ] generate doc for Response(s) with code reflection
  - [ ] all OpenAPI types are supported
  - [ ] all OpenAPI schema validations are supported
  - [ ] route constructor
  - [x] middlewares support
- [ ] Router tooling
  - [ ] Decode runtime helpers
    - [x] DecodeQueryParams runtime helper
    - [x] DecodePathParams runtime helper
    - [x] DecodeHeadersParams runtime helper
    - [x] DecodeCookiesParams runtime helper
    - [ ] all OpenAPI types are supported
  - [ ] optional runtime validation for requests based on OpenAPI schema
  - [ ] Route tags are available in runtime with a helper method
- [ ] Try to support more routers than gorilla (chi, echo)
