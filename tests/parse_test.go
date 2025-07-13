package parse_test

import (
	"os"
	"testing"

	"github.com/akshatsrivastava11/Histograph/internals/parse"
)

func TestGetChromeHistoryPath_EnvOverride(t *testing.T) {
	os.Setenv("CHROME_HISTORY_PATH", "/tmp/fake_chrome_history")
	defer os.Unsetenv("CHROME_HISTORY_PATH")

	path, err := parse.GetChromeHistoryPath()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "/tmp/fake_chrome_history" {
		t.Errorf("expected /tmp/fake_chrome_history, got %s", path)
	}
}

func TestGetFirefoxHistoryPath_EnvOverride(t *testing.T) {
	os.Setenv("FIREFOX_HISTORY_PATH", "/tmp/fake_firefox_history")
	defer os.Unsetenv("FIREFOX_HISTORY_PATH")

	path, err := parse.GetFirefoxHistoryPath()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "/tmp/fake_firefox_history" {
		t.Errorf("expected /tmp/fake_firefox_history, got %s", path)
	}
}

// Note: OS-specific path tests would require refactoring the parse package to export the path logic as a testable function and/or allow injection of runtime.GOOS and home directory. This is a basic test for env override logic.
