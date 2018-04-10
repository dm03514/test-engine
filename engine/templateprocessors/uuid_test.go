// +build unit

package templateprocessors

import (
	"github.com/satori/go.uuid"
	"reflect"
	"testing"
)

type uuidGen struct {
	callNum int
	stubIds []string
}

func (u *uuidGen) gen() uuid.UUID {
	i := u.callNum
	u.callNum++
	uid, _ := uuid.FromString(u.stubIds[i])
	return uid
}
func newUUIDGen() *uuidGen {
	return &uuidGen{
		stubIds: []string{
			"50e0ad50-7389-4e8e-9aa1-515a51ebf6d9",
			"0162c6b3-67ea-457a-9935-9dddaade380d",
		},
	}
}

func TestUUID_Process(t *testing.T) {
	testCases := []struct {
		name     string
		b        []byte
		expected []byte
	}{
		{
			"pass_through",
			[]byte("hello"),
			[]byte("hello"),
		},
		{
			"single_uuid",
			[]byte("hello, $UUID_name"),
			[]byte("hello, 50e0ad50-7389-4e8e-9aa1-515a51ebf6d9"),
		},
		{
			"single_uuid_multiple_targets",
			[]byte("hello, $UUID_name, also welcome $UUID_name"),
			[]byte("hello, 50e0ad50-7389-4e8e-9aa1-515a51ebf6d9, also welcome 50e0ad50-7389-4e8e-9aa1-515a51ebf6d9"),
		},
		{
			"multiple_uuids",
			[]byte("hello, $UUID_name, also welcome $UUID_id"),
			[]byte("hello, 50e0ad50-7389-4e8e-9aa1-515a51ebf6d9, also welcome 0162c6b3-67ea-457a-9935-9dddaade380d"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := NewUUID(newUUIDGen().gen)
			out := u.Process(tc.b)
			if !reflect.DeepEqual(out, tc.expected) {
				t.Errorf("Expected %q Received %q", tc.expected, out)
			}
		})
	}
}
