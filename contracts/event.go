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

var eventLogger = shim.NewLogger("contracts/event")

// EventContract : a struct to handle Event.
type EventContract struct {
}

// ListEvent : return a list of events.
func (ec *EventContract) ListEvent(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	eventLogger.Infof("invoke ListEvent, args=%s\n", args)
	if len(args) > 1 {
		errMsg := fmt.Sprintf("Incorrect number of arguments. Expecting = [Optional('%s'|'%s'|'%s'), Actual = %s\n", types.DepositEvent, types.RemitEvent, types.WithdrawEvent, args)
		eventLogger.Error(errMsg)
		return shim.Error(errMsg)
	}

	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"model_type": types.EventModel,
		},
	}

	if len(args) == 1 {
		switch args[0] {
		case types.DepositEvent.String():
			query["selector"].(map[string]interface{})["event_type"] = types.DepositEvent
		case types.RemitEvent.String():
			query["selector"].(map[string]interface{})["event_type"] = types.RemitEvent
		case types.WithdrawEvent.String():
			query["selector"].(map[string]interface{})["event_type"] = types.WithdrawEvent
		default:
			errMsg := fmt.Sprintf("Incorrect arguments. Expecting = [Optional('%s'|'%s'|'%s'), Actual = %s\n", types.DepositEvent, types.RemitEvent, types.WithdrawEvent, args)
			eventLogger.Error(errMsg)
			return shim.Error(errMsg)
		}
	}

	queryBytes, err := json.Marshal(query)
	if err != nil {
		eventLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	eventLogger.Infof("Query string = '%s'", string(queryBytes))
	resultsIterator, err := APIstub.GetQueryResult(string(queryBytes))
	if err != nil {
		eventLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	results := make([]*models.Event, 0)
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			eventLogger.Error(err.Error())
			return shim.Error(err.Error())
		}
		event := new(models.Event)
		if err := json.Unmarshal(queryResponse.Value, event); err != nil {
			eventLogger.Error(err.Error())
			return shim.Error(err.Error())
		}
		results = append(results, event)
	}
	jsonBytes, err := json.Marshal(results)
	if err != nil {
		eventLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	return shim.Success(jsonBytes)
}

