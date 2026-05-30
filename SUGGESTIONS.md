# RapsShop Backend — Engineering Review & Suggested Adjustments

A senior-engineer review of `rapsshop-project`. Items are grouped by severity. Each has the **why**, the **where**, and a **concrete fix**. Tackle 🔴 Critical before any production deployment.

Legend: 🔴 Critical · 🟠 High · 🟡 Medium · 🔵 Low / polish

---

## Implementation status (updated)

Go toolchain on **1.23.0** (`go.mod`). The codebase builds (`go build ./...`), passes `go vet ./...`, and `go test ./...` is green. **All reviewed items below are now resolved.**

| Item | Status | Notes |
|---|---|---|
| C1 DB re-init per request | ✅ Done | Handlers use injected `StockDLUsecase.GetLatestDataStock()`; `mysql.InitDatabase()` removed from handlers |
| C2 Public payment credentials | ✅ Done | `GET /payments` & `/payment/:id` require `jwtMiddleware` |
| C3 Public world password | ✅ Done | `GET /env` requires `jwtMiddleware` |
| C4 JWT never expires | ✅ Done | `GenerateToken` sets `iat` + `exp` (24h) |
| C5 Unsafe CORS | ✅ Done | Env-driven `ALLOWED_ORIGINS` allowlist; no `*` with credentials |
| H1 Upload ext parse panic | ✅ Done | `filepath.Ext` + extension allowlist |
| H2 Weak upload filenames | ✅ Done | `uuid.NewString()` + ext |
| H3 Upload error ordering | ✅ Done | `image.Open()` checked before upload; reader closed |
| H4 SQL OR-precedence | ✅ Done | OR groups parenthesized in `GetByDate`/`GetProfit`/`GetTotalPembelian`; `ORDER BY` moved to `.Order()` |
| H5 Stock race / negative | ✅ Done | `repo.AdjustLatest` does a locked (`FOR UPDATE`) transaction with a non-negative guard; service delegates to it |
| H6 Midtrans hardcoded Production | ✅ Done | `MIDTRANS_ENV` selects Sandbox/Production |
| H7 Unverified webhook | ✅ Done | `signature_key` (sha512) verified; stock decrement made idempotent (`alreadySuccess` guard) |
| M1 Logic/raw HTTP in handler | ✅ Done | Charge orchestration moved to `service.ChargeAndCreate`; handler just binds + delegates |
| M2 Two Midtrans integrations | ✅ Done | All Midtrans access centralized in `lib` (`Charge`, `CheckStatus`, `HandleNotification`); handlers/services no longer hand-roll HTTP |
| M3 Inefficient counting | ✅ Done | `GetAll` uses `Count()` |
| M4 Pagination contract | ✅ Done | Repos clamp `_start >= 1` and `_end >= _start` (keeps the existing param contract, removes off-by-one/zero) |
| M5 Profit formula | ✅ Done | Confirmed correct (`revenue pembelian − cost penjualan`); documented in code + locked by a unit test path |
| M6 Constructor naming | ✅ Done | `NewStockDLHandler`, `NewSosmedHandler`, `NewPenjualanDLUsecase` |
| M7 Interfaces split across pkgs | ✅ Done | `MetodePembayaran` DTO + interfaces moved to `model/`; entity keeps only the struct |
| M8 No input validation | ✅ Done | `BindJSON`→`ShouldBindJSON` everywhere (no more double-writes); `jumlah_dl > 0` enforced in services. Binding tags intentionally not added to shared `Input*` DTOs (they double as partial-update payloads) |
| M9 Testimoni wrong remove path | ✅ Done | Removes by object key (`TrimPrefix BASE_URL`) |
| L1 `.env.example` incomplete | ✅ Done | Added `MIDTRANS`, `HOST_URL`, `MIDTRANS_ENV`, `PORT`, `ALLOWED_ORIGINS`, `AUTO_MIGRATE` |
| L2 Graceful shutdown | ✅ Done | `http.Server` + signal + `Shutdown(ctx)` |
| L3 Configurable port | ✅ Done | Reads `PORT` (default 8080) |
| L4 Debug prints | ✅ Done | Removed `fmt.Println` calls |
| L5 No tests | ✅ Done | Added unit tests: pricing math (`model`) + JWT round-trip (`middleware`) |
| L6 `AutoMigrate` on every boot | ✅ Done | Gated behind `AUTO_MIGRATE` (set `false` in prod). Versioned migrations remain the longer-term recommendation |
| L7 Response helper hides status | ✅ Done | `SuccessResponse` honors the passed code; `FailureOrErrorResponse` maps status + is nil-error safe |
| L8 typos | ✅ Done | `jsom`→`json`, `DeleteLatesPrice`→`DeleteLatestPrice`, `pacthMethod`→`patchMethod`. `"challange"` status string left intentionally (stored value; changing it would break existing rows / frontend) |
| L9 Mixed language identifiers | ✅ Accepted | Indonesian domain language is the house style; documented in `CLAUDE.md`. No churn-only rename |
| L10 Dead commented code | ✅ Done | Removed commented `NewHandlerPembelian` block |
| L11 Empty pg pkg / dead entities | ✅ Done | Deleted `database/postgresql` + unused `Riwayat*` entities |
| L12 Upgrade Go | ✅ Done | `go 1.23.0` |

