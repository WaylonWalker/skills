# OpenCode Themes

Use this reference when a CLI needs interactive color themes or a theme picker.

Source of truth:

- TUI themes: `packages/opencode/src/cli/cmd/tui/context/theme.tsx`
- UI and desktop themes: `packages/ui/src/theme/default-themes.ts`
- Repo: `https://github.com/anomalyco/opencode`

## Theme IDs

Themes shared by the OpenCode TUI and UI packs:

- `aura`
- `ayu`
- `carbonfox`
- `catppuccin`
- `catppuccin-frappe`
- `catppuccin-macchiato`
- `cobalt2`
- `cursor`
- `dracula`
- `everforest`
- `flexoki`
- `github`
- `gruvbox`
- `kanagawa`
- `lucent-orng`
- `material`
- `matrix`
- `mercury`
- `monokai`
- `nightowl`
- `nord`
- `one-dark`
- `opencode`
- `orng`
- `osaka-jade`
- `palenight`
- `rosepine`
- `solarized`
- `synthwave84`
- `tokyonight`
- `vercel`
- `vesper`
- `zenburn`

Additional themes present in OpenCode's UI and desktop pack:

- `amoled`
- `oc-2`
- `onedarkpro`
- `shadesofpurple`

## What To Copy

Do not spread raw color codes through command handlers. Copy the theme IDs and map them into semantic tokens.

For the actual color codes, use the raw theme JSON files linked below. Those files are the source of truth for every hex value.

## Required TUI Theme Keys

OpenCode's TUI themes define these semantic keys under `theme`:

- `primary`
- `secondary`
- `accent`
- `error`
- `warning`
- `success`
- `info`
- `text`
- `textMuted`
- `background`
- `backgroundPanel`
- `backgroundElement`
- `border`
- `borderActive`
- `borderSubtle`
- `diffAdded`
- `diffRemoved`
- `diffContext`
- `diffHunkHeader`
- `diffHighlightAdded`
- `diffHighlightRemoved`
- `diffAddedBg`
- `diffRemovedBg`
- `diffContextBg`
- `diffLineNumber`
- `diffAddedLineNumberBg`
- `diffRemovedLineNumberBg`
- `markdownText`
- `markdownHeading`
- `markdownLink`
- `markdownLinkText`
- `markdownCode`
- `markdownBlockQuote`
- `markdownEmph`
- `markdownStrong`
- `markdownHorizontalRule`
- `markdownListItem`
- `markdownListEnumeration`
- `markdownImage`
- `markdownImageText`
- `markdownCodeBlock`
- `syntaxComment`
- `syntaxKeyword`
- `syntaxFunction`
- `syntaxVariable`
- `syntaxString`
- `syntaxNumber`
- `syntaxType`
- `syntaxOperator`
- `syntaxPunctuation`

Optional keys used by the TUI loader:

- `selectedListItemText`
- `backgroundMenu`
- `thinkingOpacity`

## Raw TUI Theme Files

Each link below contains the full color code set for that theme.

- `aura`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/aura.json`
- `ayu`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/ayu.json`
- `carbonfox`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/carbonfox.json`
- `catppuccin`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/catppuccin.json`
- `catppuccin-frappe`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/catppuccin-frappe.json`
- `catppuccin-macchiato`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/catppuccin-macchiato.json`
- `cobalt2`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/cobalt2.json`
- `cursor`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/cursor.json`
- `dracula`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/dracula.json`
- `everforest`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/everforest.json`
- `flexoki`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/flexoki.json`
- `github`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/github.json`
- `gruvbox`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/gruvbox.json`
- `kanagawa`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/kanagawa.json`
- `lucent-orng`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/lucent-orng.json`
- `material`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/material.json`
- `matrix`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/matrix.json`
- `mercury`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/mercury.json`
- `monokai`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/monokai.json`
- `nightowl`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/nightowl.json`
- `nord`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/nord.json`
- `one-dark`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/one-dark.json`
- `opencode`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/opencode.json`
- `orng`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/orng.json`
- `osaka-jade`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/osaka-jade.json`
- `palenight`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/palenight.json`
- `rosepine`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/rosepine.json`
- `solarized`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/solarized.json`
- `synthwave84`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/synthwave84.json`
- `tokyonight`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/tokyonight.json`
- `vercel`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/vercel.json`
- `vesper`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/vesper.json`
- `zenburn`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/opencode/src/cli/cmd/tui/context/theme/zenburn.json`

## Raw UI And Desktop Theme Files

These define the broader UI theme pack. Shared theme IDs have matching names. The four extra themes live only here.

- theme directory: `https://github.com/anomalyco/opencode/tree/dev/packages/ui/src/theme/themes`
- `amoled`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/ui/src/theme/themes/amoled.json`
- `oc-2`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/ui/src/theme/themes/oc-2.json`
- `onedarkpro`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/ui/src/theme/themes/onedarkpro.json`
- `shadesofpurple`: `https://raw.githubusercontent.com/anomalyco/opencode/dev/packages/ui/src/theme/themes/shadesofpurple.json`

Useful semantic slots from OpenCode's theme model:

- brand and emphasis: `primary`, `secondary`, `accent`
- status: `success`, `warning`, `error`, `info`
- text: `text`, `textMuted`
- surfaces: `background`, `backgroundPanel`, `backgroundElement`, `backgroundMenu`
- borders: `borderSubtle`, `border`, `borderActive`
- diffs: `diffAdded`, `diffRemoved`, `diffContext`, related background tokens
- markdown and syntax tokens when the CLI renders rich text or code

## Implementation Notes

- Only apply rich color themes when output is attached to a TTY.
- Always provide `--no-color`, and disable animation when not interactive.
- Keep `--json` and `--plain` free of decoration.
- Prefer a small semantic palette layer in your app over direct ANSI literals.
- If you expose theme selection, use stable IDs that match the list above.

## Python Suggestion

- Use Rich `Theme` plus one mapping layer from semantic token name to Rich style.
- Build one theme registry keyed by the OpenCode IDs above.
- Use separate consoles for `stdout` results and `stderr` diagnostics.
- Use `Console.status()` or `Progress` with `SpinnerColumn` for work expected to last more than about a second.

## Go Suggestion

- Use a central theme struct with semantic fields, then derive `lipgloss.Style` values from it.
- Keep theme registration in one package so commands share the same IDs and meanings.
- For one-shot commands, render progress on `stderr` with a spinner helper.
- For TUIs, keep theme state and progress state inside Bubble Tea models.
