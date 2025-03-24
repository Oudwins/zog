---
sidebar_position: 20
---

# Configuration

Inside the `conf` package zog provides a bunch of global configuration options. If you override these options you can **change the default behavior of all zog schemas.**

## Coercion

During [parsing](/core-concepts/parsing) zog will attempt to coerce the data into the correct type. For example if you have a `float64` field and the data is a `string` that contains a number with a comma as the decimal separator, zog will attempt to convert this to a `float64`. Zog provides a set of default coercer functions for each type, but you can override these globally.

Lets go through an example of overriding the `float64` coercer function, because we want to support floats that use a comma as the decimal separator.

```go
import (
	// import the conf package
	"github.com/Oudwins/zog/conf"
)

// we override the coercer function for float64
conf.Coercers.Float64 = func(data any) (any, error) {
	str, ok := data.(string)
	// identify the case we want to override
	if !ok && strings.Contains(str, ",") {
		return MyCustomFloatCoercer(str)
	}
	// fallback to the original function
	return conf.DefaultCoercers.Float64(data)
}
```

## Error Formatting

For information on configuring error formatting globally please refer to the [errors page](/errors#5-configure-issue-messages-globally).
