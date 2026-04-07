# Python with Typer

Use this when a Python CLI is built with Typer or when you are adding a new Python CLI.

## Why Typer

Typer gives you:

- type-hint driven arguments and options
- automatic `--help`
- automatic shell completion support
- readable subcommand trees without a lot of boilerplate
- good default errors through Click and Rich

Prefer Typer for new Python CLIs unless the repository already uses another parser.

## Default shape

Use a top-level app with `no_args_is_help=True` for multi-command tools.

```python
from __future__ import annotations

import json
import sys
from pathlib import Path
from typing import Annotated

import typer

app = typer.Typer(no_args_is_help=True)


@app.command()
def show(
    path: Annotated[Path, typer.Argument(exists=True, readable=True)],
    json_output: Annotated[bool, typer.Option("--json")] = False,
) -> None:
    data = {"path": str(path), "size": path.stat().st_size}

    if json_output:
        typer.echo(json.dumps(data))
        return

    typer.echo(f"{data['path']}\t{data['size']}")


if __name__ == "__main__":
    app()
```

For a single-command script, `typer.run(main)` is fine, but switch to `Typer()` once subcommands or shared options appear.

## stdout, stderr, and exit codes

Use `typer.echo()` and keep channels explicit.

```python
typer.echo(result_text)
typer.echo("index updated", err=True)
raise typer.Exit(code=1)
```

Guidelines:

- Send command results to `stdout`.
- Send warnings, prompts, progress notes, and errors to `stderr` with `err=True`.
- Use non-zero exit codes for failures.

For user mistakes, prefer `typer.BadParameter`, `typer.Exit`, or `typer.Abort` over raw exceptions.

## Arguments and options

Prefer explicit names and modern annotations.

```python
name: Annotated[str, typer.Argument(help="Resource name")]
force: Annotated[bool, typer.Option("--force", help="Skip confirmation")]
output: Annotated[Path | None, typer.Option("--output", "-o")]
```

Recommended patterns:

- Use positional arguments only for obvious primary inputs.
- Use options for behavior switches and optional inputs.
- Expose standard names such as `--json`, `--plain`, `--quiet`, `--dry-run`, `--force`, `--no-input`.
- Avoid one-letter flags unless the flag is genuinely common.

Common flags to include by default when they fit the command:

- `--help`
- `--version` on the root app
- `--verbose` and `--quiet` for diagnostic control
- `--json` and `--plain` for output modes
- `--no-input` for prompt-capable commands
- `--dry-run` and `--force` for side-effecting commands

Prefer `--yes` or `--assume-yes` for confirmation bypass and reserve `--force` for bypassing safety checks or overwrite protections.

For non-trivial CLIs that may be used by agents or automation, consider an `explain` command distinct from `--help`.

If the CLI depends on durable user configuration, add a `config` command group with:

- `config show`
- `config get <key>`
- `config set <key> <value>`
- `config unset <key>`
- `config edit`
- `config path`

Make bare `config` behave like `config show` when feasible.

It is fine to add `view` as a hidden or undocumented alias for `show`.

## Help and examples

Typer generates help automatically, but you still need clear command names, docstrings, and option help.

```python
@app.command(help="Upload a report and print the resulting URL.")
def upload(...):
    ...
```

Write help that explains:

- what the command does
- the important flags
- what kind of output to expect

For non-trivial CLIs, provide both:

- root `--version`
- examples in command help or README

If the CLI surface changed, update README examples too.

`--help` should stay human-oriented. If the CLI needs to describe itself to agents, add an `explain` command that returns structured data instead of formatted prose.

Recommendations for `explain`:

- default to JSON, optionally support TOML or another explicit machine format
- make the output closer to an API schema than to terminal help text
- include canonical command paths, primary aliases, arguments, options, defaults, enum values, promptable inputs, config requirements, and output modes
- include enough metadata for an agent to build a fully qualified command without scraping help text
- keep the structure stable across versions when possible

Recommended top-level fields:

- `schema_version`
- `cli`
- `commands`
- `config`
- `environment`

Recommended design basis:

- OpenAPI-style top-level metadata and versioning
- MCP-style per-command records
- JSON Schema for `input_schema` and `output_schema`
- separate completion metadata when values or suggestions are richer than validation alone can express

For agent interoperability, prefer a stable canonical JSON shape over ad hoc prose or help-text scraping.

Use one naming convention consistently. Prefer `snake_case` such as `schema_version`, `input_schema`, `output_schema`, and `external_docs`.

For secret-bearing config and environment metadata, prefer explicit fields such as `secret`, `redacted`, `reveal_requires_opt_in`, and `sources`.

Suggested shape:

