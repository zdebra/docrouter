package docrouter

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeParams(t *testing.T) {

	server := New(DefaultOptions)

	type MyParameters struct {
		StarID     int    `docrouter:"name: starid; kind: query; desc: This is int!; example: 5; required: false; schemaMin: 3"`
		Potato     bool   `docrouter:"name: potato; kind: query; desc: This is bool!; example: true; required: true"`
		FishName   string `docrouter:"name: fishName; kind: path"`
		VisitCount int    `docrouter:"name: VISIT_COUNT; kind: cookie"`
		Color      string `docrouter:"name: Color; kind: header"`
	}

	const (
		expectedStarID     = 10
		expectedPotato     = true
		expectedFishName   = "blump"
		expectedVisitCount = 999_999_999
		expectedColor      = "purple"
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
			assert.Equal(t, expectedVisitCount, inputParams.VisitCount)
			assert.Equal(t, expectedColor, inputParams.Color)
		}),
	})
	require.NoError(t, err)

	ts := httptest.NewServer(server.muxRouter)
	defer ts.Close()

	u := ts.URL + fmt.Sprintf("/example-param/%s?starid=%d&potato=%t", expectedFishName, expectedStarID, expectedPotato)
	req, err := http.NewRequest(http.MethodGet, u, nil)
	require.NoError(t, err)

	req.AddCookie(&http.Cookie{
		Name:  "VISIT_COUNT",
		Value: strconv.Itoa(expectedVisitCount),
	})

	req.Header.Set("Color", expectedColor)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

}
