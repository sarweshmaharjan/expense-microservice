package services

import (
	"io"
	"os"

	"github.com/sarweshmaharjan/json_builder/model"
	"gopkg.in/yaml.v3"
)

func Load() model.FinancialConfig {
	jsonFile, err := os.Open("../config/preferences.yml") // Adjust the path as needed
	if err != nil {
		panic("Config file not found: " + err.Error())
	}
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	// Process the config file
	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	var result model.FinancialConfig
	err = yaml.Unmarshal(bytes, &result)
	if err != nil {
		panic(err)
	}
	return result
}
