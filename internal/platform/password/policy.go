// Package password implementa política de senha seguindo NIST SP 800-63B.
//
// Regras:
//   - Mínimo de 12 caracteres (NIST recomenda ao menos 8; usamos 12 para maior segurança)
//   - Máximo de 128 caracteres
//   - Sem exigência de complexidade (maiúsculas, números, símbolos) — conforme NIST
//   - Verificação contra lista de senhas mais comuns
//   - Hash com bcrypt custo 12
package password

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	minLength  = 12
	maxLength  = 128
	bcryptCost = 12
)

// ErrTooShort é retornado quando a senha tem menos que minLength caracteres.
var ErrTooShort = errors.New("a senha deve ter ao menos 12 caracteres")

// ErrTooLong é retornado quando a senha excede maxLength caracteres.
var ErrTooLong = errors.New("a senha não pode ter mais que 128 caracteres")

// ErrCommonPassword é retornado quando a senha está na lista de senhas comuns.
var ErrCommonPassword = errors.New("a senha é muito comum — escolha uma senha mais única")

// Validate aplica a política NIST 800-63B à senha em texto claro.
// Retorna nil se a senha é aceitável.
func Validate(plain string) error {
	if len(plain) < minLength {
		return ErrTooShort
	}
	if len(plain) > maxLength {
		return ErrTooLong
	}
	if isCommon(plain) {
		return ErrCommonPassword
	}
	return nil
}

// Hash gera o hash bcrypt da senha. Valida a política antes de hashear.
func Hash(plain string) (string, error) {
	if err := Validate(plain); err != nil {
		return "", err
	}
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Verify compara a senha em texto claro com o hash armazenado.
// Retorna nil se correspondem.
func Verify(plain, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}

// isCommon verifica se a senha está na lista de senhas mais comuns.
func isCommon(plain string) bool {
	lower := strings.ToLower(strings.TrimSpace(plain))
	_, found := commonPasswords[lower]
	return found
}

// commonPasswords contém as senhas mais frequentemente usadas (NIST 800-63B §5.1.1.2).
// Lista baseada no projeto SecLists/Passwords/Common-Credentials.
var commonPasswords = map[string]struct{}{
	"password":      {},
	"123456":        {},
	"12345678":      {},
	"1234567890":    {},
	"password123":   {},
	"password1":     {},
	"qwerty":        {},
	"qwerty123":     {},
	"qwertyuiop":    {},
	"abc123":        {},
	"letmein":       {},
	"monkey":        {},
	"dragon":        {},
	"master":        {},
	"1234567":       {},
	"123123":        {},
	"654321":        {},
	"superman":      {},
	"michael":       {},
	"sunshine":      {},
	"princess":      {},
	"iloveyou":      {},
	"welcome":       {},
	"shadow":        {},
	"football":      {},
	"baseball":      {},
	"batman":        {},
	"trustno1":      {},
	"pass@word1":    {},
	"passw0rd":      {},
	"admin":         {},
	"admin123":      {},
	"administrator": {},
	"root":          {},
	"toor":          {},
	"test":          {},
	"guest":         {},
	"login":         {},
	"changeme":      {},
	"qazwsx":        {},
	"123qwe":        {},
	"zxcvbnm":       {},
	"1q2w3e4r":      {},
	"1q2w3e":        {},
	"q1w2e3r4":      {},
	"123456789":     {},
	"12345":         {},
	"1234":          {},
	"111111":        {},
	"000000":        {},
	"666666":        {},
	"888888":        {},
	"112233":        {},
	"121212":        {},
	"696969":        {},
	"11111111":      {},
	"00000000":      {},
	"iloveyou1":     {},
	"mypassword":    {},
	"mypass":        {},
	"hello":         {},
	"hello123":      {},
	"flower":        {},
	"hottie":        {},
	"loveme":        {},
	"zaq12wsx":      {},
	"aaaaaa":        {},
	"abc1234":       {},
	"abc12345":      {},
	"abcdef":        {},
	"abcdefg":       {},
	"password12":    {},
	"senha":         {},
	"senha123":      {},
	"mudar123":      {},
	"mudar@123":     {},
	"brasil":        {},
	"brasil123":     {},
	"flamengo":      {},
	"vasco":         {},
	"corinthians":   {},
	"palmeiras":     {},
}
