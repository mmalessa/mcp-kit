# mcp-kit

A collection of lightweight MCP (Model Context Protocol) servers written in Go, running over stdio.
Each server is a separate binary built from a single shared Go module.

## Binaries

- **mcp-atlassian** — read-only access to Jira and Confluence
  - `get_jira_task` — fetch an issue by key (`{"issue_key": "PROJ-123"}`)
  - `search_jira` — search issues by JQL (`{"jql": "project = PROJ"}`)
  - `get_jira_attachment` — download an attachment by ID (`{"attachment_id": "10001"}`);
    images are returned as viewable image content, text files as text (max 5 MB)
  - `get_confluence_page` — fetch a page by ID (`{"page_id": "12345"}`)

- **mcp-bitbucket** — read-only access to Bitbucket Cloud (pull requests for code review)
  - `list_pull_requests` — list PRs in a repo (`{"repo": "my-repo", "state": "OPEN"}`)
  - `get_pull_request` — PR metadata (`{"repo": "my-repo", "id": 123}`)
  - `get_pull_request_diff` — unified diff
  - `get_pull_request_diffstat` — per-file change summary
  - `get_pull_request_comments` — comments (inline + general)
  - `get_pull_request_commits` — list of commits

  All tools take `repo` (repository slug) and `id` (PR number). The workspace
  comes from `BITBUCKET_WORKSPACE`; the repo is per-call so one server instance
  can cover multiple repositories in the same workspace.

- **mcp-github** — read-only access to GitHub (pull requests for code review)
  - `list_pull_requests` — list PRs in a repo (`{"repo": "my-repo", "state": "open"}`)
  - `get_pull_request` — PR metadata (`{"repo": "my-repo", "id": 123}`)
  - `get_pull_request_diff` — unified diff
  - `get_pull_request_diffstat` — per-file change summary (lines added/removed per file)
  - `get_pull_request_comments` — review comments + general PR discussion comments
  - `get_pull_request_commits` — list of commits

  All tools take `repo` (repository name) and `id` (PR number). The owner / organization
  comes from `GITHUB_OWNER`; the repo is per-call so one server instance can cover
  multiple repositories under the same owner.

## Project structure

```
mcp-kit/
├── cmd/
│   └── mcp-<name>/
│       ├── main.go              # thin entrypoint
│       └── .env.dist            # (optional) per-binary config template
├── internal/
│   ├── mcpkit/                  # shared bootstrap: env loading, MCP server, add-to-claude/remove-from-claude
│   └── <name>/                  # handlers and logic for a given server
├── bin/                         # built binaries (gitignored)
└── Makefile
```

## Configuration

Each binary reads environment variables from its own configuration file.

**mcp-atlassian** (`~/.config/mcp-kit/mcp-atlassian.env`):

```
ATLASSIAN_EMAIL=your@email.com
ATLASSIAN_API_TOKEN=your_api_token
ATLASSIAN_BASE_URL=https://your-domain.atlassian.net
```

API token: https://id.atlassian.com/manage-profile/security/api-tokens

**mcp-bitbucket** (`~/.config/mcp-kit/mcp-bitbucket.env`):

```
BITBUCKET_EMAIL=your@email.com
BITBUCKET_API_TOKEN=your_scoped_api_token
BITBUCKET_WORKSPACE=your_workspace
```

Bitbucket Cloud requires a **scoped** Atlassian API token. The classic unscoped
token that works with Jira/Confluence (plain `Create API token` button) will
return HTTP 401/403 against the Bitbucket API — you need a separate token.

How to create one:

1. Go to https://id.atlassian.com/manage-profile/security/api-tokens
2. Click **`Create API token with scopes`** (not the plain `Create API token`)
3. Name it, set expiry, click **Next**
4. Select **Bitbucket** as the target application
5. Grant at least these scopes:
   - `read:repository:bitbucket`
   - `read:pullrequest:bitbucket`
6. Create and copy the token immediately — it is shown only once

`BITBUCKET_EMAIL` is the Atlassian account email (the address you use to sign
in at id.atlassian.com), not the Bitbucket username.

**mcp-github** (`~/.config/mcp-kit/mcp-github.env`):

```
GITHUB_TOKEN=your_github_token
GITHUB_OWNER=your_org_or_username
```

Token: GitHub Settings → Developer settings → Personal access tokens.
Fine-grained tokens are recommended; grant at least **Pull requests** (read-only)
and **Issues** (read-only) under the relevant repositories. Classic tokens need
the `repo` scope.

For GitHub Enterprise Server set `GITHUB_API_URL` to your internal API base
(e.g. `https://api.github.my-company.com`).

### `.env` loading order

At startup, the binary loads `.env` from the following locations (in order).
`godotenv` does not overwrite existing variables, so the first file found
that defines a given key wins:

1. `~/.config/mcp-kit/<binary-name>.env` — canonical location, created
   by `make install` (e.g. `~/.config/mcp-kit/mcp-atlassian.env`)
2. `./.env` — current directory (used by `go run ./cmd/<app>` in the repo)

Variables already present in the process environment always take precedence over `.env` files.

## Build & install

```sh
make help                      # list targets + detected binaries
make build                     # all binaries (linux/amd64, stripped)
make build-mcp-atlassian       # single binary
make install                   # build + copy binaries to ~/bin/ + configs in ~/.config/mcp-kit/
make clean                     # remove bin/
make tidy                      # go mod tidy
```

`make install`:
- copies each binary from `bin/` to `~/bin/<app>` (0755)
- creates `~/.config/mcp-kit/<app>.env` (0600) for each:
  - if the file already exists — skips (does not overwrite user config)
  - if the binary has `cmd/<app>/.env.dist` — copies the template
  - otherwise creates an empty file

Make sure `~/bin` is on your `PATH`.

## Registration (binary self-install)

Each binary can register itself in Claude Code or opencode for the current project.
The server name is the binary name without the `mcp-` prefix (`mcp-atlassian` → `atlassian`).

### Claude Code

```sh
cd /home/projects/some-project
mcp-atlassian add-to-claude       # register (uses its own path from os.Executable())
mcp-atlassian remove-from-claude  # unregister
mcp-atlassian help                # show available commands
```

`add-to-claude` is idempotent — if the registration already exists, it is removed first.
The binary passes its actual path (`os.Executable()` + `EvalSymlinks`) to
`claude mcp add`, regardless of CWD.

### opencode

```sh
cd /home/projects/some-project
mcp-atlassian add-to-opencode     # writes entry to ./opencode.json
mcp-atlassian remove-from-opencode
```

`add-to-opencode` writes a `local` server entry directly to `./opencode.json` (project scope).
If the file does not exist it is created. The command is idempotent — it overwrites any
existing entry for the same name.
