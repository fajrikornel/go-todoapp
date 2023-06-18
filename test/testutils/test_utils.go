package testutils

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

func ConstructUpdateFieldsMap(name *string, description *string) map[string]interface{} {
	fields := make(map[string]interface{})
	if name != nil && *name != "" {
		fields["name"] = name
	}

	if description != nil && *description != "" {
		fields["description"] = description
	}

	return fields
}