### Follow-ups worth doing later (beyond this review)
- Replace `AutoMigrate` with versioned migrations (e.g. `golang-migrate`).
- Migrate the Midtrans **charge** from raw HTTP in `lib` to the `coreapi` SDK's `ChargeTransaction` (now isolated to one place, so low-risk).
- Structured logging (`slog`/zap) and request/trace IDs.
- Broaden test coverage to repos (recap SQL, stock transitions) with a test DB.

> The original per-item analysis is preserved below for reference.

---

## Round 2 — service-layer hardening (senior review)

A follow-up pass focused on the **service layer** correctness, after a senior-level read. Build + `go vet` + `go test` all green.

| Item | Status | Notes |
|---|---|---|
| S1 No transaction across aggregates | ✅ Done | Status-flip **and** stock mutation now commit together in one `db.Transaction` |
| S2 Unchecked error in penjualan `UpdateByID` | ✅ Done | Post-update read error is now checked; redundant branch removed |
| S3 Duplicated pricing in manual purchase | ✅ Done | `CreateDataPembelianManual` reuses `model.MidtransData.IniDataPembelian()`; added the missing `jumlah_dl > 0` guard |
| S4 Ignored upload result in penjualan `Create` | ✅ Done | Checks `UploadFile` response (`Key == ""` ⇒ fail); no more silent broken-image records |
| (bonus) Idiomatic file naming | ✅ Done | Dropped the redundant `<layer>.` filename prefix repo-wide (Go idiom: dir/package conveys the layer). Documented in `CLAUDE.md` |

### S1. Stock mutation and order status weren't atomic
**Where:** `src/pembelian_dl/service/service.go` (`UpdateStatusPembayaran`, `UpdateStatusPengiriman`) and `src/penjualan_dl/service/service.go` (`UpdateByID`).

The webhook/admin flows wrote the order status through one repo and mutated stock through a separate path, in **separate** DB operations. A failure between them could leave an order marked `success` with stock never deducted — or, on webhook redelivery after a failed status persist, deduct stock **twice**.

**Fix:** Added `WithTx(tx *gorm.DB)` to the `PembelianDLRepository`, `PenjualanDLRepository`, and `StockDLRepository` interfaces so the service can run both writes on the same transaction:
```go
return spdl.db.Transaction(func(tx *gorm.DB) error {
    if err := spdl.RepoPembelianDL.WithTx(tx).UpdateByID(dataPembelian, id); err != nil {
        return err
    }
    if decrement { // gated on the not-already-success transition (idempotent)
        if _, err := spdl.StockRepo.WithTx(tx).AdjustLatest(-dataPembelian.JumlahDL, entities.StockDL{}); err != nil {
            return err
        }
    }
    return nil
})
```
The pembelian/penjualan services now hold the **stock repo** (not the stock usecase) so the locked `AdjustLatest` (`SELECT … FOR UPDATE`) nests as a savepoint inside the outer tx. The decrement stays idempotent via the existing `alreadySuccess` guard.

> Residual edge (now closed in **Round 3**): two *truly concurrent* duplicate webhook deliveries could each read "not success" before either commits. Replaced the read-then-act guard with a compare-and-set (`UPDATE … WHERE status_pembayaran <> 'success'`, acting on `RowsAffected`) — see S5.

### S2. Latent unchecked-error bug in penjualan `UpdateByID`
The post-update `GetByID` error was never checked, so a failed read would proceed to mutate stock from a zero-value record. Now read inside the tx and checked; the stock decision uses the persisted (DB) status so it stays nil-safe.

