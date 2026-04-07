# skills

A personal collection of agent skills for AI coding assistants.

This is not a public registry or shared marketplace -- it is a personal repo
for skills you create, copy, and vet yourself. You use the `skills` CLI to
apply them to your projects or globally across your machine.

Skills follow the [agentskills.io](https://agentskills.io) specification:
each skill is a directory containing a `SKILL.md` file with YAML frontmatter.

## Install

```sh
go install github.com/WaylonWalker/skills@latest
```

Or build from source:

```sh
git clone https://github.com/WaylonWalker/skills.git
cd skills
just build
```

## Quick Start

```sh
# Add a new skill to your collection
skills add

# Browse and preview available skills
skills list

# Apply a skill to the current project
skills use

# Apply a skill globally
skills use -g

# Remove a skill from the current project
skills remove

# Show current configuration
skills config
```

## Configuration

### Skills Directory

Set `SKILLS_DIR` to configure where your skills are stored:

```sh
export SKILLS_DIR="~/.config/skills"
```

Multiple directories are supported (comma-separated). The CLI searches all of
them, and new skills are created in the first directory:

```sh
export SKILLS_DIR="~/.config/skills,~/git/skills/skills,~/private/skills"
```

Default: `~/.config/skills`

### Tool Support

Set `SKILLS_TOOL` to specify which AI tools to target:

```sh
export SKILLS_TOOL="claude,copilot,opencode"
```

If not set, the CLI targets all detected tools.

Supported tools and where skills are installed:

| Tool | Project Path | Global Path |
|------|-------------|-------------|
| claude | `.claude/rules/<name>.md` | `~/.claude/rules/<name>.md` |
| copilot | `.github/instructions/<name>.instructions.md` | -- |
| cursor | `.cursor/rules/<name>.md` | -- |
| opencode | -- | `~/.config/opencode/skills/<name>/SKILL.md` |
| windsurf | `.windsurf/rules/<name>.md` | -- |
| cline | `.clinerules/<name>.md` | `~/Documents/Cline/Rules/<name>.md` |
| augment | `.augment/rules/<name>.md` | `~/.augment/rules/<name>.md` |
| roo | `.roo/rules/<name>.md` | `~/.roo/rules/<name>.md` |

## Skills Directory Structure

Skills follow the agentskills.io specification. Each skill is a directory
containing a `SKILL.md` file with YAML frontmatter:

```
~/.config/skills/
  python-best-practices/
    SKILL.md
  go-conventions/
    SKILL.md
  docker-security/
    SKILL.md
```

The `SKILL.md` file contains YAML frontmatter with `name` and `description`
fields, followed by the skill content:

```markdown
---
name: go-conventions
description: Go coding conventions and best practices
---

# Go Conventions

Your instructions here.
```

Legacy flat `.md` files (e.g. `my-skill.md`) are also supported for backward
compatibility.

## Commands

### `skills use [name]`

Apply a skill to the current project by creating symlinks in the
tool-specific directories. Without a name, opens a fuzzy picker.

```sh
skills use                     # pick from available skills
skills use go-conventions      # apply a specific skill
skills use -g                  # pick and apply globally
skills use -g go-conventions   # apply a specific skill globally
```

### `skills list`

Browse all available skills with an interactive preview.

```sh
skills list                    # browse skills, shows project install status
skills list -g                 # browse skills, shows global install status
```

### `skills add`

Create a new skill from a template (creates `<name>/SKILL.md` with frontmatter).

```sh
skills add                     # prompted for name
skills add my-new-skill        # create with given name
```

### `skills remove [name]`

Remove a skill from the current project. If the installed file is not a
symlink, confirmation is required (or use `-f`).

```sh
skills remove                  # pick from installed skills
skills remove -g               # pick from globally installed skills
skills remove -f               # force remove even if not a symlink
```

### `skills config`

Show current configuration (skills directories, tool filter, etc.).

```sh
skills config                  # show configuration
skills config show             # same as above
```

## Flags

| Flag | Description |
|------|-------------|
| `-g, --global` | Operate on global tool directories instead of project |
| `-f, --force` | Force remove even if the file is not a symlink |
| `--version` | Print version |
| `-h, --help` | Print help |

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `SKILLS_DIR` | Comma-separated list of skills directories | `~/.config/skills` |
| `SKILLS_TOOL` | Comma-separated list of tools to target | all tools |
| `NO_COLOR` | Disable colored output | unset |

## License

MIT
