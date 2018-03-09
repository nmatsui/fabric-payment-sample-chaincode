/*
 Package utils provides some utility functions.

 Copyright Nobuyuki Matsui<nobuyuki.matsui>.

 SPDX-License-Identifier: Apache-2.0
*/
package utils

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("utils/number_creator")

const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func getRandomString(n int, letterBytes string) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

// GetAccountNo : return a unique Account No.
func GetAccountNo(APIstub shim.ChaincodeStubInterface) (string, error) {
	var no string
	for {
		no = getRandomString(16, "0123456789")
		existing, err := APIstub.GetState(no)
		if err != nil {
			logger.Error(fmt.Sprintf("APIstub.GetState Error. error = %s\n", err))
			return "", err
		} else if existing != nil {
			logger.Warning(fmt.Sprintf("this no exists, no = %s\n", no))
		} else {
			break
		}
	}
	return no, nil
}

// GetEventNo : return a unique Event No.
func GetEventNo(APIstub shim.ChaincodeStubInterface) (string, error) {
	var no string
	for {
		no = getRandomString(16, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
		existing, err := APIstub.GetState(no)
		if err != nil {
			logger.Error(fmt.Sprintf("APIstub.GetState Error. error = %s\n", err))
			return "", err
		} else if existing != nil {
			logger.Warning(fmt.Sprintf("this no exists, no = %s\n", no))
		} else {
			break
		}
	}
	return no, nil
}
