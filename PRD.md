# Product Requirements Document (PRD) — RapsShop Backend

| | |
|---|---|
| **Product** | RapsShop — Growtopia Diamond Lock (DL) / Blue Gem Lock (BGL) Trading Platform |
| **Component** | Backend API (`rapsshop-project`) |
| **Stack** | Go 1.23, Gin, GORM, MySQL, JWT, Midtrans, Supabase Storage |
| **Status** | Reverse-engineered from existing codebase (as-built) |
| **Last updated** | 2026-05-30 |

> This PRD was reconstructed by reviewing the current source code. It documents what the system *does today* and frames the intended product so it can be extended deliberately. Items the code implements but that look unintended are flagged in `SUGGESTIONS.md`.

---

## 1. Overview

RapsShop is the backend for an online shop that trades **Growtopia** in-game currency — **Diamond Locks (DL)** and **Blue Gem Locks (BGL)** (1 BGL = 100 DL). The platform serves two sides of a marketplace:

- **Customers buy** DL/BGL from the shop and pay through an online payment gateway (Midtrans) or manual transfer with proof upload. Purchased items are delivered in-game to the customer's Growtopia world/GrowID.
- **Customers sell** their DL to the shop; an admin reviews the submission, confirms payment to the seller, and the stock is added to inventory.

A single **admin** role manages inventory, pricing, payment methods, the in-game delivery world credentials, social media links, and customer testimonials, and reviews transaction recaps and profit reporting.

### 1.1 Goals
- Let customers purchase DL/BGL with automated payment (Midtrans) and manual options.
- Let customers sell DL to the shop and get paid after admin verification.
- Give admins inventory, price, and transaction management plus reporting (daily recap, totals, profit).
- Maintain a public-facing storefront content surface (prices, stock, testimonials, social links).

### 1.2 Non-goals (current scope)
- No customer accounts/authentication — buyers and sellers are anonymous, identified by WhatsApp/GrowID.
- No multi-admin roles / RBAC (single implicit admin role).
- No automated in-game delivery bot (delivery is manual; status is tracked).

---

## 2. Personas

| Persona | Description | Auth |
|---|---|---|
| **Customer (Buyer)** | Wants to buy DL/BGL, pays via Midtrans or manual transfer, receives delivery in-game. | None (anonymous) |
| **Customer (Seller)** | Wants to sell DL to the shop, uploads proof, gets paid after admin review. | None (anonymous) |
| **Admin / Operator** | Manages stock, prices, payment methods, delivery world, content, and reviews transactions. | JWT (Bearer token) |

---

## 3. Domain model

| Entity | Purpose | Key fields |
|---|---|---|
| `Admin` | Operator account | id (uuid), username, password (bcrypt), nama |
| `StockDL` | Inventory + price snapshot at a point in time | stock_dl, harga_jual_dl, harga_beli_dl, harga_beli_bgl, profit, waktu |
| `HargaDL` | Current public price list | harga_jual_dl, harga_beli_dl, harga_jual_bgl, harga_beli_bgl |
| `Growtopia` | In-game delivery world credentials | world, password, owner |
| `MetodePembayaran` | Payment method config shown to buyers | index_pembayaran, jenis_pembayaran, kredensial_pembayaran, pemilik |
| `PenjualanDL` | A customer **selling** DL to the shop | nama, jumlah_dl, jumlah_transaksi, wa, transfer, nomor_transfer, status, editor, harga_jual, bukti_dl |
| `PembelianDL` | A customer **buying** DL/BGL from the shop | id (uuid), world, nama, grow_id, jenis_item, jumlah_dl, wa, metode_transfer, jumlah_transaksi, button_bayar, status_pembayaran, status_pengiriman, editor, harga_beli, bukti_pembayaran |
| `Sosmed` | Social media links | username, platform, link |
| `Testimoni` | Customer testimonials | gambar, testi, uname, title |

**Terminology note:** In this domain, `penjualan` = *customer-initiated sale to the shop* (shop acquires stock), and `pembelian` = *customer-initiated purchase from the shop* (shop sells stock). Pricing: 1–99 DL is priced per DL; multiples of 100 are priced per BGL; mixed quantities split into BGL + DL.

---

## 4. Functional requirements

