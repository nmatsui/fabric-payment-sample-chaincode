package contracts

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"

	"github.com/nmatsui/fabric-payment-sample/models"
	"github.com/nmatsui/fabric-payment-sample/types"
	"github.com/nmatsui/fabric-payment-sample/utils"
)

var logger = shim.NewLogger("contracts/account")

// AccountContract : a struct which has the methods related to manage Account
type AccountContract struct {
}

// ListAccount : return a list of all accounts
func (ac *AccountContract) ListAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	logger.Infof("invoke ListAccount, args=%s\n", args)
	if len(args) != 0 {
		errMsg := fmt.Sprintf("Incorrect number of arguments. Expecting = no argument, Actual = %s\n", args)
		logger.Error(errMsg)
		return shim.Error(errMsg)
	}

	query := map[string]interface{}{
		"selector": map[string]interface{}{
			"model_type": types.AccountModel,
		},
	}

	queryBytes, err := json.Marshal(query)
	if err != nil {
		logger.Error(err.Error())
		return shim.Error(err.Error())
	}
	logger.Infof("Query string = '%s'", string(queryBytes))
	resultsIterator, err := APIstub.GetQueryResult(string(queryBytes))
	if err != nil {
		logger.Error(err.Error())
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	results := make([]*models.Account, 0)
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			logger.Error(err.Error())
			return shim.Error(err.Error())
		}
		var account = new(models.Account)
		if err := json.Unmarshal(queryResponse.Value, account); err != nil {
			logger.Error(err.Error())
			return shim.Error(err.Error())
		}
		results = append(results, account)
	}
	jsonBytes, err := json.Marshal(results)
	if err != nil {
		logger.Error(err.Error())
		return shim.Error(err.Error())
	}
	return shim.Success(jsonBytes)
}

// CreateAccount : create a new account
func (ac *AccountContract) CreateAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	logger.Infof("invoke CreateAccount, args=%s\n", args)
	if len(args) != 1 {
		errMsg := fmt.Sprintf("Incorrect number of arguments. Expecting = ['name'], Actual = %s\n", args)
		logger.Error(errMsg)
		return shim.Error(errMsg)
	}
	name := args[0]

	no, err := utils.GetAccountNo(APIstub)
	if err != nil {
		logger.Error(err.Error())
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
		logger.Error(err.Error())
		return shim.Error(err.Error())
	}
	if err := APIstub.PutState(no, jsonBytes); err != nil {
		logger.Error(err.Error())
		return shim.Error(err.Error())
	}
	return shim.Success(jsonBytes)
}

// RetrieveAccount : return an account
func (ac *AccountContract) RetrieveAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	logger.Infof("invoke RetrieveAccount, args=%s\n", args)
	if len(args) != 1 {
		errMsg := fmt.Sprintf("Incorrect number of arguments. Expecting = ['no'], Actual = %s\n", args)
		logger.Error(errMsg)
		return shim.Error(errMsg)
	}
	no := args[0]

	account, err := getAccount(APIstub, no)
	if err != nil {
		switch e := err.(type) {
		case *utils.WarningResult:
			logger.Warning(err.Error())
			return shim.Success(e.JSONBytes())
		default:
			logger.Error(err.Error())
			return shim.Error(err.Error())
		}
	}

	jsonBytes, err := json.Marshal(account)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(jsonBytes)
}

// UpdateAccountName : update the name of an account
func (ac *AccountContract) UpdateAccountName(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	logger.Infof("invoke UpdateAccountName, args=%s\n", args)
	if len(args) != 2 {
		errMsg := fmt.Sprintf("Incorrect number of arguments. Expecting = ['no', 'name'], Actual = %s\n", args)
		logger.Error(errMsg)
		return shim.Error(errMsg)
	}
	no := args[0]
	name := args[1]

	account, err := getAccount(APIstub, no)
	if err != nil {
		switch e := err.(type) {
		case *utils.WarningResult:
			logger.Warning(err.Error())
			return shim.Success(e.JSONBytes())
		default:
			logger.Error(err.Error())
			return shim.Error(err.Error())
		}
	}

	account.Name = name

	jsonBytes, err := json.Marshal(account)
	if err != nil {
		logger.Error(err.Error())
		return shim.Error(err.Error())
	}
	if err := APIstub.PutState(no, jsonBytes); err != nil {
		logger.Error(err.Error())
		return shim.Error(err.Error())
	}
	return shim.Success(jsonBytes)
}

// DeleteAccount : delete an account
func (ac *AccountContract) DeleteAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	logger.Infof("invoke DeleteAccount, args=%s\n", args)
	if len(args) != 1 {
		errMsg := fmt.Sprintf("Incorrect number of arguments. Expecting = ['no'], Actual = %s\n", args)
		logger.Error(errMsg)
		return shim.Error(errMsg)
	}
	no := args[0]

	_, err := getAccount(APIstub, no)
	if err != nil {
		switch e := err.(type) {
		case *utils.WarningResult:
			logger.Warning(err.Error())
			return shim.Success(e.JSONBytes())
		default:
			logger.Error(err.Error())
			return shim.Error(err.Error())
		}
	}

	if err := APIstub.DelState(no); err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func getAccount(APIstub shim.ChaincodeStubInterface, no string) (*models.Account, error) {
	var account = new(models.Account)
	accountBytes, err := APIstub.GetState(no)
	if err != nil {
		return account, err
	} else if accountBytes == nil {
		msg := fmt.Sprintf("Account does not exist, no = %s", no)
		warning := &utils.WarningResult{StatusCode: 404, Message: msg}
		return account, warning
	}
	if err := json.Unmarshal(accountBytes, account); err != nil {
		return account, err
	}
	return account, nil
}
