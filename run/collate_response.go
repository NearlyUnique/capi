package run

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/NearlyUnique/capi/builder"
)

type (
	PrintableResponse struct {
		Status     string                 `json:"status"`
		StatusCode int                    `json:"statusCode"`
		Proto      string                 `json:"proto"`
		ProtoMajor int                    `json:"protoMajor"`
		ProtoMinor int                    `json:"protoMinor"`
		Header     PrintableHeader        `json:"header"`
		Body       *PrintableResponseBody `json:"body,omitempty"`
	}
	PrintableResponseBody json.RawMessage
	PrintableHeader       map[string]builder.StringOrList
)

func Collate(response *http.Response) PrintableResponse {
	p := PrintableResponse{
		Status:     response.Status,
		StatusCode: response.StatusCode,
		Proto:      response.Proto,
		ProtoMajor: response.ProtoMajor,
		ProtoMinor: response.ProtoMinor,
		Header:     convertHeaders(response.Header),
	}
	if response.Body != nil {
		defer func() { _ = response.Body.Close() }()
		buf, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Printf("download body failed: %v", err)
		}
		if len(buf) > 0 {
			rj := PrintableResponseBody(buf)
			p.Body = &rj
		}
	}
	return p
}

func convertHeaders(header http.Header) map[string]builder.StringOrList {
	h := make(map[string]builder.StringOrList)
	for k, v := range header {
		h[k] = v
	}
	return h
}
func (p *PrintableResponseBody) MarshalJSON() ([]byte, error) {
	if p == nil {
		return nil, nil
	}
	buf := []byte(*p)
	if len(buf) >= 2 {
		content := strings.TrimSpace(string(buf))
		if len(content) >= 2 {
			first, last := content[0], content[len(content)-1]
			if (first == '{' && last == '}') || (first == '[' && last == ']') {
				return []byte(content), nil
			}

		}
	}
	return json.Marshal(string(buf))
}

func (ph PrintableHeader) MarshalJSON() ([]byte, error) {
	b := map[string]builder.StringOrList(ph)
	return json.Marshal(b)
}
