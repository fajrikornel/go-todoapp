package test_utils

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func FormatNameAndDescription(name, description *string) string {
	nameString := "nil"
	if name != nil {
		nameString = *name
	}
	descriptionString := "nil"
	if description != nil {
		descriptionString = *description
	}

	return fmt.Sprintf("name:%s and description:%s", nameString, descriptionString)
}
