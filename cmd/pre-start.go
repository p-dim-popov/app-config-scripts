package cmd

import (
	"dario.cat/mergo"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"strings"
)

func PreStart(cmdArgs []string) {
	command, err := newFromCmdArgs(cmdArgs)
	if err != nil {
		panic(err)
	}

	collectedConfig := CollectFromEnv()

	if command.runningConfig == nil {
		command.runningConfig = collectedConfig
	} else {
		err := mergo.Merge(&command.runningConfig, collectedConfig)
		if err != nil {
			panic(err)
		}
	}

	command.start()
}

type preStartCommand struct {
	runningConfig map[string]any
	args          []string
}

func newFromCmdArgs(cmdArgs []string) (*preStartCommand, error) {
	if len(cmdArgs) < 1 {
		return nil, errors.New("Invalid number of args: ")
	}

	if !strings.HasPrefix(cmdArgs[0], "{") {
		return &preStartCommand{
			args: cmdArgs,
		}, nil
	}

	var baseConfig map[string]any
	if err := json.Unmarshal([]byte(cmdArgs[0]), &baseConfig); err != nil {
		return nil, err
	}

	return &preStartCommand{
		runningConfig: baseConfig,
		args:          cmdArgs[1:],
	}, nil
}

func (cmd *preStartCommand) start() {
	marshalledConfig, err := json.Marshal(cmd.runningConfig)
	if err != nil {
		panic(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	childCmd := exec.Command("sh", "-c", strings.Join(cmd.args, " "))
	childCmd.Dir = cwd
	childCmd.Stdin = os.Stdin
	childCmd.Stdout = os.Stdout
	childCmd.Stderr = os.Stderr

	err = os.Setenv("APP_CONFIG", string(marshalledConfig))
	if err != nil {
		panic(err)
	}
	childCmd.Env = nil

	err = childCmd.Start()
	if err != nil {
		panic(err)
	}
	err = childCmd.Process.Release()
	if err != nil {
		panic(err)
	}
}
