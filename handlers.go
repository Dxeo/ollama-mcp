package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// Handlers holds the shared state needed by tool handlers.
type Handlers struct {
	ollama         *OllamaClient
	reasoningModel string
	embeddingModel string
	contextSize    int
}

// HandleReasonTask performs reasoning/code generation using the chat model.
func (h *Handlers) HandleReasonTask(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	slog.Info("handling tool call", "tool", "reason_task")

	prompt, err := request.RequireString("prompt")
	if err != nil {
		slog.Error("invalid parameter", "tool", "reason_task", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("missing required parameter: %v", err)), nil
	}
	if strings.TrimSpace(prompt) == "" {
		slog.Error("empty input", "tool", "reason_task")
		return mcp.NewToolResultError("prompt must not be empty"), nil
	}

	result, err := h.ollama.Chat(ctx, h.reasoningModel, prompt, h.contextSize)
	if err != nil {
		slog.Error("ollama call failed", "tool", "reason_task", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("reasoning failed: %v", err)), nil
	}

	return mcp.NewToolResultText(result), nil
}

// HandleEmbedText generates embeddings for the given text.
func (h *Handlers) HandleEmbedText(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	slog.Info("handling tool call", "tool", "embed_text")

	text, err := request.RequireString("text")
	if err != nil {
		slog.Error("invalid parameter", "tool", "embed_text", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("missing required parameter: %v", err)), nil
	}
	if strings.TrimSpace(text) == "" {
		slog.Error("empty input", "tool", "embed_text")
		return mcp.NewToolResultError("text must not be empty"), nil
	}

	embeddings, err := h.ollama.Embed(ctx, h.embeddingModel, text)
	if err != nil {
		slog.Error("ollama call failed", "tool", "embed_text", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("embedding failed: %v", err)), nil
	}

	data, err := json.Marshal(map[string]any{"embeddings": embeddings})
	if err != nil {
		slog.Error("json marshal failed", "tool", "embed_text", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal embeddings: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

// HandleFilterDocs ranks documents by semantic similarity to a query.
func (h *Handlers) HandleFilterDocs(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	slog.Info("handling tool call", "tool", "filter_docs")

	query, err := request.RequireString("query")
	if err != nil {
		slog.Error("invalid parameter", "tool", "filter_docs", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("missing required parameter: %v", err)), nil
	}

	documents := request.GetStringSlice("documents", nil)
	if documents == nil || len(documents) == 0 {
		slog.Error("empty or missing documents", "tool", "filter_docs")
		return mcp.NewToolResultError("documents list is empty or missing"), nil
	}

	queryEmb, err := h.ollama.Embed(ctx, h.embeddingModel, query)
	if err != nil {
		slog.Error("failed to embed query", "tool", "filter_docs", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("failed to embed query: %v", err)), nil
	}

	type scored struct {
		doc   string
		score float64
	}
	results := make([]scored, 0, len(documents))

	for i, doc := range documents {
		docEmb, err := h.ollama.Embed(ctx, h.embeddingModel, doc)
		if err != nil {
			slog.Warn("skipping document, embedding failed", "tool", "filter_docs", "index", i, "error", err)
			continue
		}
		sim := CosineSimilarity(queryEmb, docEmb)
		results = append(results, scored{doc: doc, score: sim})
	}

	if len(results) == 0 {
		slog.Error("all documents failed to embed", "tool", "filter_docs", "total", len(documents))
		return mcp.NewToolResultError("all documents failed to embed"), nil
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
		slog.Error("json marshal failed", "tool", "filter_docs", "error", err)
		return mcp.NewToolResultError(fmt.Sprintf("failed to marshal results: %v", err)), nil
	}

	return mcp.NewToolResultText(string(data)), nil
}

