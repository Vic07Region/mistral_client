package mistral_client

import (
	"net/http"
	"time"
)

type MistralAi struct {
	Mistral
	HTTPClient *http.Client
}

func New(apikey string) *MistralAi {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	provider := MistralProvider{
		ApiKey:     apikey,
		HTTPClient: httpClient,
		BaseURL:    "https://api.mistral.ai/v1/chat/completions",
	}

	return &MistralAi{
		HTTPClient: httpClient,
		Mistral:    &provider,
	}
}

func (l *MistralAi) SetAPIKey(apiKey string) {
	l.Mistral.setApiKey(apiKey)
}

func (l *MistralAi) SetBaseURL(url string) {
	l.Mistral.setBaseURL(url)
}