```python
@app.command("explain")
def explain(format: str = "json") -> None:
    payload = build_cli_schema()
    emit_schema(payload, format=format)
```

Do not make agents parse human help when a real schema-like command can be provided.

Do not include live secret values in `explain` output.

For larger Typer CLIs, treat examples as part of the interface, not just docs.

- maintain a centralized examples registry keyed by canonical command path
- inject examples into command help and group help automatically
- keep examples copy-pasteable and realistic
- include both human-oriented and machine-oriented examples when the command supports formats like `--json`
- use examples in error paths too, especially after usage errors

Good example strategy:

- one or two common success paths
- one example showing machine-readable output when supported
- one example showing non-interactive or fully qualified usage when prompts are possible

## Hints and error presentation

Default Click and Typer errors are often not enough for a polished CLI. For serious CLIs, add structured hints on top of normal parsing behavior.

Recommendations:

- on usage errors, show the relevant help text and then print a short hint
- if an option is missing, name the exact flag or argument to pass
- if a subcommand is missing, suggest visible subcommands or tell the user to run `--help`
- after usage errors, show examples for the current command when available
- keep expected user mistakes out of Python tracebacks

Hints should be short and actionable, such as:

- pass this flag
- choose one of these commands
- run this help command
- use this fully qualified non-interactive command next time

For advanced CLIs, it is reasonable to wrap Typer's Click layer with custom `TyperCommand` and `TyperGroup` subclasses so help, examples, and error presentation stay consistent across the whole tree.

When you rename commands or keep compatibility aliases, it is reasonable to surface a short deprecation hint that points to the canonical command.

## Aliases

Aliases are useful, but they should not make the CLI ambiguous.

Recommendations:

- keep one canonical command name and treat aliases as compatibility or convenience layers
- use hidden or undocumented aliases for separator variants, legacy spellings, and secondary compatibility names
- store aliases in a centralized registry instead of scattering them across modules
- support aliases for both top-level groups and subcommands when the CLI is large enough to justify them
- normalize alias paths before resolving examples, hints, or destructive-action summaries
- do not make aliases the only documented surface unless the alias is becoming the new canonical name
- it is reasonable to document one primary alias inline in help, such as `list (ls)`, when that alias is short, memorable, and genuinely first-class
- do not create duplicate help entries for aliases

Good uses of aliases:

- common abbreviations for frequently used top-level groups
- underscore and hyphen compatibility
- preserving older command spellings during migrations
- adding `view` as a hidden alias for `show`

Avoid:

- too many aliases for one command
- aliases that overlap semantically with another real command
- documenting every alias in help output
- long parenthesized alias lists such as `command (a, b, c, d)`

## Root app setup

For multi-command Typer CLIs, prefer a root callback that configures the CLI runtime once.

Use the root callback to:

- attach runtime state to `typer.Context`
- resolve output mode, color, verbosity, and debug settings
- handle eager root options such as `--version`
- show root help when no subcommand was invoked

This keeps global behavior consistent and prevents every subcommand from reimplementing setup.

Useful root-level patterns:

- `invoke_without_command=True` for a helpful root entrypoint
- eager `--version` callback that prints and exits
- root-level output controls such as `--verbose`, `--quiet`, `--plain`, `--no-color`, or format defaults
- installable shell completion where the project is large enough to benefit from it

## Custom command and group classes

For small CLIs, plain Typer defaults are enough. For larger CLIs, recommend a thin custom layer.

Use custom `TyperCommand` and `TyperGroup` classes when you need consistent behavior for:

- examples appended to help output
- grouped or curated command listings
- visible-only command lists that ignore hidden aliases
- mutation or danger summaries in help output
- normalized error rendering and hint generation
- consistent handling of `KeyboardInterrupt`, `EOFError`, and broken pipes

Keep these classes lightweight. They should improve presentation and consistency, not replace Click parsing semantics.

## Organizing larger Typer CLIs

For bigger CLIs, prefer this shape:

- one root `app`
- one sub-app per domain or resource
- centralized registries for examples and aliases
- shared command and group classes for help and error formatting
- one runtime setup path in the root callback

This gives agents a repeatable structure without forcing every command module to reinvent help text, aliases, or diagnostics.

## Interactivity

Only prompt when `stdin` is interactive and the user did not disable prompts.

```python
def require_name(name: str | None, no_input: bool) -> str:
    if name:
        return name
    if no_input or not sys.stdin.isatty():
        typer.echo("Pass --name in non-interactive mode.", err=True)
        raise typer.Exit(code=2)
    return typer.prompt("Name")
```

For confirmations:

```python
if not force:
    confirmed = typer.confirm("Delete the remote environment?", default=False)
    if not confirmed:
        raise typer.Abort()
```

For severe actions, prefer typed confirmation over a simple yes/no prompt.

