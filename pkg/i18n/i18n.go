package i18n

import (
	"embed"
	"encoding/json"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.json
var localesFS embed.FS

type Manager struct {
	bundle *i18n.Bundle
}

func NewManager() (*Manager, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	if _, err := bundle.LoadMessageFileFS(localesFS, "locales/en.json"); err != nil {
		return nil, err
	}

	if _, err := bundle.LoadMessageFileFS(localesFS, "locales/ka.json"); err != nil {
		return nil, err
	}

	return &Manager{bundle: bundle}, nil
}

func (m *Manager) Localize(lang, messageID string) string {
	localizer := i18n.NewLocalizer(m.bundle, lang)

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
	if err != nil {
		return messageID
	}

	return msg
}

func (m *Manager) GetLocalizer(acceptLanguage string) *i18n.Localizer {
	langs := parseAcceptLanguage(acceptLanguage)

	return i18n.NewLocalizer(m.bundle, langs...)
}

func parseAcceptLanguage(header string) []string {
	if header == "" {
		return []string{"en"}
	}

	parts := strings.Split(header, ",")
	langs := make([]string, 0, len(parts))

	for _, part := range parts {
		lang := strings.TrimSpace(strings.Split(part, ";")[0])
		langs = append(langs, lang)
	}

	if len(langs) == 0 {
		return []string{"en"}
	}

	return langs
}
