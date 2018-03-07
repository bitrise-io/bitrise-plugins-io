package integration

import (
	"testing"

	"github.com/bitrise-io/bitrise/models"
	"github.com/bitrise-io/bitrise/plugins"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func Test_RunTest(t *testing.T) {
	t.Log("success build")
	{
		tmpDir, err := pathutil.NormalizedOSTempDirPath("")
		require.NoError(t, err)

		envs := []string{
			plugins.PluginInputDataDirKey + "=" + tmpDir,

			plugins.PluginInputPluginModeKey + "=" + string(plugins.TriggerMode),
			plugins.PluginInputFormatVersionKey + "=" + models.Version,
		}

		cmd := command.New(binPath())
		cmd.SetEnvs(envs...)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)
	}
}
