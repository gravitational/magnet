package magnet

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/containerd/console"
	"github.com/gravitational/trace"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/util/progress/progressui"
	"github.com/opencontainers/go-digest"
	"github.com/sirupsen/logrus"
)

var BuildLogDir string

func init() {
	BuildLogDir = filepath.Join("build/logs", time.Now().Format("20060102150405"))
}

type Magnet struct {
	logrus.FieldLogger

	Vertex *client.Vertex
	parent *Magnet
	status chan *client.SolveStatus
}

var root *Magnet
var rootOnce sync.Once

func Root() *Magnet {
	const statusChanSize = 128

	rootOnce.Do(func() {
		now := time.Now()
		root = &Magnet{
			status: make(chan *client.SolveStatus, statusChanSize),
			Vertex: &client.Vertex{
				Digest:    digest.FromString("root"),
				Name:      fmt.Sprint("Logs: ", BuildLogDir),
				Started:   &now,
				Completed: &now,
			},
		}
		root.FieldLogger = root.newLogger("root", digest.FromString("root"))
	})

	return root
}

// Shutdown indicates that the program is exiting, and we should shutdown the progressui
//  if it's currently running
func Shutdown() {
	if root != nil {
		// Hack: give progressui enough time to process any queues status updates
		time.Sleep(time.Second)
		close(root.status)
	}

	// Hack: give progressui enough time to shutdown
	time.Sleep(10 * time.Second)
	fmt.Fprintln(os.Stdout, "Shutdown complete")
}

func (m *Magnet) Clone(name string) *Magnet {
	started := time.Now()
	vertex := &client.Vertex{
		Digest:  digest.FromString(name),
		Name:    name,
		Started: &started,
	}

	status := &client.SolveStatus{
		Vertexes: []*client.Vertex{vertex},
	}

	select {
	case m.root().status <- status:
	default:
		fmt.Fprintln(os.Stderr, "dropped SolveStatus on Clone")
	}

	// only create loggers when cloned from root, so parallel tasks get there own log output
	// TODO: We need a better way to indicate and name individual logs
	var logger logrus.FieldLogger
	logger = m.FieldLogger.WithField("vertex", vertex.Digest.Encoded()[:5])
	if m.parent == nil {
		logger = m.newLogger(name, vertex.Digest)
	}

	return &Magnet{
		FieldLogger: logger,
		Vertex:      vertex,
		parent:      m,
	}
}

func (m *Magnet) InitOutput() error {
	if m.parent != nil {
		return trace.BadParameter("Expect output to only be run on the root of the graph")
	}

	var c console.Console
	//if cn, err := console.ConsoleFromFile(os.Stderr); err == nil {
	//	c = cn
	//}

	return trace.Wrap(progressui.DisplaySolveStatus(context.TODO(), m.Vertex.Name, c, os.Stdout, m.status))
}

func (m *Magnet) root() *Magnet {
	root := m
	for root.parent != nil {
		root = root.parent
	}

	return root
}

// Complete marks the current task as complete.
func (m *Magnet) Complete(cached bool, err error) {
	now := time.Now()
	m.Vertex.Completed = &now
	m.Vertex.Cached = cached
	m.Vertex.Error = trace.DebugReport(err)

	m.root().status <- &client.SolveStatus{
		Vertexes: []*client.Vertex{
			m.Vertex,
		},
	}
}

func (m *Magnet) newLogger(name string, vertex digest.Digest) logrus.FieldLogger {
	err := os.MkdirAll(BuildLogDir, 0755)
	if err != nil {
		panic(trace.DebugReport(trace.ConvertSystemError(err)))
	}

	filename := filepath.Join(BuildLogDir, fmt.Sprintf("%v.%v", name, vertex.Encoded()))
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		panic(trace.DebugReport(trace.ConvertSystemError(err)))
	}

	out := io.MultiWriter(f, &streamWriter{
		stream: STDOUT,
		vertex: vertex,
		status: m.root().status,
	})

	logger := logrus.New()
	logger.SetOutput(out)

	return logger.WithField("vertex", vertex.Encoded()[:5])
}
