# Discussion Notes

Working notes for product, roadmap, monetization, and strategy discussions. These are not committed requirements until they are moved into `PRD.md`, implementation tickets, or a tracked backlog.

## Product Roadmap And Monetization

### Goals

- Increase customer trust in RapsShop as a DL/BGL shop.
- Improve repeat purchase behavior and customer retention.
- Increase margin without making operations harder for admins.
- Keep operational safety high around stock, payments, and delivery credentials.

### Near-Term Roadmap

Focus: conversion, trust, and admin efficiency.

1. Order tracking page
   - Customers can enter invoice/order ID.
   - Show payment status, delivery status, estimated processing state, and support contact.

2. Public stock and live pricing
   - Show available DL/BGL stock.
   - Show buy price, sell price, and last updated time.
   - Helps reduce repeated support questions and increases trust.

3. Payment method clarity
   - Separate instant payment, manual transfer, and fast processing options.
   - Show processing expectations per payment method.

4. Admin order workflow
   - Clear queues for unpaid orders, paid orders waiting delivery, delivered orders, pending sell orders, approved sell orders, rejected sell orders, and refund/review cases.

5. WhatsApp/support integration
   - Generate a support message link with order ID already filled in.

6. Basic anti-abuse controls
   - Rate limit order creation.
   - Flag suspicious repeated orders.
   - Log useful request metadata where legally and operationally appropriate.

### Mid-Term Roadmap

Focus: repeat customers and growth mechanics.

1. Customer accounts or lightweight login
   - Customers can view order history.
   - Customers can reorder faster.
   - Customers can save GrowID/world information.

2. Loyalty system
   - Discount after buying a configured amount of DL/BGL.
   - VIP tier based on monthly volume.
   - Cashback balance usable only in the shop.

3. Promo codes
   - Support campaigns, influencers, returning customers, and first-time buyers.

4. Dynamic pricing rules
   - Admin-configured buy margin.
   - Admin-configured sell margin.
   - BGL premium or discount.
   - Stock-based pricing.
   - Event or weekend pricing.

5. Referral system
   - Give both referrer and new customer a reward after a completed paid order.
   - Prefer store credit first for simpler cash-flow control.

6. Notification system
   - Send order status updates through WhatsApp, email, or Discord webhook.

### Long-Term Roadmap

Focus: defensibility and scale.

1. Wallet/balance system
   - Customers can deposit funds or keep shop credit.
   - Requires careful accounting, audit logs, and refund rules.

2. Automated delivery workflow
   - Potentially integrate delivery automation or operational bots.
   - Requires strong security controls around world credentials and delivery state.

3. Marketplace-like expansion
   - Support other Growtopia items while keeping DL/BGL as the liquidity core.

4. Seller liquidity program
   - Trusted sellers can receive faster approval, special rates, or scheduled bulk sell arrangements.

5. Analytics and forecasting
   - Track demand by day and hour.
   - Track margins, stock velocity, payment conversion, repeat purchase rate, and profit by customer segment.

### Monetization Ideas

1. Spread-based margin
   - Core model: buy DL/BGL from customers at one price and sell at a higher price.
   - Margins can vary by small orders, bulk orders, instant payment, manual processing, or VIP tier.

2. BGL bulk pricing
   - Encourage exact BGL quantities with better pricing.
   - Keep slightly higher per-DL pricing for smaller orders.
   - Add special bulk rates above configured BGL thresholds.

3. Processing fee
   - Possible fee types: instant payment fee, QRIS/payment gateway fee, fast delivery fee.
   - Consider baking fees into displayed price when visible fees reduce conversion.

4. VIP or membership
   - Benefits can include better buy/sell rates, faster processing, priority stock reservation, and exclusive bulk deals.
   - Most useful after repeat customer behavior is already measurable.

5. Referral and affiliate margin
   - Give community sellers or influencers referral codes.
   - Reward can be fixed commission, percentage of gross margin, or store credit.
   - Store credit is safer for early cash-flow control.

6. Stock reservation
   - Free reservation for a short checkout window.
   - Longer reservation or rate lock can be a VIP or bulk buyer benefit.

7. Sell-order differentiation
   - Standard sell order: normal approval and normal rate.
   - Fast approval: faster processing with slightly worse rate.
   - Bulk seller: negotiated rate.

