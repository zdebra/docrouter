package docrouter

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocServer(t *testing.T) {
	server := New(DefaultOptions)

	type myRequest struct {
		ID string `json:"id"`
		// <here goes meta information to setup request validation, e.g. regular expression>
		FavoriteColor string `json:"favoriteColor"`
	}
	type myResponse struct {
		PotatoCount int `json:"potatoCount"`
	}
	type myQueryParams struct {
		Limit int `json:"limit"`
		// <validation info goes here as well>
		Offset int `json:"offset"`
	}
	type myHeaderParams struct {
		Authorization string
	}
	type myPathParams struct {
		x, y, z string
	}

	const expectedHandlerOutput = "Hello star!"

	err := server.AddRoute(Route{
		Path:    "/stars",
		Methods: []string{http.MethodGet},
		// RequestBody:  myRequest{},
		// ResponseBody: myResponse{},
		// QueryParams:  myQueryParams{},
		// HeaderParams: myHeaderParams{},
		// PathParams:   myPathParams{},
		Summary: "Get All Stars",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, expectedHandlerOutput)
		}),
	})
	require.NoError(t, err)

	server.muxRouter.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		methods, _ := route.GetMethods()
		path, _ := route.GetPathTemplate()
		fmt.Println("registered route", methods, path)
		return nil
	})

	ts := httptest.NewServer(server.muxRouter)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/stars")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	respBytes, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	assert.Equal(t, expectedHandlerOutput, string(respBytes))

	doc, err := server.docRoot.MarshalJSON()
	require.NoError(t, err)
	fmt.Println(string(doc))
}
