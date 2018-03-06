package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"

	"github.com/nmatsui/fabric-payment-sample-chaincode/contracts"
)

var logger = shim.NewLogger("main")

var accountContract = new(contracts.AccountContract)
var eventContract = new(contracts.EventContract)
var historyContract = new(contracts.HistoryContract)

type EntryPoint struct {
}

func (s *EntryPoint) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	logger.Info("instantiated chaincode")
	return shim.Success(nil)
}

func (s *EntryPoint) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function, args := APIstub.GetFunctionAndParameters()

	switch function {
	case "listAccount":
		return accountContract.ListAccount(APIstub, args)
	case "createAccount":
		return accountContract.CreateAccount(APIstub, args)
	case "retrieveAccount":
		return accountContract.RetrieveAccount(APIstub, args)
	case "updateAccountName":
		return accountContract.UpdateAccountName(APIstub, args)
	case "deleteAccount":
		return accountContract.DeleteAccount(APIstub, args)
	case "listEvent":
		return eventContract.ListEvent(APIstub, args)
	case "deposit":
		return eventContract.Deposit(APIstub, args)
	case "remit":
		return eventContract.Remit(APIstub, args)
	case "withdraw":
		return eventContract.Withdraw(APIstub, args)
	case "listHistory":
		return historyContract.ListHistory(APIstub, args)
	}
	msg := fmt.Sprintf("No such function. function = %s, args = %s", function, args)
	logger.Error(msg)
	return shim.Error(msg)
}

func main() {
	if err := shim.Start(new(EntryPoint)); err != nil {
		logger.Errorf("Error creating new Chaincode. Error = %s\n", err)
	}
}