### Metrics To Track Early

- Gross merchandise volume.
- Gross profit.
- Average order size.
- Repeat customer rate.
- Payment success rate.
- Paid but undelivered count.
- Average delivery time.
- Stock turnover rate.
- Buy/sell spread per DL/BGL.
- Refund and rejection rate.
- Support contact rate per order.

### Suggested Priority

1. Order tracking and public stock/pricing.
2. Better admin order queue.
3. Promo codes.
4. Loyalty and referral.
5. Dynamic pricing.
6. Wallet and VIP system.
7. Automation and analytics.

Reasoning: improve trust and operations first, then add growth mechanics, then add more complex monetization systems.

## Approved Direction: Increase Completed Orders

Decision date: 2026-06-07

Primary goal: increase completed orders by improving customer trust and reducing drop-off between order creation, payment, confirmation, and delivery.

### Current Product Position

The backend already has the base transaction engine:

1. Customer buy order through Midtrans: `POST /pembelian`.
2. Customer buy order through manual payment: `POST /new/pembelian`.
3. Payment proof upload: `PATCH /upload/:id`.
4. Customer order detail lookup: `GET /pembelian/:id`.
5. Midtrans status lookup: `GET /pembelian/status/:id`.
6. Admin manual payment confirmation: `PATCH /pembelian/confirm/:id`.
7. Admin delivery status update: `PATCH /pembelian/:id`.
8. Public stock and price lookup: `GET /stock`, `GET /price`.

The missing product layer is clearer customer confidence and a more explicit order journey.

### Phase 1 MVP: Complete Order Experience

1. Public order tracking
   - Customer can track an order using the order ID.
   - Show order ID, item type, quantity, total payment, payment status, delivery status, payment proof URL when available, and created time.

2. Checkout clarity
   - Before order creation, the storefront should show current stock, current DL/BGL price, estimated total, available payment types, and stock availability.

3. Admin order queue
   - Admin should be able to filter orders by operational state: pending payment, proof uploaded, paid and waiting delivery, delivered, failed, or rejected.

4. Manual payment trust
   - Manual transfer should have customer-friendly statuses such as `pending_upload`, `waiting_confirmation`, `dibayar`, `rejected`, and `delivered`.

5. Customer support fallback
   - Order responses or frontend views should provide a WhatsApp support action with the order ID prefilled.

### Approved Default Decisions

- Public stock visibility: show exact stock.
- Pricing model: admin-controlled fixed pricing for now.
- Payment priority: Midtrans first, manual transfer as fallback.
- Customer account system: not needed for Phase 1.
- Monetization priority: keep buy/sell spread as the core model; add promo, referral, loyalty, VIP, and dynamic pricing later.
- First success metric: completed orders per day.
- Secondary metrics: payment success rate, abandoned orders, average delivery time, and support contact rate per order.

### First Implementation Backlog

1. Improve public order tracking response.
2. Add customer-friendly order status mapping.
3. Add admin order queue filters.
4. Add stock validation before purchase creation.
5. Add support/contact metadata to order response.
6. Later: promo, referral, loyalty, VIP, and dynamic pricing.

## Phase 1 API Contract Draft

Working branch: `phase-1-complete-order-experience`

### Git Workflow

1. Base branch: `review-hardening`.
2. Feature branch: `phase-1-complete-order-experience`.
3. Keep Phase 1 commits focused on the complete order experience.
4. Suggested commit groups:
   - `docs: define phase 1 order experience contract`
   - `feat: add public order tracking response`
   - `feat: add admin pembelian queue filters`
   - `feat: validate stock before purchase creation`
   - `test: cover phase 1 purchase tracking and stock validation`
5. Before merging, run:
   - `go build ./...`
   - `go vet ./...`
   - `go test ./...`

### API Scope

Phase 1 should reuse the existing `PembelianDL` table and purchase routes unless a new response shape is needed for customer-facing tracking.

Preferred approach:

1. Keep `GET /api/v1/pembelian/:id` as the raw/order-detail route for compatibility.
2. Add a customer-friendly tracking route:
   - `GET /api/v1/pembelian/:id/tracking`
