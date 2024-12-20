package mistal_client

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type SendMessageRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ChoiceStream struct {
	Delta struct {
		Content string `json:"content"`
	} `json:"delta"`
}

type Chunk struct {
	Choices []ChoiceStream `json:"choices"`
}

type Choice struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}

type MistralResponse struct {
	Choices []Choice `json:"choices"`
}
