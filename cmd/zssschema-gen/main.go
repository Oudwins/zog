package main

// Usage:
//
//	go run ./cmd/zssschema-gen -version 0.0.1
//	go run ./cmd/zssschema-gen -version 0.0.1 -inline
//	go run ./cmd/zssschema-gen -version 0.0.1 -out /tmp/schema.json

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Oudwins/zog"
	zsscore "github.com/Oudwins/zog/pkgs/zss/core"
	zjsonschema "github.com/Oudwins/zog/pkgs/zss/jsonschema"
	zssschema "github.com/Oudwins/zog/pkgs/zss/schema"
)

const defaultBaseURL = "https://zog.dev/zss"
const defaultOutputRoot = "schemas/zss"

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("zssschema-gen", flag.ContinueOnError)
	fs.SetOutput(stderr)

	version := fs.String("version", "", "ZSS schema version to generate, e.g. 0.0.1 or 0.1.0-beta.1")
	out := fs.String("out", "", "output file path; defaults to schemas/zss/{version}/schema.json")
	inline := fs.Bool("inline", false, "write schema to stdout instead of a file")
	baseURL := fs.String("base-url", defaultBaseURL, "base URL used to build the schema $id")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *version == "" {
		return errors.New("-version is required")
	}

	id := schemaID(*baseURL, *version)
	if !zsscore.ZSS_URI_REGEX.MatchString(id) {
		return fmt.Errorf("invalid ZSS schema version %q", *version)
	}

	schema, err := generate(id, *version)
	if err != nil {
		return err
	}

	encoded, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal schema: %w", err)
	}
	encoded = append(encoded, '\n')

	if *inline {
		_, err = stdout.Write(encoded)
		return err
	}
	if *out == "" {
		*out = defaultOutputPath(*version)
	}

	if err := os.MkdirAll(filepath.Dir(*out), 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	if err := os.WriteFile(*out, encoded, 0o644); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}
	return nil
}

func schemaID(baseURL, version string) string {
	return strings.TrimRight(baseURL, "/") + "/" + version + "/schema.json"
}

func defaultOutputPath(version string) string {
	return filepath.Join(defaultOutputRoot, version, "schema.json")
}

func generate(id, version string) (zjsonschema.Schema, error) {
	doc := zog.EXPERIMENTAL_TO_ZSS[zsscore.ZSSDocument](zssschema.ZSSDocumentSchema)
	schema, err := zjsonschema.FromZSS(doc, zjsonschema.Options{})
	if err != nil {
		return nil, fmt.Errorf("convert ZSS schema to JSON Schema: %w", err)
	}
	schema["$id"] = id
	schema["title"] = "Zog Schema Specification"
	schema["description"] = "JSON Schema for ZSS documents."
	schema["version"] = version
	return schema, nil
}
