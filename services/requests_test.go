package services

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_urlWithParameters(t *testing.T) {
	baseURL := "https://some.basic.url"

	for _, tc := range []struct {
		expectedURL string
		queryParams map[string]string
	}{
		// no query params
		{
			expectedURL: baseURL,
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
		generatedURL, err := url.Parse(urlWithParameters(baseURL, tc.queryParams))
		require.NoError(t, err)

		t.Log("check params independent from order")
		for key, value := range tc.queryParams {
			require.Equal(t, value, generatedURL.Query().Get(key))
		}
	}
}
