# CLAUDE.md

Guidance for AI assistants (and humans) working in this repository.

## What this is

Backend API for **RapsShop**, a Growtopia **Diamond Lock (DL)** / **Blue Gem Lock (BGL)** trading shop (1 BGL = 100 DL). Customers buy DL/BGL (Midtrans or manual transfer) and sell DL to the shop; a single admin role manages inventory, pricing, payment methods, the in-game delivery world, content, and reporting.

See `PRD.md` for the product spec, `README.md` for setup, and `SUGGESTIONS.md` for the engineering backlog + status.

## Stack

Go 1.23 · Gin v1.9 · GORM v1.24 (MySQL) · JWT (`golang-jwt`) · bcrypt · Midtrans (coreapi SDK + raw HTTP in `lib/`) · Supabase Storage · `godotenv`.

## Architecture

Layered, one package per feature under `src/`:

```
handler (HTTP)  →  service (usecase/business logic)  →  repo (GORM)  →  MySQL
```

- `entities/` — GORM table models.
- `model/` — DTOs (`Input*`, `*Dto`) **and** repo/usecase interfaces. Interfaces live here by convention.
- `middleware/` — JWT auth + token helpers.
- `lib/` — Midtrans driver (the single place that talks to Midtrans).
- `database/mysql/` — connection + `AutoMigrate`.
- `utils/` — standardized JSON response envelope.
- `src/<feature>/{handlers,service,repo}/` — feature modules.

Wiring happens in `main.go`: `repo := New...Repository(db)` → `usecase := New...Usecase(repo)` → `handler.New...Handler(apiGroup, usecase, jwtMiddleware)`. Routes are under `/api/v1`.

Features: `admin` (auth), `stock_dl`, `harga_dl` (price), `env_growtopia` (delivery world creds), `sosmed`, `testimoni`, `penjualan_dl` (customer→shop sales), `pembelian_dl` (customer←shop purchases, Midtrans), `metode_pembayaran` (payment methods).

## Build / run / test

```bash
go build ./...        # compile everything
go vet ./...          # static checks
go test ./...         # run unit tests
go run main.go        # run locally (needs .env)
go build -o main      # production binary (what CI builds)
```

Go 1.23 is required (`go.mod`). If the local toolchain is older, keep `GOTOOLCHAIN=auto` so the 1.23 toolchain is fetched, or install Go 1.23.

## Configuration

All config is env vars (`.env`, loaded by `godotenv`). Copy `.env.example` and fill it in. Key vars: `DB_*`, `JWT_KEY`, `TOKEN` (admin-register secret), `SUPABASE_URL`/`SERVICE_TOKEN`/`STORAGE_NAME`/`BASE_URL`, `AUTHORIZATION_VALUE` (Midtrans server key), `MIDTRANS` (charge URL), `MIDTRANS_ENV` (`sandbox`/`production`), `HOST_URL`, `PORT`, `ALLOWED_ORIGINS`, `AUTO_MIGRATE`. `.env` is gitignored — never commit secrets.

## Conventions

- File names follow Go idiom: the directory/package already conveys the layer, so files are **not** prefixed with it. Single-file layer packages take the package name (`handlers/handlers.go`, `service/service.go`, `repo/repo.go`, `lib/midtrans.go`); the multi-file `model/` and `entities/` packages name files by feature (`model/pembelian_dl.go`, `entities/pembelian_dl.go`).
- Domain language is Indonesian (`penjualan` = customer sells to shop; `pembelian` = customer buys from shop; `harga` = price; `stock`/`stok`; `metode_pembayaran` = payment method). Keep new code consistent within a layer.
- Constructors are named `New<Feature><Layer>` (e.g. `NewStockDLHandler`, `NewPenjualanDLUsecase`). Don't copy-paste a constructor name from another feature.
- Handlers: bind with `c.ShouldBindJSON` (not `BindJSON`), then return early via `utils.FailureOrErrorResponse`. Keep HTTP/transport concerns here only.
- Services hold business logic; repos hold all GORM/DB access. **Never** open a DB connection or call `mysql.InitDatabase()` outside `main` / `database/mysql`.
- All responses use `utils.SuccessResponse` / `utils.FailureOrErrorResponse` (envelope: `status_code`, `status`, `message`, `data`). Pass the real HTTP status; the helper honors it.
- Money/quantity are integers. Pricing: 1–99 DL priced per DL; multiples of 100 priced per BGL; mixed splits into BGL + DL (see `model.MidtransData.IniDataPembelian`).
- `PenjualanDL.status`: `0` pending, `1` approved, `-1` rejected. Stock increments when a sale is approved and decrements when a purchase is paid/shipped.

## Gotchas / domain rules

- **Midtrans webhook** (`POST /pembelian/status`) must verify `signature_key = sha512(order_id + status_code + gross_amount + AUTHORIZATION_VALUE)` before trusting it.
- **Stock mutations** must stay atomic and non-negative — go through the repo's transactional adjust, not read-modify-write in app code.
- **Credentials** (`metode_pembayaran.kredensial_pembayaran`, `growtopia.password`) are admin-only — never expose them on public routes.
- **Image uploads**: validate extension via `filepath.Ext` against an allowlist and name files with `uuid`; check `image.Open()` before uploading.
- **Reporting SQL**: always parenthesize `OR` groups in `WHERE` (`... AND (a OR b)`), and use `.Order()` rather than embedding `ORDER BY` in a `Where` string.
- `AutoMigrate` runs at startup when `AUTO_MIGRATE != "false"`. For prod, prefer disabling it and using versioned migrations.

## Before you finish

Run `go build ./... && go vet ./... && go test ./...`. Add tests for business logic you touch (pricing, JWT, recap math). Don't introduce new dependencies without need; pin versions.
