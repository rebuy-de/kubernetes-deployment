package cmd

import "testing"

func init() {
	version = "go-test"

	// changing defaults to prevent actual running real shit
	defaultProjectConfigPath = "test-fixtures/services.yaml"
}

func runMainForTest(t *testing.T, wantedExit int, args ...string) {
	exit := Main(args...)
	if exit != wantedExit {
		t.Fatalf("got exit code %d, but wanted %d", exit, wantedExit)
	}
}

func TestVersion(t *testing.T) {
	runMainForTest(t, 0, "-version")
}

func TestUsage(t *testing.T) {
	// IMO it's nice to see the usage in the test logs
	runMainForTest(t, 2, "-h")
}
