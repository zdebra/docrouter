package docrouter

import (
	"fmt"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

type Router struct {
	opts      Options
	docRoot   *openapi3.T
	muxRouter *mux.Router
}

func New(opts Options) *Router {
	docRoot := openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:   opts.Title,
			Version: opts.Version,
		},
	}
	for _, server := range opts.Servers {
		docRoot.AddServer(&openapi3.Server{
			URL:         server.URL,
			Description: server.Description,
		})
	}
	return &Router{
		opts:      opts,
		docRoot:   &docRoot,
		muxRouter: mux.NewRouter(),
	}
}

func (srv *Router) AddRoute(route Route) error {
	if err := srv.validateRoute(&route); err != nil {
		return fmt.Errorf("route validation: %v", err)
	}
	if err := srv.addRouteToDoc(&route); err != nil {
		return fmt.Errorf("adding route do doc: %v", err)
	}
	if err := srv.registerHandler(&route); err != nil {
		return fmt.Errorf("register handler: %v", err)
	}

	return nil
}

func (srv *Router) addRouteToDoc(route *Route) error {
	params, err := route.openAPI3Params()
	if err != nil {
		return fmt.Errorf("create route params: %w", err)
	}
	for _, method := range route.Methods {
		operation := openapi3.Operation{
			Summary:     route.Summary,
			Description: route.Description,
			OperationID: uniqueOperationID(route),
			Parameters:  params,
			Responses:   openapi3.NewResponses(),
		}
		srv.docRoot.AddOperation(route.Path, method, &operation)
	}
	return nil
}

func uniqueOperationID(route *Route) string {
	// todo: ensure uniqueness
	return strings.ToLower(strings.ReplaceAll(route.Summary, " ", "-"))
}

func (*Router) validateRoute(route *Route) error {
	return validation.ValidateStruct(route,
		validation.Field(&route.Handler, validation.NotNil),
		validation.Field(&route.Path, validation.Required),
		validation.Field(&route.Summary, validation.Required),
	)
}

func (srv *Router) registerHandler(route *Route) error {
	middlewareChain := alice.New()
	for _, mw := range route.Middlewares {
		middlewareChain = middlewareChain.Append(mw)
	}
	h := middlewareChain.Then(route.Handler)
	srv.muxRouter.
		Handle(route.Path, h).
		Methods(route.Methods...)
	return nil
}
