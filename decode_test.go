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
		StarID   int    `docrouter:"name: starid; kind: query; desc: This is int!; example: 5; required: false; schemaMin: 3"`
		Potato   bool   `docrouter:"name: potato; kind: query; desc: This is bool!; example: true; required: true"`
		FishName string `docrouter:"name: fishName; kind: path"`
	}

	const (
		expectedStarID   = 10
		expectedPotato   = true
		expectedFishName = "blump"
	)

	err := server.AddRoute(Route{
		Path:       "/example-param/{fishName}",
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
			assert.Equal(t, expectedFishName, inputParams.FishName)
		}),
	})
	require.NoError(t, err)

	ts := httptest.NewServer(server.muxRouter)
	defer ts.Close()

	resp, err := http.Get(ts.URL + fmt.Sprintf("/example-param/%s?starid=%d&potato=%t", expectedFishName, expectedStarID, expectedPotato))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

}
