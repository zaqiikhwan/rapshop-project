---
name: feature-module-scaffold
description: How to add or modify a feature module (handler/service/repo/model/entity) in the RapsShop Go backend. Use when creating a new endpoint group or extending an existing feature under src/.
---

# Feature module scaffold

RapsShop uses a layered architecture, one package per feature under `src/<feature>/`:

```
handler (HTTP)  →  service (usecase)  →  repo (GORM)  →  MySQL
```

## Steps to add a feature `foo`

1. **Entity** — add the GORM model in `entities/foo.go`. Register it in `database/mysql/db.go`'s `AutoMigrate(...)` list.
2. **Model** — in `model/foo.go` define:
   - Input/Output DTOs: `InputFoo`, `FooDto` (with `json` tags; add `binding:"required"` where appropriate).
   - Interfaces: `FooRepository` and `FooUsecase`. Interfaces live in `model/`, not `entities/`.
3. **Repo** — `src/foo/repo/repo.go`: `func NewFooRepository(db *gorm.DB) model.FooRepository`. All GORM access lives here.
4. **Service** — `src/foo/service/service.go`: `func NewFooUsecase(repo model.FooRepository) model.FooUsecase`. Business logic only.
5. **Handler** — `src/foo/handlers/handlers.go`: `func NewFooHandler(r *gin.RouterGroup, uc model.FooUsecase, jwtMiddleware gin.HandlerFunc)`. Register routes; protect mutating/sensitive routes with `jwtMiddleware`.
6. **Wire in `main.go`**:
   ```go
   fooRepo := fooRepo.NewFooRepository(db)
   fooUsecase := fooSvc.NewFooUsecase(fooRepo)
   fooHandler.NewFooHandler(api, fooUsecase, jwtMiddleware)
   ```

## Conventions
- File names follow Go idiom: the directory/package conveys the layer, so files are **not** prefixed with it. Single-file layer packages are named after the package (`handlers/handlers.go`, `service/service.go`, `repo/repo.go`); the multi-file `model/` and `entities/` packages name files by feature (`model/foo.go`, `entities/foo.go`).
- Constructor names: `New<Feature><Layer>` (`NewFooHandler`, `NewFooUsecase`, `NewFooRepository`). Never reuse another feature's constructor name.
- Handlers bind with `c.ShouldBindJSON(&input)` and return early on error via `utils.FailureOrErrorResponse`.
- Never call `mysql.InitDatabase()` outside `main`/`database/mysql`. The `*gorm.DB` is injected via constructors.
- Use `utils.SuccessResponse` / `utils.FailureOrErrorResponse` for every response.
- Domain language is Indonesian; stay consistent within the feature.

## Verify
`go build ./... && go vet ./... && go test ./...`
