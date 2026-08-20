package config

import (
	"os"
	"path/filepath"
	"testing"
)

// writeTestConfig points HOME at a temp dir and writes a config file into the
// location LoadConfig will read.
func writeTestConfig(t *testing.T, contents string) string {
	t.Helper()

	home := t.TempDir()
	t.Setenv("HOME", home)

	configDir := filepath.Join(home, ".config", "reddittui")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		t.Fatalf("could not create config dir: %v", err)
	}

	configPath := filepath.Join(configDir, configFilename)
	if err := os.WriteFile(configPath, []byte(contents), 0644); err != nil {
		t.Fatalf("could not write config file: %v", err)
	}

	return configPath
}

func TestLoadConfigReadsSessionCookie(t *testing.T) {
	writeTestConfig(t, "[auth]\nsessionCookie = \"abc123\"\n")
	os.Unsetenv(sessionCookieEnvVar)

	configuration, _ := LoadConfig()

	if configuration.Auth.SessionCookie != "abc123" {
		t.Errorf("SessionCookie = %q, want %q", configuration.Auth.SessionCookie, "abc123")
	}
}

func TestLoadConfigTrimsSessionCookie(t *testing.T) {
	writeTestConfig(t, "[auth]\nsessionCookie = \"  abc123  \"\n")
	os.Unsetenv(sessionCookieEnvVar)

	configuration, _ := LoadConfig()

	if configuration.Auth.SessionCookie != "abc123" {
		t.Errorf("SessionCookie = %q, want it trimmed to %q", configuration.Auth.SessionCookie, "abc123")
	}
}

func TestLoadConfigEnvVarOverridesFile(t *testing.T) {
	writeTestConfig(t, "[auth]\nsessionCookie = \"from-file\"\n")
	t.Setenv(sessionCookieEnvVar, "from-env")

	configuration, _ := LoadConfig()

	if configuration.Auth.SessionCookie != "from-env" {
		t.Errorf("SessionCookie = %q, want the env var to win", configuration.Auth.SessionCookie)
	}
}

// An empty env var is an explicit request to run unauthenticated, so it should
// still override a configured file value.
func TestLoadConfigEmptyEnvVarOverridesFile(t *testing.T) {
	writeTestConfig(t, "[auth]\nsessionCookie = \"from-file\"\n")
	t.Setenv(sessionCookieEnvVar, "")

	configuration, _ := LoadConfig()

	if configuration.Auth.SessionCookie != "" {
		t.Errorf("SessionCookie = %q, want empty", configuration.Auth.SessionCookie)
	}
}

// The cookie grants full account access, so the file must not stay 0644.
func TestLoadConfigRestrictsPermissionsWhenCookiePresent(t *testing.T) {
	configPath := writeTestConfig(t, "[auth]\nsessionCookie = \"abc123\"\n")
	os.Unsetenv(sessionCookieEnvVar)

	LoadConfig()

	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("could not stat config file: %v", err)
	}

	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("config file mode = %04o, want 0600", perm)
	}
}

func TestLoadConfigLeavesPermissionsAloneWithoutCookie(t *testing.T) {
	configPath := writeTestConfig(t, "[server]\ndomain = \"old.reddit.com\"\n")
	os.Unsetenv(sessionCookieEnvVar)

	LoadConfig()

	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("could not stat config file: %v", err)
	}

	if perm := info.Mode().Perm(); perm != 0644 {
		t.Errorf("config file mode = %04o, want 0644 left untouched", perm)
	}
}

// Existing settings must keep working now that a new section exists.
func TestLoadConfigStillReadsOtherSections(t *testing.T) {
	writeTestConfig(t, "[server]\ndomain = \"redlib.example.com\"\ntype = \"redlib\"\n\n[auth]\nsessionCookie = \"abc123\"\n")
	os.Unsetenv(sessionCookieEnvVar)

	configuration, _ := LoadConfig()

	if configuration.Server.Domain != "redlib.example.com" {
		t.Errorf("Server.Domain = %q, want %q", configuration.Server.Domain, "redlib.example.com")
	}

	if configuration.Server.Type != "redlib" {
		t.Errorf("Server.Type = %q, want %q", configuration.Server.Type, "redlib")
	}

	if configuration.Auth.SessionCookie != "abc123" {
		t.Errorf("SessionCookie = %q, want %q", configuration.Auth.SessionCookie, "abc123")
	}
}

func TestLoadConfigDefaultsWhenNoAuthSection(t *testing.T) {
	writeTestConfig(t, "[core]\nlogLevel = \"Debug\"\n")
	os.Unsetenv(sessionCookieEnvVar)

	configuration, _ := LoadConfig()

	if configuration.Auth.SessionCookie != "" {
		t.Errorf("SessionCookie = %q, want empty", configuration.Auth.SessionCookie)
	}

	if configuration.Core.LogLevel != "Debug" {
		t.Errorf("Core.LogLevel = %q, want %q", configuration.Core.LogLevel, "Debug")
	}
}
