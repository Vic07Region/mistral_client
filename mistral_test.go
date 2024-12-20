package mistral_client_test

import (
	"github.com/Vic07Region/mistral_client"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"testing"
)

func TestMistral_SendMessage(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Error(err)
	}
	apikey := os.Getenv("API_KEY")
	type fields struct {
		apiKey     string
		HTTPClient *http.Client
		baseURL    string
	}
	type args struct {
		request mistral_client.SendMessageRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "test count 1 to 5",
			fields: fields{
				apiKey: apikey,
			},
			args: args{
				request: mistral_client.SendMessageRequest{
					Model: "mistral-large-latest",
					Messages: []mistral_client.Message{
						{
							Role:    "user",
							Content: "посчитай от 1 до 5",
						},
					},
				},
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "test wrong model",
			fields: fields{
				apiKey: apikey,
			},
			args: args{
				request: mistral_client.SendMessageRequest{
					Model: "wrong-model",
					Messages: []mistral_client.Message{
						{
							Role:    "user",
							Content: "посчитай от 1 до 5",
						},
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "test wrong apikey",
			fields: fields{
				apiKey: "apikey",
			},
			args: args{
				request: mistral_client.SendMessageRequest{
					Model: "mistral-large-latest",
					Messages: []mistral_client.Message{
						{
							Role:    "user",
							Content: "посчитай от 1 до 5",
						},
					},
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			cl := mistral_client.New(tt.fields.apiKey)

			got, err := cl.Mistral.SendMessage(tt.args.request)
			if err != nil && !tt.wantErr {
				t.Errorf("Mistral.SendMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else {
				t.Logf("Mistral.SendMessage() error = %v, wantErr %v", err, tt.wantErr)
			}

			if got != tt.want && got == "" {
				t.Errorf("Mistral.SendMessage() = %v, want %v", got, tt.want)
			} else {
				t.Logf("Mistral.SendMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
