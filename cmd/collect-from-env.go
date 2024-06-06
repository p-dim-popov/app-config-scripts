package cmd

import (
	"encoding/json"
	"github.com/mdaverde/jsonpath"
	"os"
	"strings"
)

func CollectFromEnvAsString() string {
	var marshalled, err = json.Marshal(CollectFromEnv())
	if err != nil {
		panic(err)
	}

	return string(marshalled)
}

func CollectFromEnv() map[string]any {
	var collected map[string]interface{}
	err := json.Unmarshal([]byte(os.Getenv("APP_CONFIG_JSON")), &collected)
	if err != nil {
		collected = make(map[string]interface{})
	}

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		var (
			key   = parts[0]
			value = parts[1]
		)

		separatedKey, found := strings.CutPrefix(key, "APP_CONFIG_")
		if !found || separatedKey == "JSON" {
			continue
		}

		dottedPath := strings.ReplaceAll(separatedKey, "_", ".")

		var unmarshalledValue interface{}
		err = json.Unmarshal([]byte(value), &unmarshalledValue)
		if err != nil {
			unmarshalledValue = value
		}

		err = jsonpath.Set(&collected, dottedPath, unmarshalledValue)
		if err != nil {
			panic(err)
		}
	}

	return collected
}
