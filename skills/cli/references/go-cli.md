# Go CLI

Use this when building or improving a Go CLI, especially with Cobra, the standard library `flag` package, or Charm tools for terminal UX.

Charm tools are best for terminal presentation and interaction:

- `bubbletea` for stateful TUIs
- `huh` for forms and prompts
- `lipgloss` for styling and layout
- `log` for human-readable stderr logging

Charm does not replace argument parsing. Keep the repo's parser if it already exists. For small greenfield tools, use the standard library `flag` package and explicit subcommand dispatch.

## Cobra

If the repository uses Cobra, stay with Cobra. It is a strong default for multi-command Go CLIs.

Use these patterns:

- prefer `RunE` over `Run` so command handlers return errors cleanly
- use `Use`, `Short`, `Long`, and `Example` to make help genuinely useful
- validate arguments with `Args` helpers or a custom validator
- write results to `cmd.OutOrStdout()`
- write prompts, warnings, hints, and errors to `cmd.ErrOrStderr()`
- keep local flags local unless they truly apply to the whole command tree
- reserve persistent flags for cross-cutting concerns like `--config`, `--profile`, `--json`, `--plain`, and `--no-input`
- avoid `MarkFlagRequired` unless a prompt-free non-interactive path really needs it

Common flags and commands to include when they fit:

- root `--help` and `--version`
- `help` command support from Cobra defaults
- `--verbose` and `--quiet` for diagnostics
- `--json` and `--plain` for output modes
- `--no-input` for any prompt-capable command
- `--dry-run` and `--force` for state-changing commands
- a `config` command with `show`, `get`, `set`, `unset`, `edit`, and `path` subcommands when the CLI has durable configuration

Prefer `--yes` or `--assume-yes` for skipping confirmations and reserve `--force` for bypassing safety checks or overwrite protections.

For non-trivial CLIs that may be used by agents or automation, consider an `explain` command distinct from human help.

Make bare `config` behave like `config show` when feasible.

It is fine to add `view` as a hidden or undocumented alias for `show`.

For error behavior:

- set `SilenceUsage = true` for runtime failures so usage is not dumped after ordinary execution errors
- consider `SilenceErrors = true` at the root and print normalized errors once in the execute path
- return actionable errors from `RunE` instead of printing ad hoc error text everywhere

Minimal Cobra shape:

```go
cmd := &cobra.Command{
    Use:          "deploy",
    Short:        "Deploy a service",
    Example:      "myapp deploy --env staging --region us-east-1",
    Args:         cobra.NoArgs,
    SilenceUsage: true,
    RunE: func(cmd *cobra.Command, args []string) error {
        out := cmd.OutOrStdout()
        errOut := cmd.ErrOrStderr()

        if noInput && env == "" {
            fmt.Fprintln(errOut, "Missing required value: pass --env or omit --no-input to choose interactively.")
            return fmt.Errorf("missing --env")
        }

        if env == "" {
            if !isInteractive() {
                return fmt.Errorf("missing --env in non-interactive mode")
            }

            fmt.Fprintln(errOut, "No environment provided. Options: dev, staging, prod.")
            fmt.Fprintln(errOut, "Tip: pass --env <name> or use --no-input to fail fast.")

            selected, err := promptForEnv()
            if err != nil {
                return err
            }
            env = selected
            fmt.Fprintf(errOut, "Equivalent command: myapp deploy --env %s\n", env)
        }

        fmt.Fprintln(out, "deployed")
        return nil
    },
}
```

Use Cobra completion support when it improves discoverability.

- add shell completion generation for serious CLIs
- use `ValidArgsFunction` for dynamic completions where values are discoverable ahead of time
- keep completion values stable and descriptive
- document how to install or generate completions when the CLI is intended for regular interactive use

If the CLI also uses Viper or another config layer, make precedence explicit and surface the chosen config file or profile when it affects behavior.

When the CLI has durable config, make it inspectable and editable from the command line.

