# Cross-Mod And Hot Reload

Use this file when the task mentions compatibility with other mods, optional dependencies, `Mod.Call`, or hot reload.

## Cross-mod options

There are four common approaches, from easiest to most risky:

1. Use another mod's items or tiles by lookup.
2. Use `Mod.Call` if the target mod documents an API.
3. Use strong references with `modReferences`.
4. Use weak references with `weakReferences` and JIT-safe code.

Avoid reflection unless there is no supported alternative.

## Simple optional lookups

Use this when the dependency is optional and you only need items, tiles, recipes, or similar content IDs.

Patterns:

```csharp
if (ModLoader.TryGetMod("ExampleMod", out Mod exampleMod)) {
    if (exampleMod.TryFind("ExampleWings", out ModItem exampleWings)) {
        npcShop.Add(exampleWings.Type);
    }
}
```

Or:

```csharp
if (ModContent.TryFind("ExampleMod", "ExampleTorch", out ModItem exampleTorch)) {
    npcShop.Add(exampleTorch.Type);
}
```

## `Mod.Call`

Prefer `Mod.Call` when the target mod publishes a stable call contract.

Default pattern:

1. Check whether the mod is loaded with `ModLoader.TryGetMod`.
2. If needed, check the dependency mod version.
3. Make the call in the hook the target mod expects, often `PostSetupContent`.
4. Follow the documented message format exactly.

This is usually the best optional-integration path when supported.

## Strong references

Use `modReferences` when your mod cannot function without the other mod.

Example `build.txt` entry:

```ini
modReferences = ExampleMod
```

Or pin a minimum compatible version:

```ini
modReferences = ExampleMod@2.0
```

Strong references are easy to code against but make the dependency mandatory.

## Weak references

Use `weakReferences` only when the dependency is optional and you truly need direct code-level access.

Example:

```ini
weakReferences = ExampleMod@2.0
```

This requires careful coding because the .NET runtime can crash even when the code path looks guarded.

Important rules:

- Do not assume `if (ModLoader.HasMod(...))` is enough by itself.
- Keep weak-reference code in isolated classes or members.
- Use documented attributes such as `JITWhenModsEnabled` and `ExtendsFromMod` where appropriate.
- Test by disabling the referenced mod and fully restarting tModLoader.

## Hot reload reality

There is no single universal hot reload mode for every Terraria modding task.

### Visual Studio

- Best supported path on Windows.
- Uses Edit and Continue while debugging.
- Good for changing many runtime code paths without full restart.

### VS Code

The official docs describe enabling C# Dev Kit hot reload with settings such as:

```json
"csharp.experimental.debug.hotReload": true,
"csharp.debug.hotReloadOnSave": true
```

### Rider

- Supports hot reload while debugging.
- Open the real `.csproj` path directly, not a symlinked path, if hot reload behaves strangely.

## Hot reload limits

Be explicit about what usually does not hot-reload cleanly:

- `Load`
- many `SetStaticDefaults` effects
- content registration
- type discovery and autoload behavior
- some asset changes that require reimport or full reload

If the code only runs at mod load time, expect to rebuild or reload mods.

## Best guidance when a user asks for hot reload

- Offer the best IDE-specific path they can use.
- Explain that in-game `Build + Reload` remains the reliable fallback.
- Distinguish between runtime logic edits and load-time registration edits.

## Sources

- `Expert Cross Mod Content`
- `Why Use an IDE`
- `Developing with Visual Studio`
- `Developing with Visual Studio Code`
- `Developing with Rider`
