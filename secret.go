// Package secret — случайный секрет для ссылки, cookie или токена и его хеш для базы.
//
// Соли нет намеренно: секрет случайный и длинный, словарь по нему не строится.
package secret

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
)

// Bytes — длина секрета: 32 байта перебору не поддаются, в base64url это 43 символа.
const Bytes = 32

// New даёт секрет и его sha256. Ошибки нет: crypto/rand.Read с Go 1.24 её не
// возвращает — при отказе источника случайности процесс падает, и это правильно.
func New() (string, []byte) {
	buf := make([]byte, Bytes)
	_, _ = rand.Read(buf)

	value := base64.RawURLEncoding.EncodeToString(buf)

	return value, Hash(value)
}

// NewHex — то же, что [New], с хешем строкой hex для хранилищ без bytea.
func NewHex() (string, string) {
	value, hash := New()

	return value, hex.EncodeToString(hash)
}

// Hash — то единственное, что от секрета попадает в базу.
func Hash(value string) []byte {
	sum := sha256.Sum256([]byte(value))

	return sum[:]
}

// HashHex — [Hash] строкой hex.
func HashHex(value string) string {
	return hex.EncodeToString(Hash(value))
}

// Match сверяет предъявленный секрет с записанным хешем за постоянное время.
func Match(value string, hash []byte) bool {
	return subtle.ConstantTimeCompare(Hash(value), hash) == 1
}

// MatchHex — [Match] для хеша строкой hex; неразборчивая строка не совпадает ни с чем.
func MatchHex(value, hash string) bool {
	raw, err := hex.DecodeString(hash)
	if err != nil {
		return false
	}

	return Match(value, raw)
}
