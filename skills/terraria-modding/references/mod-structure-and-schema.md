# Mod Structure And Schema

Use this file when the task involves folder layout, `build.txt`, localization, naming, or content architecture.

## Core class map

The most common tModLoader content and extension points are:

- `Mod`: the central mod entry point
- `ModSystem`: mod-wide systems, setup, integration, world-level logic
- `ModItem`: custom items
- `ModProjectile`: custom projectiles
- `ModNPC`: enemies, bosses, critters, town NPCs
- `ModTile`: placeable tiles
- `ModWall`: walls
- `ModPlayer`: player-attached custom state and hooks
- `ModConfig`: configuration exposed to users
- `GlobalItem`, `GlobalNPC`, `GlobalProjectile`: behavior attached to many existing entities

Pick the narrowest class that owns the behavior.

## Recommended content layout

```text
MyMod/
  MyMod.cs
  build.txt
  Content/
    Items/
    Projectiles/
    NPCs/
    Tiles/
    Walls/
    Buffs/
  Common/
    Systems/
    Players/
    Globals/
  Assets/
    Sounds/
    Textures/
  Localization/
```

Follow the repository's existing shape if it already has one.

## Naming rules

- Internal mod names should be unique and have no spaces.
- Class names should usually be PascalCase.
- File names should match the main class in the file.
- Asset paths must match code references exactly.
- Renaming a published mod's internal name is risky and should be treated as a migration task.

## `build.txt` schema

`build.txt` is a top-level key-value file. Capitalization matters.

Common properties:

- `displayName`
- `author`
- `version`
- `homepage`
- `dllReferences`
- `modReferences`
- `weakReferences`
- `noCompile`
- `hideCode`
- `hideResources`
- `includeSource`
- `buildIgnore`
- `side`
- `sortAfter`
- `sortBefore`
- `playableOnPreview`
- `translationMod`

Example:

```ini
author = Your Name
displayName = My Cool Mod
version = 0.1.0
homepage = https://github.com/yourname/my-cool-mod
includeSource = true
buildIgnore = *.csproj, *.user, obj\*, bin\*, .vs\*
```

Version guidance:

- tModLoader accepts 2 to 4 numeric parts.
- A 3-part semantic style such as `0.1.0` is a good default.
- Do not include letters such as `beta` or `v1.0`.

Dependency guidance:

- `modReferences`: required mods
- `weakReferences`: optional mods with careful code
- `dllReferences`: non-mod DLLs placed in a top-level `lib/` folder

Packaging guidance:

- `buildIgnore` matters for final `.tmod` output.
- Large PSDs, source exports, and tooling folders should be excluded.
- IDE builds may not reflect final `buildIgnore` packaging behavior as faithfully as in-game builds.

## Localization shape

Use HJSON files under `Localization/`, usually starting with `en-US_Mods.<ModName>.hjson`.

Typical shape:

```hjson
Items: {
  MySword: {
    DisplayName: My Sword
    Tooltip: A simple example sword
  }
}
```

The exact file often grows automatically as content is added and exported by tModLoader. Reuse the existing structure in the mod rather than inventing a new one.

## Asset conventions

- Item sprites usually live near the corresponding content file or within matching content folders.
- `icon.png` is the in-game icon and should be 80x80.
- `icon_workshop.png` can be larger for Steam Workshop.
- Asset names are part of runtime lookup. Case, folder names, and file placement matter.

## Main `Mod` class responsibilities

Use the `Mod` class for central registration and narrowly global behavior.

Common hooks and methods to know:

- `Load`
- `Unload`
- `AddRecipes`
- `PostAddRecipes`
- `PostSetupContent`
- `HandlePacket`
- `Call`

Important timing note:

- `Load` runs before all content setup is fully finished.
- `PostSetupContent` is the safer place for logic that needs all content loaded.

## Sources

- `build.txt`
- tModLoader `Mod` class docs
- `Basic tModLoader Modding Guide`
