package i18n

import (
	"github.com/Oudwins/zog/conf"
	"github.com/Oudwins/zog/internals"
	"github.com/Oudwins/zog/zconst"
)

const (
	// Default lang key used to get the language from the ParseContext
	LangKey = "lang"
)

// Takes a map[langKey]conf.LangMap
// usage is i18n.SetLanguagesErrsMap(map[string]zconst.LangMap{
// "es": es.Map, "en": en.Map,
// }, "en", i18n.WithLangKey("langKey"))
// schema.Parse(data, &dest, z.WithCtxValue("langKey", "es"))
func SetLanguagesErrsMap(m map[string]zconst.LangMap, defaultLang string, opts ...setLanguageOption) {
	langKey := LangKey

	for _, op := range opts {
		op(&langKey)
	}

	conf.IssueFormatter = func(e internals.ZogIssue, ctx internals.Ctx) {
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

// Override the default lang key used to get the language from the ParseContext
func WithLangKey(key string) setLanguageOption {
	return func(lk *string) {
		*lk = key
	}
}

// Please use the helper function this type may very well change in the future but the helper function's API will stay the same
type setLanguageOption = func(langKey *string)

// Proxy the type for easy use
type LangMap = zconst.LangMap
