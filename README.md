# Sentiary MCP

Go MCP stdio server for project string management through the Sentiary REST API.

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
      "command": "/path/to/sentiary/mcp/sentiary-mcp",
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
go build -o sentiary-mcp .
```

## Docker

```sh
docker build -t sentiary-mcp .
docker run --rm -p 8080:8080 \
  -e SENTIARY_PROJECT_ID=project-id \
  -e SENTIARY_USER_API_KEY=user-project-api-key \
  sentiary-mcp
```

The container runs streamable HTTP MCP on `$PORT` at `/` and exposes `GET /healthz`.

