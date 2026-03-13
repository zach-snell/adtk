# adtk - Azure DevOps Toolkit

A unified CLI and MCP (Model Context Protocol) server for Azure DevOps. Single Go binary, PAT-first auth, no admin approval required.

## Features

- **Dual-mode**: Use as a CLI tool or MCP server for AI agents
- **PAT-first auth**: Self-service Personal Access Tokens — no Azure AD admin approval needed
- **Single binary**: No Node.js, Python, or Docker required
- **Consolidated tools**: ~12 MCP tools covering work items, repos, PRs, pipelines, boards, wiki
- **Response flattening**: Strips `_links` and converts `System.` field prefixes to readable names
- **Rate limiting**: Built-in token bucket respecting ADO's TSTU-based limits

## Quick Start

### Install

```bash
go install github.com/zach-snell/adtk/cmd/adtk@latest
```

### Authenticate

```bash
adtk auth
```

Or use environment variables:

```bash
export AZURE_DEVOPS_ORG=myorg
export AZURE_DEVOPS_PAT=your-pat-here
```

### CLI Usage

```bash
# List projects
adtk projects list

# Get a specific project
adtk projects get MyProject
```

### MCP Server

```bash
# stdio mode (for AI agents)
adtk mcp

# HTTP Streamable mode
adtk mcp --port 8080
```

## MCP Configuration

Add to your MCP client configuration:

```json
{
  "mcpServers": {
    "adtk": {
      "command": "adtk",
      "args": ["mcp"],
      "env": {
        "AZURE_DEVOPS_ORG": "myorg",
        "AZURE_DEVOPS_PAT": "your-pat-here"
      }
    }
  }
}
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `AZURE_DEVOPS_ORG` | Azure DevOps organization name |
| `AZURE_DEVOPS_PAT` | Personal Access Token |
| `ADTK_ENABLE_WRITES` | Set to `true` to enable write operations |
| `AZURE_DEVOPS_DISABLED_TOOLS` | Comma-separated list of tools to disable |

## Competitive Comparison

| Feature | adtk | microsoft/azure-devops-mcp | Tiberriver256 |
|---------|------|---------------------------|---------------|
| Auth | PAT (self-service) | Azure AD only (admin) | PAT |
| CLI | Yes | No | No |
| Runtime | Go binary | Node.js | Node.js |
| Tools | ~12 consolidated | 82 tools | 35 tools |
| Boards | Yes | No | No |

## License

Apache 2.0 - see [LICENSE](LICENSE)
