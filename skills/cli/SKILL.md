---
name: cli
description: Build, refactor, or review command-line interfaces and terminal applications using human-first CLI design. Use this whenever the user mentions a CLI, command, subcommand, flags, options, arguments, help text, shell completion, stdout vs stderr, JSON or plain output, prompts, confirmations, config files, environment variables, exit codes, terminal UX, or a TUI. Prefer this skill for Python CLIs with Typer, for Go CLIs with Cobra or the standard library `flag` package, and for Go terminal UX using Charm tools such as Bubble Tea, Lip Gloss, Huh, and Charm Log. Use it even when the user does not say "CLI" explicitly but is clearly asking for a command-line tool, terminal workflow, or command-driven automation surface.
source: https://clig.dev/
---

# CLI

Use this skill for command-line programs, terminal commands, and terminal UX work.

Start by inspecting the existing CLI shape before editing anything:

- Find the current entrypoints, subcommands, help text, and tests.
- Detect the language and framework already in use.
- Preserve established command names and output contracts unless the user asks for a breaking change.
- Prefer the smallest change that makes the CLI clearer, safer, and more scriptable.

## Core rules

Apply these defaults from `clig.dev` unless the repository already has a deliberate pattern:

- Use a real argument parsing library when the stack provides one.
- Return exit code `0` on success and non-zero on failure.
- Send primary output and machine-readable output to `stdout`.
- Send logs, prompts, warnings, progress, and errors to `stderr`.
- Support `-h` and `--help`, including on subcommands.
- Include a root-level `--version` unless the repository already has a different version-reporting pattern.
- When a command needs input and is run with no args, show concise help or guide the user instead of failing silently.
- Prefer flags over positional arguments except for obvious, stable cases.
- Reuse standard flag names when possible: `--help`, `--version`, `--json`, `--quiet`, `--dry-run`, `--force`, `--output`, `--no-input`, `--no-color`.
- If human-friendly formatting would break pipelines, add `--plain` or `--json`.
- Only prompt interactively when `stdin` is a TTY.
- Never require prompts. Every interactive input needs a non-interactive flag, argument, or `stdin` path.
- Confirm destructive actions. Make severe actions deliberately hard to confirm.
- Respect terminal capabilities: disable color and animations for non-TTY output.
- When output is interactive, prefer rich, colorful presentation with a consistent semantic theme instead of ad hoc ANSI codes.
- Keep success output brief, but not so quiet that the command feels hung.
- If work is likely to take longer than about 1 second, show a spinner on `stderr` in interactive mode.
- For long-running interactive work, it is good to rotate short tips or context while the spinner runs.
- Rewrite expected errors into useful, actionable messages.

## Workflow

Follow this order:

1. Identify whether the task is about parsing, command structure, output, interactivity, or terminal presentation.
2. Inspect the current implementation and follow the repo's existing parser and test conventions unless the user asked to change them.
3. Choose the framework guidance below.
4. Implement the behavior and update help text, examples, or docs when the command surface changes.
5. Verify the user-visible behavior, not just the internals.

If the task includes interactive styling, themes, or color systems, also read `references/opencode-themes.md`.

Minimum verification for CLI changes:

- Run the relevant help commands.
- Run at least one success path.
- Run at least one failure path and confirm the exit code and error text make sense.
- If the command claims to support piping or machine-readable output, verify `--json` or `--plain` behavior.
- If prompts are involved, verify non-interactive behavior too.

## Framework selection

### Python

If the project already uses Typer, or the user wants a new Python CLI, read `references/python-typer.md` and follow it closely.

### Go

If the project is Go and the user wants terminal UX, prompts, forms, styled output, TUIs, Cobra command structure, or a general Go CLI implementation, read `references/go-cli.md`.

For Go argument parsing:

- Keep the repo's existing parser if one already exists.
- If the task is greenfield and small, use the standard library `flag` package with explicit subcommand dispatch.
- Use Charm libraries for terminal UX, not as a replacement for argument parsing.

## Interaction design guidance

Design for both humans and shells:

