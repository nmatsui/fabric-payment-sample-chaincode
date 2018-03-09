/*
 Package types provides the enum like type.

 Copyright Nobuyuki Matsui<nobuyuki.matsui>.

 SPDX-License-Identifier: Apache-2.0
*/
package types

import (
	"encoding/json"
)

const (
	unknownModelStr = "unknown"
	accountModelStr = "account"
	eventModelStr   = "event"
)

// ModelType : model type
type ModelType int

// concrete ModelType
const (
	UnKnownModel ModelType = iota
	AccountModel
	EventModel
)

// String : Stringer interface
func (t ModelType) String() string {
	switch t {
	case AccountModel:
		return accountModelStr
	case EventModel:
		return eventModelStr
	default:
		return unknownModelStr
	}
}

// MarshalJSON : Marshaler interface
func (t ModelType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// UnmarshalJSON : Marshaler interface
func (t *ModelType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s {
	case accountModelStr:
		*t = AccountModel
	case eventModelStr:
		*t = EventModel
	default:
		*t = UnKnownModel
	}
	return nil
}
