package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseLang(t *testing.T) {
	tests := map[string]string{
		"de_DE.UTF-8": "de",
		"de-DE":       "de",
		"en-US":       "en",
		"de:en":       "de",
		"de-DE;q=0.9": "de",
		"C":           "",
		"POSIX":       "",
		"":            "",
	}

	for input, expected := range tests {
		result := ParseLang(input)
		if result != expected {
			t.Fatalf("ParseLang(%s)=%s expected=%s", input, result, expected)
		}
	}
}

func TestRequestLanguage(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Language", "de-DE,de;q=0.9,en-US;q=0.8")

	lang := RequestLanguage(r)

	if lang != "de" {
		t.Fatalf("unexpected language: %s", lang)
	}
}

func TestRequestLanguageDefault(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)

	lang := RequestLanguage(r)

	if lang != DefaultLanguage {
		t.Fatalf("unexpected default language: %s", lang)
	}
}

func TestLanguageMiddleware(t *testing.T) {
	handler := LanguageMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lang := GetLanguageCtx(r.Context())
		if lang != "en" {
			t.Fatalf("unexpected context language: %s", lang)
		}
	}))

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Language", "en-US,en;q=0.9")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)
}
