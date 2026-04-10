# Trakt API v2 — Reference Notes

Source of truth: https://trakt.docs.apiary.io/ (Apiary, JS-rendered).
Cross-checked against `vankasteelj/trakt.tv` client source and Trakt forum/GitHub threads.

- Base URL (prod): `https://api.trakt.tv`
- Base URL (staging): `https://api-staging.trakt.tv`
- All request/response bodies are JSON.

## Required headers

Sent on **every** request:

| Header | Value | Notes |
|---|---|---|
| `Content-Type` | `application/json` | Always, even on GET |
| `trakt-api-version` | `2` | API major version |
| `trakt-api-key` | `{client_id}` | App Client ID |
| `User-Agent` | `<lib>/<version>` | Recommended; identifies your client |
| `Authorization` | `Bearer {access_token}` | Only when OAuth is needed |

- `trakt-api-key` is **always** required, even when an OAuth `Authorization` header is also present.
- No other mandatory headers.

## Authentication

Two OAuth 2.0 flows are supported. Endpoints (relative to base URL):

- Authorization code:
  - `GET /oauth/authorize` — user-facing, `response_type=code`, `client_id`, `redirect_uri`, `state`
  - `POST /oauth/token` — exchange `code` → access token (`grant_type=authorization_code`)
- Device code:
  - `POST /oauth/device/code` — returns `device_code`, `user_code`, `verification_url`, `interval`, `expires_in`
  - `POST /oauth/device/token` — poll with `code` (the `device_code`), `client_id`, `client_secret`
  - While pending, Trakt returns HTTP 400; keep polling until 200 or HTTP 404/409/410/418/429
- Refresh:
  - `POST /oauth/token` with `grant_type=refresh_token`, `refresh_token`, `client_id`, `client_secret`, `redirect_uri`
- Revoke:
  - `POST /oauth/revoke` with `token`, `client_id`, `client_secret`

Bearer token is sent as `Authorization: Bearer {access_token}`.

## Pagination

- Trigger: add `?page={n}&limit={n}` to any paginated endpoint.
  - Default `limit`: **10**
  - Practical max `limit`: high (thousands), varies per endpoint. Treat as soft-unbounded but do not assume > 1000.
  - `page` is 1-indexed.
- Response headers returned on paginated endpoints:

| Header | Meaning |
|---|---|
| `X-Pagination-Page` | Current page (1-indexed) |
| `X-Pagination-Limit` | Current page size |
| `X-Pagination-Page-Count` | Total number of pages |
| `X-Pagination-Item-Count` | Total number of items across all pages |

- Edge cases / gotchas:
  - Some endpoints return `X-Pagination-Page-Count` and `X-Pagination-Item-Count` but **omit** `X-Pagination-Page` / `X-Pagination-Limit`. Treat those two as optional on parse.
  - `X-Pagination-Item-Count` has known inaccuracies on some endpoints (e.g. `/users/{id}/history`, personal list items) — do **not** rely on it for strict equality checks.
  - Not every list endpoint is paginated. If no `X-Pagination-*` headers are present, the endpoint returns the full collection.
- Typically paginated endpoints: `/users/{id}/history`, `/users/{id}/ratings/*`, `/users/{id}/lists/{id}/items`, `/search/*`, `/shows/{id}/comments`, `/movies/{id}/comments`, `/sync/history`, trending/popular/anticipated/boxoffice listings.

## Rate limiting

Documented limits (per user / per client):

| Bucket | Limit |
|---|---|
| Unauthenticated `GET` | 1000 calls / 5 min |
| Authenticated `GET` | 1000 calls / 5 min |
| Authenticated `POST` / `PUT` / `DELETE` | 1 call / sec |

- On exceed: HTTP `429 Too Many Requests`.
- Response headers include an `X-Ratelimit` header carrying a JSON object with the bucket name (e.g. `AUTHED_API_POST_LIMIT`, `UNAUTHED_API_GET_LIMIT`), `period`, `limit`, `remaining`, `until`.
  - Parse it as JSON, not as a scalar.
