// Package secret — случайный секрет для ссылки, cookie или токена, код входа и их хеши для базы.
//
// У [Hash] соли нет намеренно: секрет случайный и длинный, словарь по нему не строится.
package secret

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
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

// maxCodeDigits — 10^18 ещё помещается в int64.
const maxCodeDigits = 18

// NewCode даёт код входа из digits цифр с ведущими нулями; digits вне 1..18 — паника.
func NewCode(digits int) string {
	if digits < 1 || digits > maxCodeDigits {
		panic(fmt.Sprintf("secret: NewCode(%d): длина кода вне 1..18", digits))
	}

	const base = 10

	limit := new(big.Int).Exp(big.NewInt(base), big.NewInt(int64(digits)), nil)
	n, _ := rand.Int(rand.Reader, limit)

	return fmt.Sprintf("%0*d", digits, n.Int64())
}

// HashWith — sha256(salt ‖ value): хеш короткого кода, привязанный к запросу.
func HashWith(salt []byte, value string) []byte {
	// Не append(salt, value...): при свободной ёмкости он пишет в массив вызывающего.
	h := sha256.New()
	h.Write(salt)
	h.Write([]byte(value))

	return h.Sum(nil)
}
