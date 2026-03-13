package utils

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const (
	languageKey contextKey = "user-language"

	DefaultLanguage = "de"
	DefaultRegion   = "DE"
	DefaultLocale   = DefaultLanguage + "_" + DefaultRegion
)

func GetLanguage() string {
	return DefaultLanguage
}

func GetRegion() string {
	return DefaultRegion
}

func ParseLang(raw string) string {
	if raw == "" || raw == "C" || strings.HasPrefix(raw, "C.") || raw == "POSIX" {
		return ""
	}

	if parts := strings.Split(raw, ":"); len(parts) > 0 {
		raw = parts[0]
	}

	if idx := strings.Index(raw, ";"); idx != -1 {
		raw = raw[:idx]
	}

	if idx := strings.IndexAny(raw, "-_.@"); idx != -1 {
		raw = raw[:idx]
	}

	raw = strings.TrimSpace(raw)

	if len(raw) < 2 {
		return ""
	}

	return raw
}

func RequestLanguage(r *http.Request) string {
	header := r.Header.Get("Accept-Language")
	if header == "" {
		return DefaultLanguage
	}

	parts := strings.Split(header, ",")
	if len(parts) == 0 {
		return DefaultLanguage
	}

	lang := ParseLang(parts[0])
	if lang == "" {
		return DefaultLanguage
	}

	return lang
}

func LanguageMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lang := RequestLanguage(r)
		ctx := context.WithValue(r.Context(), languageKey, lang)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetLanguageCtx(ctx context.Context) string {
	if val, ok := ctx.Value(languageKey).(string); ok && val != "" {
		return val
	}
	return DefaultLanguage
}
