---
sidebar_position: 10
hide_table_of_contents: false
toc_min_heading_level: 2
toc_max_heading_level: 4
---

# Common Mistakes

- schemas are optional
- creating schemas inside handlers
- creating deeply nested schemas
- Using the data after parsing even though there were errors

---

Composition and Reuse
Schema are immutable, each method call returns a new schema object. Reuse and pass them around without fear of mutating another instance.

let optionalString = string().optional();

let definedString = optionalString.defined();

let value = undefined;
optionalString.isValid(value); // true
definedString.isValid(value); // false

that is another footgun. That this isn't how it works in zog
