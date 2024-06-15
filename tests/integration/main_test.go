package integration

import (
	"os"
	"testing"

	"github.com/calvincolton/greenlight/tests/testutils"
)

func TestMain(m *testing.M) {
	testutils.SetupDatabase(nil)

	code := m.Run()

	testutils.TeardownDatabase(nil)

	os.Exit(code)
}
