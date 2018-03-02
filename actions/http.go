package actions

import (
	"fmt"
	"github.com/dm03514/test-engine/results"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type HttpResult struct {
	*http.Response

	err error
}

func (h HttpResult) Error() error {
	return h.err
}

func (h HttpResult) ValueOfProperty(property string) (results.Value, error) {
	switch property {
	case "status_code":
		return results.IntValue{V: h.Response.StatusCode}, nil
	default:
		return nil, fmt.Errorf("No property %s in %+v", property, h)
	}
}

type Http struct {
	Url       string
	Method    string
	Headers   map[string]string
	Overrides []results.Override
}

func (h Http) Execute(rs results.Results) (results.Result, error) {
	req, err := http.NewRequest(h.Method, h.Url, nil)
	if err != nil {
		return nil, err
	}

	// add headers

	// Transport Config
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	log.Infof("http.Execute() requested url: %q method: %q, received: %s",
		h.Url, h.Method, resp.StatusCode)

	return HttpResult{
		Response: resp,
	}, nil
}

func NewHttpFromMap(m map[string]interface{}) (Action, error) {
	var h Http
	err := mapstructure.Decode(m, &h)
	return h, err
}
