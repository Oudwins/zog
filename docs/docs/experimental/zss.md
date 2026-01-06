---
sidebar_position: 2
---

# ZSS - Zog Schema Specification

**ZSS (Zog Schema Specification)** is an intermediate, structured representation of Zog schemas. It serves as a bridge between Zog's runtime schema definitions and other schema formats.

> **⚠️ Experimental API**: ZSS is an experimental intermediate format that is still in flux. The structure and API are subject to breaking changes. Use with caution.

## Purpose

ZSS is designed to be an intermediate format that can be used to generate other schema formats such as:

- **JSON Schema** (planned)
- **OpenAPI/Swagger** specifications (planned)
- **TypeScript types** (planned)
- Other target formats as needed

By converting Zog schemas to ZSS first, we can generate multiple output formats from a single, well-defined intermediate representation.

## Structure

A ZSS document consists of:

- **`ZSSDocument`**: The root container with version information and a root schema
- **`ZSSSchema`**: Represents individual schema nodes with:
  - `Kind`: The schema type (string, number, bool, time, slice, struct, ptr, etc.)
  - `Processors`: Validation tests and transformers
  - `Child`: Nested schemas (for slices, structs, pointers, etc.)
  - `GoTypes`: Go type metadata (when exhaustive metadata is enabled)
  - `Required`, `DefaultValue`, `CatchValue`: Schema constraints and defaults

## Usage

Convert a Zog schema to ZSS using `EXPERIMENTAL_TO_ZSS()`:

```go
import (
    z "github.com/Oudwins/zog"
    zss "github.com/Oudwins/zog/pkgs/zss/core"
)

schema := z.String().Min(5)

// Convert to ZSS
zssDoc := z.EXPERIMENTAL_TO_ZSS(schema)

// Access the root schema
rootSchema := zssDoc.Root
```

## Exhaustive Metadata

ZSS supports an "exhaustive metadata" mode that includes additional information in the schema output. This mode must be explicitly enabled at build time.

### Enabling Exhaustive Metadata

To enable exhaustive metadata, build your application with the `zogmeta` build tag:

```bash
go build -tags zogmeta
```

Or when running tests:

```bash
go test -tags zogmeta
```

### What's Included

When exhaustive metadata is enabled, ZSS includes:

- **`GoTypes`**: Array of `ZSSGoType` containing:

  - `PkgPath`: Package path (empty for built-in types)
  - `Name`: Type name (may be empty for unnamed types)
  - `Display`: Full type string representation

- **`Format`**: For `TimeSchema`, the format string specified via `z.Time(z.Time.Format(...))` is included in the ZSS output. **Note**: The format is only included in ZSS when exhaustive metadata is enabled.

- **Custom Messages**: Custom messages set via `z.Message()` are included in the ZSS `ZSSTest.Message` field. **Note**: Custom messages are only included in ZSS output when exhaustive metadata is enabled.

For generic schemas like `PreprocessSchema[F, T]` and `BoxedSchema[B, T]`, multiple type parameters are captured in the `GoTypes` array.

### Example

```go
// With exhaustive metadata enabled
schema := z.Time(z.Time.Format(time.RFC3339))
schema := z.String().Min(5, z.Message("Custom validation message"))

zssDoc := z.EXPERIMENTAL_TO_ZSS(schema)
// Format and custom messages will be included in zssDoc
```

## Notes

- The ZSS format is versioned via the `$schema` field
- Currently, only version `0.0.1` is defined
- The format may change significantly before stabilization
