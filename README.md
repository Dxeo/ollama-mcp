# Ollama MCP Server

An MCP (Model Context Protocol) server that routes AI tasks to local [Ollama](https://ollama.com/) models. Provides private, cost-free AI capabilities — reasoning, code generation, embeddings, and semantic search — directly from your machine.

Works with any MCP-compatible agent or client, including [Claude Code](https://docs.anthropic.com/en/docs/claude-code), [Cursor](https://cursor.com/), [Cline](https://github.com/cline/cline), and others.

> This project was inspired by and built upon ideas from [jonsflow/ollama-mcp](https://github.com/jonsflow/ollama-mcp). Full credit to the original author for the concept.

## Features

- **Local-first AI** — all inference runs on your hardware via Ollama, no API keys or cloud calls
- **3 MCP tools** — reasoning, embeddings, and document filtering
- **MCP standard** — works with any agent or client that supports the Model Context Protocol
- **Configurable context window** — set the model's context size per environment
- **Structured logging** — JSON-formatted logs via Go's `slog`
- **Retry logic** — automatic retries with exponential backoff for Ollama API calls
- **Cross-platform** — builds for Linux, macOS (Intel + Apple Silicon), and Windows

## Quick Start

### 1. Install Ollama

Download and install from [ollama.com](https://ollama.com/).

### 2. Pull models

```bash
# Reasoning model (pick one)
ollama pull qwen3:1.7b

# Embedding model
ollama pull nomic-embed-text
```

### 3. Build

```bash
git clone https://github.com/Dxeo/ollama-mcp.git
cd ollama-mcp
go build -o ollama-mcp .        # Linux / macOS
go build -o ollama-mcp.exe .    # Windows
```

Or download a prebuilt binary from [Releases](https://github.com/Dxeo/ollama-mcp/releases).

## Configuration

All configuration is via environment variables:

| Variable | Required | Default | Description |
|---|---|---|---|
| `REASONING_MODEL` | Yes | — | Ollama model for reasoning and code generation (e.g. `qwen3:1.7b`) |
| `EMBEDDING_MODEL` | Yes | — | Ollama model for embeddings and document filtering (e.g. `nomic-embed-text`) |
| `CONTEXT_SIZE` | Yes | — | Context window size in tokens for the reasoning model (e.g. `16384`) |
| `OLLAMA_HOST` | No | `http://localhost:11434` | Ollama server address |

The server exits with an error if any required variable is missing.

## Integration

The server communicates over **stdio** using the MCP protocol. Any MCP-compatible client can use it by pointing to the binary and setting the required environment variables.

### Claude Code

Add to `~/.claude/settings.json` or project `.mcp.json`:

```json
{
  "mcpServers": {
    "ollama-mcp": {
      "command": "/absolute/path/to/ollama-mcp",
      "env": {
        "REASONING_MODEL": "qwen3:1.7b",
        "EMBEDDING_MODEL": "nomic-embed-text",
        "CONTEXT_SIZE": "16384"
      }
    }
  }
}
```

This repo ships a [`AGENT.md`](AGENT.md) file with instructions that tell Claude to actively use the ollama-mcp tools. Copy its contents into your project's `AGENT.md` or reference it in your Claude Code settings.

### Other MCP Clients

For any MCP-compatible client, configure a stdio server with the binary path and environment variables above. Refer to your client's documentation for the exact configuration format.

## Tools

| Tool | Description | Model |
|---|---|---|
| `reason_task` | Reasoning, code generation, summarization | `REASONING_MODEL` |
| `embed_text` | Generate text embeddings | `EMBEDDING_MODEL` |
| `filter_docs` | Rank documents by semantic similarity to a query | `EMBEDDING_MODEL` |

## License

This project is licensed under the GNU General Public License v3.0 — see the [LICENSE](LICENSE) file for details.
