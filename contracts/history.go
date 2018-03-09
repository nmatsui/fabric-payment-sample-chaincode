/*
 Package contracts provides the smart contracts for Hyperledger/fabric 1.1.

 Copyright Nobuyuki Matsui<nobuyuki.matsui>.

 SPDX-License-Identifier: Apache-2.0
*/
package contracts

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

var historyLogger = shim.NewLogger("contracts/history")

type historyType struct {
	TxID      string                 `json:"tx_id"`
	No        string                 `json:"no"`
	State     map[string]interface{} `json:"state"`
	Timestamp string                 `json:"timestamp"`
	IsDelete  bool                   `json:"is_delete"`
}

// HistoryContract : a struct to query Histories
type HistoryContract struct {
}

// ListHistory : return all histories of a state object.
func (hc *HistoryContract) ListHistory(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	historyLogger.Infof("invoke ListHistory, args=%s\n", args)
	if len(args) != 1 {
		errMsg := fmt.Sprintf("Incorrect number of arguments. Expecting = ['no'], Actual = %s\n", args)
		historyLogger.Error(errMsg)
		return shim.Error(errMsg)
	}
	no := args[0]

	resultsIterator, err := APIstub.GetHistoryForKey(no)
	if err != nil {
		historyLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	histories := make([]*historyType, 0)
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			historyLogger.Error(err.Error())
			return shim.Error(err.Error())
		}
		var state map[string]interface{}
		if !response.IsDelete {
			if err := json.Unmarshal(response.Value, &state); err != nil {
				historyLogger.Error(err.Error())
				return shim.Error(err.Error())
			}
		}
		history := &historyType{
			TxID:      response.TxId,
			No:        no,
			State:     state,
			Timestamp: time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String(),
			IsDelete:  response.IsDelete,
		}
		histories = append(histories, history)
	}
	jsonBytes, err := json.Marshal(histories)
	if err != nil {
		historyLogger.Error(err.Error())
		return shim.Error(err.Error())
	}
	return shim.Success(jsonBytes)
}
