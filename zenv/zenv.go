package zenv

import (
	"log"
	"os"

	z "github.com/Oudwins/zog"
	p "github.com/Oudwins/zog/primitives"
)

type envDataProvider struct {
}

func (e *envDataProvider) Get(key string) any {
	return os.Getenv(key)

}

// Parses environment variables into destinationStruct
func Parse(schema z.StructParser, destPtr any, panicOnError bool) p.ZogSchemaErrors {
	errs := schema.Parse(&envDataProvider{}, destPtr)

	if len(errs) > 0 && panicOnError {
		log.Fatalf("FAILED TO PARSE ENVIRONMENT VARIABLES: %+v", errs)
	}
	return errs
}
