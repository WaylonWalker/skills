# Start And Setup

Use this file when the user is starting a new Terraria mod or when the local setup is broken.

## Baseline stack

tModLoader modding is C# modding on .NET 8.

Required baseline:

- tModLoader installed
- .NET 8 SDK installed
- A real IDE or editor setup

Do not recommend .NET 9 or .NET 10 for tModLoader mod projects.

## IDE guidance

- Windows: Visual Studio 2022 v17.8+ is the official best-supported path.
- Cross-platform: Rider is a strong choice on Windows, macOS, and Linux.
- Cross-platform lightweight option: VS Code with C# Dev Kit.

If the user is inexperienced and on Windows, prefer Visual Studio.

## Save and source locations

Common save roots for current 1.4.4-era tModLoader docs:

- Windows: `%UserProfile%\Documents\My Games\Terraria\tModLoader`
- Linux: `~/.local/share/Terraria/tModLoader/` or `$XDG_DATA_HOME/Terraria/tModLoader/`
- macOS: `~/Library/Application support/Terraria/tModLoader/`

Inside that save root, the mod sources directory is typically `ModSources/`.

Common install roots:

- Steam Windows: `C:\Program Files (x86)\Steam\steamapps\common\tModLoader`
- Steam Linux: `~/.local/share/Steam/steamapps/common/tModLoader` or `~/.steam/steam/steamapps/common/tModLoader`
- Steam macOS: `~/Library/Application Support/Steam/steamapps/common/tModLoader`

## Best way to start a new mod

Prefer the built-in skeleton generator:

1. Open tModLoader.
2. Go to `Workshop`.
3. Open `Develop Mods`.
4. Click `Create Mod`.
5. Fill in the internal name, display name, author, and starter item.

This creates the mod folder, base content, `build.txt`, `.csproj`, launch settings, and starter localization.

Do not hand-roll a `.csproj` unless the user has a special need and the generated project is unavailable.

## Typical generated structure

```text
MyMod/
  MyMod.cs
  MyMod.csproj
  build.txt
  description.txt
  description_workshop.txt
  icon.png
  icon_small.png
  Properties/
    launchSettings.json
  Content/
    Items/
      StarterItem.cs
      StarterItem.png
  Localization/
    en-US_Mods.MyMod.hjson
```

## What each core file is for

- `MyMod.cs`: the main `Mod` class
- `build.txt`: packaging, versioning, dependencies, and metadata
- `description.txt`: in-game description
- `description_workshop.txt`: Steam Workshop page description
- `MyMod.csproj`: IDE build/debug project
- `Properties/launchSettings.json`: debugging launch profiles and command-line args
- `Localization/*.hjson`: localized display names, tooltips, and other user-facing text

## First build workflow

For a brand-new mod, prefer:

1. Save code and assets.
2. In tModLoader, go to `Workshop -> Develop Mods`.
3. Click `Build + Reload`.
4. Fix build errors before adding more content.

## Setup triage checklist

If setup fails, check these first:

- .NET 8 SDK is installed, not only a runtime
- The project was opened via the `.csproj`, not just a loose `.cs` file
- The IDE version matches tModLoader requirements
- The mod was created by the skeleton generator or upgraded properly
- The user is editing the actual mod source folder under `ModSources`

## Good defaults for an agent

- Reuse the generated layout.
- Keep namespaces aligned with folders.
- Keep one content type per file unless the example is intentionally tiny.
- Preserve the mod's internal name and folder name.

## Sources

- `tModLoader guide for developers`
- `Basic tModLoader Modding Guide`
- `Developing with Visual Studio`
- `Developing with Visual Studio Code`
- `Developing with Rider`
- `Basic tModLoader Usage Guide`
