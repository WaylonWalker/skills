This repo will be a skills repo for agent skills.  there will be skills md files in the skills directory.

You make a readme. make it clear that this is intended to be a personal repo for skills, skills you create, or copy, and vet yourself, then use the cli to apply them to your projects or globally.  this is not intended to be a public repo for sharing skills, but you can share your skills by sharing the md files in the skills directory.
You make a cli in go.
You make a justfile, and github actions for well tested, checked, and most up to date go.
follow best cli practice from clig.dev
use an Ayu color scheme

## Command

the command will be called skill

skill use - opens a fuzzy picker to select a skill and add it to the project via a symlink
skill use -g - opens a fuzzy picker to select a skill and add it globally via a symlink
skill use <skill-name> - adds the specified skill to the project via a symlink
skill use -g <skill-name> - adds the specified skill globally via a symlink
skill list - lists all available skills using a picker that gives a nice preview of the skill when selected
skill list -g - lists all available global skills using a picker that gives a nice preview of the skill when selected
skill add - templates a new skill to the skills directory
skill remove - opens a fuzzy picker to select a skill and remove it from the project, warns and requires confirmation or -f if the skill file is not a symlink
skill remove -g - opens a fuzzy picker to select a skill and remove it globally, warns and requires confirmation or -f if the skill file is not a symlink

common flags:
-g, --global - operate on the global skills directory instead of the project directory
-f, --force - force remove a skill even if it is not a symlink, use

## Configuration

Users can configure their skills directory by setting the `SKILLS_DIR`
environment variable. If not set, it defaults to `~/.config/skills`.

Users can have multiple comma separated directories in `SKILLS_DIR`, and the
CLI will search for skills in all of them. The first skill found in the list
will be used for adding new skills.  for example users might have
~/.config/skills, ~/git/skills, and ~/private/skills

Support for multiple tools by setting the `SKILLS_TOOL` environment variable.
This variable can be set to a comma-separated list of tools (e.g.,
`SKILLS_TOOL=tool1,tool2`). The CLI will search for skills in the specified
tools and prioritize them based on the order they are listed. If `SKILLS_TOOL`
is not set, it defaults to searching for skills in all available tools.  tools
will tell `skill` command where to put the skill file based on how the tool is
configured, we should support all of the popular known tools, opencode, pi,
copilot, claude do some research to understand where the skills go for each
tool.
