package docrouter

import (
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoute(t *testing.T) {

	type MyExamplePathParam struct {
		StarID int    `docrouter:"name:starId;desc:Star identifier in CommonMark syntax. This can be potentially issue for longer descriptions.; example: 5"`
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

	colorParam := oaParams.GetByInAndName(openapi3.ParameterInPath, "color")
	require.NotNil(t, colorParam)
	assert.Equal(t, "color", colorParam.Name)
	assert.Equal(t, "This is string value.", colorParam.Description)
	assert.Equal(t, "Ciao!", colorParam.Example)
	assert.True(t, colorParam.Required)
}
