package docserver

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gorilla/mux"
)

type Server struct {
	opts      Options
	docRoot   *openapi3.T
	muxRouter *mux.Router
}

func New(opts Options) *Server {
	docRoot := openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:   opts.Title,
			Version: opts.Version,
		},
	}
	for _, serverURL := range opts.ServerURLs {
		docRoot.AddServer(&openapi3.Server{
			URL: serverURL,
		})
	}
	return &Server{
		opts:      opts,
		docRoot:   &docRoot,
		muxRouter: mux.NewRouter(),
	}
}

func (srv *Server) AddRoute(route Route) error {
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

func (srv *Server) addRouteToDoc(route *Route) error {
	for _, method := range route.Methods {
		operation := openapi3.Operation{}
		srv.docRoot.AddOperation(route.Path, method, &operation)
	}
	return nil
}

func (*Server) validateRoute(route *Route) error {
	if route.Handler == nil {
		return fmt.Errorf("handler is nil")
	}
	if route.Path == "" {
		return fmt.Errorf("path is empty")
	}
	return nil
}

func (srv *Server) registerHandler(route *Route) error {
	// for _, serverURL := range srv.opts.ServerURLs {
	srv.muxRouter.
		Handle(route.Path, route.Handler).
		Methods(route.Methods...)
		// Host(serverURL) todo add server
	// }
	return nil
}

// Do I want to expose server or just a router?
// func (srv *Server) ListenAndServe(port string) error {
// 	httpSrv := http.Server{
// 		Handler: srv.muxRouter,
// 		Addr:    port,
// 	}
// 	return httpSrv.ListenAndServe()
// }
