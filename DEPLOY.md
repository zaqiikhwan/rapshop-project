# Linux VPS Deployment

Recommended production directory:

```bash
/opt/rapsshop
```

This is better than deploying under a personal home directory like `~/web/...` because `/opt` is commonly used for manually installed server applications. It is easier to find, easier to manage with PM2/systemd/Nginx, and cleaner for permissions and backups.

## Simple Layout

Use this layout first:

```text
/opt/rapsshop/
├── main
├── ecosystem.config.cjs
├── .env
├── logs/
└── public/
```

Optional future layout if you want release history and easier rollback:

```text
/opt/rapsshop/
├── current/
│   ├── main
│   ├── ecosystem.config.cjs
│   ├── .env
│   ├── logs/
│   └── public/
├── releases/
└── backups/
```

## Initial Server Setup

```bash
sudo mkdir -p /opt/rapsshop/logs
sudo chown -R $USER:$USER /opt/rapsshop
```

Copy production files:

```bash
rsync -avz main ecosystem.config.cjs user@server:/opt/rapsshop/
```

Create `/opt/rapsshop/.env` directly on the VPS. Do not commit real secrets.

Important production values:

```env
GIN_MODE=release
MIDTRANS_ENV=production
AUTO_MIGRATE=false
PORT=8080
HOST_URL=https://your-domain
ALLOWED_ORIGINS=https://your-frontend-domain
```

## PM2 Startup

```bash
cd /opt/rapsshop
pm2 startOrReload ecosystem.config.cjs --update-env
pm2 save
pm2 startup
```

After `pm2 startup`, PM2 prints a command that starts with `sudo env ...`. Run that exact command once. It registers PM2 with systemd so the app restarts automatically after VPS reboot.

Useful PM2 commands:

```bash
pm2 status
pm2 logs rapsshop --lines 100
pm2 monit
```

## Health Check

Local check on the VPS:

```bash
curl -fsS http://127.0.0.1:8080/ping
```

Public check through Nginx/domain:

```bash
curl -fsS https://your-domain/ping
```

Expected response:

```json
{"message":"pong"}
```

## Nginx Notes

Nginx should proxy API traffic to the Go app:

```nginx
location / {
    proxy_pass http://127.0.0.1:8080;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}
```

Use HTTPS with Certbot or your cloud provider's managed certificate.

## Deployment Checklist

- Build locally or in CI with `go build ./...`.
- Run `go vet ./...`.
- Run `go test ./...`.
- Back up MySQL before deploying.
- Upload `main` and `ecosystem.config.cjs` to `/opt/rapsshop`.
- Confirm `/opt/rapsshop/.env` has production values.
- Restart with `pm2 startOrReload ecosystem.config.cjs --update-env`.
- Run `pm2 save`.
- Check `/ping` locally and through the public domain.
- Confirm Midtrans webhook points to `https://your-domain/api/v1/pembelian/status`.
