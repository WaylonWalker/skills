# CloudNativePG And Migration Rollout

Use this file when the task touches Alembic rollout strategy, hot-table migrations, index creation, or operational safety on CloudNativePG.

## Operator Reality

CloudNativePG recovery is cluster-level. Recovery typically means restoring a new cluster from a base backup plus WAL, not undoing one migration in place.

That changes the design goal:

- expanded schemas should remain usable by old and new app versions during the safety window
- migrations should minimize WAL spikes, lock time, and rewrite-heavy DDL

## Default Rollout Pattern

Use `expand -> backfill -> validate -> contract`.

1. Expand with additive schema changes.
2. Backfill in small committed batches.
3. Validate constraints after data is in shape.
4. Contract only after the app has fully switched over and the rollback window has passed.

This is the safest default for busy self-hosted CRUD systems.

## Safe Change Patterns

- add nullable columns first
- add constant defaults when needed
- create indexes concurrently on hot tables
- add `CHECK` and `FOREIGN KEY` constraints as `NOT VALID` when supported
- validate later
- make backfills resumable and idempotent

## Operations To Treat As Dangerous

- table rewrites on large hot tables
- volatile defaults that force rewrites
- one huge transaction for a backfill
- immediate uniqueness or not-null enforcement before the data is ready
- dropping old columns before a safe compatibility window has passed

## Lock-Sensitive DDL

Assume `ALTER TABLE` is risky until you have checked its lock behavior.

For busy tables:

- prefer `CREATE INDEX CONCURRENTLY`
- use `CREATE UNIQUE INDEX CONCURRENTLY` followed by `ALTER TABLE ... ADD CONSTRAINT ... USING INDEX` when rolling out uniqueness safely
- split foreign key rollout from backfill work

Remember that `CREATE INDEX CONCURRENTLY` cannot run inside a normal transaction block and can leave an invalid index behind on failure.

## Backfill Rules

- process rows in batches
- commit often
- throttle if needed
- analyze after large backfills or index creation
- avoid long transactions that delay vacuum cleanup

Large backfills are expensive not only for the app but for replicas, WAL archiving, backup storage, and recovery time.

## Index Discipline

Every extra index adds:

- write amplification
- more WAL
- more vacuum work
- more backup and recovery cost

Add indexes from observed query paths. Do not index every column "just in case."

## Schema Patterns That Help Operations

- stable surrogate keys
- narrow hot tables
- separate audit or history tables from hot CRUD tables
- additive changes before destructive cleanup
- partitioning only when there is a real retention or scale problem

## Schema Patterns That Hurt Operations

- giant `UPDATE` jobs in one transaction
- frequent type rewrites on large tables
- deep cascade-heavy cleanup during busy periods
- late generated-column or computed-column rewrites on hot tables
- over-indexing

## Source Links

- CloudNativePG backup: https://cloudnative-pg.io/docs/1.29/backup
- CloudNativePG WAL archiving: https://cloudnative-pg.io/docs/1.29/wal_archiving
- CloudNativePG recovery: https://cloudnative-pg.io/docs/1.29/recovery
- PostgreSQL `ALTER TABLE`: https://www.postgresql.org/docs/current/sql-altertable.html
- PostgreSQL `CREATE INDEX`: https://www.postgresql.org/docs/current/sql-createindex.html
- PostgreSQL explicit locking: https://www.postgresql.org/docs/current/explicit-locking.html
- PostgreSQL routine vacuuming: https://www.postgresql.org/docs/current/routine-vacuuming.html
- PostgreSQL continuous archiving and PITR: https://www.postgresql.org/docs/current/continuous-archiving.html
- PostgreSQL populating a database: https://www.postgresql.org/docs/current/populate.html
