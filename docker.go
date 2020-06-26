package magnet

import (
	"context"
	"fmt"

	"github.com/gravitational/trace"
)

type DockerConfigCommon struct {
	magnet *Magnet

	// Env are environment variables to pass to the spawned docker command
	Env map[string]string
}

type DockerConfigBuild struct {
	DockerConfigCommon
	// Always attempt to pull a newer version of the same image (Default: true)
	Pull bool
	// Compress the build context using gzip (Default: true)
	Compress bool
	// NoCache indicated to docker to avoid caching the results (Default: false)
	NoCache bool
	// Tag Name and optionally a tag in the 'name:tag' format
	Tag []string
	// BuildArgs set build-time variables
	BuildArgs map[string]string
	// Dockerfile is the path to the Dockerfile to build
	Dockerfile string
	// Target sets the target build stage to build
	Target string
	// CacheFrom is a list of images to consider as cache sources
	// https://andrewlock.net/caching-docker-layers-on-serverless-build-hosts-with-multi-stage-builds---target,-and---cache-from/
	CacheFrom []string

	// TODO: Support custom build context behaviour (IE a whitelist type approach)
	// and possibly an implementation that scans and finds all go files ignoring common directories (like .git)
	// Or possibly make it easy to stage all required files into a temp directory, and pass that as a context
}

func (m *Magnet) DockerBuild() *DockerConfigBuild {
	return &DockerConfigBuild{
		DockerConfigCommon: DockerConfigCommon{
			magnet: m,
			Env: map[string]string{
				"DOCKER_BUILDKIT":   "1",
				"PROGRESS_NO_TRUNC": "1",
			},
		},
		Pull:     true,
		Compress: true,
	}
}

func (m *DockerConfigBuild) AddTag(tag string) *DockerConfigBuild {
	m.Tag = append(m.Tag, tag)
	return m
}

func (m *DockerConfigBuild) AddCacheFrom(from string) *DockerConfigBuild {
	m.CacheFrom = append(m.CacheFrom, from)
	return m
}

func (m *DockerConfigBuild) SetBuildArg(key, value string) *DockerConfigBuild {
	if m.BuildArgs == nil {
		m.BuildArgs = make(map[string]string)
	}

	m.BuildArgs[key] = value

	return m
}

func (m *DockerConfigBuild) SetEnv(key, value string) *DockerConfigBuild {
	if m.Env == nil {
		m.Env = make(map[string]string)
	}

	m.Env[key] = value

	return m
}

func (m *DockerConfigBuild) SetEnvs(envs map[string]string) *DockerConfigBuild {
	if m.Env == nil {
		m.Env = make(map[string]string)
	}

	for key, value := range envs {
		m.Env[key] = value
	}

	return m
}

func (m *DockerConfigBuild) SetDockerfile(dockerfile string) *DockerConfigBuild {
	m.Dockerfile = dockerfile
	return m
}

func (m *DockerConfigBuild) SetPull(pull bool) *DockerConfigBuild {
	m.Pull = pull
	return m
}

func (m *DockerConfigBuild) SetCompress(compress bool) *DockerConfigBuild {
	m.Compress = compress
	return m
}

func (m *DockerConfigBuild) SetNoCache(nocache bool) *DockerConfigBuild {
	m.NoCache = nocache
	return m
}

func (m *DockerConfigBuild) SetTarget(target string) *DockerConfigBuild {
	m.Target = target
	return m
}

func (m *DockerConfigBuild) Build(ctx context.Context, contextPath string) error {
	args := []string{"build"}

	if m.Pull {
		args = append(args, "--pull")
	}

	if m.Compress {
		args = append(args, "--compress")
	}

	if m.NoCache {
		args = append(args, "--no-cache")
	}

	if len(m.Target) > 0 {
		args = append(args, "--target", m.Target)
	}

	for key, value := range m.BuildArgs {
		args = append(args, "--build-arg", fmt.Sprint(key, "=", value))
	}

	for _, value := range m.CacheFrom {
		args = append(args, "--cache-from", value)
	}

	for _, value := range m.Tag {
		args = append(args, "-t", value)
	}

	if len(m.Dockerfile) > 0 {
		args = append(args, "-f", m.Dockerfile)
	}

	args = append(args, contextPath)

	_, err := m.magnet.Exec().SetEnvs(m.Env).Run(ctx, "docker", args...)

	return trace.Wrap(err)
}

type DockerConfigRun struct {
	DockerConfigCommon

	// Eun container in background
	Detach bool

	// User ID of spawned process
	UID string
	// Group ID of spawned process
	GID string

	// Privileged Give extended privileges to the container
	Privileged bool
	// ReadOnly mounts the containers root filesystem as read only
	ReadOnly bool
	// Automatically remove the container when it exits
	Remove bool
	// Volumes is a list of volumes to bind mount
	Volumes []string
	// Workdir sets the working directory inside the container
	WorkDir string
}

func (m *Magnet) DockerRun() *DockerConfigRun {
	return &DockerConfigRun{
		DockerConfigCommon: DockerConfigCommon{
			magnet: m,
		},
	}
}

func (m *DockerConfigRun) SetEnv(key, value string) *DockerConfigRun {
	if m.Env == nil {
		m.Env = make(map[string]string)
	}

	m.Env[key] = value

	return m
}

func (m *DockerConfigRun) SetDetach(detach bool) *DockerConfigRun {
	m.Detach = detach
	return m
}

func (m *DockerConfigRun) SetUID(uid string) *DockerConfigRun {
	m.UID = uid
	return m
}

func (m *DockerConfigRun) SetGID(gid string) *DockerConfigRun {
	m.GID = gid
	return m
}

func (m *DockerConfigRun) SetPrivileged(privileged bool) *DockerConfigRun {
	m.Privileged = privileged
	return m
}

func (m *DockerConfigRun) SetReadonly(readonly bool) *DockerConfigRun {
	m.ReadOnly = readonly
	return m
}

func (m *DockerConfigRun) SetRemove(remove bool) *DockerConfigRun {
	m.Remove = remove
	return m
}

func (m *DockerConfigRun) AddVolume(volume string) *DockerConfigRun {
	m.Volumes = append(m.Volumes, volume)
	return m
}

func (m *DockerConfigRun) SetWorkDir(workdir string) *DockerConfigRun {
	m.WorkDir = workdir
	return m
}

func (m *DockerConfigRun) Run(ctx context.Context, image, cmd string, cargs ...string) error {
	args := []string{"run"}

	if m.Detach {
		args = append(args, "-d")
	}

	if len(m.UID) > 0 {
		if len(m.GID) > 0 {
			args = append(args, "-u", fmt.Sprint(m.UID, ":", m.GID))
		} else {
			args = append(args, "-u", m.UID)
		}
	}

	if m.Privileged {
		args = append(args, "--privileged")
	}

	if m.ReadOnly {
		args = append(args, "--read-only")
	}

	if m.Remove {
		args = append(args, "--rm=true")
	}

	if len(m.WorkDir) > 0 {
		args = append(args, "-w", m.WorkDir)
	}

	for _, value := range m.Volumes {
		args = append(args, "-v", value)
	}

	for key, value := range m.Env {
		args = append(args, fmt.Sprintf("--env=%v=%v", key, value))
	}

	args = append(args, image)
	args = append(args, cmd)
	args = append(args, cargs...)

	_, err := m.magnet.Exec().Run(ctx, "docker", args...)

	return trace.Wrap(err)
}
