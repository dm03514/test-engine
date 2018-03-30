package actions

import (
	"context"
	"fmt"
	"github.com/dm03514/test-engine/ids"
	"github.com/dm03514/test-engine/results"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
)

// HTTPResult contains request result and metadata
type HTTPResult struct {
	*http.Response

	err  error
	body []byte
}

// Error gets the error associated with this result
func (h HTTPResult) Error() error {
	return h.err
}

// ValueOfProperty gets value associated with an identifier
func (h HTTPResult) ValueOfProperty(property string) (results.Value, error) {
	switch property {
	case "status_code":
		return results.IntValue{V: h.Response.StatusCode}, nil
	case "body":
		return results.StringValue{V: string(h.body)}, nil
	default:
		return nil, fmt.Errorf("No property %s in %+v", property, h)
	}
}

// HTTP action has all data necessary to make a request
type HTTP struct {
	URL       string `mapstructure:"url"`
	Method    string
	Headers   map[string]string
	Overrides []results.Override
	Body      string
	Type      string
}

// Execute makes the http request
func (h HTTP) Execute(ctx context.Context, rs results.Results) (results.Result, error) {
	req, err := http.NewRequest(h.Method, h.URL, strings.NewReader(h.Body))
	if err != nil {
		return nil, err
	}

	// add headers

	// add body

	// Transport Config
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	log.WithFields(log.Fields{
		"component":    h.Type,
		"call":         "Execute()",
		"url":          h.URL,
		"method":       h.Method,
		"status_code":  resp.StatusCode,
		"execution_id": ctx.Value(ids.Execution("execution_id")),
	}).Info("requested_url")

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return HTTPResult{
		Response: resp,
		body:     b,
	}, nil
}

// NewHTTPFromMap initializes an http action based on generic map
func NewHTTPFromMap(m map[string]interface{}) (Action, error) {
	var h HTTP
	err := mapstructure.Decode(m, &h)
	return h, err
}
