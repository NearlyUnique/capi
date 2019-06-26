package run_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NearlyUnique/capi/run"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_response_with_no_body_is_marshaled(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("h1", "v1")
		w.Header().Add("h2", "v2-1")
		w.Header().Add("h2", "v2-2")
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	resp, err := ts.Client().Get(ts.URL)
	require.NoError(t, err)

	printable := run.Collate(resp)

	buf, err := json.Marshal(printable)

	require.NoError(t, err)
	require.NotNil(t, buf)

	content := string(buf)
	assert.Contains(t, content, `"status":"200 OK"`)
	assert.Contains(t, content, `"statusCode":200`)
	assert.Contains(t, content, `"proto":"HTTP/1.1"`)
	assert.Contains(t, content, `"protoMajor":1`)
	assert.Contains(t, content, `"protoMinor":1`)
	assert.Contains(t, content, `"header"`)
	assert.Contains(t, content, `"H1":"v1"`)
	assert.Contains(t, content, `"H2":[`)
	assert.Contains(t, content, `"v2-1"`)
	assert.Contains(t, content, `"v2-2"`)
	assert.NotContains(t, content, `"body"`)
}

func Test_body_marshalling(t *testing.T) {
	var responseContent string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(responseContent))
		require.NoError(t, err)
	}))
	defer ts.Close()

	testData := []struct {
		name, content, expected string
	}{
		{"valid json object is written as literal",
			"{}", "{}"},
		{"valid json object is written as literal with trailing spaces",
			" {}  ", "{}"},
		{"valid json object is written as literal",
			`any text`, `"any text"`},
		{"valid json array is written as literal",
			`[1,2,3]`, `[1,2,3]`},
	}
	for _, td := range testData {

		t.Run(td.name, func(t *testing.T) {

			responseContent = td.content
			resp, err := ts.Client().Get(ts.URL)
			require.NoError(t, err)

			printable := run.Collate(resp)

			buf, err := json.Marshal(printable)
			s := string(buf)
			_ = s
			require.NoError(t, err)
			require.NotNil(t, buf)
			var actual struct{ Body *json.RawMessage }

			assert.NoError(t, json.Unmarshal(buf, &actual))
			assert.Equal(t, td.expected, string(*actual.Body))
		})
	}
}
