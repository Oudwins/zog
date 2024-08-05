package zog

import (
	"fmt"

	p "github.com/Oudwins/zog/primitives"
)

type TestOption func(test *p.Test)

type msgParam interface {
	string | p.ErrorFunc
}

func Message[T msgParam](msg T) TestOption {
	switch v := any(msg).(type) {
	case string:
		return func(test *p.Test) {
			test.ErrorFunc = p.DErrorFunc(v)
		}
	case p.ErrorFunc:
		return func(test *p.Test) {
			test.ErrorFunc = v
		}
	default:
		panic(fmt.Errorf("invalid message type %T", v))
	}
}
