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
	router := New(DefaultOptions)

	type MyParameters struct {
		StarID int  `docrouter:"name:starId;desc:Star identifier in CommonMark syntax. This can be potentially issue for longer descriptions.; example: 5; required: false; schemaMin: 3"`
		Potato bool `docrouter:"name:potato;desc: This is bool!; example: true; required: true"`
	}

	const expectedHandlerOutput = "Hello star!"

	err := router.AddRoute(Route{
		Path:       "/stars",
		Methods:    []string{http.MethodGet},
		Parameters: &MyParameters{},
		Summary:    "Get All Stars",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, expectedHandlerOutput)
		}),
	})
	require.NoError(t, err)

	router.muxRouter.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		methods, _ := route.GetMethods()
		path, _ := route.GetPathTemplate()
		fmt.Println("registered route", methods, path)
		return nil
	})

	ts := httptest.NewServer(router.muxRouter)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/stars")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	respBytes, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	assert.Equal(t, expectedHandlerOutput, string(respBytes))

	doc, err := router.docRoot.MarshalJSON()
	require.NoError(t, err)
	fmt.Println(string(doc))

	notFoundresp, err := http.Get(ts.URL + "/knock-knock")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, notFoundresp.StatusCode)
}
