# Baby Care Tracker - Implementation Plan

## Context

Build a baby care tracking web app from scratch (empty project directory). The app helps parents (primarily women aged 20-40) quickly log and review their baby's sleep, feeding, diaper changes, and growth. The top priority is **minimal taps to record events** with a **clean, warm UI**.

**Key decisions**: Single child per instance (no multi-child switcher). No authentication (personal/family use). Vanilla TypeScript (no framework).

## Tech Stack & Deployment

- **Frontend**: TypeScript + HTML (vanilla, no framework) bundled with esbuild
- **Backend**: Go (standard library `net/http` router — Go 1.22+ supports method+path patterns natively)
- **Database**: SQLite via `github.com/mattn/go-sqlite3`
- **Deployment**: Single Go binary with embedded static files (`embed.FS`) + SQLite file. Deploy to any $5/month VPS (DigitalOcean/Linode) — just copy one binary and run it.

## Project Structure

```
vibe_baby_care/
├── go.mod / go.sum
├── main.go                     # Entry point: flags, DB init, server start
├── embed.go                    # //go:embed static/*
├── Makefile                    # Build pipeline
├── internal/
│   ├── server/server.go        # Mux setup, API routes, static file serving + SPA fallback
│   ├── handler/
│   │   ├── handler.go          # JSON response helpers, error helpers
│   │   ├── children.go         # /api/v1/child CRUD
│   │   ├── sleep.go            # /api/v1/sleep
│   │   ├── feeding.go          # /api/v1/feeding
│   │   ├── diaper.go           # /api/v1/diaper
│   │   └── growth.go           # /api/v1/growth
│   ├── model/                  # Go structs matching DB tables
│   │   ├── child.go, sleep.go, feeding.go, diaper.go, growth.go
│   ├── store/
│   │   ├── store.go            # SQLite connection, PRAGMA setup, CREATE TABLE migrations
│   │   ├── children.go, sleep.go, feeding.go, diaper.go, growth.go
│   └── middleware/middleware.go # CORS, logging, JSON content-type
├── frontend/
│   ├── index.html              # SPA shell
│   ├── package.json            # devDeps: esbuild, typescript
│   ├── tsconfig.json
│   ├── src/
│   │   ├── main.ts             # Entry: router init, first-visit check
│   │   ├── router.ts           # Hash-based SPA router
│   │   ├── api.ts              # Typed fetch wrapper for all endpoints
│   │   ├── state.ts            # Minimal pub/sub state (currentChild, activeTimers)
│   │   ├── components/
│   │   │   ├── onboarding.ts   # First-visit child profile form
│   │   │   ├── dashboard.ts    # Today's summary cards + recent activity
│   │   │   ├── quick-add.ts    # FAB + radial menu
│   │   │   ├── sleep-modal.ts  # Start/stop sleep timer or manual entry
│   │   │   ├── feeding-modal.ts# Breast (L/R timer) or Bottle (quantity)
│   │   │   ├── diaper-modal.ts # One-tap: Wet / Dirty / Mixed = instant save
│   │   │   ├── growth-modal.ts # Weight + length input form
│   │   │   ├── history.ts      # Filtered timeline by day
│   │   │   ├── growth-chart.ts # Canvas-rendered line chart
│   │   │   └── nav.ts          # Bottom tab bar
│   │   ├── utils/
│   │   │   ├── dom.ts          # h() element creator, $, $$
│   │   │   └── date.ts         # Format, "time ago", duration
│   │   └── types/models.ts     # TS interfaces matching Go models
│   └── css/
│       ├── main.css            # CSS custom properties (design tokens), reset
│       ├── components.css, onboarding.css, dashboard.css, history.css, growth-chart.css, modal.css
└── static/                     # Build output (gitignored), embedded into Go binary
```

## Database Schema (SQLite)

**children**: `id TEXT PK`, `name TEXT`, `date_of_birth TEXT`, `gender TEXT (male|female|other)`, `photo_url TEXT`, `notes TEXT`, `created_at TEXT`, `updated_at TEXT`

**sleep_logs**: `id TEXT PK`, `child_id TEXT FK`, `start_time TEXT`, `end_time TEXT NULL` (NULL = currently sleeping), `duration_minutes INT`, `notes TEXT`, `created_at TEXT`

**feeding_logs**: `id TEXT PK`, `child_id TEXT FK`, `feed_type TEXT (breast_left|breast_right|bottle)`, `start_time TEXT`, `end_time TEXT NULL`, `duration_minutes INT`, `quantity_ml INT` (bottle only), `notes TEXT`, `created_at TEXT`

**diaper_logs**: `id TEXT PK`, `child_id TEXT FK`, `diaper_type TEXT (wet|dirty|mixed)`, `changed_at TEXT`, `notes TEXT`, `created_at TEXT`

**growth_logs**: `id TEXT PK`, `child_id TEXT FK`, `measured_on TEXT`, `weight_grams INT`, `length_mm INT`, `head_circumference_mm INT`, `notes TEXT`, `created_at TEXT`

