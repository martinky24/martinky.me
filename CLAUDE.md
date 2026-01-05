# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Context

Personal portfolio/resume website for Kyle Martin. Intentionally simple implementation - built for functionality over complexity, not as a showcase of coding ability.

**Live Site**: https://martinky.me
**Stack**: Go 1.25 (stdlib only), Bootstrap 5.3, HTML templates
**Deployment**: Docker + Fly.io (Denver region)

## Development Commands

```bash
# Run locally (port 8090 by default, configure via MARTINKY_ME_PORT)
go run server.go

# Run with Docker Compose (port 8080)
docker compose up --build

# Deploy to production
fly deploy
```

## Architecture

Single-file Go server (`server.go`) with no external dependencies. Key patterns:

**Template System**: Go `html/template` with reusable partials defined in `templates/partials.html`. Each page template uses `{{template "head"}}`, `{{template "header"}}`, and `{{template "footer"}}` for consistency.

**Routing**: Simple path-based routing with automatic `.html` extension handling in `checkExt()`. Request `/about` → serves `templates/about.html`. Security headers applied via middleware wrapper to all routes.

**Static Files**: Served from `./static/` with security headers. Includes CSS, SVG favicon, and resume icons.

**Docker**: Multi-stage build pattern:
1. Builder stage: `golang:1.25` compiles with `-ldflags="-s -w" -trimpath` for minimal binary
2. Runtime: `scratch` base image (no shell, no healthcheck possible)

**Logging**: Structured JSON via `log/slog` (not fmt.Println)

## Important Constraints

- **No external dependencies**: Only Go stdlib. Keep it that way.
- **Scratch Docker image**: No shell utilities available in production container (wget, curl, etc.). Container healthchecks won't work.
- **Bootstrap 5 quirk**: Default link underlines removed via `a { text-decoration: none; }` in main.css
- **Security headers**: CSP whitelist includes `cdn.jsdelivr.net` (Bootstrap), `fonts.googleapis.com`, `fonts.gstatic.com`

## Code Patterns

When modifying templates:
- Use partials from `templates/partials.html` for consistency
- Template names must match file paths (e.g., template name `/resume.html` in resume.html file)

When adding routes:
- Wrap handlers with `securityHeaders()` middleware
- Log requests using slog with structured fields: `slog.Info("event", "key", value)`

When updating Bootstrap CDN:
- Update both `templates/partials.html` link tag and `server.go` CSP header

## Deployment Notes

- Fly.io auto-scales based on traffic (min 0 machines)
- HTTPS enforced at edge
- Health check endpoint at `/health` returns "OK"
- No graceful shutdown implemented (not needed for this use case)
