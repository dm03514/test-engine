package templateprocessors

import (
	"bytes"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"regexp"
)

// GenUUID a function to return a new UUID
type GenUUID func() uuid.UUID

type uuidMatch struct {
	match []string
}

func (u uuidMatch) toReplace() string {
	return u.match[0]
}
func (u uuidMatch) variableName() string {
	return u.match[1]
}

// UUID contains the regex to use to find UUID's a function to generate UUIDs
type UUID struct {
	re    *regexp.Regexp
	genFn GenUUID
}

// Process applies the UUID substitutions
func (u UUID) Process(b []byte) []byte {
	matches := u.re.FindAllStringSubmatch(string(b), -1)

	log.WithFields(log.Fields{
		"component": "templateprocessors.uuid",
		"matches":   matches,
	}).Info("Process()")

	for i := 0; i < len(matches); i++ {
		m := uuidMatch{match: matches[i]}
		id := u.genFn()
		b = bytes.Replace(b, []byte(m.toReplace()), []byte(id.String()), -1)
	}

	return b
}

// NewUUID creates a new UUID template processor
func NewUUID(genFn GenUUID) UUID {
	return UUID{
		re:    regexp.MustCompile(`\$UUID\_(?P<variablename>\w+)`),
		genFn: genFn,
	}
}
