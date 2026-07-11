# Content Patterns

Use this file when implementing gameplay content or choosing hooks.

## Items

Use `ModItem` for weapons, accessories, consumables, placeables, ammo, and materials.

Common responsibilities:

- `SetDefaults` for behavior and stats
- `AddRecipes` for crafting
- `HoldoutOffset` for visual alignment of held items
- `ModifyShootStats` when an item should rewrite projectile type, damage, spawn point, or velocity
- `Shoot` when the item should spawn projectiles directly or cancel vanilla spawn behavior

Typical fields configured in `SetDefaults`:

- size: `Item.width`, `Item.height`
- rarity: `Item.rare`
- use timing: `Item.useTime`, `Item.useAnimation`, `Item.useStyle`
- auto reuse: `Item.autoReuse`
- damage settings: `Item.DamageType`, `Item.damage`, `Item.knockBack`
- projectile settings: `Item.shoot`, `Item.shootSpeed`, `Item.useAmmo`

Prefer helper methods or `CloneDefaults` only when they match the intended vanilla baseline closely.

## Projectiles

Use `ModProjectile` for anything moving or spawned: bullets, arrows, magic bolts, pets, minions, beams, lingering effects.

Common responsibilities:

- `SetDefaults` for size, time, collision, damage type, penetration, and base behavior
- `AIType = ProjectileID.SomeProjectile` when vanilla AI is the right baseline
- `OnHitNPC` to customize hit behavior and immunity

Be deliberate with projectile immunity:

- `Projectile.penetrate`
- `Projectile.usesLocalNPCImmunity`
- `Projectile.localNPCHitCooldown`
- static-ID immunity patterns when one projectile type should share cooldowns

## NPCs

Use `ModNPC` for enemies, bosses, critters, or town NPCs.

Common tasks:

- set stats and AI defaults
- assign drops
- set bestiary info
- add spawn conditions or summon logic
- synchronize boss progression with systems or world state

For broad edits to many NPCs, prefer `GlobalNPC` over rewriting individual NPC classes.

## Recipes

Prefer fluent recipe construction with `CreateRecipe()` and `.Register()`.

Pattern:

```csharp
CreateRecipe()
    .AddIngredient<ExampleItem>()
    .AddTile<Tiles.Furniture.ExampleWorkbench>()
    .Register();
```

For vanilla content, use the appropriate ID classes such as `ItemID`, `ProjectileID`, `TileID`, and `NPCID`.

## Referencing content

Use these patterns:

- `ModContent.ItemType<MyItem>()`
- `ModContent.ProjectileType<MyProjectile>()`
- `ModContent.NPCType<MyNpc>()`
- `ModContent.TileType<MyTile>()`

This is safer and easier to maintain than loose string lookups.

## Systems and world-level behavior

Use `ModSystem` when logic is not owned by one item, projectile, or NPC.

Examples:

- recipe groups
- global registration
- post-load integrations
- world state helpers
- packet routing and broader network setup

## Player-attached behavior

Use `ModPlayer` when the state lives on the player.

Examples:

- accessory effects
- custom resource flags
- save/load player state
- input handling tied to a player

## Asset and code alignment

When a content file expects an asset, verify:

- the file exists
- the path matches code usage
- the asset name matches the content name if relying on conventions

Many tModLoader content-loading issues are simple path or naming mismatches.

## Good design defaults

- Use `SetDefaults` for per-instance runtime defaults.
- Use `SetStaticDefaults` for display names, research, static sets, and metadata.
- Keep unrelated features out of one class unless the content is intentionally tiny.
- For simple learning examples, keep code close to the content class.
- For real mods, move shared logic into `Common/Systems`, `Common/Players`, or helper classes when repetition appears.

## Sources

- `Basic tModLoader Modding Guide`
- stable `ExampleMod`
- tModLoader API docs
