package config

import (
	"encoding/json"
	"os"
)

func Read(filePath string) (Configuration, error) {
	var configuration Configuration
	configText, err := os.ReadFile(filePath)
	if err != nil {
		return configuration, err
	}

	err = json.Unmarshal(configText, &configuration)
	if err != nil {
		return configuration, err
	}

	return configuration, nil
}
