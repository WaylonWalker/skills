# Build, Debug, And Publish

Use this file when the user asks how to build, test, debug, package, or publish a Terraria mod.

## Day-to-day development loop

The reliable loop is:

1. Edit code and assets.
2. Save all files.
3. Build.
4. Reload or debug.
5. Test in game.

For general correctness, in-game `Build + Reload` is the default safe workflow.

## Build paths

### In-game build

Best all-around build path:

- Open `Workshop -> Develop Mods`
- Use `Build + Reload`

Why this matters:

- It exercises the mod's normal packaging path.
- It produces the `.tmod` output used by the game.
- It respects release-oriented behavior more closely than many IDE-only builds.

### Visual Studio

Use the generated `.csproj`.

- Build from the IDE to catch compile errors quickly.
- Debug by launching the `Terraria` profile.
- Keep tModLoader closed, or ensure the mod is unloaded, before a plain IDE build.

Important caveat:

- Visual Studio build output can differ from in-game packaging, especially around `buildIgnore`.
- Before release, always do an in-game build.

### VS Code

The official docs recommend:

- open the mod folder correctly
- use C# Dev Kit
- run `dotnet msbuild`

If restore issues appear, run `dotnet restore` first.

### Rider

- Open the `.csproj` or `.sln` directly.
- Avoid symlinked paths if hot reload or IDE features behave strangely.
- Build and debug from Rider normally.

## Debugging

The best debugging path is to run the mod under the IDE debugger.

Useful behaviors:

- breakpoints
- variable inspection
- stepping through hooks and AI
- launch profiles for faster test loops

The `-skipselect` launch argument is especially useful for jumping directly into a test world.

Examples:

- `-skipselect`
- `-skipselect ":MyTestWorld"`

Launch profiles live in `Properties/launchSettings.json`.

## Sharing and local deployment

When built successfully, the mod is packaged as a `.tmod` file in the mods folder.

That file can be:

- copied to another local install for testing
- shared directly with testers
- published through Workshop

## Workshop publishing

Before publishing:

1. Set a valid numeric version in `build.txt`.
2. Update `description.txt`.
3. Update `description_workshop.txt`.
4. Update `changelog.txt` if present.
5. Confirm icons and metadata are correct.
6. Build in-game.
7. Test thoroughly.

Then:

1. Open `Workshop -> Develop Mods`
2. Click `Publish`
3. Set visibility and tags
4. Publish

## Release guidance

- Keep `displayName` user-friendly.
- Keep the internal mod name stable once published unless a migration is intentional.
- Keep source control in sync with releases.
- Exclude large art source files and junk with `buildIgnore`.

## Build issue triage

Check these quickly:

- wrong .NET SDK version
- project opened incorrectly
- tModLoader still holding the `.tmod` file open
- missing restore in VS Code
- stale references after dependency changes
- missing or mismatched asset paths

## Sources

- `Developing with Visual Studio`
- `Developing with Visual Studio Code`
- `Developing with Rider`
- `Workshop`
- `Command-Line`