- Human mode should be readable, reassuring, and suggest the next step when useful.
- Script mode should be stable, quiet, and easy to parse.
- If output differs by mode, make the switch explicit with flags such as `--json` and `--plain`.
- If a command changes remote or hidden state, tell the user what changed.
- If the command needs configuration, use clear precedence: flags, environment, project config, user config, system config.
- If the CLI has durable configuration, make it discoverable through a `config` command instead of forcing users to edit files manually.

## Interactive color and progress guidance

For interactive CLIs, plain monochrome output should be the fallback, not the design target.

- Use color to make the interface easier to scan: hierarchy, status, focus, warnings, success, and next steps.
- Use semantic theme tokens such as primary, accent, success, warning, error, muted text, and background surfaces.
- Keep theme IDs stable if the CLI exposes theme selection. Prefer the OpenCode theme set in `references/opencode-themes.md`.
- Do not rely on color alone to communicate critical state.
- Disable decorative color and animation for non-TTY output, `--plain`, `--json`, or `--no-color`.
- If the command will likely take longer than about 1 second, start a spinner quickly so the CLI never feels stalled.
- Put spinners and progress messages on `stderr`, not `stdout`.
- Spinner labels should say what the tool is doing now, in plain language.
- Prefer a braille snake spinner for interactive progress unless the repository already uses another spinner style.
- For longer tasks, it is good to rotate short tips, assumptions, or the next likely step while the spinner runs.
- Stop the spinner cleanly before printing the final result or any blocking prompt.

## Common flags and commands

Treat these as defaults for new or meaningfully revised CLIs unless the repo already has a deliberate alternative:

- Always provide `-h` and `--help` on the root command and subcommands.
- Include root-level `--version`, and add a `version` subcommand when that matches the framework or existing command style.
- Consider adding an `explain` command for non-trivial CLIs, especially when agents or automation need to understand the command surface programmatically.
- Include `--no-input` on any command that may prompt or open an interactive selector.
- Include `--json` or `--plain` whenever output may be piped, parsed, or decorated.
- Include `--verbose` and `--quiet` for non-trivial CLIs that emit diagnostics, progress, or hints.
- Include shell completion support for non-trivial CLIs when the framework supports it.
- Include `--dry-run` and `--force` on commands with meaningful side effects when those concepts fit the operation.
- Prefer `--yes` or `--assume-yes` for skipping confirmation prompts, and reserve `--force` for bypassing safety checks or overwrite protections.
- If the CLI depends on user-configurable settings, add a `config` command with `show`, `get`, `set`, `unset`, `edit`, and `path` subcommands unless the repository already has a deliberate alternative.
- Make bare `config` behave like `config show` unless the repository already has a stronger existing pattern.
- It is fine to add `view` as a hidden or undocumented alias for `show` when that improves discoverability or preserves compatibility.

Behavior expectations:

- `--verbose` and `--quiet` should affect diagnostics on `stderr`, not the primary result on `stdout`.
- Do not add both `--json` and decorative human output to `stdout` at the same time.
- `--version` output should be brief and script-friendly by default.
- Help output should include useful examples for non-trivial commands.
- If two common flags conflict, reject the combination or define precedence explicitly.
- `help` is for humans; `explain` is for agents and automation.
- `explain` output should be structured and schema-like, closer to `openapi.json` than to rendered help text.
- Prefer machine-readable formats such as JSON or TOML; if the product has a custom agent-oriented format such as `toon`, make it explicit and stable.
- `explain` should describe commands, subcommands, flags, arguments, config surfaces, output formats, and important behavioral contracts.
- `explain` should have a versioned schema and stable top-level fields.
- `config show` should make the active config easy to inspect.
- bare `config` should default to showing the active config
- `config get <key>` should print a single resolved value that works well in scripts.
- `config set <key> <value>` should update the intended config scope and say what changed.
- `config unset <key>` should remove an explicit value and make fallback behavior clear.
- `config edit` should open the correct config target in the user's editor.
- `config path` should print the path to the active config file or target.

Suggested `explain` contract:

