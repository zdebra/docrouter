package docrouter

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeParams(t *testing.T) {

	server := New(DefaultOptions)

	type MyParameters struct {
		StarID int  `docrouter:"name:starid; kind:query; desc: This is int!; example: 5; required: false; schemaMin: 3"`
		Potato bool `docrouter:"name:potato; kind:query; desc: This is bool!; example: true; required: true"`
	}

	const (
		expectedStarID = 10
		expectedPotato = true
	)

	err := server.AddRoute(Route{
		Path:       "/example-param",
		Methods:    []string{http.MethodGet},
		Parameters: &MyParameters{},
		Summary:    "Parses input params",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var inputParams MyParameters
			if err := DecodeParams(&inputParams, r); err != nil {
				t.Log("decode query params error", err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			assert.Equal(t, expectedStarID, inputParams.StarID)
			assert.Equal(t, expectedPotato, inputParams.Potato)
		}),
	})
	require.NoError(t, err)

	ts := httptest.NewServer(server.muxRouter)
	defer ts.Close()

	resp, err := http.Get(ts.URL + fmt.Sprintf("/example-param?starid=%d&potato=%t", expectedStarID, expectedPotato))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

}
