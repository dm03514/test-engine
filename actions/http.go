package actions

import (
	"fmt"
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
}

func (h Http) Execute(rs results.Results) (results.Result, error) {
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

	log.Infof("http.Execute() requested url: %q method: %q, received: %s",
		h.Url, h.Method, resp.StatusCode)

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
