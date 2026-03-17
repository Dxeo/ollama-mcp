package main

import (
	"log/slog"
	"os"

	"github.com/mark3labs/mcp-go/server"
)

const (
	envReasoningModel = "REASONING_MODEL"
	envEmbeddingModel = "EMBEDDING_MODEL"
)

func main() {
	reasoningModel := os.Getenv(envReasoningModel)
	embeddingModel := os.Getenv(envEmbeddingModel)

	if reasoningModel == "" {
		slog.Error("environment variable is required", "var", envReasoningModel)
		os.Exit(1)
	}
	if embeddingModel == "" {
		slog.Error("environment variable is required", "var", envEmbeddingModel)
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

	slog.Info("starting ollama-mcp server", "reasoning_model", reasoningModel, "embedding_model", embeddingModel)

	if err := server.ServeStdio(s); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}

	slog.Info("server stopped")
}
