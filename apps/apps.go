package apps

import (
	"os"
)

const (
	ListenPortEnv = "LISTEN_PORT"
	LogLevelEnv   = "LOG_LEVEL"

	ListenPortValue = "3000"
	LogLevelValue   = "INFO"
)

var ReplicaID = func() string {
	if h, err := os.Hostname(); err == nil {
		return h
	}
	return "undefined"
}
