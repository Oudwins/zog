package i18n

import (
	"github.com/Oudwins/zog/conf"
	"github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

// Takes a map[langKey]conf.LangMap
func SetLanguagesErrsMap(m map[string]zconst.LangMap, defaultLang string) {
	langKey := "lang"

	conf.ErrorFormatter = func(e internals.ZogError, ctx internals.ParseCtx) {
		lang := ctx.Get(langKey)
		if lang != nil {
			langM, ok := m[lang.(string)]
			if ok {
				conf.NewDefaultFormatter(langM)(e, ctx)
				return
			}
		}
		// use default lang if failed to get correct language map
		conf.NewDefaultFormatter(m[defaultLang])(e, ctx)
	}
}
