package integration

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/calvincolton/greenlight/tests/testutils"
)

func TestMain(m *testing.M) {
	testutils.SetupDatabase(nil)

	if err := waitForServer("http://localhost:8081/v1/healthcheck", 10, 2*time.Second); err != nil {
		log.Fatalf("server did not become ready: %v", err)
	}

	// Run tests
	code := m.Run()

	testutils.TeardownDatabase(nil)

	os.Exit(code)
}

func waitForServer(url string, attempts int, sleep time.Duration) error {
	for i := 0; i < attempts; i++ {
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			return nil
		}
		time.Sleep(sleep)
	}
	return fmt.Errorf("server at %s not ready after %d attempts", url, attempts)
}
