// +build !ignore_autogenerated

package model

import (
	"encoding/json"
)

// MarshalJSONBinary can marshal themselves into valid JSON.
func (obj *T1) MarshalJSONBinary() ([]byte, error) {
	return json.Marshal(obj)
}

// UnmarshalJSONBinary that can unmarshal a JSON description of themselves.
// The input can be assumed to be a valid encoding of
// a JSON value. UnmarshalJSON must copy the JSON data
// if it wishes to retain the data after returning.
func (obj *T1) UnmarshalJSONBinary(data []byte) error {
	if err := json.Unmarshal(data, &obj); err != nil {
		return err
	}

	return nil
}

// String is used to print values passed as an operand
// to any format that accepts a string or to an unformatted printer
// such as Print.
func (obj *T1) String() string {
	bs, _ := obj.MarshalJSONBinary()
	return string(bs)
}
