---
sidebar_position: 1
---

# Introduction

Zog is a schema builder for runtime value parsing and validation. Define a schema, transform a value to match, assert the shape of an existing value, or both. Zog schemas are extremely expressive and allow modeling complex, interdependent validations, or value transformations.

Killer Features:

- Concise yet expressive schema interface, equipped to model simple to complex data models
- **Zod-like API**, use method chaining to build schemas in a typesafe manner
- **Extensible**: add your own validators, schemas and data providers
- Rich error details, make debugging a breeze
- Almost no reflection when using primitive types
- **Built-in coercion** support for most types
- Zero dependencies!
- **Three Helper Packages**
  - **zenv**: parse environment variables
  - **zhttp**: parse http forms & query params
  - **i18n**: Opinionated solution to good i18n zog errors

> **API Stability:**
>
> - I will consider the API stable when we reach v1.0.0
> - However, I believe very little API changes will happen from the current implementation. The APIs are are most likely to change are the **data providers** (please don't make your own if possible use the helpers whose APIs will not change meaningfully) and the ParseCtx most other APIs should remain the same
> - Zog will not respect semver until v1.0.0 is released. Expect breaking changes (mainly in non basic apis) until then.
