/*
 Package models provides the model of state objects.

 Copyright Nobuyuki Matsui<nobuyuki.matsui>.

 SPDX-License-Identifier: Apache-2.0
*/
package models

import (
	"github.com/nmatsui/fabric-payment-sample-chaincode/types"
)

// AccountState: Holder to show the change of balance.
type AccountState struct {
	No              string `json:"no"`
	Name            string `json:"name"`
	PreviousBalance int    `json:"previous_balance"`
	CurrentBalance  int    `json:"current_balance"`
}

// Event: Event model to show deposit, remit or withdraw event.
type Event struct {
	ModelType        types.ModelType `json:"model_type"`
	EventType        types.EventType `json:"event_type"`
	No               string          `json:"no"`
	Amount           int             `json:"amount"`
	FromAccountState *AccountState   `json:"from_account"`
	ToAccountState   *AccountState   `json:"to_account"`
}
