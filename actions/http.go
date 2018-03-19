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

type HttpResult struct {
	*http.Response

	err  error
	body []byte
}

func (h HttpResult) Error() error {
	return h.err
}

func (h HttpResult) ValueOfProperty(property string) (results.Value, error) {
	switch property {
	case "status_code":
		return results.IntValue{V: h.Response.StatusCode}, nil
	case "body":
		return results.StringValue{V: string(h.body)}, nil
	default:
		return nil, fmt.Errorf("No property %s in %+v", property, h)
	}
}

type Http struct {
	Url       string
	Method    string
	Headers   map[string]string
	Overrides []results.Override
	Body      string
	Type      string
}

func (h Http) Execute(ctx context.Context, rs results.Results) (results.Result, error) {
	req, err := http.NewRequest(h.Method, h.Url, strings.NewReader(h.Body))
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
		"url":          h.Url,
		"method":       h.Method,
		"status_code":  resp.StatusCode,
		"execution_id": ctx.Value(ids.Execution("execution_id")),
	}).Info("requested_url")

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return HttpResult{
		Response: resp,
		body:     b,
	}, nil
}

func NewHttpFromMap(m map[string]interface{}) (Action, error) {
	var h Http
	err := mapstructure.Decode(m, &h)
	return h, err
}
