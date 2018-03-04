package utils

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"github.com/nmatsui/fabric-payment-sample/models"
)

// GetAccount : get an account from no
func GetAccount(APIstub shim.ChaincodeStubInterface, no string) (*models.Account, error) {
	var account = new(models.Account)
	accountBytes, err := APIstub.GetState(no)
	if err != nil {
		return account, err
	} else if accountBytes == nil {
		msg := fmt.Sprintf("Account does not exist, no = %s", no)
		warning := &WarningResult{StatusCode: 404, Message: msg}
		return account, warning
	}
	if err := json.Unmarshal(accountBytes, account); err != nil {
		return account, err
	}
	return account, nil
}

// GetAmount : convert amount to int and validate it
func GetAmount(amountStr string) (int, error) {
	var amount int
	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		msg := fmt.Sprintf("amount is not integer, amount = %s", amountStr)
		warning := &WarningResult{StatusCode: 400, Message: msg}
		return amount, warning
	}
	if amount < 0 {
		msg := fmt.Sprintf("amount is less than zero, amount = %d", amount)
		warning := &WarningResult{StatusCode: 400, Message: msg}
		return amount, warning
	}
	return amount, nil
}
