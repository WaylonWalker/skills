---
name: terraria-modding
description: Build, debug, port, review, and publish Terraria mods with tModLoader. Use this whenever the user mentions Terraria modding, tModLoader, ModItem, ModProjectile, ModNPC, ModSystem, build.txt, Workshop publishing, ExampleMod, Mod.Call, cross-mod compatibility, Build + Reload, Edit and Continue, or hot reload for a Terraria mod, even if they do not explicitly say "tModLoader."
source: https://github.com/tModLoader/tModLoader/wiki
---

# Terraria Modding

Use this skill for Terraria mods built with tModLoader.

Start by inspecting the current mod before changing anything:

- Find the mod root, `build.txt`, the `Mod` class, the `.csproj`, and `Properties/launchSettings.json` if present.
- Detect whether the task is new content, a bug fix, a port, a build problem, a Workshop release task, or cross-mod integration.
- Preserve the existing mod's internal name, namespaces, and public identifiers unless the user explicitly wants a rename or breaking change.
- Prefer the smallest correct change. Terraria mods often rely on naming, asset paths, and hook timing.

## Core rules

- Treat official tModLoader docs and the stable `ExampleMod` as the primary references.
- Prefer generated mod skeletons over hand-rolled project files.
- Keep content names stable and descriptive. Internal names usually use PascalCase and should not contain spaces.
- Follow existing folder and namespace patterns in the mod.
- Use generic `ModContent.*Type<T>()` and `ModContent.TryFind(...)` patterns instead of stringly typed lookups when possible.
- For optional mod interoperability, prefer `ModLoader.TryGetMod`, `Mod.Call`, and documented weak-reference patterns over reflection.
- For release builds, prefer in-game `Build + Reload` or the Develop Mods menu build flow because IDE-only builds can miss final packaging behavior such as `buildIgnore` handling.

## Workflow

Follow this order:

1. Identify the task type.
2. Read the relevant reference file below.
3. Inspect the existing code and asset layout.
4. Implement the change using the narrowest appropriate tModLoader hook or content class.
5. Verify with the best available build and runtime workflow.
6. If the task affects publishability, versioning, or interoperability, update `build.txt`, descriptions, or docs as needed.

## Choose the right reference

- New mod setup, required tools, save paths, skeleton generation, and base folder layout:
  Read `references/start-and-setup.md`
- Common mod structure, content classes, localization, and `build.txt` fields:
  Read `references/mod-structure-and-schema.md`
- Implementing items, projectiles, recipes, systems, and common content patterns:
  Read `references/content-patterns.md`
- Build, debug, Build + Reload, IDE usage, Workshop release, and deploy flow:
  Read `references/build-debug-publish.md`
- Cross-mod support, strong vs weak references, `Mod.Call`, and hot reload expectations:
  Read `references/cross-mod-and-hot-reload.md`
- Example snippets adapted from official docs and `ExampleMod`:
  Read `references/examples.md`

## Task guidance

### New content

For new gameplay content, prefer these shapes:

- `ModItem` for items, weapons, accessories, consumables, summon items, placeables.
- `ModProjectile` for bullets, beams, pets, minions, and other spawned moving entities.
- `ModNPC` for enemies, bosses, critters, and town NPCs.
- `ModTile` and `ModWall` for world content.
- `ModSystem` for world-level setup, integration, loading-time registration, and broader hooks.
- `ModPlayer` when the behavior belongs to player state or lifecycle.
- `GlobalItem`, `GlobalNPC`, or `GlobalProjectile` when modifying many existing entities rather than adding one new thing.

### Bug fixes

When debugging:

- Check whether the issue is a wrong hook, wrong side, wrong timing, or wrong asset path before rewriting logic.
- Inspect whether the bug appears only after reload, only in multiplayer, or only after loading an old save.
- Be careful with code that only runs during loading, such as `Load`, `SetStaticDefaults`, and some registration logic.

### Porting or updating

If the user is updating an old mod:

- Preserve the internal mod name unless they explicitly want to rename it.
- Look for outdated APIs, renamed hooks, or old dependency patterns.
- Read migration docs if the codebase clearly spans older tModLoader versions.
- Verify `build.txt`, `.csproj`, and launch settings before touching gameplay logic.

### Cross-mod support

For compatibility work:

- Prefer `ModLoader.TryGetMod` for optional presence checks.
- Prefer `Mod.Call` when the target mod documents a call-based API.
- Use `modReferences` only when the dependency is truly required.
- Use `weakReferences` only when the code is written defensively enough to avoid JIT/load crashes.

## Verification

Minimum verification depends on the task:

- Code compiles with the mod's normal build path.
- Asset names and code references agree exactly.
- The changed content loads, builds, or registers without obvious load errors.
- If the task touches multiplayer or cross-mod support, verify the relevant path specifically.

Preferred verification paths:

- In-game: `Workshop -> Develop Mods -> Build + Reload`
- Visual Studio: build or debug from the generated `.csproj`
- VS Code: `dotnet msbuild` and, if configured, debug with hot reload support
- Rider: build or debug from the `.csproj` and use its hot reload workflow when available

## Hot reload expectations

tModLoader does not offer universal "edit any file and it live-reloads everything" behavior.

- Visual Studio uses Edit and Continue while debugging.
- VS Code can use C# Dev Kit hot reload when enabled.
- Rider supports hot reload while debugging.
- In-game `Build + Reload` is still the reliable general workflow.
- Load-time hooks and static registration code often do not benefit from hot reload and may still require a rebuild or reload.

When the user asks for hot reload, set expectations clearly and prefer the best supported workflow for their IDE.

## Report what matters

When you finish, summarize:

- What changed.
- Which tModLoader classes or hooks were used.
- How to build and test the result.
- Any release or compatibility caveats, especially around `build.txt`, weak references, or hot reload limitations.
