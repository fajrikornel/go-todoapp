package test_utils

import (
	"bytes"
	"encoding/json"
	"io"
)

func ConvertStructToIoReader(body interface{}) io.Reader {
	bodyBytes, _ := json.Marshal(body)
	return bytes.NewReader(bodyBytes)
}

func CreatePointerOfString(s string) *string {
	sPointer := &s
	return sPointer
}
