package secret_test

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/lllypuk/secret"
)

// TestNew — в базу уезжает хеш, а не секрет, и два секрета подряд разные.
func TestNew(t *testing.T) {
	t.Parallel()

	value, hash := secret.New()
	if value == "" || len(hash) != sha256.Size {
		t.Fatalf("секрет %q, хеш %d байт — ожидались непустое значение и sha256", value, len(hash))
	}

	if sum := sha256.Sum256([]byte(value)); string(hash) != string(sum[:]) {
		t.Error("хеш не сходится с sha256 самого секрета")
	}

	if second, _ := secret.New(); second == value {
		t.Error("два секрета подряд совпали — источник случайности не работает")
	}
}

// TestHex — hex-форма кодирует тот же хеш, что и бинарная.
func TestHex(t *testing.T) {
	t.Parallel()

	value, hashHex := secret.NewHex()
	if hashHex != hex.EncodeToString(secret.Hash(value)) || hashHex != secret.HashHex(value) {
		t.Errorf("hex-хеш %q расходится с бинарным", hashHex)
	}
}

// TestMatch — совпадает только сам секрет, и только с годным хешем.
func TestMatch(t *testing.T) {
	t.Parallel()

	value, hash := secret.New()
	other, _ := secret.New()

	if !secret.Match(value, hash) || !secret.MatchHex(value, hex.EncodeToString(hash)) {
		t.Error("секрет не сошёлся с собственным хешем")
	}

	if secret.Match(other, hash) || secret.MatchHex(other, hex.EncodeToString(hash)) {
		t.Error("чужой секрет сошёлся с хешем")
	}

	if secret.MatchHex(value, "не hex") || secret.Match(value, nil) {
		t.Error("негодный хеш сошёлся с секретом")
	}
}
