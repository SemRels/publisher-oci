package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIntegrationDryRunResolvesOCIReferenceWithoutPublishing(t *testing.T) {
	artifact := filepath.Join(t.TempDir(), "release.tgz")
	if err := os.WriteFile(artifact, []byte("release-bytes"), 0o600); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	code := run(&stdout, &stderr, integrationEnv(map[string]string{
		"SEMREL_VERSION":         "v1.2.3",
		"SEMREL_PLUGIN_REF":      "registry.invalid/release:{version}",
		"SEMREL_PLUGIN_ARTIFACT": artifact,
		"SEMREL_DRY_RUN":         "true",
	}))
	if code != 0 {
		t.Fatalf("dry-run code = %d, stderr = %s", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "oras push registry.invalid/release:1.2.3") {
		t.Errorf("stdout = %q", stdout.String())
	}
}

func TestIntegrationInvalidOCIConfigurationFailsBeforeExternalCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run(&stdout, &stderr, integrationEnv(map[string]string{
		"SEMREL_VERSION":      "1.2.3",
		"SEMREL_PLUGIN_REF":   "registry.invalid/release",
		"SEMREL_PLUGIN_TOKEN": "integration-secret",
	}))
	if code != 1 || !strings.Contains(stderr.String(), "no artifacts configured") {
		t.Fatalf("code = %d, stderr = %q", code, stderr.String())
	}
	if strings.Contains(stdout.String(), "integration-secret") || strings.Contains(stderr.String(), "integration-secret") {
		t.Error("invalid configuration leaked its token")
	}
}

func integrationEnv(values map[string]string) func(string) string {
	return func(key string) string {
		return values[key]
	}
}
