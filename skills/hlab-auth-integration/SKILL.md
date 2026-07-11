---
name: hlab-auth-integration
description: Integrate downstream apps with `hlab-auth` using the supported v1 trust model, and guide least-privilege operational automation around apps, groups, roles, and tokens. Use this whenever the user wants to protect an app or internal tool with `hlab-auth`, Traefik `ForwardAuth`, trusted `X-Hlab-*` headers, reverse-proxy auth, app-local role mapping, homelab SSO around `hlab-auth`, or a narrowly scoped agent credential for app-related changes, even if they do not explicitly ask for an "integration skill." Also use it when the request sounds like OIDC or SAML with `hlab-auth`, so you can classify that expectation early and redirect to the supported proxy-header model or clearly explain the v1 limitation.
---

# hlab-auth Integration

Use this skill to add `hlab-auth` protection to another app without inventing a new trust model.

## Outcome

When this skill triggers, the agent should usually:

1. Inspect the target app and deployment shape.
2. Classify the request against `hlab-auth` integration tiers.
3. Implement the needed proxy and app changes when the integration is feasible.
4. Verify the trust boundary and document any remaining operator steps.

When the user wants an agent to manage `hlab-auth` state after approval, also teach the least-privilege credential path rather than assuming the agent should run with a human's broad session.

## Start Here

1. Inspect the target repo, proxy setup, app runtime, and deployment docs before proposing changes.
2. Identify whether the app already sits behind Traefik, Nginx, or another reverse proxy that can enforce external auth.
3. Identify whether the app can consume trusted upstream headers directly, auto-provision local users from them, or needs deeper adapter code.
4. If the request assumes `hlab-auth` is an OIDC or SAML provider, say clearly that this is out of scope for v1 and pivot to the supported proxy-auth pattern.

## v1 Trust Model

`hlab-auth` v1 is not "every app logs into auth directly."

- `hlab-auth` authenticates the user and decides whether the request is allowed.
- The reverse proxy calls `hlab-auth` before the app sees the request.
- The reverse proxy strips any client-supplied `X-Hlab-*` headers.
- On success, the reverse proxy injects trusted `X-Hlab-*` headers into the upstream request.
- The downstream app consumes those trusted headers and must fail closed when they are missing or malformed.

Do not replace this with browser-supplied headers, app-managed trust shortcuts, or a fake OIDC flow.

## CLI And API Reality

As of the current `hlab-auth` implementation:

- The API supports managing apps, groups, roles, route rules, personal tokens, and app-bound service tokens.
- The `hlab` CLI currently supports device login, `whoami`, session management, personal token management, and app listing.
- The `hlab` CLI does not yet expose full app, group, or role mutation commands.

That means agent-driven administration is possible today, but mutating app, group, and role resources is usually API-first after the user issues an allowed credential.

## Compatibility Tiers

Classify the target integration early and say which tier it is.

### Tier 1: Works Out of the Box

Use this when:

- The proxy can call an external auth endpoint.
- The proxy can strip spoofed inbound auth headers.
- The proxy can inject trusted `X-Hlab-*` headers.
- The app works when auth is enforced entirely at the proxy boundary.

Typical examples: Traefik `ForwardAuth`, Nginx `auth_request`.

### Tier 2: Works With Light App Config

Use this when the Tier 1 proxy behavior works and the app only needs small app-side configuration, such as:

- trusting proxy identity headers
- mapping `X-Hlab-User-Id` to a local user
- auto-provisioning a local profile from trusted headers

### Tier 3: Needs App-Specific Support

Use this when the app needs adapter code, plugin work, or local session minting from upstream identity.

Do the work only if the user wants it and the repo supports it. Be explicit about the extra moving parts.

### Tier 4: Out Of Scope For v1

Reject or reframe requests that require:

- OIDC provider behavior from `hlab-auth`
- SAML support
- direct app exposure without a trusted proxy
- an app that cannot consume proxy auth and cannot be safely wrapped

## Trusted Header Contract

Use only the stable `X-Hlab-*` contract from `references/header-contract.md`.

Default rules:

- `X-Hlab-User-Id` is the durable external subject key.
- `X-Hlab-Username` and `X-Hlab-Display-Name` are display fields, not durable identity keys.
- `X-Hlab-Groups`, `X-Hlab-Roles`, and `X-Hlab-Scopes` are canonical comma-separated lists.
- Missing trusted headers means unauthenticated.
- Direct client-supplied `X-Hlab-*` headers are never trusted.

## Delegated Agent Management

Use this section when the user wants an agent to create or manage `hlab-auth` resources after explicit approval.

The safe default is not "give the agent your full admin session." The safe default is:

1. Decide what the agent needs to manage.
2. Create the smallest practical scope set.
3. Prefer a dedicated token over reusing a browser session.
4. Prefer an app-bound service token when the work is truly app-specific.
5. Prefer a dedicated human-managed automation user only when app-bound token limits or audit needs require it.

### Credential Choice

Choose credentials in this order:

1. App-bound service token
   - Best when the agent only needs to work on one app and its related automation surface.
   - In v1, service token scopes must be a subset of that app's declared route-rule scopes.