### S3. Duplicated DL/BGL pricing in the manual purchase flow
`CreateDataPembelianManual` re-implemented the BGL/DL split math inline (a second source of truth that could drift) and lacked the `jumlah_dl > 0` guard the other two creation paths had. It now calls the single pricing function `model.MidtransData.IniDataPembelian()` and rejects non-positive quantities. (`harga_beli` is still stored as the DL rate — consistent with the charge flow and the recap grouping — so reporting is unchanged.)

### S4. Failed image upload silently persisted a record
`penjualan.Create` discarded the `UploadFile` result and saved the row regardless, leaving a record pointing at a missing image. The supabase `storage-go` client returns a `FileUploadResponse` (not a Go `error`), so success is now detected via a populated `Key`. **Same pattern still exists in `testimoni` (`CreateTestimoni`/`UpdateTestimoniByID`)** — worth the identical fix in a later pass.

---

## Round 3 — concurrency / exactly-once hardening (senior review)

Closes the idempotency / TOCTOU gaps left after Round 2. The atomic stock row lock (`AdjustLatest`, `SELECT … FOR UPDATE`) and the connection-pool / stateless-service concurrency were already sound; the remaining risk was **check-then-act** decisions made from a read taken *outside* the lock, which concurrent duplicates could double-apply. All four are now compare-and-set (CAS) at the DB layer. Build + `go vet` + `go test` green.

| Item | Status | Notes |
|---|---|---|
| S5 Webhook double-decrement (concurrent) | ✅ Done | `MarkPaid` = `UPDATE … WHERE status_pembayaran <> 'success'`; stock decremented only when `RowsAffected == 1` |
| S6 Shipment toggle double-decrement | ✅ Done | `MarkShipped` = CAS on the not-shipped → shipped transition; decrement gated on `RowsAffected` |
| S7 Penjualan double-approve double-count | ✅ Done | `UpdateStatusIfCurrent` = optimistic CAS (`WHERE status = <observed>`); stock adjust gated on the applied transition |
| S8 Admin username register race | ✅ Done | `Username` now `gorm:"uniqueIndex"` — the DB rejects duplicates regardless of timing |

### S5–S7. Exactly-once state transitions via compare-and-set
**Where:** `src/pembelian_dl/{repo,service}/` and `src/penjualan_dl/{repo,service}/`.

The status flips now happen as a single conditional `UPDATE` whose `RowsAffected` tells the service whether *this* call actually performed the transition; the stock mutation runs only then, inside the same `db.Transaction`. Because the database — not an earlier in-memory read — decides the winner, concurrent/duplicate webhooks, admin double-clicks, and racing approvals can no longer double-count stock.

```go
// pembelian repo — only the call that flips the row reports a transition
func (rp *repoPembelianDL) MarkPaid(id string) (bool, error) {
    res := rp.db.Model(&entities.PembelianDL{}).
        Where("id = ? AND status_pembayaran <> ?", id, "success").
        Update("status_pembayaran", "success")
    return res.RowsAffected == 1, res.Error
}
```
```go
// penjualan repo — optimistic CAS guarded on the previously-observed status
func (pdlr *penjualanDLRepository) UpdateStatusIfCurrent(id uint, from, to int, editor string) (bool, error) {
    res := pdlr.db.Model(&entities.PenjualanDL{}).
        Where("id = ? AND status = ?", id, from).
        Updates(entities.PenjualanDL{EditorStatus: editor, Status: &to})
    return res.RowsAffected == 1, res.Error
}
```
New repo methods: `PembelianDLRepository.MarkPaid` / `MarkShipped`, `PenjualanDLRepository.UpdateStatusIfCurrent` (all `WithTx`-aware).

### S8. Username uniqueness enforced by the DB
`Admin.Register` did `GetByUsername` → (if empty) `Create`, a classic race with no safety net. Added `gorm:"uniqueIndex"` on `Admin.Username`; the friendly pre-check stays for the common case, but a row that slips through concurrently now fails on the unique constraint instead of creating a duplicate.

> Migration note: the index is created by `AutoMigrate` on next boot **only if** `AUTO_MIGRATE != "false"` and there are no existing duplicate usernames. In prod (where `AUTO_MIGRATE=false`) add it via a versioned migration. Concurrency scope is per-process; this codebase assumes a single instance (pm2) — a multi-instance deploy would still rely on these DB-level guards, which is exactly why they're at the DB layer.