### 4.1 Authentication & Admin
- **FR-1** An admin can register with username/password, gated by a shared secret (`TOKEN`). Passwords are hashed with bcrypt.
- **FR-2** An admin can log in and receive a JWT bearer token.
- **FR-3** An admin can fetch their own profile using the token.
- **FR-4** Protected admin actions require a valid `Authorization: Bearer <token>` header.

### 4.2 Inventory (Stock)
- **FR-5** Admin can create a stock snapshot (quantity + prices).
- **FR-6** Public can read the latest stock and the full stock history.
- **FR-7** Admin can increment (`tambah`) or decrement (`kurangi`) stock; decrements are triggered automatically when a purchase is paid/shipped, increments when a customer sale is confirmed.
- **FR-8** Admin can delete the latest stock entry.

### 4.3 Pricing
- **FR-9** Admin can create/update/delete the current price list (DL & BGL buy/sell prices).
- **FR-10** Public can read the latest price.

### 4.4 Customer Purchase (Pembelian)
- **FR-11** A customer can create a purchase order specifying world, GrowID, item type, quantity, WhatsApp, and payment method.
- **FR-12** For gateway payments, the system creates a Midtrans charge (QRIS, GoPay, ShopeePay, or bank transfer for BCA/BRI/BNI) and returns the payment instructions.
- **FR-13** For manual payment, the system returns the configured payment-method credentials and the customer uploads proof of payment.
- **FR-14** The system records `jumlah_transaksi` (total price) computed from current prices and quantity.
- **FR-15** Midtrans sends a status notification (webhook); the system updates payment status and decrements stock on success.
- **FR-16** A customer/admin can query a purchase's DB record and its live Midtrans status.
- **FR-17** Admin can confirm manual payment and update delivery (`status_pengiriman`) status; the editing admin's username is recorded.

### 4.5 Customer Sale (Penjualan)
- **FR-18** A customer can submit a sale: quantity, WhatsApp, bank/transfer details, name, and an image proof. Total is computed from the latest sell price.
- **FR-19** Admin can list (paginated), view, update status, and delete sales.
- **FR-20** When a sale status transitions to *approved*, stock is incremented; reverting decrements it. The editing admin is recorded.

### 4.6 Reporting
- **FR-21** Admin can get a daily recap (`/rekapitulasi`) of sales and purchases grouped by rate for a given date.
- **FR-22** Admin can get monthly per-day totals of sales (`/penjualan/total`) and purchases (`/pembelian/total`).
- **FR-23** Admin can get a per-day profit report (`/profit`) for a given month.

### 4.7 Storefront content
- **FR-24** Admin can CRUD social media links; public can list them.
- **FR-25** Admin can CRUD testimonials (with image upload to Supabase); public can list them.
- **FR-26** Admin can set/read the Growtopia delivery world credentials.

### 4.8 Media storage
- **FR-27** Testimonial and sale-proof images are stored in Supabase Storage; purchase-proof images are stored on the local filesystem under `/public/payment` and served statically.

---

## 5. Non-functional requirements

| Category | Requirement |
|---|---|
| **Security** | Admin endpoints must require JWT. Payment credentials and in-game world passwords must not be publicly exposed. Secrets live in environment variables only. |
| **Performance** | A single shared DB connection pool must be reused across requests. Listing endpoints must paginate and count efficiently. |
| **Reliability** | Stock mutations must be atomic and must not go negative. Payment webhooks must be idempotent. |
| **Observability** | Requests and errors should be logged in a structured, queryable form. |
| **Portability** | Configuration via environment variables; runnable as a single binary. Deployed via GitHub Actions → rsync → pm2. |
| **Maintainability** | Consistent layered architecture (handler → service → repo); no business logic or direct DB access in handlers. |

> **Note:** Several of these NFRs are *targets*, not current behavior. See `SUGGESTIONS.md` for the gap analysis (e.g., per-request DB re-initialization, publicly exposed payment/world credentials, JWT without expiry).

---

## 6. API surface (as built)

Base path: `/api/v1`. 🔓 = public, 🔒 = requires JWT.

### Health
- `GET /ping` 🔓

### Admin
- `POST /admin/register` 🔓 (gated by `TOKEN` secret)
- `POST /admin/login` 🔓
- `GET /profile` 🔒

