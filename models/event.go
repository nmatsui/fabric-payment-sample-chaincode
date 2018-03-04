package models

import (
	"github.com/nmatsui/fabric-payment-sample/types"
)

type AccountState struct {
	No              string `json:"no"`
	Name            string `json:"name"`
	PreviousBalance int    `json:"previous_balance"`
	CurrentBalance  int    `json:"current_balance"`
}

type Event struct {
	ModelType        types.ModelType `json:"model_type"`
	EventType        types.EventType `json:"event_type"`
	No               string          `json:"no"`
	Amount           int             `json:"amount"`
	FromAccountState *AccountState   `json:"from_account"`
	ToAccountState   *AccountState   `json:"to_account"`
}
