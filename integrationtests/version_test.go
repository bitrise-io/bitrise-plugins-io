// +build integration

package integrationtests

import (
	"testing"

	"github.com/bitrise-io/bitrise-plugins-io/version"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("")
	require.NoError(t, err)

	cmd := command.New(binPath(), "version")
	cmd.SetDir(tmpDir)
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	require.NoError(t, err, out)
	require.Equal(t, version.VERSION, out)
}
