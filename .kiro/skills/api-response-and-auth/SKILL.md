---
name: api-response-and-auth
description: RapsShop HTTP conventions — the JSON response envelope, JWT auth middleware, and which routes must be protected. Use when adding/reviewing handlers or auth.
---

# API responses & auth

## Response envelope
Every response goes through `utils`:
```go
utils.SuccessResponse(c, http.StatusOK, "message", data)
utils.FailureOrErrorResponse(c, http.StatusBadRequest, "message", err)
```
Shape: `{ "status_code", "status", "message", "data" }`. Pass the real HTTP status — the helper honors it. Never pass a `nil` error to `FailureOrErrorResponse` (it calls `err.Error()`); use `errors.New(...)`.

## Auth
- Login (`POST /admin/login`) returns a JWT. Protected routes use the `jwtMiddleware` (from `middleware.NewAuthMiddleware`) which sets `id` in the context.
- Tokens carry `id`, `iat`, `exp` (24h). Generate via `middleware.GenerateToken(id)`.
- Admin registration is gated by the shared `TOKEN` env secret, not JWT.

## Which routes must be protected
- All create/update/delete (mutating) routes → `jwtMiddleware`.
- Any route returning **secrets** → `jwtMiddleware`: payment credentials (`/payments`, `/payment/:id`) and the delivery-world password (`/env`).
- Public reads are OK for storefront content (prices, stock, testimonials, social links) and customer-initiated create flows (`/penjualan`, `/pembelian`, proof upload).

## Binding & validation
Bind with `c.ShouldBindJSON(&input)` (not `BindJSON`, which double-writes the response). Add `binding:"required"` / `binding:"gt=0"` to dedicated `Input*` DTOs — not to shared entity structs that are also used for partial updates. Enforce business rules (e.g. `jumlah_dl > 0`) in the service.
