/*
 Package models provides the model of state objects.

 Copyright Nobuyuki Matsui<nobuyuki.matsui>.

 SPDX-License-Identifier: Apache-2.0
*/
package models

import (
	"github.com/nmatsui/fabric-payment-sample-chaincode/types"
)

// Account: Account model
type Account struct {
	ModelType types.ModelType `json:"model_type"`
	No        string          `json:"no"`
	Name      string          `json:"name"`
	Balance   int             `json:"balance"`
}
