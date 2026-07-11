---
name: shots-api
description: Integrate with a running Shot Scraper API instance to request screenshots, poll queued jobs, fetch image responses, inspect queue and URL stats, or delete cached shots. Use this whenever the user mentions shots, shot-scraper-api, screenshot generation over HTTP, `/shot`, `/shot/async`, `/shot/blocking`, `/trigger/shot`, `/job/{job_id}`, or wants to wire another service, script, or app to screenshot endpoints, even if they only describe the workflow and do not name the API.
---

# Shots API

Use this skill when an agent needs to consume a running Shot Scraper API service.

Start by identifying the base URL before making assumptions. Prefer a user-provided URL. If none is provided, inspect the repo, docs, env files, or local dev config to find the running server address.

## Endpoint selection

Choose the smallest endpoint that matches the task:

- Use `GET /shot` when the caller wants the image response directly.
- Use `GET /shot/blocking` when you want the same blocking behavior but a more explicit endpoint.
- Use `GET /shot/async` when the caller wants JSON queue status instead of image bytes.
- Use `POST /trigger/shot` when the caller wants to enqueue work and receive `job_id`, `job_url`, and `result_url`.
- Use `GET /job/{job_id}` to poll queued or active work.
- Use `GET /queue/stats` or `GET /urls/stats` for operational visibility.
- Use `DELETE /shot` or `DELETE /shot/{filename}` to remove cached shots.

## Request rules

Use the API's real parameter names and validation rules:

- `url` is required and must start with `http`.
- `width` defaults to `800`.
- `height` defaults to `450`.
- `scaled_width` and `scaled_height` are optional and default to the unscaled dimensions.
- `selectors` is a comma-separated string.
- `format` supports `webp`, `png`, `jpg`, and `jpeg`.
- `v` must be a positive integer when provided.
- `timeout` is in milliseconds and cannot exceed `60000`.
- `theme` must be `light` or `dark`.
- `wait` is used on blocking endpoints, must be positive, and cannot exceed `120000`.
- `priority` is only for `POST /trigger/shot` and must be between `0` and `10`.

If the server returns relative paths such as `job_url` or `result_url`, join them with the base URL before presenting them as complete links.

## Response handling

Handle binary and JSON responses differently:

- Never dump image bytes into the terminal.
- Save images to a file with `curl -o` or an equivalent client feature.
- Use `HEAD` when the caller only wants availability or headers.
- Surface HTTP status, response body, and relevant headers when debugging.
- Pay attention to `Content-Type` and `X-Screenshot-Status` on image responses.

Common JSON statuses include queued work, already queued work, existing cached screenshots, and deletion results. Treat these statuses as part of the contract instead of collapsing them into a generic success message.

## Recommended workflow

1. Confirm or discover the API base URL.
2. Pick the endpoint based on whether the caller wants image bytes, async status, queueing, polling, stats, or deletion.
3. Build the request with the exact query parameter names above.
4. Execute the request with a tool that preserves headers and body visibility.
5. If the response is an image, save it and report the file path plus key headers.
6. If the response is JSON, summarize the returned status fields exactly.
7. If the request fails, show the method, URL, status code, and response body before proposing changes.

## curl examples

Direct blocking screenshot:

```bash
curl -G "$BASE_URL/shot" \
  --data-urlencode "url=https://example.com" \
  --data-urlencode "width=1280" \
  --data-urlencode "height=720" \
  --data-urlencode "format=webp" \
  -o example.webp
```

Async status response:

```bash
curl -G "$BASE_URL/shot/async" \
  --data-urlencode "url=https://example.com" \
  --data-urlencode "theme=dark"
```

Queued job response:

```bash
curl -X POST -G "$BASE_URL/trigger/shot" \
  --data-urlencode "url=https://example.com" \
  --data-urlencode "priority=5"
```

Poll a job:

```bash
curl "$BASE_URL/job/$JOB_ID"
```

Check headers without downloading the body:

```bash
curl -I -G "$BASE_URL/shot" \
  --data-urlencode "url=https://example.com"
```

Delete by URL-derived parameters:

```bash
curl -X DELETE -G "$BASE_URL/shot" \
  --data-urlencode "url=https://example.com" \
  --data-urlencode "format=webp"
```

## Integration guidance

When writing integration code for another app or service:

- Preserve the API's method and query parameter names exactly.
- Treat image endpoints as binary responses, not JSON.
- Treat queue endpoints as JSON contracts that may require polling.
- Prefer explicit timeout handling in the client.
- Return or log enough context to debug bad `url`, `timeout`, `theme`, `wait`, or `priority` values quickly.
- Do not invent undocumented request fields or response fields.