When prompting, be generous with hints:

- explain which config or default caused the prompt
- mention the valid options when they are not obvious
- remind the user about `--no-input` for fail-fast behavior
- if the prompt corresponds to flags, print the equivalent full command at the end

```python
def select_env(env: str | None, no_input: bool) -> str:
    if env:
        return env
    if no_input or not sys.stdin.isatty():
        typer.echo("Pass --env in non-interactive mode, or omit --no-input to choose interactively.", err=True)
        raise typer.Exit(code=2)

    typer.echo("No environment was provided. Available options: dev, staging, prod.", err=True)
    value = typer.prompt("Environment", default="dev")
    typer.echo(f"Equivalent command: myapp deploy --env {value}", err=True)
    return value
```

If config affects behavior, say so explicitly before prompting.

```python
typer.echo("Using default region from ~/.config/myapp/config.toml: us-east-1", err=True)
```

## Machine-readable output

If humans need pretty output, add an explicit machine mode.

```python
@app.command()
def list_items(
    json_output: Annotated[bool, typer.Option("--json")] = False,
    plain: Annotated[bool, typer.Option("--plain")] = False,
) -> None:
    items = [{"name": "alpha"}, {"name": "beta"}]

    if json_output:
        typer.echo(json.dumps(items))
        return
    if plain:
        for item in items:
            typer.echo(item["name"])
        return

    typer.echo("NAME")
    for item in items:
        typer.echo(item["name"])
```

Do not mix logs into `stdout` when `--json` is active.

## Errors

Rewrite expected failures into clear messages.

```python
if not path.exists():
    raise typer.BadParameter(f"File not found: {path}")
```

Prefer messages that tell the user what to do next.

## Config and environment

Use clear precedence:

1. CLI flags
2. environment variables
3. project config
4. user config
5. system config

Expose environment-backed defaults when they make sense, but keep them visible in help text.

Avoid accepting secrets as raw flag values.

Redaction guidance:

- redact secrets in `config show` by default
- avoid printing secret values in hints, examples, or debug output
- if `config get` may touch secret keys, prefer masked output or require an explicit reveal path
- `config path` should reveal only the path, not the contents
- `explain` should describe secret-bearing keys and env vars without printing live values

For large human-readable output, consider `PAGER`, but avoid paging in `--json`, `--plain`, non-interactive, or piped modes.

When a Python CLI has real config, make it manageable from the CLI instead of requiring manual file edits.

```python
config_app = typer.Typer(help="Inspect and manage configuration.")
app.add_typer(config_app, name="config")


@config_app.command("show")
def config_show() -> None:
    typer.echo(json.dumps(load_config(), indent=2))


@config_app.command("get")
def config_get(key: str) -> None:
    value = resolve_config_value(key)
    typer.echo(value)


@config_app.command("set")
def config_set(key: str, value: str) -> None:
    write_config_value(key, value)
    typer.echo(f"Updated {key}", err=True)


@config_app.command("unset")
def config_unset(key: str) -> None:
    remove_config_value(key)
    typer.echo(f"Unset {key}", err=True)


@config_app.command("edit")
def config_edit() -> None:
    open_in_editor(active_config_path())


@config_app.command("path")
def config_path() -> None:
    typer.echo(str(active_config_path()))
```

Guidance:

- `config show` should help the user see the active config file or profile
- bare `config` should route to the same behavior as `config show`
- `view` may exist as a hidden alias for `show`
- `config get` should print only the resolved value unless the user asks for more detail
- `config set` should confirm what changed and where it was written
- `config unset` should explain the resulting fallback source or default when useful
- `config edit` should respect `EDITOR` and fail clearly if no editor is available
- `config path` should print just the path for scriptability
- if config supports scopes such as project vs user, make the scope explicit with flags

## Testing

Use Typer's testing helpers for parsing and output behavior.

```python
from typer.testing import CliRunner

runner = CliRunner()


def test_show_json() -> None:
    result = runner.invoke(app, ["show", "README.md", "--json"])
    assert result.exit_code == 0
    assert '"path"' in result.stdout
```

Cover at least:

- top-level help
- one success path
- one failure path
- `--json` or `--plain` behavior when supported
- prompt bypass behavior in non-interactive mode
- helpful hint text when prompting or inferring values
- redaction behavior for secret-bearing config
- explain output structure if an `explain` command exists

## Good defaults to reach for

- `Typer(no_args_is_help=True)` for multi-command apps
- `Annotated[...]` with `typer.Argument` and `typer.Option`
- `typer.echo(..., err=True)` for stderr
- `typer.confirm()` and `typer.prompt()` only behind TTY checks
- `typer.Exit(code=...)` for deliberate exits
