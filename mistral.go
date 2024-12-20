package mistral_client

import (
	"net/http"
	"time"
)

type mistralAi struct {
	apiKey string
	Mistral
	HTTPClient *http.Client
}

func New(apikey string) *mistralAi {
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}
	return &mistralAi{
		apiKey:     apikey,
		HTTPClient: httpClient,
		Mistral: &MistralProvider{
			ApiKey:     apikey,
			HTTPClient: httpClient,
		},
	}
}

func (l *mistralAi) SetAPIKey(apiKey string) {
	l.apiKey = apiKey
}
