# AGENTS.md — Agentic Coding Guidelines

## Primary Role

This repository is maintained by **Prae**, Tech Lead. All coding tasks follow these rules.

---

## Task Handling

1. Read and understand the task fully
2. If unclear → ask exactly **1 question**, the most important one only
3. If clear → start immediately, no need to ask for permission again
4. Format: briefly state what I will do → deliver code → summarize what was done

---

## Build / Lint / Test Commands

Since this is a static HTML project with no package.json or build system:

- **No standard build commands exist** - HTML files are served directly
- **No lint commands configured** - Consider adding a linter if the project grows
- **No test framework** - Consider adding tests if JavaScript logic is added

If the project grows beyond a single HTML file, consider adding:

```bash
# Example for a future Node.js project
npm install           # Install dependencies
npm run dev          # Start development server
npm run build        # Build for production
npm run lint         # Run linter
npm test             # Run tests
npm test -- --testNamePattern="specific test"  # Run single test
```

---

## Tech Stack

### Backend — Go
- **Language:** Go
- **Framework:** Fiber v2 or Gin
- **Config:** Viper (`github.com/spf13/viper`) — all config via `.env` + `config.yaml`
- **Database:** PostgreSQL + pgvector (`github.com/pgvector/pgvector-go`)
- **ORM/Query:** sqlc or pgx directly (`github.com/jackc/pgx/v5`)
- **Real-time:** WebSocket (`github.com/gorilla/websocket`)
- **Auth:** JWT (`github.com/golang-jwt/jwt/v5`)
- **Architecture:** Clean — handler / service / repository separated
- **Project layout:**
  ```
  /cmd/server/main.go
  /internal/handler/
  /internal/service/
  /internal/repository/
  /internal/model/
  /config/
  ```

### Frontend — Vue 3
- **Framework:** Vue 3 + Vite
- **State:** Pinia
- **Routing:** Vue Router 4
- **HTTP:** axios or native fetch
- **Real-time:** native WebSocket or `@vueuse/core useWebSocket`
- **Styling:** Tailwind CSS or CSS Variables per design spec
- **Components:** Composition API + `<script setup>` always, no Options API

### DevOps
- **Container:** Docker + Docker Compose
- **Config:** `.env` loaded by Viper on backend, `.env` for Vite on frontend

---

## Coding Rules

### Go
- Always run `gofmt`
- Handle every error — never ignore with `_` on important calls
- Use structs for request/response, not `map[string]interface{}`
- Keep business logic out of HTTP handlers
- Load all config through Viper, never hardcode values
- Name: camelCase (unexported), PascalCase (exported)

### Vue 3
- `<script setup>` on every component, no Options API
- Move all API calls to `composables/` or `services/`
- No logic in templates — use `computed` or methods
- Component filenames: PascalCase

### General
- Readable code over clever code
- Error handling on every I/O and network call
- Never hardcode credentials — always use `.env`
- Sanitize all external input

---

## API Standard

```json
// Success
{ "data": {...}, "error": null }

// Error
{ "data": null, "error": "message here" }
```

WebSocket event format:
```json
{ "event": "agent_update", "payload": {...} }
```

---

## Docker Standard

Every project must have:
- `Dockerfile` for backend (Go multi-stage build)
- `Dockerfile` for frontend (node build + nginx)
- `docker-compose.yml` to run the full stack with one command
- `.env.example` as template

---

## Code Style Guidelines

### HTML Structure
- Use semantic HTML5 elements (`<header>`, `<main>`, `<section>`, `<nav>`)
- Keep attributes in consistent order (class, id, data-*, src, alt)
- Use lowercase for tags and attributes
- Quote attribute values

### CSS Guidelines
- Use CSS custom properties (variables) for theming
- Follow existing color naming convention (`--accent`, `--green`, `--red`)
- Use flexbox and grid for layout
- Group related styles together

### JavaScript Guidelines
- Use ES6+ syntax (const/let, arrow functions, template literals)
- Declare variables with `const` by default, `let` when reassignment needed
- Wrap async operations in try/catch
- Handle errors on every I/O and network call

### Naming Conventions

| Type | Convention | Example |
|------|------------|---------|
| Variables | camelCase | `userName`, `isActive` |
| Functions | camelCase | `fetchData()`, `calculateTotal()` |
| Constants | UPPER_SNAKE_CASE | `MAX_RETRY_COUNT`, `API_BASE_URL` |
| Classes | PascalCase | `DashboardPanel`, `UserProfile` |
| CSS classes | kebab-case | `.btn-primary`, `.card-header` |
| File names | kebab-case | `dashboard-v3.html`, `utils.js` |

### Error Handling
- Always wrap async operations in try/catch
- Provide meaningful error messages
- Log errors appropriately (console.error for errors, console.warn for warnings)
- Never expose sensitive information in error messages
- Never ignore errors with `_` on important calls

---

## File Organization

```
project-root/
├── cmd/                    # Entry points
├── internal/               # Private application code
│   ├── handler/           # HTTP handlers
│   ├── service/           # Business logic
│   ├── repository/        # Database queries
│   └── model/             # Data models
├── config/                # Configuration files
├── web/                   # Frontend (Vue 3)
├── Dockerfile
├── docker-compose.yml
├── .env.example
└── AGENTS.md
```

---

## Scope

✅ Accept: Go backend, Vue 3 frontend, REST API, WebSocket, PostgreSQL, pgvector, Docker, architecture, bug fixes, code review

❌ Decline: non-technical tasks (marketing, graphic design unrelated to UI)

If a task is out of scope → say so directly and suggest the appropriate agent.

---

## Recommendations for Future Development

1. **Add a package.json** - For dependency management and scripts
2. **Add a linter** - ESLint for JavaScript, HTMLHint for HTML
3. **Add a formatter** - Prettier for consistent code formatting
4. **Add tests** - Vitest or Jest for JavaScript testing
5. **Consider TypeScript** - For better type safety as the project grows
