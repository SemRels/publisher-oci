//go:build integration

package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestIntegrationPublishesArtifactToDistributionRegistry(t *testing.T) {
	ref := os.Getenv("SEMREL_TEST_OCI_REF")
	if ref == "" {
		t.Fatal("SEMREL_TEST_OCI_REF is required for the OCI integration test")
	}
	if _, err := exec.LookPath("oras"); err != nil {
		t.Fatalf("oras is required for the OCI integration test: %v", err)
	}

	dir := t.TempDir()
	artifactName := "release.txt"
	artifactBody := []byte("semrel OCI integration artifact\n")
	if err := os.WriteFile(filepath.Join(dir, artifactName), artifactBody, 0o600); err != nil {
		t.Fatal(err)
	}
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})

	var stdout, stderr bytes.Buffer
	code := run(&stdout, &stderr, integrationEnv(map[string]string{
		"SEMREL_VERSION":         "v1.2.3",
		"SEMREL_PLUGIN_REF":      ref,
		"SEMREL_PLUGIN_ARTIFACT": artifactName,
	}))
	if code != 0 {
		t.Fatalf("publication code = %d, stdout = %q, stderr = %q", code, stdout.String(), stderr.String())
	}
	if !strings.Contains(stdout.String(), "published 1 artifact") {
		t.Fatalf("publication output = %q", stdout.String())
	}

	output, err := exec.Command("oras", "manifest", "fetch", ref).CombinedOutput()
	if err != nil {
		t.Fatalf("fetch published manifest: %v: %s", err, output)
	}
	var manifest struct {
		SchemaVersion int `json:"schemaVersion"`
		Layers        []struct {
			Size        int64             `json:"size"`
			Annotations map[string]string `json:"annotations"`
		} `json:"layers"`
	}
	if err := json.Unmarshal(output, &manifest); err != nil {
		t.Fatalf("decode published manifest: %v: %s", err, output)
	}
	if manifest.SchemaVersion != 2 || len(manifest.Layers) != 1 {
		t.Fatalf("unexpected published manifest: %s", output)
	}
	layer := manifest.Layers[0]
	if layer.Annotations["org.opencontainers.image.title"] != artifactName || layer.Size != int64(len(artifactBody)) {
		t.Fatalf("manifest layer does not describe the artifact: %s", output)
	}

	pullDir := filepath.Join(dir, "pulled")
	pullOutput, err := exec.Command("oras", "pull", "--output", pullDir, ref).CombinedOutput()
	if err != nil {
		t.Fatalf("pull published artifact: %v: %s", err, pullOutput)
	}
	pulledBody, err := os.ReadFile(filepath.Join(pullDir, artifactName))
	if err != nil {
		t.Fatalf("read pulled artifact: %v", err)
	}
	if !bytes.Equal(pulledBody, artifactBody) {
		t.Fatalf("pulled artifact = %q, want %q", pulledBody, artifactBody)
	}
}
