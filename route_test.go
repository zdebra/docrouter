package docrouter

import (
	"fmt"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoute(t *testing.T) {

	t.Run("params", func(t *testing.T) {
		type MyParameters struct {
			StarID    int    `docrouter:"name:starId; kind:path; desc:Star identifier in CommonMark syntax. This can be potentially issue for longer descriptions.; example: 5; schemaMin: 3"`
			Color     string `docrouter:"name:color; kind:path; desc:This is string value.; example: Ciao!"`
			Limit     int    `docrouter:"name:limit; kind: query; desc:Star limit; example: 93; required: false"`
			Potato    bool   `docrouter:"name:potato; kind: query; desc: This is bool!; example: true; required: true"`
			StarName  string `docrouter:"name:Star-Name; kind: header; desc: This is star name header param!; example: Sun; required: true"`
			SessionID string `docrouter:"name:sessionId; kind: cookie; desc: session identifier; example: bananas"`
		}

		r := Route{
			Parameters: &MyParameters{},
		}

		oaParams, err := r.openAPI3Params()
		require.NoError(t, err)

		tests := []struct {
			kind     string
			name     string
			desc     string
			example  interface{}
			required bool
			schema   *openapi3.Schema
		}{
			{
				kind:     openapi3.ParameterInPath,
				name:     "starId",
				desc:     "Star identifier in CommonMark syntax. This can be potentially issue for longer descriptions.",
				example:  5,
				required: true,
				schema:   openapi3.NewIntegerSchema().WithMin(3),
			},
			{
				kind:     openapi3.ParameterInPath,
				name:     "color",
				desc:     "This is string value.",
				example:  "Ciao!",
				required: true,
				schema:   openapi3.NewStringSchema(),
			},
			{
				kind:     openapi3.ParameterInQuery,
				name:     "limit",
				desc:     "Star limit",
				example:  93,
				required: false,
				schema:   openapi3.NewIntegerSchema(),
			},
			{
				kind:     openapi3.ParameterInQuery,
				name:     "potato",
				desc:     "This is bool!",
				example:  true,
				required: true,
				schema:   openapi3.NewBoolSchema(),
			},
			{
				kind:     openapi3.ParameterInHeader,
				name:     "Star-Name",
				desc:     "This is star name header param!",
				example:  "Sun",
				required: true,
				schema:   openapi3.NewStringSchema(),
			},
			{
				kind:     openapi3.ParameterInCookie,
				name:     "sessionId",
				desc:     "session identifier",
				example:  "bananas",
				required: false,
				schema:   openapi3.NewStringSchema(),
			},
		}

		for _, test := range tests {
			testName := fmt.Sprintf("%s/%s", test.kind, test.name)
			t.Run(testName, func(t *testing.T) {
				oaParam := oaParams.GetByInAndName(test.kind, test.name)
				require.NotNil(t, oaParam)
				assert.Equal(t, test.name, oaParam.Name)
				assert.Equal(t, test.desc, oaParam.Description)
				assert.Equal(t, test.example, oaParam.Example)
				assert.Equal(t, test.required, oaParam.Required)
				require.NotNil(t, oaParam.Schema)
				require.NotNil(t, oaParam.Schema.Value)
				require.Equal(t, test.schema, oaParam.Schema.Value)
			})
		}

		t.Run("not-allowed header names", func(t *testing.T) {
			t.Skip("todo: not-allowed header names Content-Type, Accept, Authorization https://swagger.io/docs/specification/describing-parameters/#header-parameters")
		})
	})

	t.Run("request", func(t *testing.T) {
		type MyRequestBody struct {
			Name                     string
			SurfaceTemperatureKelvin int
			Mass                     float64
			OlderThanSun             bool
		}
		r := Route{
			RequestBody: &MyRequestBody{},
		}

		oaParams, err := r.openAPI3Params()
		require.NoError(t, err)
	})
}
