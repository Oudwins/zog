package zsscore // Zog Schema Specification

import "regexp"

type ZSSSchemaVersion = string

const (
	ZSS_VERSION_0_0_1  ZSSSchemaVersion = "https://zog.dev/zss/0.0.1/schema.json" // no file here yet
	ZSS_VERSION_LATEST ZSSSchemaVersion = ZSS_VERSION_0_0_1
)

// This accepts any url following the pattern of
// https://{domain}/..../{version}/schema.json
// Allows:
// https://example.com/money/1.0.0/schema.json
// https://example.com/money/0.1.0/schema.json
// https://example.com/zss/extensions/money/2.3.4-beta/schema.json
// https://example.com/zss/extensions/money/2.3.4-beta.1/schema.json
// Rejects:
// https://example.com/money/1/schema.json
// https://example.com/money/1.0/schema.json
// https://example.com/money/v1.0.0/schema.json
// https://example.com/money/1.0.0-alpha/schema.json
// https://example.com/money/1.0.0-rc.1/schema.json
// https://example.com/money/1.0.0-beta.01/schema.json
// Extracts values like so
// https://example.com/money/1.0.0/schema.json
// - Id = https://example.com/money
// - Version = 1.0.0
const ZSS_URI_REGEX_PATTERN = `^(?P<id>https?://[A-Za-z0-9.-]+(?::[0-9]+)?(?:/[A-Za-z0-9._~!$&'()*+,;=:@%-]+)+)/(?P<version>(?:0|[1-9][0-9]*)\.(?:0|[1-9][0-9]*)\.(?:0|[1-9][0-9]*)(?:-beta(?:\.(?:0|[1-9][0-9]*))?)?)/schema\.json$`

var ZSS_URI_REGEX = regexp.MustCompile(ZSS_URI_REGEX_PATTERN)
