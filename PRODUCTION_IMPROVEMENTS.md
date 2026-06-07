# Production Improvements

This plan is based on `PRD.md`, `AGENTS.md`, and the current transaction/deployment shape. The target is not a literal guarantee of 99% payment completion, because Midtrans, banks/e-wallets, Growtopia, and the VPS can still fail. The goal is to make successful customer intent become a correct recorded transaction at least 99% of the time under normal operations.

## Reliability Target

| Flow | Success definition |
|---|---|
| Midtrans purchase | Order is created, payment instructions are returned, webhook/live status eventually marks paid, stock is decremented once, admin can deliver. |
| Manual purchase | Order is created, proof upload is stored, admin confirmation marks paid once, stock is decremented once. |
| Customer sale | Proof is stored, admin approval changes status once, stock is incremented once. |
| VPS runtime | Process restarts automatically after crashes/reboots, deploys are repeatable, rollback is possible. |

## Already Good In Current Codebase

- Payment status transitions and stock mutations are transactional and idempotent.
- Stock mutation uses DB-level guards so stock cannot go negative.
- Midtrans webhook signatures are verified before trust.
- Admin-only data is protected by JWT.
- `main.go` has graceful shutdown and `/ping` health check.
- CI runs build, vet, tests, and lint for new issues.

## Priority 1: VPS Auto-Restart And Repeatable Deploy

Use the checked-in `ecosystem.config.cjs` on any manual/cloud VPS deployment.

Initial VPS setup:

```bash
sudo mkdir -p /opt/rapsshop/logs
sudo chown -R $USER:$USER /opt/rapsshop
cd /opt/rapsshop
pm2 startOrReload ecosystem.config.cjs --update-env
pm2 save
pm2 startup
```

After `pm2 startup`, run the exact command printed by PM2. That registers PM2 with systemd so the app returns after a VPS reboot.

Manual deploy:

```bash
go build -o main
rsync -avz main ecosystem.config.cjs user@server:/opt/rapsshop/
ssh user@server 'cd /opt/rapsshop && mkdir -p logs && pm2 startOrReload ecosystem.config.cjs --update-env && pm2 save'
```

Production env requirements:

- `GIN_MODE=release`
- `MIDTRANS_ENV=production`
- `AUTO_MIGRATE=false`
- `PORT=8080` or the port used by your reverse proxy
- `HOST_URL=https://your-domain`
- `ALLOWED_ORIGINS=https://your-frontend-domain`

Recommended reverse proxy checks:

- Nginx proxies `/api/v1/*` and `/ping` to `127.0.0.1:$PORT`.
- TLS is enabled with Certbot or a managed certificate.
- Midtrans notification URL points to `https://your-domain/api/v1/pembelian/status`.
- Server firewall opens only `22`, `80`, and `443`; MySQL is private.

## Priority 2: Transaction Reconciliation

Even with webhooks, payment providers can miss or delay notifications. Add a reconciliation job so paid orders are recovered without manual DB edits.

Recommended behavior:

- Every 2-5 minutes, find recent `pembelian_dl` rows whose `status_pembayaran` is not final.
- Call the existing Midtrans status check for each row.
- Apply the same service path used by the webhook, so stock decrement remains idempotent.
- Log failures and retry later; never update stock directly from the job.

Final states should stop polling:

- `success`
- `settlement`
- `capture`
- `deny`
- `cancel`
- `expire`
- `failure`

## Priority 3: Versioned Database Migrations

Disable `AutoMigrate` in production and replace schema changes with explicit migrations.

Minimum first migration:

- Create all current tables.
- Add unique index for admin username.
- Add indexes for transaction lookup fields:
  - `pembelian_dl.id`
  - `pembelian_dl.status_pembayaran`
  - `pembelian_dl.created_at`
  - `penjualan_dl.status`
  - `penjualan_dl.created_at`

Before each deploy:

```bash
mysqldump --single-transaction --routines --triggers DB_NAME > backup-$(date +%Y%m%d-%H%M%S).sql
```

## Priority 4: Observability And Alerts

Add structured logs and basic monitoring before increasing traffic.

Must-have alerts:

- `/ping` fails for 2 checks.
- PM2 process restarts more than 3 times in 10 minutes.
- MySQL connection errors occur.
- Midtrans charge/status calls fail repeatedly.
- Orders stay pending for more than 30 minutes.
- Stock is below the minimum operating threshold.

Useful commands on VPS:

```bash
pm2 status
pm2 logs rapsshop --lines 100
pm2 monit
curl -fsS http://127.0.0.1:8080/ping
```

## Priority 5: Operational Transaction Playbook

When a customer reports a paid order that still shows pending:

1. Check the order by `GET /api/v1/pembelian/:id`.
2. Check live Midtrans status with `GET /api/v1/pembelian/status/:id`.
3. If Midtrans is paid but local status is pending, rerun reconciliation or resend the Midtrans notification from the Midtrans dashboard.
4. Do not manually edit stock in MySQL. Use the app's admin confirmation/status flow so idempotency rules stay intact.

When stock is wrong:

1. Export purchases, sales, and stock history for the affected date.
2. Recalculate expected stock from approved sales minus paid/shipped purchases.
3. Apply one admin stock correction entry and document the reason.

## Priority 6: Next Code Improvements

These are the highest-value code changes after the current hardening:

1. Add a background reconciliation command or worker mode, for example `go run main.go reconcile` or a separate `cmd/reconciler`.
2. Add request IDs and structured `slog` logging around order creation, webhook handling, stock adjustment, and Midtrans API calls.
3. Add integration tests for stock transitions using a test MySQL container.
4. Move Midtrans charge creation fully onto the official SDK now that access is centralized in `lib`.
5. Store local payment proofs in object storage too, so deploys/rebuilds cannot lose uploaded files.
6. Add rate limiting to public transaction and upload endpoints.

## Deployment Checklist

- CI is green: `go build ./...`, `go vet ./...`, `go test ./...`.
- `.env` on VPS has production values and no sandbox Midtrans credentials.
- `AUTO_MIGRATE=false` in production.
- Database backup exists before deployment.
- `ecosystem.config.cjs` exists on VPS.
- `pm2 startOrReload ecosystem.config.cjs --update-env` succeeds.
- `pm2 save` has been run.
- `pm2 startup` has been installed for reboot recovery.
- `/ping` returns `{"message":"pong"}` through the public domain.
- Midtrans webhook URL is correct and reachable over HTTPS.
