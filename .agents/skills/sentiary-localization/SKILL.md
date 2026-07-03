---
name: sentiary-localization
description: Use Sentiary MCP tools to add, update, remove, and sync project localization strings between local string files and a Sentiary project.
---

# Sentiary Localization

Use this skill when the user asks to manage project localization strings through Sentiary, add new strings, update translations, remove strings, or sync local string files with a Sentiary project.

## Requirements

- The Sentiary MCP server must be configured in the client as `sentiary`.
- The MCP server uses these environment variables:
  - `SENTIARY_PROJECT_ID`
  - `SENTIARY_USER_API_KEY`
- Do not ask for or use bearer tokens.
- Treat Sentiary as the source of truth unless the user explicitly says the local file should override it.

## MCP Tools

Use these Sentiary MCP tools:

- `list_project_strings` — list strings with pagination and optional translations.
- `search_project_string_by_name` — search by string key/name only.
- `get_project_string` — fetch one string by internal ID with translations.
- `get_project_string_by_key` — fetch one string by exact key/name with translations.
- `add_project_string` — create one string and optionally set an initial translation.
- `edit_project_string` — update string name, description, and context.
- `remove_project_string` — remove one string and its translations.
- `set_project_string_translation` — create or update one translation.
- `remove_project_string_translation` — remove one translation.
- `get_project_strings_info` — read project languages and batch metadata.
- `export_project_strings` — export strings in `json`, `android`, `apple`, or `compose` format.
- `import_project_strings` — import strings in `json`, `android`, `apple`, or `compose` format.

## Adding a String

1. Identify the string key, description/context, source language, and translation text.
2. Call `get_project_string_by_key` for the key when checking if it already exists.
3. If the string exists, ask before overwriting unless the user explicitly requested an update.
4. If it does not exist, call `add_project_string`.
5. Add or update translations with `set_project_string_translation`.
6. Report the string key and languages changed.

## Editing a String

1. Find the string with `get_project_string_by_key` when the user provides a key/name.
2. Fetch details with `get_project_string` only when you already have the internal ID.
3. Use `edit_project_string` for name, description, or context changes.
4. Use `set_project_string_translation` for translation text changes.
5. Use `remove_project_string_translation` only when the user clearly requested removal.

## Removing a String

1. Fetch by exact key first with `get_project_string_by_key`.
2. Confirm ambiguity if multiple matches exist.
3. Use `remove_project_string` only for the intended string.
4. Mention that all translations were removed with the string.

## Syncing a Local Strings File

Use this workflow when the user asks to sync a file such as Android XML, Apple `.strings`, Compose resources XML, or JSON localization files.

1. Inspect the local file format and language.
2. Call `get_project_strings_info` to identify project languages when language IDs are unknown.
3. Never search each local key one-by-one. Avoid workflows that issue many database calls.
4. For local file → Sentiary sync, prefer `import_project_strings` with the full file content and correct format.
5. For Sentiary → local file sync, prefer `export_project_strings` and write the exported content to the local file.
6. For small targeted updates, load existing project strings in bulk with `list_project_strings` using pagination and the needed `includeLanguageIds`, then compare locally in memory.
7. Use single-string tools only for a small number of explicit user-requested changes, not for whole-file sync.
8. Never remove Sentiary strings during a sync unless the user explicitly asks for deletion or pruning.
9. Preserve placeholders and formatting exactly, including `%s`, `%d`, `{name}`, `${name}`, XML escapes, newlines, and ICU message syntax.
10. After syncing, summarize:
   - imported or exported file
   - created keys, if known
   - updated translations, if known
   - skipped keys
   - removals, if any

## Conflict Rules

- If the same key exists locally and in Sentiary with different text, ask which source wins unless the user already specified direction.
- If a local key has no matching Sentiary language, stop and ask which project language to use.
- If duplicate local keys are found, resolve duplicates before writing to Sentiary.
- Do not invent translations and do not copy source text into missing translations. If a translation is missing, leave it empty and ask whether the user wants translations generated.

## Local File Handling

- Read the relevant local strings file before changing it.
- Keep output format consistent with the file being edited.
- Make minimal diffs.
- Do not reorder keys unless the existing file is already sorted or the user asks for sorting.
- Validate XML/JSON syntax after edits when possible.