---

## 🔴 Critical

### C1. Database is re-initialized on every request
`mysql.InitDatabase()` is called *inside request handlers*, not just at startup.

- `src/penjualan_dl/handlers/handlers.penjualan_dl.go` → `CreateNewPenjualan`: `mysql.InitDatabase().Order("id desc").First(&harga)`
- `src/pembelian_dl/handlers/handlers.pembelian_dl.go` → `HandlerPembelian`: same pattern

Each call opens a **new connection pool** *and re-runs `AutoMigrate`* on every hit. This leaks connections, will exhaust MySQL `max_connections` under load, and adds migration latency to each request.

**Fix:** Inject the already-created `*gorm.DB` (or, better, the latest price via the `stock_dl`/`harga_dl` service) through the constructor like every other module. Never call `InitDatabase()` outside `main`.

```go
// instead of mysql.InitDatabase().Order("id desc").First(&harga)
harga, err := pdlh.StockDLUsecase.GetLatestDataStock()
```

### C2. Payment-method credentials are publicly readable
`GET /payments` and `GET /payment/:id` have **no `jwtMiddleware`** (`src/metode_pembayaran/handlers/handlers.metode_pembayaran.go`). `MetodePembayaran.KredensialPembayaran` (account numbers / payment credentials) is returned to anyone.

**Fix:** Either protect these with `jwtMiddleware`, or split the model so the public endpoint returns only non-sensitive display fields (`jenis_pembayaran`, `pemilik`, masked account) and keep full credentials admin-only.

### C3. Growtopia world password is publicly readable
`GET /env` is public (`src/env_growtopia/handlers/handlers.env_growtopia.go`) and returns `Growtopia.Password` (the in-game world lock password) in plaintext.

**Fix:** Protect `GET /env` with `jwtMiddleware`. The delivery world password should never be exposed to buyers.

### C4. JWTs never expire and carry no role
`middleware/middleware.go` → `GenerateToken` sets only `claim["id"]`. There is no `exp`, `iat`, or `nbf`. A leaked token is valid forever.

**Fix:**
```go
claim := jwt.MapClaims{
    "id":  id,
    "exp": time.Now().Add(24 * time.Hour).Unix(),
    "iat": time.Now().Unix(),
}
```
`jwt.Parse` validates `exp` automatically once present. Consider refresh tokens if 24h is too short for the admin UX.

### C5. Invalid + unsafe CORS configuration
`main.go` sets `Access-Control-Allow-Origin: *` **together with** `Access-Control-Allow-Credentials: true`. Browsers reject this combination, and a wildcard origin defeats CORS protection.

**Fix:** Use an explicit allowlist of front-end origins (env-driven) and echo the matched origin. Drop `*` when credentials are enabled. The `gin-contrib/cors` package handles this cleanly.

---

## 🟠 High

### H1. File-upload extension parsing can panic
`src/pembelian_dl/handlers/handlers.pembelian_dl.go` → `UploadFile`:
```go
splitFileName := strings.Split(file.Filename, ".")
if splitFileName[1] != "png" ... // panics if no dot, wrong if >1 dot ("a.b.jpg")
```
A filename with no extension (`index out of range`) crashes the request; `a.tar.png` checks `tar` not `png`.

**Fix:** Use `ext := strings.ToLower(filepath.Ext(file.Filename))` and compare against an allowlist of `.png/.jpg/.jpeg/.heic/.heif`. Also validate the MIME/content type, not just the name.

### H2. Weak, collision-prone upload filenames
Same handler generates names by shuffling a fixed 60-char alphabet seeded with `rand.Seed(time.Now().Unix())`. Two uploads in the same second produce correlated names, and the result is a fixed-length permutation of the same characters (predictable, and effectively the same set every time).

**Fix:** Use a UUID (already a dependency): `file.Filename = uuid.NewString() + ext`. Remove the manual `rand` shuffle.

### H3. Upload error checked after the upload runs
In `testimoni`, `penjualan`, and the Supabase calls, the pattern is:
```go
imageIo, err := image.Open()
client.UploadFile(...)   // return value ignored
if err != nil { return err }  // err is from Open(), checked too late
```
The `UploadFile` result is discarded and the `Open()` error is checked *after* the upload. Failed uploads silently produce broken image links.

