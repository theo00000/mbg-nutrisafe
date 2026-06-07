# NutriSafe Backend

REST API for the NutriSafe platform — connecting schools, SPPG partners, and admins (BGN) around the MBG (Makan Bergizi Gratis) program. Built with Go (Fiber) + PostgreSQL.

## Stack

- **Go** 1.25 (see `back-end-go/go.mod`)
- **Fiber** v2 — HTTP framework
- **GORM** + **PostgreSQL** (pgx driver) — ORM and database
- **JWT** (HS256) — auth
- **SMTP** — transactional emails (credentials, password reset)

## Prerequisites

- Go ≥ 1.25
- PostgreSQL ≥ 14 (15-alpine used in Docker)
- SMTP credentials (for credential emails and password reset)
- Git

## Quick Start

```bash
git clone <repo-url>
cd mbg-nutrisafe/back-end-go

# 1. Install dependencies
go mod tidy

# 2. Create .env (see Environment Variables below)
cp .env.example .env   # or create manually

# 3. Create database
createdb nutrisafe      # GORM AutoMigrate handles tables on first run

# 4. Run
go run main.go
```

Server listens on `:3000` (override with `PORT`). Verify with:

```bash
curl http://localhost:3000/ping
# {"status":"success","message":"API Nutrisafe successfully running!"}
```

## Environment Variables

Create `back-end-go/.env`:

| Variable | Required | Default | Description |
|---|---|---|---|
| `DATABASE_URL` | optional | — | Full Postgres DSN. If set, overrides individual `DB_*` vars |
| `DB_HOST` | yes* | — | PostgreSQL host |
| `DB_USER` | yes* | — | PostgreSQL user |
| `DB_PASSWORD` | yes* | — | PostgreSQL password |
| `DB_NAME` | yes* | — | Database name |
| `DB_PORT` | yes* | — | PostgreSQL port |
| `DB_SSLMODE` | no | `disable` | `disable` \| `require` \| `verify-full` |
| `JWT_SECRET` | **yes** | — | HMAC secret for JWT signing |
| `PORT` | no | `3000` | HTTP listen port |
| `ALLOWED_ORIGINS` | no | `localhost:{4173,4174,5173,5174}` | CORS allowlist (comma-separated) |
| `SMTP_HOST` | yes | — | SMTP server hostname |
| `SMTP_PORT` | yes | — | SMTP port (587 for STARTTLS) |
| `SMTP_USER` | yes | — | SMTP username |
| `SMTP_PASSWORD` | yes | — | SMTP password |
| `SMTP_FROM` | yes | — | Sender address shown to recipients |

`*` Required only when `DATABASE_URL` is not provided.

## Project Layout

```text
back-end-go/
├── app/
│   ├── controllers/        # HTTP handlers, grouped by role
│   │   ├── admin/          # BGN admin endpoints
│   │   ├── auth/           # register, login, profile, password reset
│   │   ├── school/         # student, teacher, class, delivery, reports
│   │   ├── sppg/           # SPPG-side operations
│   │   ├── umum/           # public user endpoints
│   │   ├── dashboard/      # public stats
│   │   └── allergy/        # allergy reference data
│   ├── cron/               # background jobs (student account generator)
│   ├── middleware/         # JWT, role guards
│   ├── models/             # GORM models (single source of truth for schema)
│   ├── routes/             # Route wiring per role
│   └── services/           # cross-cutting domain logic
├── cmd/seed/               # seed CLI (admin user, reference data)
├── config/database.go      # DB connection + AutoMigrate + role seeding
├── utils/                  # mailer, slug/password helpers
├── uploads/                # runtime file uploads (gitignored)
├── Dockerfile
├── docker-compose.yml
└── main.go
```

## API Overview

All routes are flat (no `/api` prefix). Auth-protected routes expect `Authorization: Bearer <jwt>`.

**Public**
- `GET  /ping` — healthcheck
- `GET  /dashboard/stats` — public counters
- `POST /register` — register umum/admin
- `POST /register/school` — register school (login credentials emailed)
- `POST /register/sppg` — register SPPG (requires approval; multipart with proposal PDF + kitchen photo)
- `POST /login`, `/logout`
- `POST /forgot-password/request`, `/forgot-password/verify`

**Authenticated** (role-gated via middleware)
- `/profile/detail`, `/change-password`
- `/admin/*` — registrations, schools, SPPG, reports, dashboard
- `/school/*` — students, teachers, classes, delivery receipts, reports
- `/sppg/*` — schools, scan, reports, food problems, status updates
- `/umum/*` — reports

Roles seeded on first boot: `admin`, `school`, `sppg`, `umum`, `siswa`.

## Background Jobs

A cron loop runs every 5 minutes (`app/cron/student_account.go`):

- **`syncAllStudentCounts`** — reconciles `school_profiles.student_count` against actual `students` rows
- **`processStudentAccounts`** — generates login accounts for students added > 15 minutes ago, then emails the credentials list to the school's contact email

## Database

`GORM AutoMigrate` runs on every boot and applies any model changes additively (no destructive migrations). If you add a model, register it in `config/database.go` inside the `AutoMigrate(...)` call.

For a fresh admin user, run the seed CLI:

```bash
go run ./cmd/seed
```

## Docker

```bash
cd back-end-go
docker compose up --build
```

`docker-compose.yml` provisions Postgres, the backend, and (optionally) an ML service from a sibling `../NutritionModel/` directory. Backend reads its config from `.env` in the same directory.

## File Uploads

Uploaded files land under `back-end-go/uploads/`:

- `uploads/sppg/` — SPPG registration proposals + kitchen photos
- `uploads/delivery/` — SPPG delivery photos
- `uploads/delivery-receipt/` — school delivery receipts
- `uploads/reports/` — school + umum food reports
- `uploads/sppg-problems/` — SPPG food-problem evidence

There is no static file server wired up — paths are intentionally not exposed in API responses. Serve them through a dedicated authenticated endpoint when needed.

## Contributing

- New models → register in `config/database.go` `AutoMigrate(...)`
- New routes → wire under the appropriate role file in `app/routes/`
- Secrets (`.env`, credentials, uploaded files) must stay out of git — see `.gitignore`
