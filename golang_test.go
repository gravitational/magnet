package magnet

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildsContainerPackagePathFromGopath(t *testing.T) {
	testCases := []struct {
		comment           string
		hostPath          string
		containerRootPath string
		expectedPath      string
		srcDirs           []string
	}{
		{
			comment:           "input path outside of gopath dirs, returns default",
			hostPath:          "/root/go/github.com/repository/package",
			containerRootPath: "/host",
			expectedPath:      "/host",
			srcDirs:           []string{"/a/b/c/d", "/foo/bar"},
		},
		{
			comment:           "input path is within one of gopath dirs, returns it relative to container root",
			hostPath:          "/root/go/src/github.com/repository/package",
			containerRootPath: "/host",
			expectedPath:      "/host/src/github.com/repository/package",
			srcDirs:           []string{"/root/go", "/root/another/gopath"},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.comment, func(t *testing.T) {
			path := dockerSrcPathFromGopath(tc.hostPath, tc.containerRootPath, tc.srcDirs)
			require.Equal(t, tc.expectedPath, path)
		})
	}
}
