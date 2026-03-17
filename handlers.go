package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/mark3labs/mcp-go/mcp"
)

// Handlers holds the shared state needed by tool handlers.
type Handlers struct {
	ollama         *OllamaClient
	reasoningModel string
	embeddingModel string
}

// HandleReasonTask performs reasoning/code generation using the chat model.
func (h *Handlers) HandleReasonTask(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	prompt, err := request.RequireString("prompt")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("missing required parameter: %v", err)), nil
	}

	result, err := h.ollama.Chat(ctx, h.reasoningModel, prompt)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("reasoning failed: %v", err)), nil
	}

	return mcp.NewToolResultText(result), nil
}

// HandleEmbedText generates embeddings for the given text.
func (h *Handlers) HandleEmbedText(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	text, err := request.RequireString("text")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("missing required parameter: %v", err)), nil
	}

	embeddings, err := h.ollama.Embed(ctx, h.embeddingModel, text)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("embedding failed: %v", err)), nil
	}

	data, err := json.Marshal(map[string]any{"embeddings": embeddings})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal embeddings: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

// HandleFilterDocs ranks documents by semantic similarity to a query.
func (h *Handlers) HandleFilterDocs(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := request.RequireString("query")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("missing required parameter: %v", err)), nil
	}

	documents := request.GetStringSlice("documents", nil)
	if documents == nil {
		return mcp.NewToolResultError("missing required parameter: documents"), nil
	}

	queryEmb, err := h.ollama.Embed(ctx, h.embeddingModel, query)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to embed query: %v", err)), nil
	}

	type scored struct {
		doc   string
		score float64
	}
	results := make([]scored, 0, len(documents))

	for _, doc := range documents {
		docEmb, err := h.ollama.Embed(ctx, h.embeddingModel, doc)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("failed to embed document: %v", err)), nil
		}
		sim := CosineSimilarity(queryEmb, docEmb)
		results = append(results, scored{doc: doc, score: sim})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	ranked := make([]string, len(results))
	for i, r := range results {
		ranked[i] = r.doc
	}

	data, err := json.Marshal(map[string]any{"results": ranked})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal results: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

// HandlePreprocessCode cleans and formats code using the generate model.
func (h *Handlers) HandlePreprocessCode(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	code, err := request.RequireString("code")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("missing required parameter: %v", err)), nil
	}

	prompt := fmt.Sprintf("Format and clean this code. Return only the cleaned code:\n%s", code)
	result, err := h.ollama.Generate(ctx, h.reasoningModel, prompt)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("preprocessing failed: %v", err)), nil
	}

	data, err := json.Marshal(map[string]any{"processed_code": result})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal result: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}
