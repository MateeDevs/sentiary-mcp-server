# Sentiary MCP

Go MCP stdio server for project string management through the Sentiary REST API.

## Install

Install the latest release binary:

```sh
curl -fsSL https://raw.githubusercontent.com/MateeDevs/sentiary-mcp-server/main/install.sh | sh
```

Or install from source with Go:

```sh
go install github.com/MateeDevs/sentiary-mcp-server@latest
```

`go install` places `sentiary-mcp-server` in `$GOBIN` or `$GOPATH/bin`. Use `sentiary-mcp-server` as the MCP command when installing this way. The release installer installs the binary as `sentiary-mcp-server`.

## Configuration

Configure credentials through environment variables in the MCP client config:

- `SENTIARY_USER_API_KEY` — required. Use the current user's project API key from `/project/{projectId}/userApiKey`; sent as `Authorization: Ribbon <key>`.
- `SENTIARY_PROJECT_ID` — optional default project ID.

The API base URL is fixed to `https://api.sentiary.com/`.

Example:

```json
{
  "mcpServers": {
    "sentiary": {
      "command": "/path/to/sentiary/mcp/sentiary-mcp-server",
      "env": {
        "SENTIARY_PROJECT_ID": "project-id",
        "SENTIARY_USER_API_KEY": "user-project-api-key"
      }
    }
  }
}
```

## Tools

String management:

- `list_project_strings`
- `search_project_string_by_name`
- `get_project_string`
- `get_project_string_by_key`
- `add_project_string`
- `edit_project_string`
- `remove_project_string`
- `set_project_string_translation`
- `remove_project_string_translation`

Batch operations:

- `get_project_strings_info`
- `export_project_strings`
- `import_project_strings`

All tools use `Authorization: Ribbon <api-key>`. For string-management tools, use a user API key so Sentiary resolves the request to the normal user principal. Project-wide API keys still exist for batch-style automation but do not carry user identity. Supported batch formats: `json`, `android`, `apple`, `compose`.

## Agent Skill

This repository includes a reusable agent skill for localization workflows:

- `sentiary-localization` — use Sentiary MCP tools to add, edit, remove, import, export, and sync project strings with local localization files.

The skill lives at `.agents/skills/sentiary-localization/SKILL.md` and is intended to be versioned with the MCP server repository. It can be adapted by any MCP-capable agent workflow that supports reusable instruction files.

### Add the skill to an agent

1. Configure this MCP server in your AI client with `SENTIARY_PROJECT_ID` and `SENTIARY_USER_API_KEY`.
2. Copy or reference `.agents/skills/sentiary-localization/SKILL.md` in your agent's reusable instructions/skills directory.
3. Tell the agent to use the `sentiary-localization` skill when working with project strings or localization files.

For agents without native skill support, paste the contents of `SKILL.md` into the agent's project instructions or custom system instructions.

## Build

```sh
go build -o sentiary-mcp-server .
```

## Docker

```sh
docker build -t sentiary-mcp-server .
docker run --rm -p 8080:8080 \
  -e SENTIARY_PROJECT_ID=project-id \
  -e SENTIARY_USER_API_KEY=user-project-api-key \
  sentiary-mcp-server
```

The container runs streamable HTTP MCP on `$PORT` at `/` and exposes `GET /healthz`.

### Per-user API keys over HTTP

A single hosted HTTP instance can serve many users, each supplying their own
Sentiary API key from their MCP client instead of sharing the server's
environment key. On session initialize the server reads the key from the request
in this order:

1. `X-Sentiary-User-Api-Key: <key>`
2. `Authorization: Bearer <key>` (also accepts `Ribbon <key>` or a bare `<key>`)
3. Falls back to `SENTIARY_USER_API_KEY` from the server environment when no
   per-request key is sent.

This means `SENTIARY_USER_API_KEY` is optional when clients provide their own
key, and `SENTIARY_PROJECT_ID` can be omitted so each user targets their own
project via the per-tool `projectId` argument (or their own default).

Example MCP client config against a hosted server:

```json
{
  "mcpServers": {
    "sentiary": {
      "type": "http",
      "url": "https://sentiary-mcp.example.com/",
      "headers": {
        "Authorization": "Bearer user-project-api-key"
      }
    }
  }
}
```

> Note: the key is only read when the MCP session is established. Each user gets
> an isolated session bound to their own key.