- use `snake_case` field names consistently
- include a top-level schema version such as `schema_version`
- include the CLI name and version
- include canonical commands keyed by canonical command path
- include primary aliases separately from hidden compatibility aliases
- include arguments, options, defaults, enum values, required inputs, output modes, and interactivity behavior
- include config commands, config keys, config scopes, and environment variable mappings when relevant
- include important behavioral contracts such as whether a command mutates state, prompts, or supports machine output
- keep the payload stable enough that agents do not need to scrape help output

Recommended shape:

- use an OpenAPI-like top-level envelope for metadata and versioning
- use MCP-like per-command records with names, descriptions, and schemas
- use JSON Schema for input and output structures
- keep completion metadata separate from validation metadata when needed

Suggested top-level fields:

- `schema_version`
- `kind`
- `metadata`
- `commands`
- `components`
- `annotations`
- `completion`
- `examples`
- `external_docs`

For secret-bearing config and environment metadata, prefer explicit fields such as:

- `secret: true`
- `redacted: true`
- `reveal_requires_opt_in: true`
- `sources: ["env", "user_config"]`

Canonical JSON example:

```json
{
  "schema_version": "1.0",
  "kind": "cli",
  "metadata": {
    "name": "myapp",
    "version": "1.4.2",
    "summary": "Deployment and operations CLI"
  },
  "commands": [
    {
      "name": "deploy",
      "path": ["myapp", "deploy"],
      "summary": "Deploy a service",
      "description": "Create or update a deployment.",
      "aliases": {
        "primary": [],
        "hidden": []
      },
      "input_schema": {
        "type": "object",
        "properties": {
          "env": {
            "type": "string",
            "enum": ["dev", "staging", "prod"],
            "description": "Deployment environment"
          },
          "no_input": {
            "type": "boolean",
            "description": "Fail instead of prompting"
          }
        },
        "required": ["env"]
      },
      "output_schema": {
        "type": "object",
        "properties": {
          "deployment_id": {"type": "string"},
          "status": {"type": "string"}
        },
        "required": ["deployment_id", "status"]
      },
      "annotations": {
        "interactive": true,
        "destructive": false,
        "requires_network": true,
        "safe_for_automation": true
      },
      "completion": {
        "options": {
          "--env": {
            "values": ["dev", "staging", "prod"]
          }
        }
      },
      "examples": [
        "myapp deploy --env prod --no-input",
        "myapp deploy --env staging --json"
      ]
    }
  ],
  "config": {
    "commands": {
      "show": true,
      "get": true,
      "set": true,
      "unset": true,
      "edit": true,
      "path": true
    },
    "keys": [
      {
        "name": "profile",
        "secret": false,
        "scopes": ["user", "project"]
      },
      {
        "name": "api_token",
        "secret": true,
        "redacted": true,
        "reveal_requires_opt_in": true,
        "sources": ["env", "user_config"],
        "env_var": "MYAPP_API_TOKEN"
      }
    ]
  }
}
```

Canonical TOML example:

```toml
schema_version = "1.0"
kind = "cli"

[metadata]
name = "myapp"
version = "1.4.2"
summary = "Deployment and operations CLI"

[[commands]]
name = "deploy"
path = ["myapp", "deploy"]
summary = "Deploy a service"
description = "Create or update a deployment."
examples = [
  "myapp deploy --env prod --no-input",
  "myapp deploy --env staging --json"
]

[commands.aliases]
primary = []
hidden = []

[commands.input_schema]
type = "object"
required = ["env"]

[commands.input_schema.properties.env]
type = "string"
enum = ["dev", "staging", "prod"]
description = "Deployment environment"

[commands.input_schema.properties.no_input]
type = "boolean"
description = "Fail instead of prompting"

[commands.output_schema]
type = "object"
required = ["deployment_id", "status"]

[commands.output_schema.properties.deployment_id]
type = "string"

[commands.output_schema.properties.status]
type = "string"

[commands.annotations]
interactive = true
destructive = false
requires_network = true
safe_for_automation = true

[commands.completion.options."--env"]
values = ["dev", "staging", "prod"]

[config.commands]
show = true
get = true
set = true
unset = true
edit = true
path = true

[[config.keys]]
name = "profile"
secret = false
scopes = ["user", "project"]

[[config.keys]]
name = "api_token"
secret = true
redacted = true
reveal_requires_opt_in = true
sources = ["env", "user_config"]
env_var = "MYAPP_API_TOKEN"
```

