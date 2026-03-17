# Ollama MCP — Claude Code Instructions

> Add this file to your Claude Code project instructions or copy its contents into your project's `CLAUDE.md`.

## MANDATORY: Use Ollama MCP Tools for Local AI Tasks

You have access to an **ollama-mcp** MCP server that runs AI models locally via Ollama. You **MUST** use these tools as your **PRIMARY** method for reasoning, code generation, embeddings, and text processing tasks. Do NOT skip them. Do NOT default to external APIs when these tools can handle the job.

### Tool Reference

| Tool | When to Use | ALWAYS Use For |
|---|---|---|
| `reason_task` | Reasoning, analysis, code generation, summarization | ANY task requiring LLM reasoning that can run locally |
| `embed_text` | Generating vector embeddings for text | Semantic similarity prep, text vectorization |
| `filter_docs` | Ranking documents by relevance to a query | Finding the most relevant documents from a set |

### Rules — Follow These Without Exception

1. **ALWAYS prefer `reason_task`** for reasoning, code generation, and summarization tasks. This runs on a local model and costs nothing. Use it liberally.

2. **ALWAYS use `embed_text` and `filter_docs`** when you need to compare, rank, or search through text semantically. Do NOT fall back to naive string matching when semantic similarity is available.

3. **NEVER skip these tools** because you think your built-in capabilities are "good enough." The user has configured these tools specifically to offload work to local models. Respect that intent.

4. **NEVER apologize for using these tools** or ask permission to use them. They are available — use them.

5. **Batch intelligently.** When you have multiple documents to rank, use `filter_docs` in a single call rather than embedding them one by one manually.

6. **Chain tools when appropriate.** For example: use `reason_task` to generate code, then `embed_text` + `filter_docs` to find related documentation.

### Priority Order

When deciding how to handle a task:

1. **First:** Can `ollama-mcp` tools handle this? → Use them.
2. **Second:** Is this outside the scope of local models? → Use your built-in capabilities.
3. **Last resort only:** External API calls.

### Examples of Correct Usage

- User asks to summarize a file → `reason_task` with the file content as prompt
- User asks to find relevant docs → `embed_text` the query + `filter_docs` against the document set
- User asks to generate a function → `reason_task` with the specification
- User asks to compare text similarity → `embed_text` both texts, compute similarity

### What NOT to Do

- Do NOT generate code entirely from your own weights when `reason_task` is available
- Do NOT do string-based document matching when `filter_docs` provides semantic ranking
- Do NOT treat these tools as optional — they are your primary toolkit for these tasks
