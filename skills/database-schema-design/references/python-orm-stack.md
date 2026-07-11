# Python ORM Stack Guidance

Use this file when the task is about SQLAlchemy, SQLModel, Alembic, Pydantic, or the boundary between schema design and Python model design.

## Layer Boundaries

Keep these concerns separate:

- table or ORM models for persistence
- create models for required input
- patch models for partial updates
- response models for output

Do not force one class to serve all four roles.

Why:

- pre-insert ORM objects often allow `id = None`
- response models should usually require IDs and derived fields
- create and patch models have different optionality
- response models often omit secrets and internal fields

## SQLAlchemy Defaults

Prefer SQLAlchemy 2.x annotated declarative style with `Mapped[...]` and `mapped_column()`.

Benefits:

- clearer nullability and typing
- direct alignment between Python annotations and schema intent
- easier review of persistence models

Be deliberate about mismatches. If Python allows `None` but the column is `NOT NULL`, do that only when object creation flow truly needs it.

## Naming Convention And Alembic

Add a metadata naming convention on day one and use the same metadata in Alembic.

```python
from sqlalchemy import MetaData
from sqlalchemy.orm import DeclarativeBase

convention = {
    "ix": "ix_%(column_0_label)s",
    "uq": "uq_%(table_name)s_%(column_0_name)s",
    "ck": "ck_%(table_name)s_%(constraint_name)s",
    "fk": "fk_%(table_name)s_%(column_0_name)s_%(referred_table_name)s",
    "pk": "pk_%(table_name)s",
}

metadata = MetaData(naming_convention=convention)


class Base(DeclarativeBase):
    metadata = metadata
```

Why this matters:

- stable migration diffs
- easier constraint changes and drops
- clearer database errors
- fewer autogenerate surprises

## Relationships And Delete Semantics

Align ORM behavior and database behavior.

Database side:

- use `ForeignKey(..., ondelete=...)`
- use `UNIQUE` for one-to-one
- use real link tables for many-to-many

ORM side:

- use `cascade` intentionally
- use `delete-orphan` only when child rows have no life outside the parent
- use `passive_deletes=True` when the database is doing the cascade work

Do not let ORM cascade settings accidentally define lifecycle rules that the database does not enforce.

## Optimistic Locking

For rows that users commonly edit concurrently, add a version column and wire it into SQLAlchemy versioning.

- use a `NOT NULL` version column
- prefer explicit integer versioning for normal CRUD apps
- do not depend on system columns as a durable app contract

## Alembic Rules

- Use autogenerate as a draft, not as truth.
- Review every migration manually.
- Renames usually need hand-written migrations.
- Separate schema rollout from data backfill when possible.
- Prefer additive change first, cleanup later.

Good habits:

- set `target_metadata` correctly
- keep type comparison enabled unless you have a reason not to
- run `alembic check` or an equivalent migration sanity check in CI

## SQLModel Guidance

SQLModel is useful for straightforward CRUD services, but it does not remove the need for schema design.

Use SQLModel when:

- the project benefits from lower ceremony
- the mapping is conventional
- the team wants one typed stack for ORM plus data models

Still keep distinct models for:

- table models
- create input
- patch input
- read output

Drop to plain SQLAlchemy when the mapping or migration behavior becomes more advanced than SQLModel makes comfortable.

## Pydantic Guidance

For pure SQLAlchemy plus Pydantic:

- use explicit Pydantic models
- use `from_attributes=True` for ORM-to-API conversion
- keep validation rules close to the API contract, not as a replacement for database constraints

Use strict validation where bad coercion would hide client errors.

## Common Anti-Patterns

- one class shared by DB table, request body, and response body
- unnamed constraints and indexes
- assuming autogenerate detects renames
- one-to-one relationships without a database unique constraint
- many-to-many link tables without a composite uniqueness rule
- database `ondelete` rules without matching ORM configuration
- bulk updates treated as if ORM version checks or cascades still apply

## Source Links

- SQLAlchemy declarative tables: https://docs.sqlalchemy.org/en/20/orm/declarative_tables.html
- SQLAlchemy relationships: https://docs.sqlalchemy.org/en/20/orm/basic_relationships.html
- SQLAlchemy cascades: https://docs.sqlalchemy.org/en/20/orm/cascades.html
- SQLAlchemy constraints and naming conventions: https://docs.sqlalchemy.org/en/20/core/constraints.html
- SQLAlchemy version counters: https://docs.sqlalchemy.org/en/20/orm/versioning.html
- Alembic tutorial: https://alembic.sqlalchemy.org/en/latest/tutorial.html
- Alembic naming conventions: https://alembic.sqlalchemy.org/en/latest/naming.html
- Alembic autogenerate: https://alembic.sqlalchemy.org/en/latest/autogenerate.html
- Pydantic ORM example: https://docs.pydantic.dev/latest/examples/orms/
- SQLModel features: https://sqlmodel.tiangolo.com/features/
- SQLModel create tables and migrations note: https://sqlmodel.tiangolo.com/tutorial/create-db-and-table/
- SQLModel multiple models: https://sqlmodel.tiangolo.com/tutorial/fastapi/multiple-models/
- SQLModel cascade delete relationships: https://sqlmodel.tiangolo.com/tutorial/relationship-attributes/cascade-delete-relationships/
