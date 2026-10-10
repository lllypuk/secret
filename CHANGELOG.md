# Changelog

## v0.2.0

Только добавляет API.

- `NewCode` — код входа из 1–18 цифр с ведущими нулями, без смещения от остатка.
- `HashWith` — `sha256(salt ‖ value)`: байты совпадают с прежними копиями в DPP, fitness-app и
  family-money, выданные до перехода коды остаются действительными.

## v0.1.0

`New`, `NewHex`, `Hash`, `HashHex`, `Match`, `MatchHex` — случайный секрет и его sha256.
