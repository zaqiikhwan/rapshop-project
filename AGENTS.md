# AGENTS.md

Operational guidance for AI agents working in this repository.

## Project Snapshot

RapsShop is a Go backend API for a Growtopia Diamond Lock (DL) / Blue Gem Lock (BGL) shop. Customers buy DL/BGL from the shop through Midtrans or manual transfer, and sell DL back to the shop. Admin APIs manage stock, pricing, payment methods, delivery-world credentials, content, and reports.

Read these before larger changes:

- `README.md` for setup, route reference, deployment, and current status.
- `PRD.md` for product behavior.
- `SUGGESTIONS.md` for backlog and historical risk notes.
- `CLAUDE.md` for an existing assistant-oriented overview.

## Stack

- Go 1.23
- Gin
- GORM with MySQL
- JWT auth with `github.com/golang-jwt/jwt`
- bcrypt via `golang.org/x/crypto`
- Midtrans integration in `lib/`
- Supabase Storage for images
- `.env` loading via `godotenv`

## Repository Layout

- `main.go`: loads env, opens DB, constructs dependencies, registers routes, starts HTTP server.
- `database/mysql/`: DB connection and optional `AutoMigrate`.
- `entities/`: GORM table models.
- `model/`: DTOs plus repository/usecase interfaces.
- `middleware/`: JWT middleware and token helpers.
- `lib/`: external service drivers, especially Midtrans.
- `utils/`: JSON response helpers.
- `src/<feature>/handlers`: HTTP layer.
- `src/<feature>/service`: business logic/usecase layer.
- `src/<feature>/repo`: GORM data access layer.

The architecture is:

```text
handler -> service/usecase -> repo -> MySQL
```

Keep new code in the same layer pattern unless there is a strong reason to change it.

## Common Commands

Run these from the repo root:

```bash
go build ./...
go vet ./...
go test ./...
```

Local server:

```bash
go run main.go
```

Production-style binary:

```bash
go build -o main
```

Linting is configured in `.golangci.yml`; CI runs `golangci-lint` only on new issues.

## Configuration

Runtime config comes from environment variables loaded from `.env`. Use `.env.example` as the template. Never commit real secrets.

Important variables include:

- `DB_USER`, `DB_PASS`, `DB_HOST`, `DB_NAME`
- `JWT_KEY`
- `TOKEN`
- `SUPABASE_URL`, `SERVICE_TOKEN`, `STORAGE_NAME`, `BASE_URL`
- `AUTHORIZATION_VALUE`, `MIDTRANS`, `MIDTRANS_ENV`
- `HOST_URL`
- `PORT`
- `ALLOWED_ORIGINS`
- `AUTO_MIGRATE`

`AutoMigrate` runs unless `AUTO_MIGRATE=false`.

## Coding Conventions

- Use existing Indonesian domain naming consistently: `pembelian` means customer buys from shop, `penjualan` means customer sells to shop, `harga` is price, `metode_pembayaran` is payment method.
- Constructors follow `New<Feature><Layer>`, for example `NewStockDLHandler`, `NewStockDLUsecase`, `NewStockDLRepository`.
- Keep repository and usecase interfaces in `model/`, matching existing patterns.
- Handlers should handle HTTP binding and responses only.
- Services should contain business rules.
- Repos should contain all GORM queries and transactions.
- Do not call `mysql.InitDatabase()` outside `main.go` or `database/mysql/`.
- Use `c.ShouldBindJSON` in handlers and return early on bind errors.
- Use `utils.SuccessResponse` and `utils.FailureOrErrorResponse` for API responses.
- Keep money and DL/BGL quantities as integers unless a requirement explicitly changes that.
- Run `gofmt` on changed Go files.

## Domain Rules That Must Stay Intact

- Pricing rule: 1-99 DL uses per-DL pricing; exact multiples of 100 use BGL pricing; mixed quantities split into BGL plus DL remainder. See `model/pembelian_dl_test.go`.
- `PenjualanDL.status`: `0` pending, `1` approved, `-1` rejected.
- Stock increases when a customer sale to the shop is approved.
- Stock decreases when a customer purchase from the shop is paid/shipped.
- Stock changes must remain atomic and non-negative; use the repo transaction helpers instead of read-modify-write logic in services.
- Midtrans webhook status updates must verify the signature before trusting the payload.
- Payment credentials and Growtopia world passwords are admin-only data. Do not expose them through public routes.
- Upload handling should validate file extension, inspect image content, and generate safe names instead of trusting user-provided filenames.
- Reporting queries with mixed `AND`/`OR` conditions need explicit parentheses.

## Testing Guidance

Before finishing a code change, run:

```bash
go build ./...
go vet ./...
go test ./...
```

Add or update tests when touching:

- Pricing calculations.
- JWT/token behavior.
- Stock mutation logic.
- Midtrans webhook behavior.
- Recap/profit/reporting calculations.
- Route authorization changes.

Prefer focused unit tests for service/domain logic. Use integration-style tests only when the behavior genuinely depends on GORM/MySQL interaction.

## Deployment Notes

CI runs on pull requests and pushes to `main`:

- `go build ./...`
- `go vet ./...`
- `go test ./...`
- `golangci-lint` for new issues

Deployment is handled by `.github/workflows/deploy.yml` on `main` only when the commit message contains `DEPLOY`, or by manual dispatch. It builds `main`, rsyncs it to the server, then restarts `pm2`.

## Agent Workflow

1. Inspect existing code in the target feature before editing.
2. Keep edits scoped to the requested behavior.
3. Preserve public API response envelopes unless intentionally changing an API contract.
4. Avoid introducing dependencies without a clear need.
5. Do not modify `.env` or commit secrets.
6. Be careful with the current dirty worktree; do not revert unrelated user changes.
7. After edits, run the relevant verification commands and report any command that could not be run.
