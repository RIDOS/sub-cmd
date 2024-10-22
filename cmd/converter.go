package cmd

import (
	"encoding/json"
	"errors"
)

func parseAndReformatJson(jsonData string) ([]byte, error) {
	if jsonData == "" {
		return nil, errors.New("input JSON string is empty")
	}

	var parsedData map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &parsedData)
	if err != nil {
		return nil, err
	}

	formattedJson, err := json.Marshal(parsedData)
	if err != nil {
		return nil, err
	}

	return formattedJson, nil
}
