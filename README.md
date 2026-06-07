# RapsShop Backend

Backend API for **RapsShop**, an online shop that trades Growtopia **Diamond Locks (DL)** and **Blue Gem Locks (BGL)** (1 BGL = 100 DL). It supports customers **buying** DL/BGL (with Midtrans online payment or manual transfer) and **selling** DL to the shop, plus an admin panel API for inventory, pricing, payment methods, reporting, and storefront content.

Built with **Go + Gin + GORM (MySQL)**, JWT auth, Midtrans payments, and Supabase Storage for images.

> See [`PRD.md`](./PRD.md) for the product spec, [`SUGGESTIONS.md`](./SUGGESTIONS.md) for the engineering backlog, and [`PRODUCTION_IMPROVEMENTS.md`](./PRODUCTION_IMPROVEMENTS.md) for VPS deployment and transaction reliability hardening.

---

## Tech stack

| Concern | Choice |
|---|---|
| Language | Go 1.23 |
| HTTP framework | Gin v1.9 |
| ORM / DB | GORM v1.24 / MySQL |
| Auth | JWT (`golang-jwt`), bcrypt password hashing |
| Payments | Midtrans (Core API + raw HTTP) |
| Object storage | Supabase Storage |
| Config | `.env` via `godotenv` |
| Deploy | GitHub Actions → rsync → pm2 |

---

## Architecture

Layered (clean-ish) architecture, one package per feature under `src/`:

```
handler (HTTP)  →  service (business logic / usecase)  →  repo (GORM data access)  →  MySQL
```

```
rapsshop-project/
├── main.go                 # entrypoint: load env, init DB, wire deps, register routes
├── entities/               # GORM models (DB tables)
├── model/                  # DTOs, request/response types, repo & usecase interfaces
├── middleware/             # JWT auth middleware + token helpers
├── lib/                    # Midtrans Core API driver
├── database/mysql/         # DB connection + AutoMigrate
├── utils/                  # standardized JSON response helpers
├── public/payment/         # locally stored purchase-proof uploads (served statically)
└── src/
    ├── admin/              # auth (register/login/profile)
    ├── stock_dl/           # inventory
    ├── harga_dl/           # pricing
    ├── env_growtopia/      # in-game delivery world credentials
    ├── sosmed/             # social media links
    ├── testimoni/          # testimonials (image upload)
    ├── penjualan_dl/       # customer → shop sales
    ├── pembelian_dl/       # customer ← shop purchases (Midtrans)
    └── metode_pembayaran/  # payment methods
        ├── handlers/       # HTTP layer
        ├── service/        # business logic
        └── repo/           # data access
```

Each feature is wired in `main.go`: `repo := New...Repository(db)` → `usecase := New...Usecase(repo)` → `handler.New...Handler(apiGroup, usecase, jwtMiddleware)`.

---

## Getting started

### Prerequisites
- Go 1.23+
- MySQL 5.7+/8.0 (a database created and reachable)
- A Midtrans account (server key) for payments
- A Supabase project + storage bucket for images

### 1. Clone & install dependencies
```bash
git clone <repo-url>
cd rapsshop-project
go mod download
```

### 2. Configure environment
Copy the example file and fill in values:
```bash
cp .env.example .env
```

| Variable | Required | Description |
|---|---|---|
| `DB_USER` | ✅ | MySQL username |
| `DB_PASS` | ✅ | MySQL password |
| `DB_HOST` | ✅ | MySQL host, e.g. `127.0.0.1:3306` |
| `DB_NAME` | ✅ | MySQL database name |
| `JWT_KEY` | ✅ | Secret used to sign/verify admin JWTs |
| `TOKEN` | ✅ | Shared secret required to register an admin |
| `SUPABASE_URL` | ✅ | Supabase Storage API URL |
| `SERVICE_TOKEN` | ✅ | Supabase service token |
| `STORAGE_NAME` | ✅ | Supabase bucket name |
| `BASE_URL` | ✅ | Public base URL prefix for Supabase-hosted images |
| `AUTHORIZATION_VALUE` | ✅ | Midtrans server key (used as Basic auth, base64-encoded at runtime) |
| `MIDTRANS` | ✅ | Midtrans charge endpoint URL (e.g. `https://api.midtrans.com/v2/charge`) |
| `MIDTRANS_STATUS_URL` | optional | Midtrans status base URL; defaults to `https://api.midtrans.com/v2` |
| `HOST_URL` | ✅ | Public base URL of this server, used to build links to uploaded proof images |
| `GIN_MODE` | optional | `debug` (default) or `release` |
| `AUTO_MIGRATE` | optional | Set to `false` in production; defaults to enabled for local development |

> `.env` is gitignored — never commit real secrets.

### 3. Run
```bash
go run main.go
```
The server starts on the port from `PORT` (default `:8080`) and shuts down gracefully on SIGINT/SIGTERM. `AutoMigrate` creates/updates tables on startup. Health check:
```bash
curl http://localhost:8080/ping   # -> {"message":"pong"}
```

### 4. Build
```bash
go build -o main
./main
```

---

## Configuration notes