- `config show` should display the active config and where it came from when useful
- bare `config` should route to the same behavior as `config show`
- `view` may exist as a hidden alias for `show`
- `config get <key>` should print a single resolved value for scripts
- `config set <key> <value>` should update config predictably and report the scope or file changed
- `config unset <key>` should remove an explicit value and explain fallback behavior when useful
- `config edit` should open the config in `EDITOR` or fail clearly
- `config path` should print the active config file path only

Redaction guidance:

- redact secrets in `config show` by default
- do not print live secrets in logs, hints, or explain output
- if `config get` can access secret-bearing keys, prefer masked output or require an explicit reveal path
- treat environment variables for secrets as metadata in `explain`, not as values

For Cobra, model this as a `config` parent command with `show`, `get`, `set`, `unset`, `edit`, and `path` child commands rather than hidden flags scattered across unrelated commands.

## Larger Cobra CLI architecture

For larger Go CLIs, do not let each command invent its own help, aliases, and runtime wiring. Use shared infrastructure.

Recommended structure:

- one root command that owns global flags, config loading, output mode, and context setup
- one package or file per domain command group
- shared helpers for output, error rendering, prompt fallback, and config resolution
- centralized registries or metadata for examples, aliases, and command annotations when the tree is large

### Examples and help

Treat examples as part of the command surface.

- populate `Example` on non-trivial commands
- keep examples copy-pasteable and realistic
- include machine-readable examples when `--json` is supported
- include fully qualified non-interactive examples when prompts are possible
- for very large CLIs, it is reasonable to keep examples in a centralized registry and inject them during command construction

If Cobra's default help is too sparse, customize help templates so examples, grouped commands, environment notes, or mutation summaries are easy to find.

Keep `help` human-oriented. If the CLI needs to describe itself to agents, add an `explain` command that emits structured data.

Recommendations for `explain`:

- default to JSON, optionally support TOML or another explicit machine format
- make the output closer to an API schema than to rendered help
- describe canonical command paths, aliases, flags, args, defaults, required inputs, config surfaces, output modes, and interaction behavior
- include enough information for an agent to construct valid non-interactive invocations
- keep the shape stable and versionable

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
- separate completion metadata when discovery hints are richer than validation alone

For agent interoperability, prefer a stable canonical JSON shape over scraping Cobra help output.

Use one naming convention consistently. Prefer `snake_case` such as `schema_version`, `input_schema`, `output_schema`, and `external_docs`.

For secret-bearing config and environment metadata, prefer explicit fields such as `secret`, `redacted`, `reveal_requires_opt_in`, and `sources`.

For Cobra, model this as a normal command such as `myapp explain --format json`, backed by a shared schema builder rather than by scraping command help output.

### Aliases

Use aliases as a compatibility layer, not a second command taxonomy.

- keep one canonical command name
- use `Aliases` for obvious short forms and common compatibility spellings
- use hidden compatibility commands when you need alias behavior that should not clutter help
- keep aliases centralized for larger CLIs instead of scattering them across packages
- normalize command paths before looking up examples, hints, or mutation metadata
- it is reasonable to document one primary alias inline in help, such as `list (ls)`, when that alias is short, memorable, and genuinely first-class
- do not create duplicate help entries for aliases
- when a deprecated alias is used, it is reasonable to print a short hint pointing to the canonical command

Good uses of aliases:

- common abbreviations for heavily used command groups
- underscore and hyphen compatibility
- migration compatibility for renamed commands
- `view` as a hidden alias for `show`

Avoid:

- exposing every alias in command listings
- long alias lists in help output
- using aliases to create a second public command taxonomy

### Root runtime setup

For multi-command CLIs, configure runtime behavior once near the root.

Good places for this:

- `PersistentPreRunE` on the root command
- a shared setup function called from command constructors
- a central `Execute()` path that normalizes errors and exit behavior

Use that setup to:

- load config and decide precedence
- resolve output mode such as human, plain, or JSON
- set color and verbosity behavior
- attach shared clients or runtime values to `context.Context`
- determine whether interactive prompting is allowed

### Error and hint rendering

Cobra defaults are fine for small tools, but larger CLIs should present errors consistently.

Recommendations:

