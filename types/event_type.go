package types

import (
	"encoding/json"
)

const (
	unknownEventStr  = "unknown"
	depositEventStr  = "deposit"
	remitEventStr    = "remit"
	withdrawEventStr = "withdraw"
)

// EventType : event type
type EventType int

// concrete EventType
const (
	UnKnownEvent EventType = iota
	DepositEvent
	RemitEvent
	WithdrawEvent
)

// String : Striner interface
func (t EventType) String() string {
	switch t {
	case DepositEvent:
		return depositEventStr
	case RemitEvent:
		return remitEventStr
	case WithdrawEvent:
		return withdrawEventStr
	default:
		return unknownEventStr
	}
}

// MarshalJSON : Marshaler interface of ObjectType
func (t EventType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// UnmarshalJSON : Marshaler interface of ObjectType
func (t *EventType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	switch s {
	case depositEventStr:
		*t = DepositEvent
	case remitEventStr:
		*t = RemitEvent
	case withdrawEventStr:
		*t = WithdrawEvent
	default:
		*t = UnKnownEvent
	}
	return nil
}
