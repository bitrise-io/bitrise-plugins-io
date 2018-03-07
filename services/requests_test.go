package services

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_urlWithParameters(t *testing.T) {
	url := "https://some.basic.url"

	for _, tc := range []struct {
		expectedURL string
		queryParams map[string]string
	}{
		// no query params
		{
			expectedURL: url,
		},
		// one parameter
		{
			expectedURL: "https://some.basic.url?param=value",
			queryParams: map[string]string{
				"param": "value",
			},
		},
		// more parameters
		{
			expectedURL: "https://some.basic.url?param=value&param2=value2",
			queryParams: map[string]string{
				"param":  "value",
				"param2": "value2",
			},
		},
		// empty parameters don't get into the URL
		{
			expectedURL: "https://some.basic.url?param=value&param2=value2",
			queryParams: map[string]string{
				"param":  "value",
				"param2": "value2",
				"param3": "",
			},
		},
	} {
		generatedURL := urlWithParameters(url, tc.queryParams)
		require.Equal(t, tc.expectedURL, generatedURL)
	}
}
