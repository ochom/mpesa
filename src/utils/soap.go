package utils

import (
	"encoding/xml"
	"strings"
)

func ParseXml[T any](data string) (*T, error) {
	var res T

	decoder := xml.NewDecoder(strings.NewReader(data))
	err := decoder.Decode(&res)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
