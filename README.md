# Ollama MCP Server

A Model Context Protocol (MCP) server written in Go that routes AI tasks to local Ollama models. Provides cost-optimized, privacy-focused AI capabilities for Claude Code integration.

## Prerequisites

- [Go](https://go.dev/dl/) 1.22+
- [Ollama](https://ollama.com/) installed and running locally

## Build

```bash
# Clone the repository
git clone https://github.com/your-username/ollama-mcp.git
cd ollama-mcp

# Build the binary
go build -o ollama-mcp .      # Linux/macOS
go build -o ollama-mcp.exe .  # Windows
```

## Configuration

The server is configured via environment variables:

| Variable | Required | Default | Description |
|---|---|---|---|
| `REASONING_MODEL` | Yes | — | Ollama model for reasoning, code gen, and preprocessing |
| `EMBEDDING_MODEL` | Yes | — | Ollama model for embeddings and document filtering |
| `OLLAMA_HOST` | No | `http://localhost:11434` | Ollama server address |

The server will exit with an error if `REASONING_MODEL` or `EMBEDDING_MODEL` are not set.

## Claude Code Integration

Add the following to your Claude Code MCP settings (`~/.claude/settings.json` or project `.mcp.json`):

```json
{
  "mcpServers": {
    "ollama-mcp": {
      "command": "/path/to/ollama-mcp",
      "env": {
        "REASONING_MODEL": "your-reasoning-model",
        "EMBEDDING_MODEL": "your-embedding-model"
      }
    }
  }
}
```

Replace `/path/to/ollama-mcp` with the absolute path to your built binary, and set the model names to whichever Ollama models you have pulled locally.

## Tools

| Tool | Purpose | Model Used |
|---|---|---|
| `reason_task` | Reasoning, code generation, summarization | `REASONING_MODEL` |
| `embed_text` | Generate text embeddings | `EMBEDDING_MODEL` |
| `filter_docs` | Rank documents by semantic similarity to a query | `EMBEDDING_MODEL` |
| `preprocess_code` | Clean and format code | `REASONING_MODEL` |

## License

MIT
