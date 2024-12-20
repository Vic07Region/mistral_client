package mistal_client

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var (
	url = "https://api.mistral.ai/v1/chat/completions"
)

type Mistral interface {
	SendMessage(request SendMessageRequest) (string, error)
	SendMessageStream(ctx context.Context, request SendMessageRequest) (*StreamIterator, error)
}

type MistralProvider struct {
	ApiKey     string
	HTTPClient *http.Client
	model      string
}

func (m *MistralProvider) SendMessage(request SendMessageRequest) (string, error) {
	data, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("Error marshaling JSON: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return "", fmt.Errorf("Error creating request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+m.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := m.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response: %v", err)
	}

	var response MistralResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("Error Unmarshal response: %v", err)
	}
	if response.Choices != nil {
		result := ""
		for _, choice := range response.Choices {
			result += choice.Message.Content
		}
		return result, nil
	} else {
		return "", fmt.Errorf("No result")
	}
}

type SendMessageStreamData struct {
	Stream   bool      `json:"stream"`
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type StreamIterator struct {
	dataChan <-chan string
	errChan  <-chan error
	current  string
	err      error
}

func (it *StreamIterator) Next() bool {
	if it.err != nil {
		return false
	}

	select {
	case current, ok := <-it.dataChan:
		if !ok {
			return false
		}
		it.current = current
		return true
	case err, ok := <-it.errChan:
		it.err = err
		if !ok {
			return false
		}
		return true
	}
}

func (it *StreamIterator) Value() string {
	return it.current
}

func (it *StreamIterator) Err() error {
	return it.err
}

func (m *MistralProvider) SendMessageStream(ctx context.Context, request SendMessageRequest) (*StreamIterator, error) {
	responseCh := make(chan string)
	errCh := make(chan error, 1)

	asyncData := SendMessageStreamData{
		Stream:   true,
		Model:    request.Model,
		Messages: request.Messages,
	}

	go func() {
		defer close(responseCh)
		defer close(errCh)
		data, err := json.Marshal(asyncData)
		if err != nil {
			errCh <- fmt.Errorf("Error marshaling JSON: %v", err)
			return
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(data))
		if err != nil {
			errCh <- fmt.Errorf("MError creating request: %v", err)
			return
		}

		req.Header.Set("Authorization", "Bearer "+m.ApiKey)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		resp, err := m.HTTPClient.Do(req)
		if err != nil {
			fmt.Println(err)
			errCh <- fmt.Errorf("Error making request: %v", err)
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			errCh <- fmt.Errorf("Status not ok: %d", resp.Status)
			return
		}

		scanner := bufio.NewScanner(resp.Body)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			line := scanner.Text()

			if strings.HasPrefix(line, "data: ") {
				line = strings.TrimPrefix(line, "data: ")
			}

			if !strings.HasPrefix(strings.TrimSpace(line), "{") {
				continue
			}

			var chunk Chunk
			if err := json.Unmarshal([]byte(line), &chunk); err != nil {
				errCh <- fmt.Errorf("Error Decode JSON data: %v", err)
				return
			}

			for _, choice := range chunk.Choices {
				if choice.Delta.Content == "" {
					continue
				}
				select {
				case <-ctx.Done():
					errCh <- fmt.Errorf("Context canceled: %v", ctx.Err())
					return
				case responseCh <- choice.Delta.Content:
				}
			}
		}

		if err := scanner.Err(); err != nil {
			errCh <- fmt.Errorf("Scanner error: %v", err)
		}
	}()

	return &StreamIterator{dataChan: responseCh, errChan: errCh}, nil
}