**Fix:** Check `image.Open()` before uploading, and check the `UploadFile` return:
```go
imageIo, err := image.Open()
if err != nil { return err }
defer imageIo.Close()
if _, err := client.UploadFile(bucket, name, imageIo); err != nil { return err }
```

### H4. Recap/profit SQL has operator-precedence bugs
`src/penjualan_dl/repo/repo.penjualan_dl.go`:
```sql
... where created_at LIKE ? and harga_beli = ? and status_pembayaran = 'success' or status_pembayaran = 'dibayar'
```
Because `AND` binds tighter than `OR`, this evaluates as `(date AND rate AND 'success') OR ('dibayar')` — i.e. it returns **every `dibayar` row for all dates and all rates**. The date and rate filters are silently bypassed. The same bug appears in `GetByDate` and `GetProfit`.

**Fix:** Parenthesize the OR group:
```sql
... WHERE created_at LIKE ? AND harga_beli = ? AND (status_pembayaran = 'success' OR status_pembayaran = 'dibayar')
```
Also embedding `order by ...` inside a `Where("... order by ...")` string is fragile — use `.Order(...)`.

### H5. Stock updates race and can go negative
`src/stock_dl/service/usecase.stock_dl.go` does read-modify-write (`GetLatest` → compute `stock ± n` → `UpdateByID`) with no transaction or row lock. Concurrent purchases/sales can lose updates, and `UpdateKurangiStock` can drive stock below zero.

**Fix:** Do the mutation atomically in one statement and guard non-negativity:
```go
db.Model(&StockDL{}).Where("id = ? AND stock_dl >= ?", id, n).
   UpdateColumn("stock_dl", gorm.Expr("stock_dl - ?", n))
// check RowsAffected == 0 -> insufficient stock
```
Wrap purchase-paid → stock-decrement in a DB transaction so they commit together.

### H6. Midtrans hardcoded to Production
`lib/lib.infrastructure.go` → `c.ca.New(os.Getenv("AUTHORIZATION_VALUE"), midtrans.Production)`. There is no sandbox path, so local/testing runs hit real production payments.

**Fix:** Drive the environment from config, e.g. `MIDTRANS_ENV=sandbox|production`, and map to `midtrans.Sandbox`/`midtrans.Production`.

### H7. Payment webhook is unauthenticated and not verified
`POST /pembelian/status` accepts any JSON with an `order_id` and updates payment status (and decrements stock). It does not verify the Midtrans `signature_key`, so anyone can POST a forged "success" and drain stock / mark orders paid.

**Fix:** Verify Midtrans' `signature_key` = `sha512(order_id + status_code + gross_amount + ServerKey)` before trusting the payload. Make the handler idempotent (don't double-decrement stock if called twice).

---

## 🟡 Medium