- **Database tables** are auto-migrated from `entities/` on every startup.
- **Midtrans** environment is selected via `MIDTRANS_ENV` (`sandbox` default, `production` when set).
- **CORS** uses an explicit allowlist from `ALLOWED_ORIGINS` (comma-separated); credentials are only echoed for matched origins.

---

## API reference

Base path: **`/api/v1`**. Protected routes require a header `Authorization: Bearer <token>` obtained from `POST /api/v1/admin/login`. 🔓 public · 🔒 admin (JWT).

### Auth & admin
| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/admin/register` | 🔓\* | Register admin (body must include valid `token`) |
| POST | `/admin/login` | 🔓 | Login, returns JWT |
| GET | `/profile` | 🔒 | Current admin profile |

\* Gated by the shared `TOKEN` secret, not JWT.

### Stock
| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/stock` | 🔒 | Create stock snapshot |
| GET | `/stocks` | 🔓 | List all stock history |
| GET | `/stock` | 🔓 | Latest stock |
| PATCH | `/stock` | 🔒 | Add to latest stock |
| DELETE | `/stock` | 🔒 | Delete latest stock |

### Price
| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/price` | 🔒 | Create price |
| GET | `/price` | 🔓 | Latest price |
| PATCH | `/price` | 🔒 | Update latest price |
| DELETE | `/price` | 🔒 | Delete latest price |

### Growtopia delivery env
| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/env` | 🔒 | Set delivery world credentials |
| GET | `/env` | 🔒 | Get latest world credentials (now protected — exposes world password) |
| PATCH | `/env` | 🔒 | Update world credentials |

### Social media
| Method | Path | Auth |
|---|---|---|
| POST | `/platform` | 🔒 |
| GET | `/platforms` | 🔓 |
| GET | `/platform/:id` | 🔒 |
| PATCH | `/platform/:id` | 🔒 |
| DELETE | `/platform/:id` | 🔒 |

### Testimonials
| Method | Path | Auth | Notes |
|---|---|---|---|
| POST | `/testimoni` | 🔒 | multipart: `gambar`, `testimoni`, `username`, `title` |
| GET | `/testimonis` | 🔓 | |
| GET | `/testimoni/:id` | 🔒 | |
| PATCH | `/testimoni/:id` | 🔒 | multipart |
| DELETE | `/testimoni/:id` | 🔒 | |

### Payment methods
| Method | Path | Auth | Notes |
|---|---|---|---|
| POST | `/payment` | 🔒 | |
| GET | `/payments` | 🔒 | now protected (exposes payment credentials) |
| GET | `/payment/:id` | 🔒 | now protected (exposes payment credentials) |
| GET | `/checkout/options` | 🔓 | public checkout payment options; excludes stored credentials |
| PATCH | `/payment/:id` | 🔒 | |
| DELETE | `/payment/:id` | 🔒 | |

### Sales — customer sells to shop (Penjualan)
| Method | Path | Auth | Notes |
|---|---|---|---|
| POST | `/penjualan` | 🔓 | multipart: `image`, `jumlah_dl`, `whatsapp`, `transfer`, `nomor_transfer`, `nama` |
| GET | `/penjualans?_start=&_end=` | 🔒 | paginated list |
| GET | `/penjualan/:id` | 🔒 | |
| GET | `/rekapitulasi?_date=` | 🔒 | daily recap (sales + purchases by rate) |
| GET | `/profit?_date=` | 🔒 | per-day profit for the month |
| GET | `/penjualan/total?_date=` | 🔒 | per-day sale totals |
| PATCH | `/penjualan/:id` | 🔒 | update status (0 pending / 1 approved / -1 rejected) |
| DELETE | `/penjualan/:id` | 🔒 | |

### Purchases — customer buys from shop (Pembelian)
| Method | Path | Auth | Notes |
|---|---|---|---|
| GET | `/checkout/preview?jumlah_dl=` | 🔓 | stock availability, price breakdown, and checkout eligibility |
| POST | `/pembelian` | 🔓 | create order + Midtrans charge |
| POST | `/new/pembelian` | 🔓 | create order (manual payment) |
| POST | `/pembelian/status` | 🔓 | Midtrans webhook (order status notification) |
| GET | `/pembelians?_start=&_end=&queue=` | 🔒 | paginated list; optional queue filter |
| GET | `/pembelian/total?_date=` | 🔒 | per-day purchase totals |
| GET | `/pembelian/:id` | 🔓 | DB record |
| GET | `/pembelian/:id/tracking` | 🔓 | customer-facing tracking response |
| GET | `/pembelian/status/:id` | 🔓 | live Midtrans status |
| PATCH | `/pembelian/:id` | 🔒 | update delivery status |
| PATCH | `/pembelian/button/:id` | 🔓 | toggle "pay" button state |
| PATCH | `/pembelian/confirm/:id` | 🔒 | admin confirm manual payment |
| PATCH | `/upload/:id` | 🔓 | upload payment proof (multipart `file`) |
| GET | `/public/*` | 🔓 | static proof images |

### Response envelope
All responses use a common shape:
```json
{
  "status_code": 200,
  "status": "success, request OK!",
  "message": "human readable message",
  "data": {}
}
```

