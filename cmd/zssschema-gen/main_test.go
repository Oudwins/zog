package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunWritesSchemaToDefaultLocation(t *testing.T) {
	workdir := t.TempDir()
	previousWorkdir, err := os.Getwd()
	require.NoError(t, err)
	require.NoError(t, os.Chdir(workdir))
	t.Cleanup(func() { require.NoError(t, os.Chdir(previousWorkdir)) })
	var stdout, stderr bytes.Buffer

	err = run([]string{"-version", "0.1.0-beta.1"}, &stdout, &stderr)

	require.NoError(t, err)
	assert.Empty(t, stdout.String())
	assert.Empty(t, stderr.String())

	data, err := os.ReadFile(filepath.Join(workdir, defaultOutputPath("0.1.0-beta.1")))
	require.NoError(t, err)

	var schema map[string]any
	require.NoError(t, json.Unmarshal(data, &schema))
	assert.Equal(t, "https://json-schema.org/draft/2020-12/schema", schema["$schema"])
	assert.Equal(t, "https://zog.dev/zss/0.1.0-beta.1/schema.json", schema["$id"])
	assert.Equal(t, "0.1.0-beta.1", schema["version"])
	assert.Equal(t, "Zog Schema Specification", schema["title"])
	assert.Contains(t, schema, "properties")
}

func TestRunWritesSchemaToStdoutWithInline(t *testing.T) {
	var stdout, stderr bytes.Buffer

	err := run([]string{"-version", "0.1.0-beta.1", "-inline"}, &stdout, &stderr)

	require.NoError(t, err)
	assert.Empty(t, stderr.String())

	var schema map[string]any
	require.NoError(t, json.Unmarshal(stdout.Bytes(), &schema))
	assert.Equal(t, "https://json-schema.org/draft/2020-12/schema", schema["$schema"])
	assert.Equal(t, "https://zog.dev/zss/0.1.0-beta.1/schema.json", schema["$id"])
	assert.Equal(t, "0.1.0-beta.1", schema["version"])
	assert.Equal(t, "Zog Schema Specification", schema["title"])
	assert.Contains(t, schema, "properties")
	assert.True(t, strings.HasPrefix(stdout.String(), "{\n  \"$id\":"))
}

func TestRunWritesSchemaToFile(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "nested", "schema.json")
	var stdout, stderr bytes.Buffer

	err := run([]string{"-version", "0.0.1", "-out", out}, &stdout, &stderr)

	require.NoError(t, err)
	assert.Empty(t, stdout.String())
	assert.Empty(t, stderr.String())

	data, err := os.ReadFile(out)
	require.NoError(t, err)
	assert.Contains(t, string(data), "https://zog.dev/zss/0.0.1/schema.json")
}

func TestRunRejectsInvalidVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer

	err := run([]string{"-version", "1.0.0-alpha"}, &stdout, &stderr)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid ZSS schema version")
	assert.Empty(t, stdout.String())
}

func TestRunRequiresVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer

	err := run(nil, &stdout, &stderr)

	require.Error(t, err)
	assert.Equal(t, "-version is required", err.Error())
	assert.Empty(t, stdout.String())
}
