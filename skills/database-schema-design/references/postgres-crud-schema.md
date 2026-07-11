# PostgreSQL CRUD Schema Defaults

Use this file when the task is mainly about table shape, data integrity, indexing, delete behavior, or long-term PostgreSQL durability.

## Durable Defaults

- Give every table a stable primary key.
- Usually prefer `bigint GENERATED ... AS IDENTITY` for long-lived CRUD apps.
- Keep natural keys as `UNIQUE` business constraints, not as mutable primary keys.
- Most business columns should be `NOT NULL`.
- Use `timestamptz` for lifecycle timestamps.
- Normalize first. Denormalize only when a measured read path needs it.

## When To Normalize Versus Use `jsonb`

Use regular columns and tables when the data is:

- filtered
- sorted
- joined
- constrained
- reported on
- updated independently

Use `jsonb` for:

- sparse metadata
- third-party payload snapshots
- extension points that do not define the main domain model

Avoid giant `jsonb` blobs for core entity state. CRUD-heavy systems pay for that later through poor constraints, harder indexing, larger row rewrites, and more lock contention.

## Constraints To Prefer Early

- `NOT NULL` for most business columns
- `UNIQUE` for durable business identifiers such as slug, external ID, or `(tenant_id, code)`
- `CHECK` for same-row invariants such as positive amounts, bounded percentages, or start-before-end rules
- `FOREIGN KEY` for ownership and referential integrity
- partial unique indexes for active-row-only uniqueness, such as soft-delete patterns

Do not leave integrity only in app code. If a bad import, one-off script, or later bug would matter, the rule belongs in PostgreSQL.

## Relationship Defaults

- Index foreign keys on the referencing side. PostgreSQL does not create these automatically.
- Pick `ON DELETE` behavior intentionally:
  - `CASCADE` when the child cannot outlive the parent
  - `RESTRICT` or `NO ACTION` when the entities are independent
  - `SET NULL` only when the relationship is truly optional
- One-to-one means a foreign key plus a `UNIQUE` constraint.
- Many-to-many means a real link table, usually with a composite primary key.

## Timestamp And Time Zone Rules

- Use `timestamptz` for events, lifecycle timestamps, and anything that crosses systems or users.
- Avoid `time with time zone`.
- Use `timestamp without time zone` only when the value is intentionally local and naive.

## Soft Delete Guidance

Use soft delete only when you need recoverability, legal hold, or a first-class restore workflow.

If you do use it:

- add `deleted_at timestamptz null`
- default application queries to `WHERE deleted_at IS NULL`
- enforce active-row uniqueness with partial unique indexes
- add a purge or archive plan

Soft delete is not free. It increases dead tuples, index churn, and autovacuum pressure.

## Enums, Lookup Tables, Generated Columns

Use PostgreSQL enums only for small, stable sets such as a handful of internal lifecycle states.

Use a lookup table instead when values may:

- be renamed
- be retired
- be extended by operators or tenants
- carry labels, sort order, or permissions

Use generated columns only for deterministic same-row derivations. Do not use them for cross-row rules or workflow logic.

## Concurrency, Audit, And Tenancy

- Add an explicit `version` column for rows that humans may edit concurrently.
- Prefer append-only history or audit tables over stuffing change history into hot entity tables.
- Make tenancy boundaries explicit with `tenant_id` in schema, constraints, and uniqueness. Do not rely only on app filters.

## PostgreSQL-Specific Tools Worth Knowing

- partial indexes: great for `deleted_at IS NULL`, `status = 'active'`, and conditional uniqueness
- exclusion constraints: the right tool for overlap rules such as reservations and time ranges
- named constraints: easier migrations, better errors, easier operations

## Common Anti-Patterns

- mutable natural keys as the only primary key
- one giant `jsonb` column for most business data
- no foreign keys because the ORM "knows" the relationships
- no index on referencing foreign keys
- soft delete without partial uniqueness or purge strategy
- `timestamp without time zone` for real-world events
- using enums for values admins will want to manage
- using `CHECK` for cross-row rules that should be `UNIQUE`, `EXCLUDE`, or `FOREIGN KEY`

## Source Links

- PostgreSQL constraints: https://www.postgresql.org/docs/current/ddl-constraints.html
- PostgreSQL identity columns: https://www.postgresql.org/docs/current/ddl-identity-columns.html
- PostgreSQL generated columns: https://www.postgresql.org/docs/current/ddl-generated-columns.html
- PostgreSQL date and time types: https://www.postgresql.org/docs/current/datatype-datetime.html
- PostgreSQL enum types: https://www.postgresql.org/docs/current/datatype-enum.html
- PostgreSQL JSON and JSONB: https://www.postgresql.org/docs/current/datatype-json.html
- PostgreSQL partial indexes: https://www.postgresql.org/docs/current/indexes-partial.html
- PostgreSQL range types: https://www.postgresql.org/docs/current/rangetypes.html
- PostgreSQL row security: https://www.postgresql.org/docs/current/ddl-rowsecurity.html
- PostgreSQL MVCC intro: https://www.postgresql.org/docs/current/mvcc-intro.html
- PostgreSQL routine vacuuming: https://www.postgresql.org/docs/current/routine-vacuuming.html
- PostgreSQL audit trigger notes: https://wiki.postgresql.org/wiki/Audit_trigger