// Deposit : deposit to an account.
func (ec *EventContract) Deposit(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	eventLogger.Infof("invoke Deposit, args=%s\n", args)
	if len(args) != 2 {
		errMsg := fmt.Sprintf("Incorrect number of arguments. Expecting = ['to_account_no', 'amount'], Actual = %s\n", args)
		eventLogger.Error(errMsg)
		return shim.Error(errMsg)
	}
	toAccountNo := args[0]
	amountStr := args[1]

	amount, err := utils.GetAmount(amountStr)
	if err != nil {
		switch e := err.(type) {
		case *utils.WarningResult:
			eventLogger.Warning(err.Error())
			return shim.Success(e.JSONBytes())
		default:
			eventLogger.Error(err.Error())
			return shim.Error(err.Error())
		}
	}

	toAccount, err := utils.GetAccount(APIstub, toAccountNo)
	if err != nil {
		switch e := err.(type) {
		case *utils.WarningResult:
			eventLogger.Warning(err.Error())
			return shim.Success(e.JSONBytes())
		default:
			eventLogger.Error(err.Error())
			return shim.Error(err.Error())
		}
	}

	eventNo, err := utils.GetEventNo(APIstub)
	if err != nil {
		eventLogger.Error(err.Error())
		return shim.Error(err.Error())
	}

	toAccountPreviousBalance := toAccount.Balance
	toAccount.Balance += amount

	toAccountState := &models.AccountState{
		No:              toAccount.No,
		Name:            toAccount.Name,
		PreviousBalance: toAccountPreviousBalance,
		CurrentBalance:  toAccount.Balance,
	}

	event := &models.Event{
		ModelType:        types.EventModel,
		EventType:        types.DepositEvent,
		No:               eventNo,
		Amount:           amount,
		FromAccountState: nil,
		ToAccountState:   toAccountState,
	}

	toAccountBytes, err := json.Marshal(toAccount)
	if err != nil {
		eventLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	if err := APIstub.PutState(toAccount.No, toAccountBytes); err != nil {
		eventLogger.Error(err.Error())
		return shim.Error(err.Error())
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return shim.Error(err.Error())
	}
	if err := APIstub.PutState(event.No, eventBytes); err != nil {
		eventLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	return shim.Success(eventBytes)
}

// Remit : remit from an account to another account
func (ec *EventContract) Remit(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	eventLogger.Infof("invoke Remit, args=%s\n", args)
	if len(args) != 3 {
		errMsg := fmt.Sprintf("Incorrect number of arguments. Expecting = ['from_account_no, 'to_account_no', 'amount'], Actual = %s\n", args)
		eventLogger.Error(errMsg)
		return shim.Error(errMsg)
	}
	fromAccountNo := args[0]
	toAccountNo := args[1]
	amountStr := args[2]

	amount, err := utils.GetAmount(amountStr)
	if err != nil {
		switch e := err.(type) {
		case *utils.WarningResult:
			eventLogger.Warning(err.Error())
			return shim.Success(e.JSONBytes())
		default:
			eventLogger.Error(err.Error())
			return shim.Error(err.Error())
		}
	}

	fromAccount, err := utils.GetAccount(APIstub, fromAccountNo)
	if err != nil {
		switch e := err.(type) {
		case *utils.WarningResult:
			eventLogger.Warning(err.Error())
			return shim.Success(e.JSONBytes())
		default:
			eventLogger.Error(err.Error())
			return shim.Error(err.Error())
		}
	}

	toAccount, err := utils.GetAccount(APIstub, toAccountNo)
	if err != nil {
		switch e := err.(type) {
		case *utils.WarningResult:
			eventLogger.Warning(err.Error())
			return shim.Success(e.JSONBytes())
		default:
			eventLogger.Error(err.Error())
			return shim.Error(err.Error())
		}
	}

	if fromAccount.Balance < amount {
		msg := fmt.Sprintf("amount is grator than the fromAccount.Balance, amount = %d, fromAccount.Balance = %d", amount, fromAccount.Balance)
		warning := &utils.WarningResult{StatusCode: 400, Message: msg}
		eventLogger.Warning(warning.Error())
		return shim.Success(warning.JSONBytes())
	}

	eventNo, err := utils.GetEventNo(APIstub)
	if err != nil {
		eventLogger.Error(err.Error())
		return shim.Error(err.Error())
	}

	fromAccountPreviousBalance := fromAccount.Balance
	fromAccount.Balance -= amount

	toAccountPreviousBalance := toAccount.Balance
	toAccount.Balance += amount

	fromAccountState := &models.AccountState{
		No:              fromAccount.No,
		Name:            fromAccount.Name,
		PreviousBalance: fromAccountPreviousBalance,
		CurrentBalance:  fromAccount.Balance,
	}

	toAccountState := &models.AccountState{
		No:              toAccount.No,
		Name:            toAccount.Name,
		PreviousBalance: toAccountPreviousBalance,
		CurrentBalance:  toAccount.Balance,
	}

	event := &models.Event{
		ModelType:        types.EventModel,
		EventType:        types.RemitEvent,
		No:               eventNo,
		Amount:           amount,
		FromAccountState: fromAccountState,
		ToAccountState:   toAccountState,
	}

	fromAccountBytes, err := json.Marshal(fromAccount)
	if err != nil {
		eventLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	if err := APIstub.PutState(fromAccount.No, fromAccountBytes); err != nil {
		eventLogger.Error(err.Error())
		return shim.Error(err.Error())
	}

	toAccountBytes, err := json.Marshal(toAccount)
	if err != nil {
		eventLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	if err := APIstub.PutState(toAccount.No, toAccountBytes); err != nil {
		eventLogger.Error(err.Error())
		return shim.Error(err.Error())
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return shim.Error(err.Error())
	}
	if err := APIstub.PutState(event.No, eventBytes); err != nil {
		eventLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	return shim.Success(eventBytes)
}

// Withdraw : withdraw from an account
func (ec *EventContract) Withdraw(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	eventLogger.Infof("invoke Withdraw, args=%s\n", args)
	if len(args) != 2 {
		errMsg := fmt.Sprintf("Incorrect number of arguments. Expecting = ['from_account_no, 'amount'], Actual = %s\n", args)
		eventLogger.Error(errMsg)
		return shim.Error(errMsg)
	}
	fromAccountNo := args[0]
	amountStr := args[1]

	amount, err := utils.GetAmount(amountStr)
	if err != nil {
		switch e := err.(type) {
		case *utils.WarningResult:
			eventLogger.Warning(err.Error())
			return shim.Success(e.JSONBytes())
		default:
			eventLogger.Error(err.Error())
			return shim.Error(err.Error())
		}
	}

	fromAccount, err := utils.GetAccount(APIstub, fromAccountNo)
	if err != nil {
		switch e := err.(type) {
		case *utils.WarningResult:
			eventLogger.Warning(err.Error())
			return shim.Success(e.JSONBytes())
		default:
			eventLogger.Error(err.Error())
			return shim.Error(err.Error())
		}
	}

	if fromAccount.Balance < amount {
		msg := fmt.Sprintf("amount is grator than the fromAccount.Balance, amount = %d, fromAccount.Balance = %d", amount, fromAccount.Balance)
		warning := &utils.WarningResult{StatusCode: 400, Message: msg}
		eventLogger.Warning(warning.Error())
		return shim.Success(warning.JSONBytes())
	}

	eventNo, err := utils.GetEventNo(APIstub)
	if err != nil {
		eventLogger.Error(err.Error())
		return shim.Error(err.Error())
	}

	fromAccountPreviousBalance := fromAccount.Balance
	fromAccount.Balance -= amount

	fromAccountState := &models.AccountState{
		No:              fromAccount.No,
		Name:            fromAccount.Name,
		PreviousBalance: fromAccountPreviousBalance,
		CurrentBalance:  fromAccount.Balance,
	}

	event := &models.Event{
		ModelType:        types.EventModel,
		EventType:        types.WithdrawEvent,
		No:               eventNo,
		Amount:           amount,
		FromAccountState: fromAccountState,
		ToAccountState:   nil,
	}

	fromAccountBytes, err := json.Marshal(fromAccount)
	if err != nil {
		eventLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	if err := APIstub.PutState(fromAccount.No, fromAccountBytes); err != nil {
		eventLogger.Error(err.Error())
		return shim.Error(err.Error())
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return shim.Error(err.Error())
	}
	if err := APIstub.PutState(event.No, eventBytes); err != nil {
		eventLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	return shim.Success(eventBytes)
}
