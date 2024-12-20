package mistral_client

import (
	"net/http"
	"time"
)

type mistralAi struct {
	Mistral
	HTTPClient *http.Client
}

func New(apikey string) *mistralAi {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	provider := MistralProvider{
		ApiKey:     apikey,
		HTTPClient: httpClient,
		BaseURL:    "https://api.mistral.ai/v1/chat/completions",
	}

	return &mistralAi{
		HTTPClient: httpClient,
		Mistral:    &provider,
	}
}

func (l *mistralAi) SetAPIKey(apiKey string) {
	l.Mistral.setApiKey(apiKey)
}

func (l *mistralAi) SetBaseURL(url string) {
	l.Mistral.setBaseURL(url)
}