- keep one normalized error printing path near `Execute()`
- return actionable errors from `RunE` rather than printing ad hoc messages everywhere
- add short hints for common problems such as missing required flags, unknown subcommands, or disabled interactivity
- after usage errors, point users to the exact help command to run
- when prompting could have been avoided with flags, print the equivalent full command after the interactive flow completes

Hints should stay short and practical:

- which flag to pass
- which subcommand to try
- whether `--no-input` would fail fast
- what exact full command would reproduce the selected options

### Command metadata

For large CLIs, it is useful to attach lightweight metadata to commands.

Useful metadata includes:

- examples
- mutation or danger summaries
- config requirements
- visibility or support status

In Cobra this can live in:

- command annotations
- shared registries keyed by canonical command path
- helper constructors that stamp commands with consistent fields

Use this metadata to improve help output and hinting, not to create a second parallel routing system.

### Interrupts and broken pipes

Handle common terminal edge cases consistently.

- return cleanly on interrupted input when possible
- treat Ctrl-C as a normal CLI interruption, not a stack trace
- handle broken pipes without dumping noisy errors into pipelines

For large human-readable output, consider respecting `PAGER`, but avoid paging in `--json`, `--plain`, non-interactive, or piped modes.

### Recommended approach

For larger Cobra CLIs, prefer a small internal framework made of shared constructors and helpers rather than hand-built per-command behavior. The goal is a command tree that feels consistent in:

- help output
- examples
- aliases
- error text and hints
- output streams
- config and runtime setup

## Choose the right tool

Pick the lightest tool that fits:

- One-shot command with simple flags: `flag` plus `fmt` and maybe `log`
- Interactive prompt or confirmation: `huh`
- Styled terminal output: `lipgloss`
- Full-screen or stateful workflow: `bubbletea`

Do not reach for Bubble Tea when a normal command plus a prompt is enough.

## Minimal subcommand shape

For small CLIs, explicit `flag.FlagSet` subcommands are often enough.

```go
package main

import (
    "encoding/json"
    "flag"
    "fmt"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Fprintln(os.Stderr, "Usage: myapp <command> [options]")
        os.Exit(2)
    }

    switch os.Args[1] {
    case "show":
        show(os.Args[2:])
    case "help", "-h", "--help":
        fmt.Fprintln(os.Stdout, "Usage: myapp <command> [options]")
    default:
        fmt.Fprintf(os.Stderr, "Unknown command %q\n", os.Args[1])
        os.Exit(2)
    }
}

func show(args []string) {
    fs := flag.NewFlagSet("show", flag.ContinueOnError)
    fs.SetOutput(os.Stderr)

    jsonOutput := fs.Bool("json", false, "print JSON")
    if err := fs.Parse(args); err != nil {
        os.Exit(2)
    }

    result := map[string]string{"name": "alpha"}
    if *jsonOutput {
        _ = json.NewEncoder(os.Stdout).Encode(result)
        return
    }

    fmt.Fprintln(os.Stdout, result["name"])
}
```

If the repo already uses Cobra, Kong, or another parser, stay consistent with that choice and use Charm only for terminal UX.

## stdout, stderr, and logging

Keep channels clean:

- `stdout` for results and machine-readable output
- `stderr` for logs, prompts, warnings, progress, and human-facing errors

Charm Log is a good default for stderr logs.

```go
logger := log.NewWithOptions(os.Stderr, log.Options{
    ReportTimestamp: false,
})

logger.Info("index updated", "count", 12)
```

For machine-readable modes, keep logs off `stdout`. If the command emits JSON on `stdout`, diagnostics still belong on `stderr`.

Charm Log supports text, JSON, and logfmt formatters. Its styling is disabled automatically when output is not a TTY.

## Styled output with Lip Gloss

Use Lip Gloss for readable human output, not for machine mode.

- Keep theme definitions in one package and style from semantic tokens, not raw literals scattered through commands.
- If you want a larger built-in theme catalog, start with the OpenCode theme IDs in `references/opencode-themes.md`.
- Map each theme into semantic slots such as `primary`, `accent`, `success`, `warning`, `error`, `muted`, and surface colors.

```go
title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))
fmt.Fprintln(os.Stdout, title.Render("Deployment complete"))
```

