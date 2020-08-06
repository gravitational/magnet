package magnet

import (
	"context"
	"os"
	"strings"

	"github.com/gravitational/trace"
)

var EnvVars map[string]EnvVar
var ImportEnvVars map[string]string

type EnvVar struct {
	Key     string
	Value   string
	Default string
	Short   string
	Long    string
	Secret  bool
}

func E(e EnvVar) string {
	if e.Key == "" {
		panic("Key shouldn't be empty")
	}

	if e.Secret && len(e.Default) > 0 {
		panic("Secrets shouldn't be embedded with defaults")
	}

	if EnvVars == nil {
		EnvVars = make(map[string]EnvVar)
		ImportEnvVars = make(map[string]string)
	}
	e.Value = os.Getenv(e.Key)

	EnvVars[e.Key] = e

	return GetEnv(e.Key)
}

func GetEnv(key string) string {
	if EnvVars == nil {
		EnvVars = make(map[string]EnvVar)
		ImportEnvVars = make(map[string]string)
	}

	if v, ok := EnvVars[key]; ok {
		if v.Value != "" {
			return v.Value
		}

		if v, ok := ImportEnvVars[key]; ok {
			return v
		}

		return v.Default
	}

	panic("Requested environment variable hasn't been registered")
}

func importMakeEnv(target string) error {
	out, err := Output(context.TODO(), "make", target)
	if err != nil {
		return trace.Wrap(err)
	}

	lines := strings.Split(out, "\n")
	for _, line := range lines {
		cols := strings.SplitN(line, "=", 2)
		ImportEnvVars[cols[0]] = cols[1]
	}

	return nil
}
