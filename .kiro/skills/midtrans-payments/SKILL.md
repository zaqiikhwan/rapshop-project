---
name: midtrans-payments
description: Patterns for working with Midtrans payments in RapsShop (charge, status check, webhook signature verification, sandbox/production). Use when touching pembelian_dl payment flows or the lib/ Midtrans driver.
---

# Midtrans payments

All Midtrans access goes through `lib/midtrans.go` (`CoreAPI`). Handlers/services must not hand-roll Midtrans HTTP calls.

## Environment
- `AUTHORIZATION_VALUE` — Midtrans **server key** (base64-encoded into the `Basic` auth header at runtime).
- `MIDTRANS` — charge endpoint URL (e.g. `https://api.midtrans.com/v2/charge`).
- `MIDTRANS_ENV` — `sandbox` (default) or `production`; selects the SDK environment for status checks.

Never default to production. Local/dev must use sandbox.

## Charge
Build the payload via `model.MidtransData.IniDataPembelian()` (handles DL/BGL pricing) and send it through the `lib` driver. Branch on the parsed `status_code` field (>= 400 = failure) — never on response body length.

## Webhook (`POST /pembelian/status`)
Verify the signature **before** trusting the payload:
```
signature_key == sha512( order_id + status_code + gross_amount + AUTHORIZATION_VALUE )
```
Reject with 401 on mismatch. Make status handling idempotent: only decrement stock once per order (guard on the stored status so a repeated "success" notification doesn't double-decrement).

## Status mapping
`capture`+`accept` / `settlement` → `success` (decrement stock); `capture`+`challenge` → `challange`; `deny` → `deny`; `cancel`/`expire` → `failure`; `pending` → `pending`.

## Pricing rule (keep in sync with `model.MidtransData`)
- 1–99 DL → priced per DL (`HargaBeliDL`).
- multiple of 100 → priced per BGL (`HargaBeliBGL`), qty = DL/100.
- mixed (>100, not multiple) → BGL portion + DL remainder.
