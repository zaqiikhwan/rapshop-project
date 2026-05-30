---
name: gorm-mysql-data-access
description: GORM + MySQL conventions for RapsShop repos — atomic stock mutation, pagination, safe WHERE/OR clauses, and counting. Use when writing or reviewing repo code under src/**/repo.
---

# GORM / MySQL data access

All DB access lives in `src/<feature>/repo`. The `*gorm.DB` is injected via the repo constructor.

## Atomic, non-negative stock
Never read-modify-write stock in app code. Do it in one transaction with a row lock and a guard:
```go
err := db.Transaction(func(tx *gorm.DB) error {
    var s entities.StockDL
    if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Order("id desc").First(&s).Error; err != nil {
        return err
    }
    newStock := s.StockDL + delta
    if newStock < 0 {
        return errors.New("insufficient stock")
    }
    if err := tx.Model(&entities.StockDL{}).Where("id = ?", s.ID).Updates(prices).Error; err != nil { // prices: zero fields skipped
        return err
    }
    return tx.Model(&entities.StockDL{}).Where("id = ?", s.ID).UpdateColumn("stock_dl", newStock).Error // force even when 0
})
```
`Updates(struct)` skips zero-value fields (so partial updates don't clobber prices); use `UpdateColumn`/a map when you must write a zero.

## WHERE / OR precedence
`AND` binds tighter than `OR`. Always parenthesize OR groups, or the filters get bypassed:
```sql
WHERE created_at LIKE ? AND harga_beli = ? AND (status_pembayaran = 'success' OR status_pembayaran = 'dibayar')
```
Use `.Order("col asc")` — never embed `ORDER BY` inside a `Where(...)` string.

## Counting & pagination
Count with `db.Model(&Entity{}).Count(&total)`, not by loading all ids. Guard the existing `_start`/`_end` contract: clamp `_start` to `>= 1` and ensure the limit isn't negative before `Offset(_start-1).Limit(_end-_start+1)`.

## Migrations
Schema is created by `AutoMigrate` at startup, gated by `AUTO_MIGRATE` (set `false` in prod and use versioned migrations).
