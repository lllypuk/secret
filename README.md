# secret

Случайный секрет для ссылки входа, cookie сессии или токена устройства и его sha256 для базы.

```go
value, hash := secret.New()        // value — в письмо или cookie, hash — в bytea
value, hashHex := secret.NewHex()  // hash строкой hex для хранилищ без bytea

secret.Match(presented, hash)      // сверка за постоянное время
secret.MatchHex(presented, hashHex)
```

Соли нет намеренно: секрет 32 случайных байта, словарь по нему не строится. Ошибки у `New` нет:
`crypto/rand.Read` с Go 1.24 её не возвращает, при отказе источника случайности процесс падает.

Формат хранения выбирает потребитель: и `bytea`, и `text` в hex — один и тот же хеш.
