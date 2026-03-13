# adtk - Azure DevOps Toolkit

A unified CLI and MCP (Model Context Protocol) server for Azure DevOps. Single Go binary, PAT-first auth, no admin approval required.

## Features

- **Dual-mode**: Use as a CLI tool or MCP server for AI agents
- **PAT-first auth**: Self-service Personal Access Tokens — no Azure AD admin approval needed
- **Single binary**: No Node.js, Python, or Docker required
- **11 consolidated MCP tools** covering work items, repos, PRs, pipelines, boards, wiki, and more
- **Full CLI** with table-formatted output and `--json` flag
- **Response flattening**: Strips `_links` and converts `System.` field prefixes to readable names
- **Rate limiting**: Built-in token bucket respecting ADO's TSTU-based limits
- **Write protection**: Write operations gated behind `ADTK_ENABLE_WRITES=true`

## Quick Start

### Install

```bash
go install github.com/zach-snell/adtk/cmd/adtk@latest
```

Or build from source:

```bash
git clone https://github.com/zach-snell/adtk.git
cd adtk
go build -o adtk ./cmd/adtk
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
# Projects & teams
adtk projects list
adtk projects get MyProject
adtk projects teams MyProject

# Work items
adtk work-items list -p MyProject
adtk work-items get 42
adtk work-items types -p MyProject

# Repositories
adtk repos list -p MyProject
adtk repos get myrepo -p MyProject
adtk repos branches myrepo -p MyProject
adtk repos tree myrepo /src -p MyProject

# Pull requests
adtk pull-requests list myrepo -p MyProject
adtk pull-requests get myrepo 1
adtk pull-requests reviewers myrepo 1

# Pipelines
adtk pipelines list -p MyProject
adtk pipelines runs 42 -p MyProject

# Iterations & boards
adtk iterations list -p MyProject
adtk iterations current -p MyProject
adtk boards list -p MyProject
adtk boards columns Stories -p MyProject

# Wiki
adtk wiki list -p MyProject
adtk wiki get ProjectWiki /Home -p MyProject

# Search
adtk search code "func main" -p MyProject
adtk search work-items "login bug" -p MyProject
adtk search wiql "SELECT [System.Id] FROM WorkItems WHERE [System.State] = 'Active'"

# Attachments
adtk attachments list 42 -p MyProject

# All commands support --json for raw JSON output
adtk projects list --json
```

### MCP Server

```bash
# stdio mode (for AI agents)
adtk mcp

# HTTP Streamable mode
adtk mcp --port 8080
```

## MCP Tools

adtk exposes 11 MCP tools with 60+ actions:

| Tool | Actions | Description |
|------|---------|-------------|
| `manage_work_items` | get, batch_get, create, update, delete, add_comment, list_comments, get_links, list_types, get_history | Full work item lifecycle |
| `manage_projects` | list, get, list_teams, get_team, create | Projects and teams |
| `manage_users` | get_current, search | Identity and user lookup |
| `manage_search` | wiql, code, work_items, wiki | Multi-domain search |
| `manage_repos` | list, get, list_branches, get_file, get_tree | Git repositories |
| `manage_pull_requests` | list, get, create, update, add_comment, list_comments, vote, list_reviewers | Pull request management |
| `manage_iterations` | list, get, get_current | Sprint/iteration tracking |
| `manage_boards` | list, get, get_columns | Kanban board management |
| `manage_wiki` | list, get_page, create_page, update_page, delete_page | Markdown-native wiki |
| `manage_pipelines` | list, get, list_runs, get_run, trigger, get_logs, get_log | CI/CD pipeline management |
| `manage_attachments` | list, upload, download | Work item attachments |

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
| `ADTK_ENABLE_WRITES` | Set to `true` to enable write operations (create, update, delete, etc.) |
| `AZURE_DEVOPS_DISABLED_TOOLS` | Comma-separated list of tools to disable (e.g., `manage_wiki,manage_pipelines`) |

## Architecture

- **Custom HTTP client** — direct REST API calls, no third-party SDK
- **Multi-base-URL routing** — `dev.azure.com`, `vssps.dev.azure.com`, `almsearch.dev.azure.com`, `vsrm.dev.azure.com`
- **PAT auth** — `Authorization: Basic base64(":" + pat)` (empty username)
- **JSON Patch** for work item writes — `Content-Type: application/json-patch+json`
- **WIQL 2-step** — query IDs then batch fetch fields (max 200)
- **Team-scoped APIs** — iterations/boards use `{project}/{team}` in URL
- **Response flattener** — `System.Title` → `title`, `Microsoft.VSTS.Common.Priority` → `priority`

## Competitive Comparison

| Feature | adtk | microsoft/azure-devops-mcp | Tiberriver256 |
|---------|------|---------------------------|---------------|
| Auth | PAT (self-service) | Azure AD only (admin) | PAT |
| CLI | Full CLI with 10 command groups | No | No |
| Runtime | Go binary | Node.js | Node.js |
| MCP Tools | 11 consolidated (60+ actions) | 82 tools | 35 tools |
| Boards/Iterations | Yes | No | No |
| Wiki | Yes (markdown-native) | Yes | Yes |
| Attachments | Yes | No | No |
| Write Protection | ADTK_ENABLE_WRITES gate | N/A | N/A |
| Rate Limiting | Built-in token bucket | No | No |

## License

Apache 2.0 - see [LICENSE](LICENSE)
