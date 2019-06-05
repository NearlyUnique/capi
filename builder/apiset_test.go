package builder_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/NearlyUnique/capi/builder"
	"github.com/stretchr/testify/require"
)

func Test_body_can_be_json(t *testing.T) {
	body := `{
          "body": {
            "literal_string": "string",
            "literal_number": 42,
            "literal_bool": true,
            "object" : {"isObject": true},
            "with_replacement": "{attribute}"
          }
        }`

	var cmd builder.Command
	err := json.Unmarshal([]byte(body), &cmd)

	require.NoError(t, err)
	assert.Equal(t, byte('{'), cmd.Body.Data[0])
	assert.Equal(t, byte('}'), cmd.Body.Data[len(cmd.Body.Data)-1])
}

func Test_body_can_be_a_string(t *testing.T) {
	body := `{"body": "some string content"}`

	var cmd builder.Command
	err := json.Unmarshal([]byte(body), &cmd)

	require.NoError(t, err)
	assert.Equal(t, "some string content", string(cmd.Body.String()))
}

func Test_body_can_be_a_number(t *testing.T) {
	body := `{"body": 12.34}`

	var cmd builder.Command
	err := json.Unmarshal([]byte(body), &cmd)

	require.NoError(t, err)
	assert.Equal(t, "12.34", string(cmd.Body.String()))
}

func Test_body_can_be_an_array(t *testing.T) {
	body := `{"body": [1,2,3,4]}`

	var cmd builder.Command
	err := json.Unmarshal([]byte(body), &cmd)

	require.NoError(t, err)
	assert.Equal(t, "[1,2,3,4]", string(cmd.Body.String()))
}
