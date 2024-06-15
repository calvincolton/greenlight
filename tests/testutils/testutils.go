package testutils

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	_ "github.com/lib/pq"
)

var (
	DB                *sql.DB
	_, b, _, _        = runtime.Caller(0)
	projectRoot       = filepath.Join(filepath.Dir(b), "../..")
	dockerComposeFile = filepath.Join(projectRoot, "docker-compose.test.yml")
)

func runCommandWithProject(name string, project string, args ...string) error {
	allArgs := append([]string{"-p", project}, args...)
	cmd := exec.Command(name, allArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func StartDockerCompose() error {
	return runCommandWithProject("docker-compose", "greenlight_test", "-f", dockerComposeFile, "up", "--build", "-d")
}

func StopDockerCompose() error {
	return runCommandWithProject("docker-compose", "greenlight_test", "-f", dockerComposeFile, "down")
}

func SetupDatabase(t *testing.T) {
	var err error
	var databaseDSN string

	flag.StringVar(&databaseDSN, "test-db-dsn", os.Getenv("TEST_DATABASE_DSN"), "PostgreSQL DSN")
	flag.Parse()

	if err := StartDockerCompose(); err != nil {
		if t != nil {
			t.Fatalf("could not start Docker Compose: %v", err)
		} else {
			log.Fatalf("could not start Docker Compose: %v", err)
		}
	}

	DB, err = sql.Open("postgres", databaseDSN)
	if err != nil {
		if t != nil {
			t.Fatalf("could not connect to test database: %v", err)
		} else {
			log.Fatalf("could not connect to test database: %v", err)
		}
	}
}

func TeardownDatabase(t *testing.T) {
	if err := StopDockerCompose(); err != nil {
		if t != nil {
			t.Logf("could not stop docker-compose: %v", err)
		} else {
			log.Printf("could not stop docker-compose: %v", err)
		}
	}
}

func MakeRequest(
	t *testing.T,
	method string,
	url string,
	body interface{},
	expectedStatus int,
) *http.Response {
	t.Helper()

	var jsonData []byte
	var err error

	if body != nil {
		jsonData, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("could not marshal request body: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("could not make request: %v", err)
	}

	if resp.StatusCode != expectedStatus {
		t.Fatalf("expected status %d, got %d", expectedStatus, resp.StatusCode)
	}

	return resp
}
