package main

import (
	"app-config-scripts/cmd"
	"fmt"
	"os"
)

const CmdCollectFromEnv = "collect-from-env"
const CmdPreStart = "pre-start"
const CmdPreServe = "pre-serve"

// TODO: extract error factories?
func main() {
	if len(os.Args) < 2 {
		panic(
			fmt.Sprintf(
				"Please specify command: '%s', '%s', '%s'",
				CmdCollectFromEnv,
				CmdPreStart,
				CmdPreServe,
			),
		)
	}

	command := os.Args[1]

	switch command {
	case CmdCollectFromEnv:
		fmt.Println(cmd.CollectFromEnvAsString())
		return
	case CmdPreStart:
		cmd.PreStart(os.Args[2:])
		return
	case CmdPreServe:
		cmd.PreServe(os.Args[2])
		return
	default:
		panic(
			fmt.Sprintf(
				"Unknown command '%s'. Available commands: '%s', '%s', '%s'",
				command,
				CmdCollectFromEnv,
				CmdPreStart,
				CmdPreServe,
			),
		)
	}
}
