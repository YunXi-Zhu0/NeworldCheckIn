package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseReadsConfigYAML(t *testing.T) {
	t.Setenv("NEWORD_EMAIL", "")
	t.Setenv("NEWORD_PASSWD", "")
	withWorkingDir(t, writeConfig(t, "email: yaml@example.com\npasswd: yaml-pass\n"))

	cfg, err := Parse(nil)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Email != "yaml@example.com" || cfg.Passwd != "yaml-pass" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestParseEnvironmentOverridesYAML(t *testing.T) {
	t.Setenv("NEWORD_EMAIL", "env@example.com")
	t.Setenv("NEWORD_PASSWD", "env-pass")
	withWorkingDir(t, writeConfig(t, "email: yaml@example.com\npasswd: yaml-pass\n"))

	cfg, err := Parse(nil)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Email != "env@example.com" || cfg.Passwd != "env-pass" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestParseCLIOverridesEnvironmentAndYAML(t *testing.T) {
	t.Setenv("NEWORD_EMAIL", "env@example.com")
	t.Setenv("NEWORD_PASSWD", "env-pass")
	withWorkingDir(t, writeConfig(t, "email: yaml@example.com\npasswd: yaml-pass\n"))

	cfg, err := Parse([]string{"--email", "cli@example.com", "--passwd", "cli-pass"})
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Email != "cli@example.com" || cfg.Passwd != "cli-pass" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestParseMissingConfigFileIsAllowed(t *testing.T) {
	t.Setenv("NEWORD_EMAIL", "env@example.com")
	t.Setenv("NEWORD_PASSWD", "env-pass")
	withWorkingDir(t, t.TempDir())

	cfg, err := Parse(nil)
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Email != "env@example.com" || cfg.Passwd != "env-pass" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestParseInvalidYAMLReturnsError(t *testing.T) {
	t.Setenv("NEWORD_EMAIL", "")
	t.Setenv("NEWORD_PASSWD", "")
	withWorkingDir(t, writeConfig(t, "email: [broken\n"))

	_, err := Parse(nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "invalid config.yaml") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseMissingCredentials(t *testing.T) {
	t.Setenv("NEWORD_EMAIL", "")
	t.Setenv("NEWORD_PASSWD", "")
	withWorkingDir(t, writeConfig(t, ""))

	_, err := Parse(nil)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "config.yaml") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func withWorkingDir(t *testing.T, dir string) {
	t.Helper()

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldWD); err != nil {
			t.Fatal(err)
		}
	})
}

func writeConfig(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, defaultConfigPath)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}
