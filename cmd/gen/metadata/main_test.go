package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestGeneratedMetadataIsCurrent(t *testing.T) {
	root := filepath.Join("..", "..", "..")
	generated, err := generate(root)
	if err != nil {
		t.Fatal(err)
	}

	current, err := os.ReadFile(filepath.Join(root, defaultOutput))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(generated, current) {
		t.Fatal("metadata_generated.go is stale; run go generate .")
	}
}
