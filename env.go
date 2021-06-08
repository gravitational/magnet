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
	"sync"

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

// ImportEnvFrom sets the environment importer
func ImportEnvFrom(importer EnvImporterFunc) {
	env.importer = importer
}

// EnvImporterFunc defines a function type to import environment from external source
type EnvImporterFunc func() map[string]string

// E defines a new environment variable specified with e.
// Returns the current value of the variable with precedence
// given to previously imported environment variables.
// If the variable was not previously imported and no value
// has been specified, the default is returned
func E(e EnvVar) string {
	if e.Key == "" {
		panic("key shouldn't be empty")
	}
	if e.Secret && len(e.Default) > 0 {
		panic("secrets shouldn't be embedded with defaults")
	}

	env.init()

	if v, ok := env.imported[e.Key]; ok {
		e.Value = v
	} else {
		e.Value = os.Getenv(e.Key)
	}
	env.env[e.Key] = e

	return MustGetEnv(e.Key)
}

// MustGetEnv returns the value of the environment variable given with key.
// The variable is assumed to have been registered either with E or
// imported from existing environment - otherwise the function will panic.
// For non-panicking version use GetEnv
func MustGetEnv(key string) (value string) {
	if v, ok := GetEnv(key); ok {
		return v
	}
	panic(fmt.Sprintf("Requested environment variable %q hasn't been registered", key))
}

// GetEnv returns the value of the environment variable given with key.
// The variable is assumed to have been registered either with E or
// imported from existing environment
func GetEnv(key string) (value string, exists bool) {
	env.init()
	var v EnvVar
	if v, exists = env.env[key]; !exists {
		return "", false
	}
	if v.Value != "" {
		return v.Value, true
	}
	return v.Default, true
}

// Env returns the complete environment
func Env() map[string]EnvVar {
	m := make(map[string]EnvVar, len(env.env))
	for key, value := range env.env {
		def := value.Default
		if def == "" {
			def = env.imported[key]
		}
		value.Default = def
		m[key] = value
	}
	return m
}

// ImportEnvFromMakefile invokes `make` to generate configuration for this mage script.
// The makefile target is assumed to be named `magnet-vars`.
// Assumes the makefile is named `Makefile`
//
// The script outputs a set of environment variables prefixed with `MAGNET_` which
// are used as default values for the configuration variables defined by the script.
// Any errors are ignored since this is a best-effort operation.
func ImportEnvFromMakefile() (env map[string]string) {
	cmd := exec.Command("make", "magnet-vars")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil
	}
	env, _ = ImportEnvFromReader(bytes.NewReader(out))
	return env
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

var env = environ{
	env:      make(map[string]EnvVar),
	importer: ImportEnvFromMakefile,
}

func (r *environ) init() {
	r.once.Do(func() {
		r.imported = r.importer()
	})
}

type environ struct {
	// env specifies the builder's configuration from environment
	env map[string]EnvVar
	// imported optionally specifies environment overrides
	imported map[string]string
	importer EnvImporterFunc
	once     sync.Once
}

var debianFrontend = E(EnvVar{
	Key:   "DEBIAN_FRONTEND",
	Short: "Set to noninteractive or stderr to null to enable non-interactive output",
})

var cacheDir = E(EnvVar{
	Key:     "XDG_CACHE_HOME",
	Short:   "Location to store/cache build assets",
	Default: "_build/cache",
})
