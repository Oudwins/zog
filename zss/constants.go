package zss // Zog Schema Specification

const (
	ZSS_TYPE_KEY     = "typeName"
	ZSS_TYPE_FLOAT64 = "float64"
	ZSS_TYPE_FLOAT32 = "float32"
	ZSS_TYPE_INT     = "int"
	ZSS_TYPE_INT64   = "int64"
	ZSS_TYPE_INT32   = "int32"
	ZSS_TYPE_UINT    = "uint"
	ZSS_TYPE_UINT64  = "uint64"
	ZSS_TYPE_UINT32  = "uint32"
)

type ZSSVersion = string

const (
	ZSS_VERSION_0_0_1 ZSSVersion = "0.0.1"
)
