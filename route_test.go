package docrouter

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoute(t *testing.T) {

	t.Run("params", func(t *testing.T) {
		t.Run("path", func(t *testing.T) {

			type MyExamplePathParam struct {
				StarID int    `docrouter:"name:starId;desc:Star identifier in CommonMark syntax. This can be potentially issue for longer descriptions.; example: 5; schemaMin: 3"`
				Color  string `docrouter:"name:color;desc:This is string value.; example: Ciao!"`
			}

			r := Route{
				PathParams: &MyExamplePathParam{},
			}

			oaParams, err := r.openAPI3Params()
			require.NoError(t, err)
			starParam := oaParams.GetByInAndName(openapi3.ParameterInPath, "starId")
			require.NotNil(t, starParam)
			assert.Equal(t, "starId", starParam.Name)
			assert.Equal(t, "Star identifier in CommonMark syntax. This can be potentially issue for longer descriptions.", starParam.Description)
			assert.Equal(t, 5, starParam.Example)
			assert.True(t, starParam.Required)
			require.NotNil(t, starParam.Schema)
			require.NotNil(t, starParam.Schema.Value)
			assert.Equal(t, float64(3), *starParam.Schema.Value.Min)
			assert.Equal(t, "integer", starParam.Schema.Value.Type)

			colorParam := oaParams.GetByInAndName(openapi3.ParameterInPath, "color")
			require.NotNil(t, colorParam)
			assert.Equal(t, "color", colorParam.Name)
			assert.Equal(t, "This is string value.", colorParam.Description)
			assert.Equal(t, "Ciao!", colorParam.Example)
			assert.True(t, colorParam.Required)
			require.NotNil(t, colorParam.Schema)
			require.NotNil(t, colorParam.Schema.Value)
			assert.Equal(t, "string", colorParam.Schema.Value.Type)
		})

		t.Run("query", func(t *testing.T) {
			type MyExampleQueryParam struct {
				StarID int  `docrouter:"name:starId;desc:Star identifier in CommonMark syntax. This can be potentially issue for longer descriptions.; example: 5; required: false"`
				Potato bool `docrouter:"name:potato;desc: This is bool!; example: true; required: true"`
			}

			r := Route{
				QueryParams: &MyExampleQueryParam{},
			}

			oaParams, err := r.openAPI3Params()
			require.NoError(t, err)
			starParam := oaParams.GetByInAndName(openapi3.ParameterInQuery, "starId")
			require.NotNil(t, starParam)
			assert.Equal(t, "starId", starParam.Name)
			assert.Equal(t, "Star identifier in CommonMark syntax. This can be potentially issue for longer descriptions.", starParam.Description)
			assert.Equal(t, 5, starParam.Example)
			assert.False(t, starParam.Required)
			require.NotNil(t, starParam.Schema)
			require.NotNil(t, starParam.Schema.Value)
			assert.Equal(t, "integer", starParam.Schema.Value.Type)

			potatoParam := oaParams.GetByInAndName(openapi3.ParameterInQuery, "potato")
			require.NotNil(t, potatoParam)
			assert.Equal(t, "potato", potatoParam.Name)
			assert.Equal(t, "This is bool!", potatoParam.Description)
			assert.Equal(t, true, potatoParam.Example)
			assert.True(t, potatoParam.Required)
			require.NotNil(t, potatoParam.Schema)
			require.NotNil(t, potatoParam.Schema.Value)
			assert.Equal(t, "boolean", potatoParam.Schema.Value.Type)
		})

		t.Run("headers", func(t *testing.T) {
			type MyExampleHeadersParam struct {
				StarName string `docrouter:"name:Star-Name;desc: This is star name header param!; example: Sun; required: true"`
			}

			r := Route{
				HeaderParams: &MyExampleHeadersParam{},
			}

			oaParams, err := r.openAPI3Params()
			require.NoError(t, err)
			starParam := oaParams.GetByInAndName(openapi3.ParameterInHeader, "Star-Name")
			require.NotNil(t, starParam)
			assert.Equal(t, "Star-Name", starParam.Name)
			assert.Equal(t, "This is star name header param!", starParam.Description)
			assert.Equal(t, "Sun", starParam.Example)
			assert.True(t, starParam.Required)
			require.NotNil(t, starParam.Schema)
			require.NotNil(t, starParam.Schema.Value)
			assert.Equal(t, "string", starParam.Schema.Value.Type)

			t.Run("not-allowed names", func(t *testing.T) {
				// todo: Content-Type, Accept, Authorization https://swagger.io/docs/specification/describing-parameters/#header-parameters
			})
		})
	})
}
