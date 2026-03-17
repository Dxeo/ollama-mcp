package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/ollama/ollama/api"
)

// OllamaClient wraps the Ollama API client with synchronous helpers.
type OllamaClient struct {
	client *api.Client
}

// NewOllamaClient creates a new OllamaClient using environment configuration.
// It reads OLLAMA_HOST from the environment (defaults to http://localhost:11434).
func NewOllamaClient() (*OllamaClient, error) {
	client, err := api.ClientFromEnvironment()
	if err != nil {
		return nil, fmt.Errorf("failed to create ollama client: %w", err)
	}
	return &OllamaClient{client: client}, nil
}

// Chat sends a chat message and returns the full accumulated response.
func (o *OllamaClient) Chat(ctx context.Context, model, prompt string) (string, error) {
	var b strings.Builder
	stream := false

	req := &api.ChatRequest{
		Model:  model,
		Stream: &stream,
		Messages: []api.Message{
			{Role: "user", Content: prompt},
		},
	}

	err := o.client.Chat(ctx, req, func(resp api.ChatResponse) error {
		b.WriteString(resp.Message.Content)
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("chat failed: %w", err)
	}

	return b.String(), nil
}

// Generate sends a generate request and returns the full accumulated response.
func (o *OllamaClient) Generate(ctx context.Context, model, prompt string) (string, error) {
	var b strings.Builder
	stream := false

	req := &api.GenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: &stream,
	}

	err := o.client.Generate(ctx, req, func(resp api.GenerateResponse) error {
		b.WriteString(resp.Response)
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("generate failed: %w", err)
	}

	return b.String(), nil
}

// Embed generates embeddings for the given text and returns the vector.
func (o *OllamaClient) Embed(ctx context.Context, model, text string) ([]float32, error) {
	resp, err := o.client.Embed(ctx, &api.EmbedRequest{
		Model: model,
		Input: text,
	})
	if err != nil {
		return nil, fmt.Errorf("embed failed: %w", err)
	}

	if len(resp.Embeddings) == 0 {
		return nil, fmt.Errorf("no embeddings returned")
	}

	return resp.Embeddings[0], nil
}