### Stock
- `POST /stock` 🔒 · `GET /stocks` 🔓 · `GET /stock` 🔓 · `PATCH /stock` 🔒 · `DELETE /stock` 🔒

### Price
- `POST /price` 🔒 · `GET /price` 🔓 · `PATCH /price` 🔒 · `DELETE /price` 🔒

### Growtopia env
- `POST /env` 🔒 · `GET /env` 🔓 *(exposes world password — see SUGGESTIONS)* · `PATCH /env` 🔒

### Social media
- `POST /platform` 🔒 · `GET /platforms` 🔓 · `GET /platform/:id` 🔒 · `PATCH /platform/:id` 🔒 · `DELETE /platform/:id` 🔒

### Testimonials
- `POST /testimoni` 🔒 (multipart) · `GET /testimonis` 🔓 · `GET /testimoni/:id` 🔒 · `PATCH /testimoni/:id` 🔒 · `DELETE /testimoni/:id` 🔒

### Payment methods
- `POST /payment` 🔒 · `GET /payments` 🔓 *(exposes credentials)* · `GET /payment/:id` 🔓 *(exposes credentials)* · `PATCH /payment/:id` 🔒 · `DELETE /payment/:id` 🔒

### Sales (Penjualan)
- `POST /penjualan` 🔓 (multipart) · `GET /penjualans` 🔒 · `GET /penjualan/:id` 🔒 · `GET /rekapitulasi` 🔒 · `GET /profit` 🔒 · `GET /penjualan/total` 🔒 · `PATCH /penjualan/:id` 🔒 · `DELETE /penjualan/:id` 🔒

### Purchases (Pembelian)
- `POST /pembelian` 🔓 (Midtrans charge) · `POST /new/pembelian` 🔓 (manual) · `POST /pembelian/status` 🔓 (Midtrans webhook) · `GET /pembelians` 🔒 · `GET /pembelian/total` 🔒 · `GET /pembelian/:id` 🔓 · `GET /pembelian/status/:id` 🔓 · `PATCH /pembelian/:id` 🔒 · `PATCH /pembelian/button/:id` 🔓 · `PATCH /pembelian/confirm/:id` 🔒 · `PATCH /upload/:id` 🔓 (proof upload) · `GET /public/*` 🔓 (static files)

---

## 7. External integrations

| Service | Use | Config |
|---|---|---|
| **MySQL** | Primary datastore (GORM, AutoMigrate) | `DB_USER`, `DB_PASS`, `DB_HOST`, `DB_NAME` |
| **Midtrans** | Payment gateway (charge + status) | `AUTHORIZATION_VALUE` (server key, base64), `MIDTRANS` (charge URL) — hardcoded to Production |
| **Supabase Storage** | Image hosting for testimonials & sale proofs | `SUPABASE_URL`, `SERVICE_TOKEN`, `STORAGE_NAME`, `BASE_URL` |
| **Local FS** | Purchase-proof images | `HOST_URL`, served from `/public/payment` |

---

## 8. Payment & order lifecycle

**Purchase (gateway):** create order → Midtrans charge → customer pays → Midtrans webhook (`/pembelian/status`) → status set to `success`/`pending`/`deny`/`failure`/`challange` → on success, stock decremented → admin marks `status_pengiriman` after in-game delivery.

**Purchase (manual):** create order (`/new/pembelian`) → return payment-method credentials → customer uploads proof (`/upload/:id`) → admin confirms (`/pembelian/confirm/:id`).

**Sale:** customer submits with proof image → admin reviews → status approved (`status = 1`) → stock incremented and recorded with editor.

Status codes for `PenjualanDL.status`: `0` = pending (default), `1` = approved, `-1` = rejected.

---

## 9. Open questions / assumptions
- Profit is currently computed as `total_pembelian − total_penjualan` per day; the intended sign/formula should be confirmed (see SUGGESTIONS).
- Midtrans is hardcoded to Production — there is no sandbox toggle.
- There is no customer-facing authentication; abuse protection (rate limiting, captcha) is out of scope today.
- Manual vs gateway purchase flows are both active (`/pembelian` and `/new/pembelian`); product should confirm which is canonical.
