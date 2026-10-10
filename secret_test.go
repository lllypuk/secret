package secret_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"strings"
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

// TestNewCode — ровно digits цифр, ведущие нули не теряются.
func TestNewCode(t *testing.T) {
	t.Parallel()

	for _, digits := range []int{1, 8, 18} {
		code := secret.NewCode(digits)
		if len(code) != digits || strings.Trim(code, "0123456789") != "" {
			t.Errorf("NewCode(%d) = %q", digits, code)
		}
	}

	for range 1000 {
		if code := secret.NewCode(2); code[0] == '0' {
			return
		}
	}

	t.Error("1000 кодов из 2 цифр без ведущего нуля")
}

// TestNewCodePanics — длина вне 1..18 — ошибка программиста.
func TestNewCodePanics(t *testing.T) {
	t.Parallel()

	for _, digits := range []int{0, 19} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("NewCode(%d) не запаниковал", digits)
				}
			}()

			secret.NewCode(digits)
		}()
	}
}

// TestHashWith — байты совпадают с копиями в потребителях, salt не портится.
func TestHashWith(t *testing.T) {
	t.Parallel()

	salt := make([]byte, 16, 64)
	salt[0] = 1
	code := "01234567"

	want := sha256.Sum256(append(bytes.Clone(salt), code...))
	if got := secret.HashWith(salt, code); !bytes.Equal(got, want[:]) {
		t.Errorf("HashWith = %x, ожидался %x", got, want)
	}

	if spare := salt[:cap(salt)]; !bytes.Equal(spare[len(salt):], make([]byte, cap(salt)-len(salt))) {
		t.Error("HashWith записал в свободную ёмкость salt")
	}

	other := bytes.Clone(salt)
	other[0] = 2

	if bytes.Equal(secret.HashWith(salt, code), secret.HashWith(other, code)) {
		t.Error("разная соль дала одинаковый хеш")
	}
}
