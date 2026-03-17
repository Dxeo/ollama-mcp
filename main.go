package main

import (
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	reasoningModel := os.Getenv("REASONING_MODEL")
	embeddingModel := os.Getenv("EMBEDDING_MODEL")

	if reasoningModel == "" {
		fmt.Fprintln(os.Stderr, "error: REASONING_MODEL environment variable is required")
		os.Exit(1)
	}
	if embeddingModel == "" {
		fmt.Fprintln(os.Stderr, "error: EMBEDDING_MODEL environment variable is required")
		os.Exit(1)
	}

	ollamaClient, err := NewOllamaClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to initialize ollama client: %v\n", err)
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

	fmt.Fprintln(os.Stderr, "ollama-mcp server starting...")
	fmt.Fprintf(os.Stderr, "  reasoning model: %s\n", reasoningModel)
	fmt.Fprintf(os.Stderr, "  embedding model: %s\n", embeddingModel)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "error: server failed: %v\n", err)
		os.Exit(1)
	}
}