- `Retry-After` header **may** be present on 429; always respect it when present.
- Note (undocumented): Trakt periodically applies stricter limits and Cloudflare bot-detection; 429s can occur below documented thresholds. Implement exponential backoff on 429.

## Errors

- All error responses use standard HTTP status codes. Common ones:

| Status | Meaning |
|---|---|
| `400` | Bad Request (bad body or query) |
| `401` | Unauthorized (missing/invalid OAuth token) |
| `403` | Forbidden (valid token, insufficient scope / VIP-only) |
| `404` | Not Found |
| `405` | Method Not Allowed |
| `409` | Conflict (e.g. resource already exists) |
| `412` | Precondition Failed |
| `422` | Unprocessable Entity (validation failed) |
| `429` | Rate Limit Exceeded |
| `500` / `502` / `503` / `504` | Server errors / maintenance |

- Error response body:
  - OAuth endpoints (`/oauth/*`) return a JSON body shaped like `{"error": "...", "error_description": "..."}` (OAuth 2.0 spec).
  - Other endpoints often return an **empty body** or a minimal `{"error": "..."}`; do not rely on a structured body being present. Prefer using HTTP status + endpoint context for error classification.
- On `401`, Trakt sets the `WWW-Authenticate` header with a human-readable description — worth surfacing in typed errors.

## Extended info

- Query param: `?extended={level}`
- Supported levels:
  - *(omitted)* — minimal: title/name, year, standard ID set only
  - `full` — full metadata (overview, rating, votes, runtime, genres, air dates, certification, etc.)
  - `metadata` — sync/collection metadata (resolution, audio, etc.) on collection endpoints
  - `episodes` — expand season responses to include all episode objects
  - `noseasons` — on show summary, skip seasons expansion
  - `guest_stars` — on episode/season endpoints, include guest stars on people lookups
  - `vip` — VIP-only fields where applicable
- Multiple values are comma-separated: `?extended=full,metadata`.
- Not every endpoint honors every level; unknown levels are silently ignored.
- The historical `images` extended level is **deprecated / removed** — Trakt no longer serves images via the API. Ignore any old references.

## Filters

- Sent as query-string params alongside `extended` / `page` / `limit`. Multi-value params are comma-delimited.
- Common filters (apply to `/search/*`, `/movies/{trending,popular,anticipated,boxoffice,updates}`, `/shows/{trending,popular,anticipated,updates}`, calendars, etc.):

| Param | Type | Example | Notes |
|---|---|---|---|
| `query` | string | `query=batman` | Text search; `/search/*` only |
| `years` | int or range | `years=2016` / `years=2010-2020` | 4-digit year |
| `genres` | slug list | `genres=action,adventure` | Comma = OR |
| `languages` | 2-char list | `languages=en,fr` | ISO 639-1 |
| `countries` | 2-char list | `countries=us,gb` | ISO 3166-1 alpha-2 |
| `runtimes` | int range | `runtimes=30-90` | Minutes |
| `ratings` | int range | `ratings=75-100` | 0–100 (Trakt rating %) |
| `votes` | int range | `votes=5000-100000` | Total vote count |
| `certifications` | slug list | `certifications=pg,pg-13` | Movies only |
| `networks` | slug list | `networks=hbo,netflix` | Shows only |
| `status` | slug list | `status=returning%20series,ended` | Shows only; URL-encode spaces |
| `studio_ids` | id list | `studio_ids=1,2` | Movies/shows (VIP) |

- All values must be URL-encoded.
- Range filters use `min-max`; open ranges are not supported — pass a concrete min and max.
- Applying an unsupported filter to an endpoint is silently ignored.

## Misc conventions

- IDs: every media object carries an `ids` sub-object with `trakt`, `slug`, `imdb`, `tmdb`, `tvdb` (where applicable). Endpoints accept any of `{trakt_id | slug | imdb_id}` in path position.
- Dates: ISO-8601 UTC with trailing `Z` (e.g. `2024-01-15T12:34:56.000Z`).
- POST/PUT bodies are JSON; always send `Content-Type: application/json` even when the body is empty (`{}`).
- Sync endpoints (`/sync/*`) return per-item add/update/not_found breakdowns, not flat arrays — keep that in mind when modeling responses.
