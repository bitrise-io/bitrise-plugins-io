package configs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SetAPIToken_and_readAPIToken(t *testing.T) {
	apiToken := "s0m3-ap1-t0k3n"
	DataDir = "."

	t.Log("OK - token is empty if not set")
	{
		tokenFromConfig, err := readAPIToken()
		require.NoError(t, err)
		require.Equal(t, "", tokenFromConfig)
	}

	t.Log("OK - set and read api token")
	{
		err := SetAPIToken(apiToken)
		require.NoError(t, err)

		tokenFromConfig, err := readAPIToken()
		require.NoError(t, err)
		require.Equal(t, "s0m3-ap1-t0k3n", tokenFromConfig)
	}

	t.Log("clean up - remove config file")
	err := os.Remove("config.yml")
	require.NoError(t, err)
}
