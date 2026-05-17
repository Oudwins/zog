package zjsonschema

import (
	"fmt"

	zsscore "github.com/Oudwins/zog/pkgs/zss/core"
	"github.com/Oudwins/zog/pkgs/zss/jsonschema/draft2020_12"
	"github.com/Oudwins/zog/pkgs/zss/jsonschema/shared"
	"github.com/Oudwins/zog/zconst"
)

type Draft string

const Draft2020_12 Draft = "https://json-schema.org/draft/2020-12/schema"

type Options struct {
	Draft Draft

	UnknownKindConverter UnknownKindConverter
	TestConverter        TestConverter
}

type Schema = shared.Schema

type UnknownKindConverter = shared.UnknownKindConverter

type TestConverter = shared.TestConverter

func ConvertTest(out Schema, kind zconst.ZogType, test *zsscore.ZSSTest) error {
	return draft2020_12.ConvertTest(out, kind, test)
}

func FromZSS(doc zsscore.ZSSDocument, opts Options) (Schema, error) {
	draft := opts.Draft
	if draft == "" {
		draft = Draft2020_12
	}
	if doc.Root == nil {
		return nil, fmt.Errorf("zss document root is nil")
	}

	switch draft {
	case Draft2020_12:
		return draft2020_12.FromZSS(doc, draft2020_12.Options{
			UnknownKindConverter: opts.UnknownKindConverter,
			TestConverter:        opts.TestConverter,
		})
	default:
		return nil, fmt.Errorf("unsupported JSON Schema draft %q", draft)
	}
}
