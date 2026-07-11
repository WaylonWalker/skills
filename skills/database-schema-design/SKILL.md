---
name: database-schema-design
description: Design or review relational database schemas and related Python persistence models for self-hosted CRUD applications on PostgreSQL or CloudNativePG. Use this whenever the user mentions tables, models, SQLAlchemy, SQLModel, Alembic, Pydantic, entities, relationships, constraints, indexes, soft deletes, audit logs, multitenancy, or schema migrations. Use it before creating new CRUD features, endpoints, admin screens, or API models that will need storage, even if the user does not ask for "schema design" explicitly, because early database choices are expensive to reverse.
---

# Database Schema Design

Use this skill to make durable schema choices for Python CRUD apps.

## Start Here

1. Inspect the existing models, migrations, constraints, and query paths.
2. Identify the entity lifecycle: create, update, delete, ownership, uniqueness, tenancy, and audit needs.
3. Put integrity in PostgreSQL first, then align the ORM and Pydantic models with it.
4. Design the migration path before adding fields or tables.

## One-Screen Defaults

- Give every table a stable primary key. Usually use `bigint` identity; use UUID only when offline or distributed ID generation is a real need.
- Keep mutable business identifiers such as email, slug, and external IDs as `UNIQUE`, not as the primary key.
- Default to normalized tables for data you filter, join, constrain, sort, or report on. Use `jsonb` only for sparse or semi-structured metadata.
- Most columns in business tables should be `NOT NULL`. Missing data should be an explicit business choice, not the default.
- Use `timestamptz` for lifecycle timestamps such as `created_at`, `updated_at`, `deleted_at`, and `expires_at`.
- Enforce invariants with `NOT NULL`, `CHECK`, `UNIQUE`, `FOREIGN KEY`, and partial unique indexes when needed.
- Index foreign keys on the referencing side.
- Choose `ON DELETE` behavior explicitly. Do not let ORM cascade settings silently define lifecycle rules.
- Add a `version` column for high-contention editable rows.
- Use soft delete only if recovery or audit requirements justify it, and pair it with partial unique indexes plus a purge or archive plan.
- Keep audit and history tables separate from hot OLTP tables.
- Use Alembic migrations, not `create_all()`, for production schema evolution.

## Workflow

1. Inspect the current tables, models, migrations, and the feature's likely reads and writes.
2. Decide the storage shape:
   - core entity tables
   - lookup and reference tables
   - join tables
   - history and audit tables
3. Define row identity, business uniqueness, foreign keys, and delete behavior.
4. Define indexes from actual query patterns, not guesswork.
5. Align the Python stack:
   - SQLAlchemy or SQLModel table models for persistence
   - Pydantic or SQLModel data models for create, read, and update contracts
6. Plan migrations with `expand -> backfill -> validate -> contract`.
7. Call out risks and anti-patterns before editing code.

## Rules

### PostgreSQL Is The Source Of Truth

- Pydantic validates input shape. PostgreSQL enforces durable correctness.
- If a rule still matters after a bug, deploy, import, or manual SQL change, encode it in the database.

### Keep Model Layers Separate

- Do not reuse one model for table mapping, create input, patch input, and response output.
- SQLModel can reduce duplication, but you still need distinct models when optionality, secrets, or response shape differ.

### Relationship Rules

- One-to-one needs a database `UNIQUE` constraint on the foreign key.
- Many-to-many needs a real link table with a composite primary key or composite unique constraint.
- If the link table has extra columns, treat it as a first-class entity.
- Use composite foreign keys only when the domain boundary truly depends on them, such as tenant-scoped references.

### JSON And Enums

- Do not store core relational data in `jsonb`.
- Use PostgreSQL enums only for small, stable value sets. If admins or tenants may manage the values, use a lookup table instead.

### Soft Delete And History

- Soft delete is not free. It adds bloat, index churn, and uniqueness complexity.
- If you soft-delete, default application queries to active rows and enforce active-row uniqueness with partial indexes.
- If the real need is audit, prefer append-only history or audit tables over keeping everything forever in the hot table.

### Migrations And CNPG

- Favor additive migrations.
- On busy tables, prefer `CREATE INDEX CONCURRENTLY`, staged constraint rollout, and resumable backfills.
- Avoid rewrite-heavy DDL and giant one-transaction backfills.
- CloudNativePG recovery is cluster-level and WAL-driven, so rollback safety depends on old and new app versions tolerating the expanded schema.

## Output

When using this skill, return:

1. Recommended schema shape.
2. Constraints and indexes to add.
3. ORM and API model boundaries.
4. Migration rollout notes.
5. The main anti-patterns or risks.

## Read More

- For durable PostgreSQL defaults and anti-patterns, read `references/postgres-crud-schema.md`.
- For SQLAlchemy, SQLModel, Alembic, and Pydantic alignment, read `references/python-orm-stack.md`.
- For CloudNativePG and migration rollout concerns, read `references/cnpg-migrations.md`.