### Example: admin login
```bash
curl -X POST http://localhost:8080/api/v1/admin/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"secret"}'
```
```bash
# use the returned token
curl http://localhost:8080/api/v1/profile \
  -H "Authorization: Bearer <token>"
```

### Example: checkout preview
```bash
curl "http://localhost:8080/api/v1/checkout/preview?jumlah_dl=150"
```
```json
{
  "status_code": 200,
  "status": "success, request OK!",
  "message": "success create checkout preview",
  "data": {
    "jumlah_dl": 150,
    "stock_dl": 1000,
    "stock_enough": true,
    "can_checkout": true,
    "harga_beli_dl": 100,
    "harga_beli_bgl": 9500,
    "breakdown": {
      "bgl_quantity": 1,
      "dl_quantity": 50,
      "bgl_subtotal": 9500,
      "dl_subtotal": 5000
    },
    "total_payment": 14500
  }
}
```

### Example: checkout options
```bash
curl http://localhost:8080/api/v1/checkout/options
```
```json
{
  "status_code": 200,
  "status": "success, request OK!",
  "message": "success fetch checkout payment options",
  "data": {
    "gateway": [
      {
        "index_pembayaran": 1,
        "jenis_pembayaran": "QRIS",
        "checkout_type": "gateway",
        "provider": "midtrans",
        "requires_payment_proof": false,
        "enabled": true
      }
    ],
    "manual": [
      {
        "index_pembayaran": 7,
        "jenis_pembayaran": "BCA Manual",
        "checkout_type": "manual",
        "provider": "manual_transfer",
        "pemilik": "RapsShop",
        "requires_payment_proof": true,
        "enabled": true
      }
    ]
  }
}
```

### Example: gateway purchase
```bash
curl -X POST http://localhost:8080/api/v1/pembelian \
  -H "Content-Type: application/json" \
  -d '{
    "world": "BUYDL",
    "nama": "Customer",
    "grow_id": "GrowID",
    "jumlah_dl": 150,
    "wa": "628123456789",
    "metode_transfer": 1
  }'
```

### Example: manual purchase
```bash
curl -X POST http://localhost:8080/api/v1/new/pembelian \
  -H "Content-Type: application/json" \
  -d '{
    "world": "BUYDL",
    "nama": "Customer",
    "grow_id": "GrowID",
    "jumlah_dl": 150,
    "wa": "628123456789",
    "metode_transfer": 7
  }'
```
```json
{
  "status_code": 201,
  "status": "success, request OK!",
  "message": "transaction successfully created",
  "data": {
    "id_transaksi": "order-uuid",
    "payment": {
      "index_pembayaran": 7,
      "jenis_pembayaran": "BCA Manual",
      "checkout_type": "manual",
      "provider": "manual_transfer",
      "pemilik": "RapsShop",
      "requires_payment_proof": true
    },
    "upload_proof_path": "/api/v1/upload/order-uuid",
    "tracking_path": "/api/v1/pembelian/order-uuid/tracking"
  }
}
```

### Example: upload proof and track order
```bash
curl -X PATCH http://localhost:8080/api/v1/upload/order-uuid \
  -F "file=@payment-proof.jpg"

curl http://localhost:8080/api/v1/pembelian/order-uuid/tracking
```

Supported admin purchase queues:

| Queue | Meaning |
|---|---|
| `pending_payment` | Waiting for gateway payment or manual proof upload |
| `proof_uploaded` | Manual proof uploaded and waiting for admin confirmation |
| `waiting_delivery` | Paid and waiting for in-game delivery |
| `delivered` | Delivered orders |
| `failed` | Denied or failed payments |
| `review` | Gateway challenge/review status |

---

## Deployment

CI/CD via `.github/workflows/deploy.yml`:
1. Triggered on push to `main` **only when the commit message contains `DEPLOY`** (or via manual `workflow_dispatch`).
2. Builds the Go binary (`go build -o main`).
3. Rsyncs the binary and `ecosystem.config.cjs` to the server.
4. Restarts or starts the process with `pm2 startOrReload ecosystem.config.cjs --update-env`.
5. Saves the PM2 process list with `pm2 save`.

For manual/cloud VPS deployment and automatic restart setup, follow [`PRODUCTION_IMPROVEMENTS.md`](./PRODUCTION_IMPROVEMENTS.md).

Required GitHub secrets: `SSH_HOST`, `SSH_USERNAME`, `SSH_KEY`.

---

## Project status & known issues

This codebase had several security and reliability issues; all items in [`SUGGESTIONS.md`](./SUGGESTIONS.md) (🔴 Critical → 🔵 Low) have now been addressed — per-request DB re-init removed, payment/world credentials protected, JWT expiry, env-driven CORS, safe uploads, webhook signature verification + idempotency, atomic non-negative stock mutation, Midtrans access centralized in `lib`, and standardized responses.

Automated tests cover the pricing math (`model`) and JWT round-trip (`middleware`); run them with `go test ./...`. Recommended follow-ups (versioned migrations, moving the Midtrans charge to the SDK, structured logging, repo-level tests) are listed at the bottom of `SUGGESTIONS.md`.
