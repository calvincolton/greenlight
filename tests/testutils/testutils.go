package testutils

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
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
	if err := runCommandWithProject("docker-compose", "greenlight_test", "-f", dockerComposeFile, "down", "-v"); err != nil {
		return err
	}
	return runCommandWithProject("docker", "", "volume", "rm", "greenlight_test_test_db_data")
}

func MakeRequest(
	t *testing.T,
	method string,
	url string,
	body interface{},
	expectedStatus int,
) (*http.Response, error) {
	t.Helper()

	var jsonData []byte
	var err error

	if body != nil {
		jsonData, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("could not marshal request body: %v", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("could not make request: %v", err)
	}

	if resp.StatusCode != expectedStatus {
		return resp, fmt.Errorf("expected status %d, got %d", expectedStatus, resp.StatusCode)
	}

	return resp, nil
}

func Equal(a, b map[string]any) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if vb, ok := b[k]; !ok || !deepEqual(v, vb) {
			return false
		}
	}
	return true
}

// Helper function to deeply compare two values
func deepEqual(a, b any) bool {
	jsonA, errA := json.Marshal(a)
	if errA != nil {
		return false
	}
	jsonB, errB := json.Marshal(b)
	if errB != nil {
		return false
	}

	return bytes.Equal(jsonA, jsonB)
}