2. Narrow personal token
   - Best when the agent needs admin API access that service tokens do not cover.
   - The token scopes must be a subset of the issuing user's effective scopes.
3. Dedicated automation user plus narrow personal token
   - Best when the user wants separation from their daily admin identity, explicit audit ownership, or tighter role assignment.

Do not tell the user to hand an agent a broad long-lived admin token when a smaller token would work.

### Scope Rules

- For app inventory or route inspection, prefer `apps:read`.
- For app creation, route-rule changes, or app-token rotation, use `apps:write`.
- For group inspection or mutation, use `groups:read` and only add `groups:write` if needed.
- For role inspection or mutation, use `roles:read` and only add `roles:write` if needed.
- For token self-management, use `tokens:self`.
- Avoid unrelated scope families.

If the user says "the agent only works on Grafana," keep the credential limited to the smallest app-related surface that can satisfy that task.

### Recommended User Instructions

When teaching a user to authorize an agent, instruct them to:

1. Log in with `hlab login` or a browser session.
2. Confirm current identity and scopes with `hlab whoami` when useful.
3. Create a dedicated token with only the needed scopes and an expiry.
4. Store that token in a tool-specific secret store or environment variable, not in source control.
5. Use the token only for the target app or management workflow.
6. Revoke and rotate the token when the task is done or if exposure is suspected.

### Usage Guidance For Agents

- Prefer `Authorization: Bearer <token>` against `/api/v1/...` for machine actions.
- Use the CLI for login and simple self-service token workflows when it helps the user bootstrap access.
- Use the API directly for mutating apps, groups, roles, memberships, route rules, or app tokens.
- Remind the user that token secrets are shown only once at creation time.
- Recommend expiries for agent-issued tokens unless the user has a concrete need for a longer-lived credential.

Read `references/agent-credentials.md` for the practical workflow and example API calls.

## Integration Workflow

1. Inspect the current deployment path.
   - Where does traffic enter?
   - Which proxy owns auth enforcement?
   - Can the app be reached directly outside that path?
2. Classify the app.
   - Tier 1, 2, 3, or 4.
   - State the reason in one or two sentences.
3. Identify the app's auth boundary.
   - Does the app need only upstream identity?
   - Does it need a local user record?
   - Does it have dangerous behavior keyed off usernames or emails?
4. Implement the integration.
   - Add or update reverse-proxy config.
   - Strip spoofable inbound headers.
   - Forward only the stable `X-Hlab-*` headers.
   - Update app middleware or request handling to read trusted headers.
   - Fail closed when proxy provenance or trusted headers are absent.
5. Configure authorization.
   - Let `hlab-auth` decide whether a request may reach the app.
   - Keep any app-local checks narrower than the forwarded scope set.
   - Never widen privilege based on `X-Hlab-Username=admin` or similar shortcuts.
6. Verify.
   - Confirm protected requests succeed through the proxy.
   - Confirm direct requests without trusted headers are denied or treated as anonymous.
   - Confirm spoofed `X-Hlab-*` headers from a client do not grant access.

## Implementation Rules

- Prefer Traefik `ForwardAuth` when the environment already uses Traefik.
- For Nginx, keep the trust model identical even if the directives differ.
- Keep the app behind the proxy rather than exposing it directly.
- If direct network access must exist for operations, make the app reject requests that do not come from the trusted proxy path.
- Recompute request identity from trusted headers on each request unless the app has a strong reason to mint a short-lived local session.
- If the app persists identity locally, link on `X-Hlab-User-Id`.
- Log the upstream user ID where useful for audit correlation.
- Do not introduce alternate unofficial header names.
- Do not quietly promise unsupported v1 protocols.

## Output

Return these items when the skill is used:

1. The compatibility tier and why.
2. The concrete proxy changes.
3. The concrete app changes.
4. The verification steps you ran or recommend.
5. Any remaining operator actions, risks, or out-of-scope requirements.

If delegated management is part of the task, also include:

6. The recommended credential type.
7. The exact minimal scope set.
8. The command or API flow the user should use to mint and revoke the credential.

## Short Examples

**Good pattern:**

- Traefik calls `https://auth.home.arpa/api/v1/forward-auth`.
- Traefik strips inbound `X-Hlab-*` headers and injects only the trusted auth response headers.
- The app reads `X-Hlab-User-Id` and `X-Hlab-Scopes` from the proxied request.

**Bad pattern:**

- The app trusts `X-Hlab-Username` from any request on the LAN.
- The app grants admin rights because the username string equals `admin`.
- The integration claim says "use hlab-auth as OIDC" with no supported provider flow.

## Ask Only When Needed

Ask a short clarifying question only when a missing detail would change the implementation materially, such as:

- which reverse proxy is actually in front of the app
- whether direct access to the app exists
- which routes or hostnames must be protected
- which scopes should gate access

Otherwise, inspect the repo and make the integration.

## Read More

- Read `references/header-contract.md` for the stable header contract and auth-order rules.
- Read `references/agent-credentials.md` for least-privilege user instructions, token selection, and API-first management examples.
