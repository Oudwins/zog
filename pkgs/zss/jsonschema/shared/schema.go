package shared

import (
	zsscore "github.com/Oudwins/zog/pkgs/zss/core"
	"github.com/Oudwins/zog/zconst"
)

type Schema map[string]any

type UnknownKindConverter func(schema *zsscore.ZSSSchema) (Schema, error)

type TestConverter func(out Schema, kind zconst.ZogType, test *zsscore.ZSSTest) error
