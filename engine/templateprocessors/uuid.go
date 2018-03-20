package templateprocessors

import (
	"bytes"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"regexp"
)

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

type UUID struct {
	re    *regexp.Regexp
	genFn GenUUID
}

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

func NewUUID(genFn GenUUID) UUID {
	return UUID{
		re:    regexp.MustCompile(`\$UUID\_(?P<variablename>\w+)`),
		genFn: genFn,
	}
}
