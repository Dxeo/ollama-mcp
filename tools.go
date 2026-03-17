package main

import "github.com/mark3labs/mcp-go/mcp"

var ReasonTaskTool = mcp.NewTool("reason_task",
	mcp.WithDescription("Perform reasoning, code generation, or summarization using local Ollama model"),
	mcp.WithString("prompt",
		mcp.Required(),
		mcp.Description("The task prompt for reasoning or code generation"),
	),
)

var EmbedTextTool = mcp.NewTool("embed_text",
	mcp.WithDescription("Generate embeddings for text using local embedding model"),
	mcp.WithString("text",
		mcp.Required(),
		mcp.Description("Text to generate embeddings for"),
	),
)

var FilterDocsTool = mcp.NewTool("filter_docs",
	mcp.WithDescription("Filter and rank documents by semantic similarity to query"),
	mcp.WithString("query",
		mcp.Required(),
		mcp.Description("Query text to match against"),
	),
	mcp.WithArray("documents",
		mcp.Required(),
		mcp.Description("List of documents to filter and rank"),
		mcp.WithStringItems(),
	),
)

