package cmd

import (
	"dario.cat/mergo"
	"encoding/json"
	"os"
)

func PreServe(distDir string) {
	appConfigVarName := "__app_config__"
	appConfigJsonFile := distDir + "/" + appConfigVarName + ".json"

	var appConfig string

	baseConfig, err := os.ReadFile(appConfigJsonFile)
	if err == nil {
		var unmarshalled map[string]any

		if err := json.Unmarshal(baseConfig, &unmarshalled); err != nil {
			panic(err)
		}

		if err = mergo.Merge(&unmarshalled, CollectFromEnv()); err != nil {
			panic(err)
		}

		marshalled, err := json.Marshal(unmarshalled)
		if err != nil {
			panic(err)
		}

		appConfig = string(marshalled)
	} else if os.IsNotExist(err) {
		appConfig = CollectFromEnvAsString()
	} else {
		panic(err)
	}

	println("Applying app config,", appConfig)

	err = os.WriteFile(appConfigJsonFile, []byte(appConfig), 0644)
	if err != nil {
		panic(err)
	}

	appConfigJsFile := distDir + "/" + appConfigVarName + ".js"
	appConfigJsFileContents := "var " + appConfigVarName + " = " + appConfig

	if err = os.WriteFile(appConfigJsFile, []byte(appConfigJsFileContents), 0644); err != nil {
		panic(err)
	}
}