## Interactivity guidance

Treat prompts as an enhancement, not the only path:

- Prompt only when `stdin` is interactive and `--no-input` is not set.
- If required input is missing in non-interactive mode, fail fast and say which flag or argument to pass.
- For secrets, accept `stdin` or a file path, not raw CLI flags.
- For severe destructive actions, ask for an explicit typed confirmation or a clear `--confirm` value.

## Generous hints

Make the CLI feel helpful and conversational, especially when the result depends on context the user may not know.

- Surface assumptions before acting when config, environment variables, defaults, or auto-detected state can affect the result.
- Call out config knobs that may explain surprising output or behavior.
- If you use fuzzy pickers, interactive selection, or inferred defaults, explain what happened and what other options were available.
- Always provide a `--no-input` path so scripts and careful users can fail early instead of getting stuck in prompts.
- In interactive mode, do your best to guide the user through valid options instead of forcing them to restart.
- If you prompt for something that could have been supplied via flags or arguments, end by showing the fully qualified equivalent command.

Examples of the kind of hints to prefer:

- "Using profile `staging` from `~/.config/myapp/config.toml`. Pass `--profile prod` to override."
- "No project was specified, so I picked the closest match: `payments-api`. Pass `--project payments-api` to skip selection next time."
- "Interactive mode is available here, but `--no-input` will fail fast if required values are missing."
- "Equivalent non-interactive command: `myapp deploy --env staging --region us-east-1 --service payments-api`"

## Output guidance

Use these conventions unless the repo already has a stronger contract:

- `stdout`: command results, machine-readable payloads, line-oriented output for pipes.
- `stderr`: diagnostics, warnings, progress, confirmations, and human-facing errors.
- `--json`: structured output with stable field names.
- `--plain`: stable, line-oriented output without decoration.
- Human formatting may use color and layout only when writing to a TTY.
- Spinners, progress notes, and rotating tips belong on `stderr` and should never contaminate `stdout` payloads.

For large human-oriented output:

- respect `PAGER` when the repository already does so or when paging is clearly useful
- do not invoke a pager for `--json`, `--plain`, non-TTY output, or when output is being piped
- keep paging optional and predictable

## Secrets and redaction

Never leak secrets through config, explain output, or diagnostics.

- `config show` should redact secret values by default
- `config get` should avoid printing secrets unless the command is explicitly designed for secure retrieval
- `config edit` and `config path` should avoid printing secret contents accidentally
- `explain` should describe secret-bearing fields as metadata, not emit live secret values
- diagnostics and hints should refer to secret sources safely, for example `Using token from environment variable` rather than printing the token
- if the CLI supports revealing secret values, require an explicit and clearly named opt-in

## Aliases and deprecation

Use aliases to preserve compatibility, but make migration behavior clear.

- keep canonical commands primary in docs and examples
- when renaming a command, it is reasonable to keep the old name as a hidden alias for a transition period
- when a hidden deprecated alias is used, prefer a short hint pointing users to the canonical command
- do not keep compatibility aliases forever without a reason

## What to avoid

- Do not print stack traces for expected user mistakes.
- Do not send log chatter to `stdout`.
- Do not add prompts without a non-interactive alternative.
- Do not introduce full-screen terminal UI for simple one-shot commands.
- Do not redesign command names, flags, or output formats unless necessary for the task.

## Deliverables

When you finish, make sure the result includes the user-visible pieces that matter:

- updated CLI code
- updated help text and examples when the interface changed
- tests for parsing, exit behavior, and key output modes when the repo has CLI tests
- a concise explanation of the behavior change and any new flags or subcommands
