package integrationtests

import "os"

func binPath() string {
	return os.Getenv("INTEGRATION_TEST_BINARY_PATH")
}