3. Extend the admin listing route with queue filters:
   - `GET /api/v1/pembelians?_start=1&_end=20&queue=waiting_delivery`

### Public Tracking Response

Endpoint:

```text
GET /api/v1/pembelian/:id/tracking
```

Response data:

```json
{
  "id": "order uuid",
  "world": "customer world",
  "nama": "customer name",
  "grow_id": "customer grow id",
  "jenis_item": "DL",
  "jumlah_dl": 100,
  "metode_transfer": 1,
  "jumlah_transaksi": 1000,
  "status_pembayaran": "pending",
  "status_pembayaran_label": "Waiting for payment",
  "status_pengiriman": false,
  "status_pengiriman_label": "Waiting for delivery",
  "bukti_pembayaran": "https://example.com/api/v1/public/file.jpg",
  "support": {
    "channel": "whatsapp",
    "message": "Halo admin, saya ingin bertanya tentang order <id>"
  },
  "created_at": "2026-06-07T00:00:00+07:00",
  "updated_at": "2026-06-07T00:00:00+07:00"
}
```

Notes:

- Do not expose admin-only data.
- Keep `wa` out of the public tracking response unless the customer already provided and the frontend needs it.
- `jenis_item` should be display-friendly even though the database stores it as `bool`.
- Tracking should return `404` when the order ID does not exist.

### Customer-Friendly Status Mapping

Payment status mapping:

| Stored status | Customer label | Meaning |
|---|---|---|
| `belum_dibayar` | Waiting for payment | Order exists but no payment has been completed. |
| `pending` | Payment pending | Midtrans/payment provider is still processing. |
| `success` | Payment confirmed | Gateway payment succeeded. |
| `dibayar` | Payment confirmed | Manual payment was approved by admin. |
| `deny` | Payment denied | Gateway rejected the payment. |
| `failure` | Payment failed | Payment failed, expired, or was cancelled. |
| `challange` | Under review | Gateway marked the payment for review. |

Delivery status mapping:

| Stored status | Customer label | Meaning |
|---|---|---|
| `false` | Waiting for delivery | Payment may still be pending or admin has not delivered yet. |
| `true` | Delivered | Admin marked the order as delivered. |

Manual payment derived state:

| Condition | Customer label |
|---|---|
| `status_pembayaran = belum_dibayar` and `bukti_pembayaran = ""` | Waiting for proof upload |
| `status_pembayaran = belum_dibayar` and `bukti_pembayaran != ""` | Waiting for admin confirmation |
| `status_pembayaran = dibayar` and `status_pengiriman = false` | Payment confirmed, waiting for delivery |
| `status_pembayaran = dibayar` and `status_pengiriman = true` | Delivered |

### Admin Order Queue Filters

Extend:

```text
GET /api/v1/pembelians?_start=1&_end=20&queue=<queue>
```

Supported queues:

| Queue | Filter |
|---|---|
| `pending_payment` | `status_pembayaran IN ('belum_dibayar', 'pending')` and no uploaded proof |
| `proof_uploaded` | `status_pembayaran = 'belum_dibayar'` and `bukti_pembayaran != ''` |
| `waiting_delivery` | `status_pembayaran IN ('success', 'dibayar')` and `status_pengiriman = false` |
| `delivered` | `status_pengiriman = true` |
| `failed` | `status_pembayaran IN ('deny', 'failure')` |
| `review` | `status_pembayaran = 'challange'` |

Default behavior without `queue` should stay compatible with the current paginated list.

### Checkout Validation

Before creating a purchase through `POST /pembelian` or `POST /new/pembelian`:

1. Reject `jumlah_dl <= 0`.
2. Reject if latest stock is missing.
3. Reject if requested `jumlah_dl` exceeds latest `stock_dl`.
4. Return a clear error message such as `insufficient stock`.
5. Keep final stock mutation atomic when payment is confirmed or delivery is marked, as currently required by the stock transaction rules.

### Open Implementation Check

Manual purchase creation should be reviewed before coding. The service has `CreateDataPembelianManual`, but the current handler path for `POST /new/pembelian` calls `CreateDataPembelian`. Phase 1 should confirm whether this is intentional or switch the manual route to the manual creation method so total pricing is computed consistently.
