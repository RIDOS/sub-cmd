package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
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

// strToParams converts a string of "name1=value1 name2=value2" into slice of strings ["name1", "value1", "name2", "value2"].
// If input string is empty or invalid, it returns an error.
func strToParams(str string) ([]string, error) {
	if str == "" {
		return nil, errors.New("input string is empty")
	}

	var result []string
	pairs := strings.Split(str, " ")
	for _, pair := range pairs {
		parts := strings.Split(pair, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid form data: %s", str)
		}
		result = append(result, parts[0], parts[1])
	}

	return result, nil
}
