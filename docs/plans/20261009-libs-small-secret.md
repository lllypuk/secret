# `secret` v0.2.0: код входа

## Overview

Выпуск только добавляет API.

- **`secret`**: код из 8 цифр и его хеш `sha256(request_id ‖ code)` написаны трижды —
  DPP `internal/app/identity/auth.go:379-395`, fitness-app `internal/app/identity/auth.go:569-584`,
  family-money `internal/auth/code.go:29-42`. Копии побайтно совпадают по результату (DPP и fitness
  `sha256.Sum256(append(id[:], code...))`, family-money — два `Write` в один хеш).

Не делаем:
- `secret.MatchWith`: потребители сверяют хеш кода в SQL (`where code_hash = $1`).

## Решения (приняты 2026-10-09, не переигрывать)

1. **`secret.NewCode(digits int) string`** — `rand.Int(rand.Reader, 10^digits)` с ведущими нулями:
   без смещения от остатка, 10^18 помещается в int64. `digits` вне `1..18` — паника: это ошибка
   программиста, а ошибки у `New` нет по той же причине (`crypto/rand` с Go 1.24 не отказывает).
2. **`secret.HashWith(salt []byte, value string) []byte`** = `sha256(salt ‖ value)` — байты
   совпадают со всеми тремя копиями, поэтому выданные до перехода коды остаются действительными.
   Не `append(salt, …)`: при свободной ёмкости он пишет в массив вызывающего; `sha256.New()` и два
   `Write`, как в family-money.
   `HashWithHex` не заводится: все три потребителя хранят `bytea`.
3. README `secret`: соль не спасает 8 цифр от перебора по утёкшей базе — защита кода в сроке и лимите
   попыток потребителя. Одной строкой.
4. **CHANGELOG** заводится: `## v0.2.0` и строка о v0.1.0; раздел `Ломает:` —
   с первой ломающей версии (правило — как в `llm`).

## Context (from discovery)

- `secret.go` — `New`, `NewHex`, `Hash`, `HashHex`, `Match`, `MatchHex`; README с примерами.
- `go 1.26`, stdlib only, CI — `.github/workflows/ci.yml`, общий `.golangci.yml`.

## Development Approach

- **testing approach**: TDD.
- Перед каждым коммитом: `go vet ./... && go test -race ./...` и
  `go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run` — 0 замечаний.
- Коммиты на русском: `feat:`, `docs:`; комментарии — по `~/.claude/CLAUDE.md`.

## Testing Strategy

- Юнит-тесты без сети и файлов вне `testdata`.

## Progress Tracking

`[x]` — сделано, ➕ — найдено по ходу, ⚠️ — блокер.

## Implementation Steps

### Task 1: `NewCode` и `HashWith`

**Files:** `secret.go`, `secret_test.go`

- [x] `NewCode`, `HashWith` по Решениям 1–2
- [x] тесты: длина и только цифры для 1, 8, 18; ведущие нули встречаются (1000 кодов из 2 цифр
  дают значение `< 10`); паника на 0 и 19; `HashWith(id, code)` совпадает с
  `sha256.Sum256(append(id, code...))`; разная соль — разные хеши; `salt` после вызова не меняется
  (срез со свободной ёмкостью)
- [x] `go test -race ./...` — зелёный
- [x] коммит `feat: NewCode и HashWith — код входа и его хеш с солью`

### Task 2: Verify acceptance criteria

- [ ] импорты — только stdlib; `go vet`, тесты, линтер — зелёные

### Task 3: [Final] Документы

**Files:** `README.md`, `CHANGELOG.md`

- [ ] README: пример кода входа (`NewCode(8)`, `HashWith(requestID[:], code)`) и строка Решения 3
- [ ] `CHANGELOG.md` по Решению 4
- [ ] коммит `docs: v0.2.0 — код входа в README и CHANGELOG`
