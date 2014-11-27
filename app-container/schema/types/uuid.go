package types

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	ErrNoEmptyUUID = errors.New("UUID cannot be empty")
)

// UUID encodes an RFC4122-compliant UUID, marshaled to/from a string
// TODO(jonboulle): vendor a package for this?
// TODO(jonboulle): consider more flexibility in input string formats.
// Right now, we only accept:
//   "6733C088-A507-4694-AABF-EDBE4FC5266F"
//   "6733C088A5074694AABFEDBE4FC5266F"
type UUID [16]byte

func (u UUID) String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:16])
}

// NewUUID generates a new UUID from the given string. If the string does not
// represent a valid UUID, nil and an error are returned.
func NewUUID(s string) (*UUID, error) {
	s = strings.Replace(s, "-", "", -1)
	if len(s) != 32 {
		return nil, errors.New("bad UUID length != 32")
	}
	dec, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	var u UUID
	for i, b := range dec {
		u[i] = b
	}
	return &u, nil
}

func (u UUID) Empty() bool {
	return reflect.DeepEqual(u, UUID{})
}

func (u *UUID) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	uu, err := NewUUID(s)
	if uu.Empty() {
		return ErrNoEmptyUUID
	}
	if err == nil {
		*u = *uu
	}
	return err
}

func (u UUID) MarshalJSON() ([]byte, error) {
	if u.Empty() {
		return nil, ErrNoEmptyUUID
	}
	return json.Marshal(u.String())
}