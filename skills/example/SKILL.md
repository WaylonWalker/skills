---
name: example
description: An example skill to demonstrate the agentskills.io format
---

# example

An example skill to demonstrate the format.

## Instructions

This is an example skill file. Skills are directories containing a
`SKILL.md` file with YAML frontmatter, following the agentskills.io
specification.

When applied to a project or globally, the skill file is symlinked
into the directories each tool expects:

- Claude: `.claude/rules/`
- Copilot: `.github/instructions/`
- Cursor: `.cursor/rules/`
- OpenCode: `~/.config/opencode/skills/`
- Windsurf: `.windsurf/rules/`
- And more.

Replace this content with your own instructions.