Useful properties:

- color downsampling happens automatically
- layouts are easier with padding, borders, width, and joining helpers
- tables, lists, and trees are available for readable terminal output

If output may be piped, add `--plain` or `--json` and bypass decorative rendering.

## Prompts and forms with Huh

Use `huh` when a prompt improves the experience, but keep it optional.

```go
var confirm bool

form := huh.NewForm(
    huh.NewGroup(
        huh.NewConfirm().
            Title("Delete the remote environment?").
            Value(&confirm),
    ),
)

if err := form.Run(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
}
```

Rules:

- only run prompts when `stdin` is interactive
- add `--no-input` to disable prompts
- provide flags for required input in non-interactive mode
- consider `WithAccessible(true)` behind config or env for screen-reader support

For severe actions, use a typed confirmation with `huh.NewInput()` instead of a simple yes/no prompt.

Be generous with hints around prompts and inferred choices:

- explain when config, env, or auto-detection influenced the current choice
- mention the available options if you are asking the user to choose
- remind the user that `--no-input` can fail fast instead of prompting
- if the answer could have been passed with flags, print the equivalent full command after the prompt completes

```go
fmt.Fprintln(os.Stderr, "No environment provided. Options: dev, staging, prod.")
fmt.Fprintln(os.Stderr, "Tip: pass --env <name> or use --no-input to fail fast.")
```

After interactive selection, show the reproducible command:

```go
fmt.Fprintln(os.Stderr, "Equivalent command: myapp deploy --env staging --region us-east-1")
```

If fuzzy selection or auto-picking is used, explain the match and how to override it.

## Bubble Tea for stateful terminal apps

Use Bubble Tea only when the interface is truly stateful: multi-step workflows, rich navigation, live status, or mixed input and rendering.

Core structure:

```go
type model struct {
    cursor int
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyPressMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return m, tea.Quit
        }
    }
    return m, nil
}

func (m model) View() string {
    return "Press q to quit\n"
}
```

Use Bubble Tea with Lip Gloss for layout and styling. If you need form-like inputs inside a Bubble Tea app, embed `huh.Form`; it is already a `tea.Model`.

Important:

- Bubble Tea owns terminal rendering, so do not print regular logs to `stdout` while it is running.
- For debugging Bubble Tea, log to a file instead of stdout.

## Progress and long-running work

If work takes noticeable time:

- say something quickly on `stderr`
- show progress only when writing to a TTY
- avoid animations in non-interactive contexts
- if the work will likely take longer than about 1 second, start a spinner
- prefer a braille snake spinner when the library allows custom frames
- keep the spinner text specific about the current step
- for longer waits, surface short tips or context on `stderr`

For prompt-driven workflows, `huh/spinner` is a good fit after submission. For TUIs, keep progress inside Bubble Tea.

For one-shot commands, a small spinner wrapper around the blocking operation is usually enough. Keep it off `stdout`, stop it before prompts, and fall back to a single status line when stderr is not a TTY.

## Errors and exit behavior

Prefer clear, actionable errors.

```go
if err != nil {
    fmt.Fprintf(os.Stderr, "cannot open %s: %v\n", path, err)
    os.Exit(1)
}
```

Do not dump stack traces for ordinary user mistakes.

Use:

- exit code `2` for usage and argument problems
- exit code `1` for execution failures
- exit code `0` for success

If the tool changes state, tell the user what changed.

## Config and environment

Use predictable precedence:

1. flags
2. environment variables
3. project config
4. user config
5. system config

Respect common environment variables where relevant:

- `NO_COLOR`
- `TERM`
- `DEBUG`
- `PAGER`
- `EDITOR`

Do not accept secrets as raw flag values.

## Verification

For Go CLI changes, verify at least:

- top-level help or usage output
- one success path
- one failure path and exit code
- `--json` or `--plain` behavior when supported
- prompt-free behavior in non-interactive mode
- redaction behavior for secret-bearing config
- explain output structure if an `explain` command exists

If you add Bubble Tea or Huh flows, verify they degrade cleanly when the command cannot interact with a TTY.
