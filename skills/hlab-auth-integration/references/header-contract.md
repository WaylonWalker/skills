# hlab-auth Header Contract

This reference bundles the parts of the `hlab-auth` integration contract that a downstream app agent needs most often.

## Stable Headers

Use this exact header set:

- `X-Hlab-User-Id`
- `X-Hlab-Username`
- `X-Hlab-Display-Name`
- `X-Hlab-Email` when available and allowed
- `X-Hlab-Groups`
- `X-Hlab-Roles`
- `X-Hlab-Scopes`

Rules:

- Treat header names as case-insensitive HTTP headers.
- Treat values as proxy assertions, not browser assertions.
- Treat `X-Hlab-Groups`, `X-Hlab-Roles`, and `X-Hlab-Scopes` as canonical comma-separated lists.
- Do not accept alternate unofficial header names.
- If the app stores identity locally, use `X-Hlab-User-Id` as the durable linkage key.

## Responsibility Split

### `hlab-auth`

- Authenticates the user.
- Evaluates whether the request may reach the app.
- Returns the trusted identity header set.

### Reverse Proxy

- Calls the auth endpoint before the app sees the request.
- Strips spoofed inbound `X-Hlab-*` headers.
- Injects the trusted `X-Hlab-*` headers into the upstream request.

### Downstream App

- Trusts `X-Hlab-*` only when the request arrived through the protected proxy path.
- Treats missing trusted headers as unauthenticated.
- Applies only equal or narrower authorization than the forwarded scope set.
- Never widens privilege based on a display field such as username.

## Authorization Order

1. Let `hlab-auth` decide whether the request may reach the app at all.
2. Let the app apply narrower app-local checks if needed.
3. Never let the app widen access beyond the forwarded scopes.

## Deployment Rules

- Prefer keeping the app reachable only behind the protected proxy.
- If direct access must exist, reject requests that lack trusted proxy provenance.
- Local development shortcuts must not become the default production path.