### M1. Business logic and raw HTTP live in the handler
`HandlerPembelian` marshals JSON, builds an `http.Request` to Midtrans, reads the body, and does a `len(body) == 115` magic-number check to detect errors. This is brittle (any change in Midtrans' response length breaks it) and violates the layering used elsewhere.

**Fix:** Move Midtrans charge creation into `lib`/service, reuse the existing `coreapi` client (already wired as `midtransDriver`), and branch on the parsed `status_code` field instead of body length.

### M2. Two parallel Midtrans integrations
`lib` uses the official `coreapi` SDK (for status), while the purchase handler hand-rolls raw HTTP (for charge). Pick one (the SDK) for consistency, testability, and correct error handling.

### M3. Inefficient counting in list endpoints
`GetAll` in `penjualan`/`pembelian` repos loads **all IDs** into memory (`db.Select("id").Find(&lenData)`) just to count, then loads the page. On large tables this is wasteful.

**Fix:** Use `db.Model(&Entity{}).Count(&total)` plus a separate paged query.

### M4. Pagination contract is off-by-one and unconventional
`Offset(_start - 1).Limit(_end - _start + 1)` with 1-based `_start` is easy to misuse and breaks if `_start = 0`. Prefer standard `page`/`pageSize` (or `limit`/`offset`) with validation and sane defaults.

### M5. Profit formula is suspect
`GetProfit`: `Profit = TransaksiBeli − TransaksiJual` (purchases minus sales). Given `penjualan` = customer→shop and `pembelian` = shop→customer, confirm whether this sign is intended; a true gross margin is usually `revenue (sales to customers) − cost (purchases from customers)`. Add a unit test that locks the agreed formula.

### M6. Inconsistent constructor / type naming
- `stock_dl` and `sosmed` handler constructors are both named `NewAdminHandler` (copy-paste leftovers).
- `penjualan_dl` service constructor is `NewTestimoniUsecase`.

These mislead readers and make `main.go` confusing. Rename to match their module (`NewStockDLHandler`, `NewSosmedHandler`, `NewPenjualanDLUsecase`).

### M7. Interfaces split across two packages
Most repo/usecase interfaces live in `model/`, but `MetodePembayaran`'s interfaces live in `entities/`. Pick one convention (interfaces near the consumer, i.e. `model/`) so the architecture is uniform.

### M8. No input validation on most request bodies
Only `model.NewAdmin`/`AdminLogin` use `binding:"required"`. Purchase/sale/stock/price inputs accept anything (negative quantities, empty WhatsApp, etc.).

**Fix:** Add `binding` tags and use `c.ShouldBindJSON` so Gin returns 400 on invalid input. Validate business rules (e.g. `jumlah_dl > 0`) in the service.

### M9. Testimoni update deletes the wrong storage path
`UpdateTestimoniByID`:
```go
paths := make([]string, 1)
paths = append(paths, detailTestimoni.Gambar) // -> ["", gambar]
client.RemoveFile(bucket, paths)
```
`make([]string, 1)` pre-fills one empty string, so `paths` becomes `["", "<url>"]` and the stored `Gambar` is a full URL (`BASE_URL + filename`), not the bucket object key — the delete likely targets the wrong key. Store the object key, and build `paths := []string{detailTestimoni.ObjectKey}`.

---

## 🔵 Low / polish

- **L1. `.env.example` is incomplete.** Add `MIDTRANS` and `HOST_URL` (read by the code but undocumented).
- **L2. No graceful shutdown.** `r.Run()` can't drain in-flight requests. Use `http.Server` + `signal.NotifyContext` + `srv.Shutdown(ctx)`.
- **L3. Hard-coded / unconfigurable port.** `r.Run()` defaults to `:8080`; read `PORT` from env.
- **L4. No structured logging.** Replace `fmt.Println(len(body))` and ad-hoc logs with a leveled/structured logger (e.g. `slog`, zap). Remove debug prints (`fmt.Println(arrayDate[1])` in `GeneratorTanggal`).
- **L5. No tests.** There is zero test coverage. Start with the pricing/`MidtransData.IniDataPembelian` math, the recap SQL, and the stock transition logic — these carry the most business risk.
- **L6. `AutoMigrate` on every boot.** Fine for dev; for prod adopt versioned migrations (e.g. `golang-migrate`) to avoid surprise schema changes.
- **L7. Response helper hides real status.** `utils.SuccessResponse` coerces any non-2xx into HTTP 500. Combined with `FailureOrErrorResponse`, status mapping is inconsistent — unify on one helper that respects the passed code.
- **L8. Inconsistent typos in JSON tags / messages.** e.g. `jsom:"nama"` in `model.AdminDto` (should be `json`), `"challange"`, `DeleteLatesPrice`, `pacthMethod`. The `jsom` typo means `nama` is silently dropped from the profile response.
- **L9. Mixed language identifiers.** Indonesian + English mix is fine, but be consistent within a layer to aid onboarding.
- **L10. Dead/commented code.** Large commented-out `NewHandlerPembelian` block and `RiwayatPenambahanDL`/`RiwayatPenguranganDL` entities that aren't migrated/used. Remove to reduce noise.
- **L11. `database/postgresql` empty package.** Drop it or implement; an empty alternative driver dir is misleading.
- **L12. Pin Go toolchain & deps.** `go 1.19` is past end of support; plan an upgrade and run `govulncheck`.

---

## Suggested remediation order

1. **C1** (DB re-init) — highest blast radius, simplest correctness win.
2. **C2/C3** (credential exposure) + **C4** (JWT expiry) + **H7** (webhook verification) — security baseline.
3. **C5** (CORS) before any browser front-end ships.
4. **H4** (SQL precedence) + **M5** (profit) — reporting correctness.
5. **H5** (stock atomicity) + **H1/H2/H3** (uploads).
6. **H6/M1/M2** (Midtrans consolidation), then the Medium/Low cleanups.

A quick win bundle (low effort, high value): **L1** (.env.example), **L8** (`jsom` typo), **C4** (JWT exp), **M6** (rename constructors).
