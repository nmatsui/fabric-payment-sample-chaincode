package utils

import (
	"encoding/json"
	"fmt"
)

// WarningResult : use this type when chaincode was invoked successfully but the expected result did not obtained
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
