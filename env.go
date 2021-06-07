package magnet

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/gravitational/trace"
)

// EnvVar represents a configuration with optional defaults
// obtained from environment
type EnvVar struct {
	Key     string
	Value   string
	Default string
	Short   string
	Long    string
	Secret  bool
}

// E defines a new environment variable specified with e.
// Returns the current value of the variable with precedence
// given to previously imported environment variables.
// If the variable was not previously imported and no value
// has been specified, the default is returned
func (m *Magnet) E(e EnvVar) string {
	if e.Key == "" {
		panic("Key shouldn't be empty")
	}

	if e.Secret && len(e.Default) > 0 {
		panic("Secrets shouldn't be embedded with defaults")
	}

	if v, ok := m.ImportEnv[e.Key]; ok {
		e.Value = v
	} else {
		e.Value = os.Getenv(e.Key)
	}
	m.env[e.Key] = e

	return m.MustGetEnv(e.Key)
}

// MustGetEnv returns the value of the environment variable given with key.
// The variable is assumed to have been registered either with E or
// imported from existing environment - otherwise the function will panic.
// For non-panicking version use GetEnv
func (m *Magnet) MustGetEnv(key string) (value string) {
	if v, ok := m.GetEnv(key); ok {
		return v
	}
	panic(fmt.Sprintf("Requested environment variable %q hasn't been registered", key))
}

// GetEnv returns the value of the environment variable given with key.
// The variable is assumed to have been registered either with E or
// imported from existing environment
func (m *Magnet) GetEnv(key string) (value string, exists bool) {
	var v EnvVar
	if v, exists = m.env[key]; !exists {
		return "", false
	}
	if v.Value != "" {
		return v.Value, true
	}
	return v.Default, true
}

// Env returns the complete environment
func (m *Magnet) Env() map[string]EnvVar {
	return m.env
}

// ImportEnvFromMakefile invokes `make` to generate configuration for this mage script.
// The makefile target is assumed to be named `magnet-vars`.
// The script outputs a set of environment variables prefixed with `MAGNET_` which
// are used as default values for the configuration variables defined by the script.
// Assumes the Makefile is named `Makefile`
func ImportEnvFromMakefile() (env map[string]string, err error) {
	cmd := exec.Command("make", "magnet-vars")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Failed to import environ from makefile: %v", err)
		return nil, trace.Wrap(err)
	}
	return ImportEnvFromReader(bytes.NewReader(out))
}

// ImportEnvFromReader consumes configuration for this mage script from the specified reader.
// Expects the reader to produce a list of environment variables as key=value pairs with a single
// variable per line.
// Only the environment variables prefixed with `MAGNET_` are considered which
// are used as default values for the configuration variables defined by the script itself.
func ImportEnvFromReader(r io.Reader) (env map[string]string, err error) {
	env = make(map[string]string)

	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}
		cols := strings.SplitN(line, "=", 2)
		if len(cols) != 2 || !strings.HasPrefix(cols[0], "MAGNET_") {
			log.Printf("Skip line that does not look like magnet envar: %q\n", line)
			continue
		}
		key, value := strings.TrimPrefix(cols[0], "MAGNET_"), cols[1]
		env[key] = value
	}
	if s.Err() != nil {
		return nil, trace.Wrap(s.Err())
	}
	return env, nil
}
