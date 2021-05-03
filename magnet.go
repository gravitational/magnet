package magnet

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/mod/modfile"

	"github.com/gravitational/magnet/pkg/progressui"
	"github.com/gravitational/trace"

	"github.com/containerd/console"
	"github.com/opencontainers/go-digest"
)

// Config defines logger's configuration
type Config struct {
	// LogDir optionally specifies the logging directory root
	LogDir string
	// CacheDir optionally specifies the artifact cache directory root
	CacheDir string

	// ModulePath specifies the path of the Go module being built
	ModulePath string
	// Version specifies the module version
	Version string

	// PrintConfig configures whether magnet will output its configuration
	PrintConfig bool
	// PlainProgress specifies whether the logger uses fancy progress reporting.
	// Set to true to see streaming output (e.g. on CI)
	PlainProgress bool
	// ImportEnv optionally specifies the external configuration as a set of
	// environment variables
	ImportEnv map[string]string
}

func (c *Config) checkAndSetDefaults() error {
	if c.Version == "" {
		c.Version = DefaultVersion()
	}

	if c.LogDir == "" {
		c.LogDir = DefaultLogDir()
	}

	if c.ModulePath == "" {
		buf, err := ioutil.ReadFile("go.mod")
		// TODO (knisbet), silently discarding the error for now
		// detection only works if using go modules
		if err == nil {
			c.ModulePath = modfile.ModulePath(buf)
		} else {
			wd, err := os.Getwd()
			if err != nil {
				return trace.Wrap(err)
			}

			gopath := os.Getenv("GOPATH")
			if gopath != "" {
				c.ModulePath, err = filepath.Rel(filepath.Join(gopath, "src"), wd)
				if err != nil {
					return trace.Wrap(err, "invalid working directory %s in GOPATH mode", wd)
				}
			}
		}
	}
	return nil
}

func (m *Magnet) printHeader() {
	fmt.Printf("Logs:    %v (%v)\n", m.statusLogger.dirLink(), m.statusLogger.dirReal())

	if m.Version != "" {
		fmt.Println("Version: ", m.Version)
	}
	if m.CacheDir != "" {
		fmt.Println("Cache:   ", m.cacheDir())
	}
}

// Magnet describes the root logger
type Magnet struct {
	Config

	status       chan *progressui.SolveStatus
	statusLogger *SolveStatusLogger
	root         MagnetTarget

	env map[string]EnvVar

	wg  sync.WaitGroup
	ctx context.Context
	// cancel cancels the logger process
	cancel         context.CancelFunc
	initOutputOnce sync.Once
}

// MagnetTarget describes a child logging target
type MagnetTarget struct {
	root   *Magnet
	vertex *progressui.Vertex
	cached bool
}

// Root creates a root vertex for executing and capturing status of each build target.
func Root(c Config) (*Magnet, error) {
	if err := c.checkAndSetDefaults(); err != nil {
		return nil, trace.Wrap(err)
	}

	statusLogger, err := newSolveStatusLogger(c.LogDir)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	now := time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	root := &Magnet{
		Config: c,
		root: MagnetTarget{
			vertex: &progressui.Vertex{
				Digest:    digest.FromString("root"),
				Started:   &now,
				Completed: &now,
			},
		},
		status:       statusLogger.source,
		statusLogger: statusLogger,
		env:          make(map[string]EnvVar),
		ctx:          ctx,
		cancel:       cancel,
	}
	// :-)
	root.root.root = root
	return root, nil
}

func newSecretsRedactor(env map[string]EnvVar) secretsRedactor {
	var secrets []EnvVar
	for _, value := range env {
		if value.Secret {
			secrets = append(secrets, value)
		}
	}
	return secretsRedactor{secrets: secrets}
}

// secretsRedactor redacts literal secrets in a text stream.
// Implements redactor
type secretsRedactor struct {
	secrets []EnvVar
}

func (r secretsRedactor) redact(s string) string {
	for _, secret := range r.secrets {
		if len(secret.Value) > 0 {
			s = strings.ReplaceAll(s, secret.Value, "<redacted>")
		}
	}
	return s
}

// Shutdown indicates that the program is exiting, and we should shutdown the progressui
//  if it's currently running
func (m *Magnet) Shutdown() {
	close(m.status)
	m.cancel()
	m.wg.Wait()
}

func (m *Magnet) Target(name string) *MagnetTarget {
	m.initOutput()
	return m.root.newTarget(&progressui.Vertex{
		Digest: digest.FromString(name),
		Name:   name,
	})
}

func (m *MagnetTarget) Target(name string) *MagnetTarget {
	m.root.initOutput()
	return m.newTarget(&progressui.Vertex{
		Digest: digest.FromString(name),
		Name:   name,
		// the root vertex doesn't get fully added to the progress ui. So only add a parent if we're not root
		Inputs: []digest.Digest{m.vertex.Digest},
	})
}

func (m *MagnetTarget) newTarget(vertex *progressui.Vertex) *MagnetTarget {
	now := time.Now()
	vertex.Started = &now

	status := &progressui.SolveStatus{
		Vertexes: []*progressui.Vertex{vertex},
	}

	m.root.status <- status

	return &MagnetTarget{
		vertex: vertex,
		root:   m.root,
	}
}

// AbsCacheDir is the configured cache directory as an absolute path.
func (c Config) AbsCacheDir() (path string, err error) {
	if filepath.IsAbs(c.cacheDir()) {
		return c.cacheDir(), nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", trace.Wrap(err)
	}
	return filepath.Join(wd, c.cacheDir()), nil
}

// initOutput starts the internal progress logging process
func (m *Magnet) initOutput() {
	m.initOutputOnce.Do(func() {
		redactor := newSecretsRedactor(m.env)
		m.statusLogger.start(redactor)

		if m.PrintConfig {
			m.printHeader()
		}

		var c console.Console

		if !m.PlainProgress {
			if cn, err := console.ConsoleFromFile(os.Stderr); err == nil {
				c = cn
			}
		}

		m.wg.Add(1)
		go func() {
			progressui.DisplaySolveStatus(
				m.ctx,
				m.root.vertex.Name,
				c,
				os.Stdout,
				m.statusLogger.destination,
			)
			m.wg.Done()
		}()
	})
}

func (c Config) cacheDir() string {
	return filepath.Join(c.CacheDir, "magnet", c.ModulePath)
}

// Complete marks the current task as complete.
func (m *MagnetTarget) Complete(err error) {
	now := time.Now()
	m.vertex.Completed = &now
	m.vertex.Cached = m.cached
	m.vertex.Error = trace.DebugReport(err)
	m.root.status <- &progressui.SolveStatus{
		Vertexes: []*progressui.Vertex{
			m.vertex,
		},
	}
}

// SetCached marks the current task as cached when it's completed.
func (m *MagnetTarget) SetCached(cached bool) {
	m.cached = cached
}

// Println allows writing log entries to the log output for the target.
func (m *MagnetTarget) Println(args ...interface{}) {
	msg := fmt.Sprintln(args...)

	m.root.status <- &progressui.SolveStatus{
		Logs: []*progressui.VertexLog{
			{
				Vertex:    m.vertex.Digest,
				Stream:    STDOUT,
				Data:      []byte(msg),
				Timestamp: time.Now(),
			},
		},
	}
}

// Printlnf allows writing log entries to the log output for the target.
func (m *MagnetTarget) Printlnf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	m.Println(msg)
}
