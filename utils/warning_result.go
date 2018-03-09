/*
 Package utils provides some utility functions.

 Copyright Nobuyuki Matsui<nobuyuki.matsui>.

 SPDX-License-Identifier: Apache-2.0
*/
package utils

import (
	"encoding/json"
	"fmt"
)

// WarningResult : a struct implements Error Interface.
//    use this type when chaincode was invoked successfully but the expected result did not obtained
type WarningResult struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

// Error : error interface
func (e *WarningResult) Error() string {
	return fmt.Sprintf("Error: StatusCode=%d, Message=%s\n", e.StatusCode, e.Message)
}

// JSONBytes : return error json bytes
func (e *WarningResult) JSONBytes() []byte {
	errBytes, _ := json.Marshal(e)
	return errBytes
}