All IDs are UUID v4 generated in Go. Indexes on `(child_id, start_time/changed_at/measured_on)` for each log table. Migrations run as idempotent `CREATE TABLE IF NOT EXISTS` on startup.

## API Endpoints

Base: `/api/v1`. All JSON. Go 1.22+ `net/http` pattern routing (no third-party router).

| Area | Endpoints |
|------|-----------|
| Child | `POST /child` (create, first visit), `GET /child` (get the child), `PUT /child` (update profile) |
| Sleep | `POST/GET /sleep`, `GET /sleep/active`, `PUT/DELETE /sleep/{logId}` |
| Feeding | `POST/GET /feeding`, `GET /feeding/active`, `PUT/DELETE /feeding/{logId}` |
| Diaper | `POST/GET /diaper`, `PUT/DELETE /diaper/{logId}` |
| Growth | `POST/GET /growth`, `PUT/DELETE /growth/{logId}` |
| Summary | `GET /summary?date=YYYY-MM-DD` — aggregated day stats |

## UI/UX Design

### Screens
1. **Onboarding** (`#/onboarding`) — shown on first visit. Single form: baby name, DOB (date picker), gender (pill toggles), optional photo. Large "Get Started" button.
2. **Dashboard** (`#/dashboard`) — default screen. Header with baby name + age. 2x2 summary cards (sleep/feeding/diaper/growth) in pastel category colors. Recent activity list. Active timer banner if sleep/feeding in progress.
3. **History** (`#/history`) — date navigator, filter pills (All/Sleep/Feeding/Diaper), color-coded vertical timeline.
4. **Growth Chart** (`#/growth`) — toggle weight/length, Canvas 2D line chart, data table below.

### Quick-Add (Floating Action Button)
Persistent FAB (bottom-right) opens radial menu with 4 options: Sleep, Feed, Diaper, Growth.

**Taps to record:**
| Event | Taps | Flow |
|-------|------|------|
| Diaper | **2** | FAB → tap Wet/Dirty/Mixed (instant save with current time) |
| Sleep start | **2** | FAB → "Start Sleep" (auto-timestamps) |
| Sleep stop | **1** | Tap "Stop" on active timer banner |
| Breast feed start | **2** | FAB → tap Left/Right (starts timer) |
| Bottle feed | **3** | FAB → Bottle → Save (quantity pre-filled from last use) |
| Growth | **5** | FAB → fill weight/length → Save |

### Design Tokens (CSS Custom Properties)
- **Primary**: `#E8507A` (rose) — action buttons, active states
- **Category colors**: Sleep `#8B5CF6` (purple), Feeding `#EC4899` (pink), Diaper `#10B981` (mint), Growth `#F59E0B` (amber)
- **Background**: `#FEFCFB` (warm off-white)
- **Font**: Inter (Google Fonts), 16px base, never smaller on mobile
- **Touch targets**: min 44x44px, buttons 48px tall
- **Layout**: Mobile-first, max-width 480px centered on desktop (phone-like container)

## Build & Deploy Pipeline

```bash
# Development (two terminals)
make frontend-watch    # esbuild watch mode → static/app.js
make run               # go run . --port 8080

# Production build
make all               # frontend build → Go embed → single binary "baby-care"

# Deploy to VPS
scp baby-care user@server:/home/app/
ssh user@server "systemctl restart baby-care"
```

SQLite DB stored at `~/.baby-care/data.db` by default (overridable via `--db` flag). WAL mode + foreign keys enabled on connection.

## Implementation Order

| Step | What | Testable result |
|------|------|----------------|
| 1 | Go scaffolding: `go.mod`, `main.go`, health-check endpoint | `curl localhost:8080/health` returns OK |
| 2 | SQLite store: connection, migrations, child CRUD | Unit tests for store layer |
| 3 | Child API handlers + middleware | `curl` POST/GET child |
| 4 | Frontend shell: `index.html`, CSS tokens, router, onboarding screen | Browser shows onboarding form |
| 5 | Connect onboarding → API, dashboard skeleton + summary endpoint | Create child → see dashboard |
| 6 | Sleep tracking (store → handler → modal → timer banner) | Record and view sleep |
| 7 | Feeding tracking (breast timer + bottle form) | Record both feed types |
| 8 | Diaper tracking (one-tap modal) | 2-tap diaper logging works |
| 9 | Growth tracking (form + Canvas chart) | Record and visualize growth |
| 10 | History view with filters | Browse past events by day |
| 11 | Polish: empty states, error toasts, loading states, active timer | Smooth UX end-to-end |
| 12 | Build pipeline: Makefile, embed verification, single binary test | `./baby-care` serves everything |

## Verification

1. `make all` produces a single `baby-care` binary
2. Run `./baby-care` — opens browser to `http://localhost:8080`
3. Complete onboarding (create child profile)
4. Log one of each event type via the FAB quick-add
5. Verify dashboard summary cards update
6. Check history timeline shows all events
7. Add 3+ growth measurements, verify chart renders
8. Kill and restart the binary — all data persists (SQLite file)
9. Verify diaper logging takes exactly 2 taps
