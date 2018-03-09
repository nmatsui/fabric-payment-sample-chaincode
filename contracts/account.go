/*
 Package contracts provides the smart contracts for Hyperledger/fabric 1.1.

 Copyright Nobuyuki Matsui<nobuyuki.matsui>.

 SPDX-License-Identifier: Apache-2.0
*/
package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"

	"github.com/nmatsui/fabric-payment-sample-chaincode/models"
	"github.com/nmatsui/fabric-payment-sample-chaincode/types"
	"github.com/nmatsui/fabric-payment-sample-chaincode/utils"
)

var accountLogger = shim.NewLogger("contracts/account")

// AccountContract : a struct to handle Account.
type AccountContract struct {
}

// ListAccount : return a list of all accounts.
func (ac *AccountContract) ListAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	accountLogger.Infof("invoke ListAccount, args=%s\n", args)
	if len(args) != 0 {
		errMsg := fmt.Sprintf("Incorrect number of arguments. Expecting = no argument, Actual = %s\n", args)
		accountLogger.Error(errMsg)
		return shim.Error(errMsg)
	}

	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"model_type": types.AccountModel,
		},
	}

	queryBytes, err := json.Marshal(query)
	if err != nil {
		accountLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	accountLogger.Infof("Query string = '%s'", string(queryBytes))
	resultsIterator, err := APIstub.GetQueryResult(string(queryBytes))
	if err != nil {
		accountLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	results := make([]*models.Account, 0)
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			accountLogger.Error(err.Error())
			return shim.Error(err.Error())
		}
		account := new(models.Account)
		if err := json.Unmarshal(queryResponse.Value, account); err != nil {
			accountLogger.Error(err.Error())
			return shim.Error(err.Error())
		}
		results = append(results, account)
	}
	jsonBytes, err := json.Marshal(results)
	if err != nil {
		accountLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	return shim.Success(jsonBytes)
}

// CreateAccount : create a new account.
func (ac *AccountContract) CreateAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	accountLogger.Infof("invoke CreateAccount, args=%s\n", args)
	if len(args) != 1 {
		errMsg := fmt.Sprintf("Incorrect number of arguments. Expecting = ['name'], Actual = %s\n", args)
		accountLogger.Error(errMsg)
		return shim.Error(errMsg)
	}
	name := args[0]

	no, err := utils.GetAccountNo(APIstub)
	if err != nil {
		accountLogger.Error(err.Error())
		return shim.Error(err.Error())
	}

	account := models.Account{
		ModelType: types.AccountModel,
		No:        no,
		Name:      name,
		Balance:   0,
	}
	jsonBytes, err := json.Marshal(account)
	if err != nil {
		accountLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	if err := APIstub.PutState(no, jsonBytes); err != nil {
		accountLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	return shim.Success(jsonBytes)
}

// RetrieveAccount : return an account.
func (ac *AccountContract) RetrieveAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	accountLogger.Infof("invoke RetrieveAccount, args=%s\n", args)
	if len(args) != 1 {
		errMsg := fmt.Sprintf("Incorrect number of arguments. Expecting = ['no'], Actual = %s\n", args)
		accountLogger.Error(errMsg)
		return shim.Error(errMsg)
	}
	no := args[0]

	account, err := utils.GetAccount(APIstub, no)
	if err != nil {
		switch e := err.(type) {
		case *utils.WarningResult:
			accountLogger.Warning(err.Error())
			return shim.Success(e.JSONBytes())
		default:
			accountLogger.Error(err.Error())
			return shim.Error(err.Error())
		}
	}

	jsonBytes, err := json.Marshal(account)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(jsonBytes)
}

// UpdateAccountName : update the name of an account.
func (ac *AccountContract) UpdateAccountName(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	accountLogger.Infof("invoke UpdateAccountName, args=%s\n", args)
	if len(args) != 2 {
		errMsg := fmt.Sprintf("Incorrect number of arguments. Expecting = ['no', 'name'], Actual = %s\n", args)
		accountLogger.Error(errMsg)
		return shim.Error(errMsg)
	}
	no := args[0]
	name := args[1]

	account, err := utils.GetAccount(APIstub, no)
	if err != nil {
		switch e := err.(type) {
		case *utils.WarningResult:
			accountLogger.Warning(err.Error())
			return shim.Success(e.JSONBytes())
		default:
			accountLogger.Error(err.Error())
			return shim.Error(err.Error())
		}
	}

	account.Name = name

	jsonBytes, err := json.Marshal(account)
	if err != nil {
		accountLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	if err := APIstub.PutState(no, jsonBytes); err != nil {
		accountLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	return shim.Success(jsonBytes)
}

// DeleteAccount : delete an account.
func (ac *AccountContract) DeleteAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	accountLogger.Infof("invoke DeleteAccount, args=%s\n", args)
	if len(args) != 1 {
		errMsg := fmt.Sprintf("Incorrect number of arguments. Expecting = ['no'], Actual = %s\n", args)
		accountLogger.Error(errMsg)
		return shim.Error(errMsg)
	}
	no := args[0]

	_, err := utils.GetAccount(APIstub, no)
	if err != nil {
		switch e := err.(type) {
		case *utils.WarningResult:
			accountLogger.Warning(err.Error())
			return shim.Success(e.JSONBytes())
		default:
			accountLogger.Error(err.Error())
			return shim.Error(err.Error())
		}
	}

	if err := APIstub.DelState(no); err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}
