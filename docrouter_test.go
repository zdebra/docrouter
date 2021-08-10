package docrouter

import (
	"context"
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

func TestMiddlewares(t *testing.T) {
	const ctxKey = "ctxkey"
	const ctxVal = "ctxval"
	const ctxVal2 = "ctxval2"

	router := New(DefaultOptions)
	myHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := r.Context().Value(ctxKey)
		require.NotNil(t, v)
		vs, ok := v.(string)
		require.True(t, ok)
		fmt.Fprintf(w, "%s:%s", ctxKey, vs)
	})

	setCtxVal := func(k, v string) func(next http.Handler) http.Handler {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), k, v)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
			})
		}
	}

	err := router.AddRoute(Route{
		Path:    "/",
		Methods: []string{http.MethodGet},
		Middlewares: []func(http.Handler) http.Handler{
			setCtxVal(ctxKey, ctxVal),
		},
		Handler: myHandler,
		Summary: "testing middlewares",
	})
	require.NoError(t, err)

	err = router.AddRoute(Route{
		Path:    "/another-route",
		Methods: []string{http.MethodGet},
		Middlewares: []func(http.Handler) http.Handler{
			setCtxVal(ctxKey, ctxVal2),
		},
		Handler: myHandler,
		Summary: "testing middlewares another route",
	})
	require.NoError(t, err)

	ts := httptest.NewServer(router.muxRouter)
	defer ts.Close()

	t.Run("first route", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/")
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		respBytes, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		expectedHandlerOutput := fmt.Sprintf("%s:%s", ctxKey, ctxVal)
		assert.Equal(t, expectedHandlerOutput, string(respBytes))
	})

	t.Run("another route", func(t *testing.T) {
		resp, err := http.Get(ts.URL + "/another-route")
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)

		respBytes, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		expectedHandlerOutput := fmt.Sprintf("%s:%s", ctxKey, ctxVal2)
		assert.Equal(t, expectedHandlerOutput, string(respBytes))
	})

}
