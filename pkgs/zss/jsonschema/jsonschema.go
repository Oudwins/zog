package zjsonschema

import (
	"fmt"

	zsscore "github.com/Oudwins/zog/pkgs/zss/core"
	"github.com/Oudwins/zog/pkgs/zss/jsonschema/draft2020_12"
	"github.com/Oudwins/zog/pkgs/zss/jsonschema/internal"
)

type Draft string

const Draft2020_12 Draft = "https://json-schema.org/draft/2020-12/schema"

type Options struct {
	Draft Draft
}

type Schema = internal.Schema

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
		return draft2020_12.FromZSS(doc)
	default:
		return nil, fmt.Errorf("unsupported JSON Schema draft %q", draft)
	}
}
