package models

import (
	"github.com/nmatsui/fabric-payment-sample-chaincode/types"
)

type Account struct {
	ModelType types.ModelType `json:"model_type"`
	No        string          `json:"no"`
	Name      string          `json:"name"`
	Balance   int             `json:"balance"`
}
