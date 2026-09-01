---
name: sentiary-cli
description: Use sentiary-cli for direct terminal workflows that manage Sentiary projects, languages, members, invitations, strings, and localization files. Do not use it for hosted MCP calls.
---

# Sentiary CLI

Use `sentiary-cli` to manage Sentiary data through the Sentiary API.

## Requirements

- Make sure that `sentiary-cli` is available on `PATH`.
- Use `SENTIARY_USER_API_KEY` for authentication.
- Use `SENTIARY_PROJECT_ID` as the default project when it is available.
- Use `--project-id` when the command must select a different project.
- Never print, store, or commit the API key.
- If the API key is missing, ask the user to configure it. Do not ask the user to paste it.

Run `sentiary-cli <command> --help` when a command option is uncertain.

## Select a Command

Projects:

- Use `list-projects` to list projects for the current user.
- Use `get-project` to read one project and the current user's membership.
- Use `create-project` to create a project.
- Use `edit-project` to rename a project.
- Use `remove-project` to remove a project and its data.

Languages:

- Use `list-languages` to list all languages that Sentiary supports.
- Use `list-project-languages` to list the languages in a project.
- Use `add-project-language` to add a language to a project.
- Use `remove-project-language` to remove a language and its translations.

Members and invitations:

- Use `list-project-members` to list members and roles.
- Use `set-project-member-role` to set a member role.
- Use `remove-project-member` to remove a member.
- Use `list-invitations` to list invitations for the current user.
- Use `create-invitation` to invite a user to a project.
- Use `accept-invitation` or `decline-invitation` to respond to an invitation.
- Use `list-project-invitations` to list outstanding project invitations.
- Use `remove-project-invitation` to revoke a project invitation.

Strings and files:

- Use `list` to read project strings with pagination and optional translation filters.
- Use `search` to find text in keys, translations, descriptions, or context.
- Use `get --key` for an exact key lookup.
- Use `get --string-id` when the internal string ID is known.
- Use `add` to create one string and an optional first translation.
- Use `edit` to change a string name, description, or context.
- Use `remove` to remove one string and all its translations.
- Use `set-translation` to create or replace one translation.
- Use `remove-translation` to remove one translation.
- Use `info` to read the project name, languages, and modification date.
- Use `export` to get a localization file from Sentiary.
- Use `import` to send a localization file to Sentiary.

Use `administrator`, `developer`, or `translator` for a member or invitation role.

## Manage Access and Project Settings

List project members before you change a role or remove a member. List project invitations before you revoke one.

Only remove a project, project language, member, or invitation when the user explicitly requests that change. A project removal deletes project data. A project language removal also deletes translations in that language. A member role change can change project permissions.

When a user wants to leave a project, use `remove-project-member` with the user's ID. The API returns an empty member list after self-removal.

## Read Strings

Search with a general query:

```sh
sentiary-cli search \
  --project-id="project-id" \
  --query="Ahoj" \
  --include-language="cs_CZ,en_US"
```

The search can match translation content even when the string key does not match.

Use the returned `pagination.next` value as `--cursor` to read the next page. Continue until `pagination.next` is empty.

Use an exact key lookup before a targeted change:

```sh
sentiary-cli get \
  --project-id="project-id" \
  --key="welcome_message" \
  --include-language="cs_CZ,en_US"
```

## Change Strings

Resolve an exact key before an edit, removal, or translation change. These commands require the internal string ID.

Set a translation:

```sh
sentiary-cli set-translation \
  --project-id="project-id" \
  --string-id="string-id" \
  --language="cs_CZ" \
  --text="Ahoj"
```

Only remove a string or translation when the user explicitly requests removal. Report that string removal also removes every translation.

Do not use `--override-existing` during an import unless the user explicitly requests replacement.

## Synchronize Files

Use batch export and import for full localization files. Do not run one command for each key.

Export to a file:

```sh
sentiary-cli export \
  --project-id="project-id" \
  --format="json" \
  --language="cs_CZ" \
  --output="strings.json"
```

Import from a file:

```sh
sentiary-cli import \
  --project-id="project-id" \
  --format="json" \
  --language="cs_CZ" \
  --input="strings.json"
```

Supported formats are `json`, `android`, `apple`, and `compose`.

Before an import, inspect the input file and confirm its format and language. Preserve placeholders, escapes, newlines, and message syntax.

Use `-` as the input or output path for a standard stream. Use a file path when the user requests a file change.

## Report Results

Report the project, command, affected users, keys or file, and languages. Include role changes, invitation decisions, skipped imports, and removals when they occur.
