package templateprocessors

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"regexp"
)

// LookupEnv fetches a key from the environment
type LookupEnv func(key string) (string, bool)

// Env contains regular expression to look for and function to pull variables from env
type Env struct {
	re        *regexp.Regexp
	lookupEnv LookupEnv
}

// Process replaces templated variables with variables from the environment
// TODO not sure about the type conversions, favoring strings for ease of logging
// but probably an easier way anyhow
func (e Env) Process(b []byte) []byte {
	matches := e.re.FindAllStringSubmatch(string(b), -1)

	log.WithFields(log.Fields{
		"component": "templateprocessors.env",
		"matches":   matches,
	}).Info("Process()")

	for i := 0; i < len(matches); i++ {
		toReplace := matches[i][0]
		m := matches[i][1]
		envValue, ok := e.lookupEnv(m)
		if ok {
			b = bytes.Replace(b, []byte(toReplace), []byte(envValue), -1)
		} else {
			log.WithFields(log.Fields{
				"component":       "templateprocessors.env",
				"to_replace":      toReplace,
				"looking_for_env": m,
			}).Warn("not_found_in_env")
		}
	}

	return b
}

// NewEnv creates an env parser
func NewEnv(lookupEnv LookupEnv) Env {
	return Env{
		re:        regexp.MustCompile(`\$ENV\_(?P<variablename>\w+)`),
		lookupEnv: lookupEnv,
	}
}
