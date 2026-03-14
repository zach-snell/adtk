# adtk — Azure DevOps Toolkit

[![CI](https://github.com/zach-snell/adtk/actions/workflows/ci.yml/badge.svg)](https://github.com/zach-snell/adtk/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/zach-snell/adtk)](https://goreportcard.com/report/github.com/zach-snell/adtk)
[![Docs](https://img.shields.io/badge/docs-starlight-blue)](https://zach-snell.github.io/adtk)
[![License](https://img.shields.io/badge/license-Apache%202.0-green)](LICENSE)

A dual-mode Go CLI & MCP server for Azure DevOps. Single binary, PAT-first auth, 14 MCP tools with 91 actions, 4 MCP prompts.

> **The most comprehensive Azure DevOps MCP server.** A single Go binary — CLI for humans, MCP for AI agents.

## Comparison: adtk vs Microsoft azure-devops-mcp

| Feature | adtk | microsoft/azure-devops-mcp |
|---------|------|---------------------------|
| **Language** | Go (single static binary) | TypeScript (Node.js) |
| **MCP Tools / Actions** | 14 tools / 91 actions | ~75 individual tools |
| **MCP Prompts** | 4 built-in prompts | None |
| **CLI mode** | Full CLI with 15 command groups | No |
| **Auth** | PAT (self-service, no admin) | Azure AD (requires admin consent) |
| **Startup** | ~50ms | ~2s |
| **Response flattening** | `System.Title` → `title` | Raw API responses |
| **Boards & Iterations** | Full support | No |
| **Test Plans** | Full support | No |
| **Advanced Security** | Alert listing & details | No |
| **Metrics** | Cycle time, lead time, time-in-status | No |
| **Git branch detection** | Auto-detect work items from branch | No |
| **Branch policies & tags** | List policies, list/create tags | No |
| **Variable groups & environments** | List/get variable groups, environments | No |
| **Saved queries** | Get and run saved queries | No |
| **Attachments** | Upload, download, list | No |
| **Write protection** | `ADTK_ENABLE_WRITES` gate | None |
| **Rate limiting** | Built-in token bucket | None |
| **Binary size** | ~15 MB | `npm install` (~200+ MB) |

## Features

- **Dual-mode** — Full CLI with table output + MCP server for AI agents
- **14 consolidated MCP tools** with 91 actions covering every Azure DevOps domain
- **4 MCP prompts** — sprint_summary, pr_review_digest, pipeline_health, release_readiness
- **Git branch detection** — Auto-detect work item IDs from branch names (e.g., `feature/12345-description`)
- **Work item metrics** — Cycle time, lead time, time-in-status computed from revision history
- **PAT-first auth** — Self-service Personal Access Tokens, no Azure AD admin approval
- **Single binary** — No Node.js, Python, or Docker required
- **Response flattening** — Strips `_links` and converts `System.*` fields to readable names
- **Write protection** — All mutations gated behind `ADTK_ENABLE_WRITES=true`
- **Rate limiting** — Built-in token bucket respecting Azure DevOps TSTU limits
- **Token-optimized** — AI agents get clean, concise payloads (40-60% fewer tokens)

## Installation

### go install

```bash
go install github.com/zach-snell/adtk/cmd/adtk@latest
```

### Build from source

```bash
git clone https://github.com/zach-snell/adtk.git
cd adtk
./install.sh
```

### Pre-built binaries

Download from the [Releases](https://github.com/zach-snell/adtk/releases) page.

## Quick Start

### Authenticate

```bash
adtk auth
# Or use environment variables:
export AZURE_DEVOPS_ORG=myorg
export AZURE_DEVOPS_PAT=your-pat-here
```

### CLI Usage

```bash
# Projects
adtk projects list
adtk projects get MyProject
adtk projects teams MyProject

# Work items
adtk work-items get 42
adtk work-items list -p MyProject

# Repositories
adtk repos list -p MyProject
adtk repos branches myrepo -p MyProject
adtk repos tree myrepo /src -p MyProject
adtk repos policies myrepo -p MyProject
adtk repos tags myrepo -p MyProject

# Pull requests
adtk pull-requests list myrepo -p MyProject
adtk pull-requests get myrepo 1

# Pipelines
adtk pipelines list -p MyProject
adtk pipelines runs 42 -p MyProject
adtk pipelines var-groups -p MyProject
adtk pipelines var-group 1 -p MyProject
adtk pipelines environments -p MyProject

# Iterations & boards
adtk iterations current -p MyProject
adtk boards list -p MyProject
adtk boards columns Stories -p MyProject

# Wiki
adtk wiki list -p MyProject
adtk wiki get ProjectWiki /Home -p MyProject

# Work item metrics
adtk work-items metrics 42 -p MyProject

# Work item auto-detect from git branch
adtk work-items get          # auto-detects work item ID from branch name

# Search
adtk search code "func main" -p MyProject
adtk search work-items "login bug" -p MyProject
adtk search wiql "SELECT [System.Id] FROM WorkItems WHERE [System.State] = 'Active'"
adtk search query "My Saved Queries/Active Bugs" -p MyProject

# Test plans
adtk test-plans list -p MyProject

# Security alerts
adtk security alerts myrepo -p MyProject

# Attachments
adtk attachments list 42 -p MyProject

# All commands support --json for raw output
adtk projects list --json
```

## MCP Server

### stdio mode (for AI agents)

```bash
adtk mcp
```

### HTTP Streamable mode

```bash
adtk mcp --port 8080
```

### MCP Client Configuration

**Claude Desktop / Cursor / Claude Code:**

```json
{
  "mcpServers": {
    "adtk": {
      "command": "adtk",
      "args": ["mcp"],
      "env": {
        "AZURE_DEVOPS_ORG": "myorg",
        "AZURE_DEVOPS_PAT": "your-pat-here",
        "ADTK_ENABLE_WRITES": "true"
      }
    }
  }
}
```

## Environment Variables

| Variable | Required | Description |
|----------|----------|-------------|
| `AZURE_DEVOPS_ORG` | Yes | Azure DevOps organization name |
| `AZURE_DEVOPS_PAT` | Yes | Personal Access Token |
| `ADTK_ENABLE_WRITES` | No | Set to `true` to enable write operations (default: `false`) |
| `AZURE_DEVOPS_DISABLED_TOOLS` | No | Comma-separated list of tools to disable |

## MCP Tools

adtk exposes 14 MCP tools with 91 actions and 4 MCP prompts:

| Tool | Actions | Description |
|------|---------|-------------|
| `manage_work_items` | get, batch_get, create, update, delete, add_comment, update_comment, list_comments, get_links, list_types, get_history, batch_update, add_children, link, unlink, add_artifact_link, my_items, iteration_items | Full work item lifecycle (18 actions) |
| `manage_projects` | list, get, list_teams, get_team, create | Projects and teams (5 actions) |
| `manage_users` | get_current, search | Identity and user lookup (2 actions) |
| `manage_search` | wiql, code, work_items, wiki, get_query, run_query | Multi-domain search + saved queries (6 actions) |
| `manage_repos` | list, get, list_branches, get_file, get_tree, create_branch, search_commits, list_policies, list_tags, create_tag | Git repositories, policies, and tags (10 actions) |
| `manage_pull_requests` | list, get, create, update, add_comment, list_comments, vote, list_reviewers, update_reviewers, create_thread, update_thread, reply_to_comment | Pull request management (12 actions) |
| `manage_iterations` | list, get, get_current, create, get_team_settings | Sprint/iteration tracking (5 actions) |
| `manage_boards` | list, get, get_columns | Kanban board management (3 actions) |
| `manage_wiki` | list, get_page, list_pages, create_page, update_page, delete_page | Markdown-native wiki (6 actions) |
| `manage_pipelines` | list, get, list_runs, get_run, trigger, get_logs, get_log, get_build_changes, list_definitions, list_variable_groups, get_variable_group, list_environments | CI/CD pipelines, variable groups, and environments (12 actions) |
| `manage_test_plans` | list_plans, create_plan, list_suites, create_suite, list_cases, get_test_results | Test plan management (6 actions) |
| `manage_advanced_security` | list_alerts, get_alert | Security alert management (2 actions) |
| `manage_metrics` | get_metrics | Work item lifecycle metrics — cycle time, lead time, time-in-status (1 action) |
| `manage_attachments` | list, upload, download | Work item attachments (3 actions) |

### MCP Prompts

| Prompt | Arguments | Description |
|--------|-----------|-------------|
| `sprint_summary` | `project` (required), `team`, `iteration` | Generate a sprint/iteration status report |
| `pr_review_digest` | `project` (required), `repo` (required) | Generate a PR review digest for a repository |
| `pipeline_health` | `project` (required), `pipeline_id` | Analyze CI/CD pipeline health and failure trends |
| `release_readiness` | `project` (required), `iteration` | Assess release readiness with go/no-go recommendation |

## Security

### Write Protection

adtk is **read-only by default**. All create, update, delete, and trigger operations require:

```bash
export ADTK_ENABLE_WRITES=true
```

### PAT Auth

Uses HTTP Basic Auth with empty username: `Authorization: Basic base64(":" + pat)`. No Azure AD admin consent required.

### Rate Limiting

Built-in token bucket rate limiter (30 tokens, refill 1/2s) prevents hitting Azure DevOps TSTU throttling limits.

## Architecture

- **Custom HTTP client** — Direct REST API calls, no third-party SDK
- **Multi-base-URL routing** — `dev.azure.com`, `vssps.dev.azure.com`, `almsearch.dev.azure.com`, `vsrm.dev.azure.com`
- **JSON Patch** for work item writes — `Content-Type: application/json-patch+json`
- **WIQL 2-step** — Query IDs, then batch fetch fields (max 200 per request)
- **Response flattener** — `System.Title` → `title`, `Microsoft.VSTS.Common.Priority` → `priority`
- **ETag concurrency** — Wiki updates use `If-Match` headers for optimistic concurrency
- **Team-scoped APIs** — Iterations/boards use `{project}/{team}` URL routing

## Development

```bash
# Build
go build -o adtk ./cmd/adtk

# Test
go test -race ./...

# Lint
golangci-lint run ./...

# Vulnerability check
govulncheck ./...
```

See the [development guide](https://zach-snell.github.io/adtk/advanced/development/) for full details.

## Documentation

Full documentation: **[zach-snell.github.io/adtk](https://zach-snell.github.io/adtk)**

## License

Apache 2.0 — see [LICENSE](LICENSE)
