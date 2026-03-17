package main

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/mark3labs/mcp-go/server"
)

const (
	envReasoningModel = "REASONING_MODEL"
	envEmbeddingModel = "EMBEDDING_MODEL"
	envContextSize    = "CONTEXT_SIZE"
)

func main() {
	reasoningModel := os.Getenv(envReasoningModel)
	embeddingModel := os.Getenv(envEmbeddingModel)
	contextSizeStr := os.Getenv(envContextSize)

	if reasoningModel == "" {
		slog.Error("environment variable is required", "var", envReasoningModel)
		os.Exit(1)
	}
	if embeddingModel == "" {
		slog.Error("environment variable is required", "var", envEmbeddingModel)
		os.Exit(1)
	}
	if contextSizeStr == "" {
		slog.Error("environment variable is required", "var", envContextSize)
		os.Exit(1)
	}
	contextSize, err := strconv.Atoi(contextSizeStr)
	if err != nil {
		slog.Error("environment variable must be a valid integer", "var", envContextSize, "value", contextSizeStr, "error", err)
		os.Exit(1)
	}

	ollamaClient, err := NewOllamaClient()
	if err != nil {
		slog.Error("failed to initialize ollama client", "error", err)
		os.Exit(1)
	}

	h := &Handlers{
		ollama:         ollamaClient,
		reasoningModel: reasoningModel,
		embeddingModel: embeddingModel,
		contextSize:    contextSize,
	}

	s := server.NewMCPServer(
		"ollama-mcp",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	s.AddTool(ReasonTaskTool, h.HandleReasonTask)
	s.AddTool(EmbedTextTool, h.HandleEmbedText)
	s.AddTool(FilterDocsTool, h.HandleFilterDocs)
	s.AddTool(PreprocessCodeTool, h.HandlePreprocessCode)

	slog.Info("starting ollama-mcp server", "reasoning_model", reasoningModel, "embedding_model", embeddingModel, "context_size", contextSize)

	if err := server.ServeStdio(s); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped")
}
