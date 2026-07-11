# Examples

Use these examples as templates. They are adapted from official tModLoader docs and the stable `ExampleMod`.

## Example 1: Simple ranged `ModItem`

```csharp
using Terraria;
using Terraria.ID;
using Terraria.ModLoader;

namespace MyMod.Content.Items.Weapons;

public class MyBlaster : ModItem
{
    public override void SetDefaults() {
        Item.width = 40;
        Item.height = 20;
        Item.useStyle = ItemUseStyleID.Shoot;
        Item.useTime = 12;
        Item.useAnimation = 12;
        Item.autoReuse = true;
        Item.DamageType = DamageClass.Ranged;
        Item.damage = 18;
        Item.knockBack = 3f;
        Item.noMelee = true;
        Item.shoot = ProjectileID.PurificationPowder;
        Item.shootSpeed = 10f;
        Item.useAmmo = AmmoID.Bullet;
        Item.rare = ItemRarityID.Green;
    }

    public override void AddRecipes() {
        CreateRecipe()
            .AddIngredient(ItemID.IllegalGunParts)
            .AddIngredient(ItemID.FallenStar, 5)
            .AddTile(TileID.Anvils)
            .Register();
    }
}
```

## Example 2: Item that swaps projectile type

```csharp
public override void ModifyShootStats(Player player, ref Vector2 position, ref Vector2 velocity, ref int type, ref int damage, ref float knockback) {
    if (type == ProjectileID.Bullet) {
        type = ModContent.ProjectileType<MySpecialBullet>();
    }
}
```

Use this when the item should reuse normal ammo but change the fired projectile.

## Example 3: Simple `ModProjectile`

```csharp
using Terraria;
using Terraria.ID;
using Terraria.ModLoader;

namespace MyMod.Content.Projectiles;

public class MyBullet : ModProjectile
{
    public override void SetDefaults() {
        Projectile.width = 10;
        Projectile.height = 10;
        Projectile.friendly = true;
        Projectile.DamageType = DamageClass.Ranged;
        Projectile.timeLeft = 120;
        Projectile.penetrate = 1;
        Projectile.ignoreWater = true;
        AIType = ProjectileID.Bullet;
    }
}
```

## Example 4: Optional cross-mod lookup

```csharp
using Terraria.ModLoader;

if (ModLoader.TryGetMod("ExampleMod", out Mod exampleMod)) {
    if (exampleMod.TryFind("ExampleWings", out ModItem exampleWings)) {
        shop.Add(exampleWings.Type);
    }
}
```

Use this for recipes, shops, drops, and other optional content references.

## Example 5: `Mod.Call` integration in `PostSetupContent`

```csharp
using Terraria.ModLoader;

public override void PostSetupContent() {
    if (!ModLoader.TryGetMod("BossChecklist", out Mod bossChecklistMod)) {
        return;
    }

    bossChecklistMod.Call(
        "LogBoss",
        Mod,
        "MyBoss",
        5.2f,
        (Func<bool>)(() => DownedBossSystem.downedMyBoss),
        ModContent.NPCType<MyBossNPC>()
    );
}
```

Use this only after confirming the target mod's documented call signature.

## Example 6: `build.txt`

```ini
author = Your Name
displayName = My Mod
version = 0.1.0
homepage = https://github.com/yourname/my-mod
includeSource = true
buildIgnore = *.csproj, *.user, obj\*, bin\*, .vs\*, ArtSource\*
```

## Example 7: Localization entry

```hjson
Items: {
  MyBlaster: {
    DisplayName: My Blaster
    Tooltip: Fires a custom round
  }
}
```

## Example 8: Command-line debugging shortcut

Use `-skipselect` in a launch profile to get into a test world faster.

Examples:

```text
-skipselect
-skipselect ":MyTestWorld"
```

## Notes

- `ExampleGun.cs`, `ExamplePiercingProjectile.cs`, and `ModIntegrationsSystem.cs` in stable `ExampleMod` are excellent real references.
- When adapting an example, remove unrelated teaching code and keep only the behavior the user needs.
