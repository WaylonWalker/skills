# Agent Credentials And Least-Privilege Workflow

Use this reference when a user wants an agent to manage `hlab-auth` resources after explicit approval.

## What Is Possible Today

Current implementation shape:

- `hlab` CLI supports device login, `whoami`, sessions, personal tokens, and app listing.
- The HTTP API supports create/list/read flows for apps, groups, roles, scopes, and users.
- The HTTP API also supports user token revocation, app-bound service token creation, app-bound service token listing, and app-bound service token revocation.

Practical consequence:

- Use the CLI to help the user authenticate and mint a narrow personal token.
- Use the API for app, group, and role mutations.

## Recommended Credential Order

### 1. App-Bound Service Token

Prefer this when the automation is truly about one app.

- Endpoint family: `/api/v1/apps/{app_id}/tokens`
- Constraint: token scopes must be a subset of the app's declared route-rule scopes
- Good for: app-specific automation, app-owned integrations, and narrow machine use

### 2. Narrow Personal Token

Prefer this when the agent needs admin API endpoints that are not naturally app-bound.

- CLI issuance path: `hlab tokens create`
- API issuance path: `POST /api/v1/me/tokens`
- Constraint: token scopes must be a subset of the issuing user's effective scopes

### 3. Dedicated Automation User + Narrow Personal Token

Prefer this when the user wants:

- clean separation from a day-to-day admin identity
- more explicit audit ownership
- a dedicated role bundle for agent tasks

## Scope Selection

Start with the smallest useful set.

Examples:

- app inspection only: `apps:read`
- app creation or route updates: `apps:write`
- group inspection only: `groups:read`
- group membership or group-role changes: `groups:write`
- role inspection only: `roles:read`
- role creation or updates: `roles:write`
- self-managed token creation and revocation: `tokens:self`

Avoid combining unrelated admin scopes unless the task genuinely requires them.

## User Workflow

### Option A: User Mints A Narrow Personal Token With The CLI

1. Login:

```bash
hlab login
```

2. Optionally inspect current identity:

```bash
hlab whoami
```

3. Create a narrow token:

```bash
hlab tokens create \
  --name agent-grafana-admin \
  --scope apps:write \
  --scope apps:read \
  --expires-at 2026-07-31T00:00:00Z
```

4. Store the returned secret in an environment variable or secret store.

Example environment variable for one shell session:

```bash
export HLAB_TOKEN='replace-with-issued-secret'
```

The token secret is shown only once.

### Option B: User Creates An App-Bound Service Token Through The API

This is useful when the user has already authenticated and wants a token tied to one app.

```bash
curl -X POST "$HLAB_BASE_URL/api/v1/apps/$APP_ID/tokens" \
  -H "Authorization: Bearer $HLAB_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "agent-grafana-service",
    "scopes": ["grafana:access"],
    "expires_at": "2026-07-31T00:00:00Z"
  }'
```

Only do this when the requested scopes fit inside the app's declared route-rule scopes.

## API Examples

### Create A Group

```bash
curl -X POST "$HLAB_BASE_URL/api/v1/groups" \
  -H "Authorization: Bearer $HLAB_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "grafana-editors",
    "description": "Editors for Grafana"
  }'
```

### Create A Role

```bash
curl -X POST "$HLAB_BASE_URL/api/v1/roles" \
  -H "Authorization: Bearer $HLAB_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "grafana-admin",
    "description": "Grafana management role",
    "scopes": ["apps:write", "apps:read"]
  }'
```

### Create An App

```bash
curl -X POST "$HLAB_BASE_URL/api/v1/apps" \
  -H "Authorization: Bearer $HLAB_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "grafana",
    "enabled": true,
    "route_rules": []
  }'
```

### Add An App Route Rule

```bash
curl -X POST "$HLAB_BASE_URL/api/v1/apps/$APP_ID/route-rules" \
  -H "Authorization: Bearer $HLAB_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "host": "grafana.home.arpa",
    "path_prefix": "/",
    "required_scopes": ["grafana:access"],
    "priority": 0
  }'
```

## Safety Rules To Repeat To Users

- Do not paste long-lived broad admin tokens into prompts if a smaller token can do the job.
- Do not commit tokens to the repo.
- Recommend expiries by default.
- Revoke the token after the task if the user does not need to keep it.
- Use a dedicated automation user when audit separation matters.
