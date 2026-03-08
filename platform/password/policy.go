// Package password implementa política de senha seguindo NIST SP 800-63B.
//
// Regras:
//   - Mínimo de 12 caracteres (NIST recomenda ao menos 8; usamos 12 para maior segurança)
//   - Máximo de 128 caracteres
//   - Sem exigência de complexidade (maiúsculas, números, símbolos) — conforme NIST
//   - Verificação contra lista de senhas mais comuns
//   - Hash com Argon2id (Custo generoso de 1GB, otimizado para performance concorrente)
package password

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	minLength = 12
	maxLength = 128
)

// Parâmetros Argon2id (Custo generoso de 1GB, otimizado para performance concorrente).
const (
	argon2Time    = 16         // Otimizado para 1 iteração (performance)
	argon2Memory  = 512 * 1024 // 0.5 GB (custo generoso de memória contra ataques de GPU)
	argon2Threads = 16         // Paralelismo agressivo (usa goroutines internamente para acelerar a criação)
	argon2KeyLen  = 32
	argon2SaltLen = 16
)

// GlobalPepper é o segredo do lado do servidor usado para aumentar a segurança do hash.
// Deve ser carregado de uma variável de ambiente ou configuração segura.
var GlobalPepper string

// ErrTooShort é retornado quando a senha tem menos que minLength caracteres.
var ErrTooShort = errors.New("a senha deve ter ao menos 12 caracteres")

// ErrTooLong é retornado quando a senha excede maxLength caracteres.
var ErrTooLong = errors.New("a senha não pode ter mais que 128 caracteres")

// ErrCommonPassword é retornado quando a senha está na lista de senhas comuns.
var ErrCommonPassword = errors.New("a senha é muito comum — escolha uma senha mais única")

// ErrInvalidHash é retornado quando o formato do hash é inválido ou incompatível.
var ErrInvalidHash = errors.New("o hash armazenado é inválido")

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

// Hash gera o hash Argon2id da senha incluindo o Pepper global.
// O Argon2id utiliza o parâmetro argon2Threads para disparar goroutines e processar em paralelo,
// garantindo que o custo generoso de memória seja processado de forma performática.
func Hash(plain string) (string, error) {
	if err := Validate(plain); err != nil {
		return "", err
	}

	salt := make([]byte, argon2SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	saltedPlain := plain + GlobalPepper

	// O processamento pesado é paralelizado via goroutines internas do pacote argon2.
	hash := argon2.IDKey([]byte(saltedPlain), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLen)

	// Encode salt e hash para base64.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Formata o hash final: $argon2id$v=19$m=1048576,t=1,p=16$<salt>$<hash>
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, argon2Memory, argon2Time, argon2Threads, b64Salt, b64Hash)

	return encodedHash, nil
}

// Verify compara a senha em texto claro com o hash Argon2id armazenado, usando o Pepper.
func Verify(plain, encodedHash string) error {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return ErrInvalidHash
	}

	if parts[1] != "argon2id" {
		return ErrInvalidHash
	}

	var version int
	_, err := fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil || version != argon2.Version {
		return ErrInvalidHash
	}

	var memory uint32
	var iterations uint32
	var parallelism uint8
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return ErrInvalidHash
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return ErrInvalidHash
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return ErrInvalidHash
	}

	saltedPlain := plain + GlobalPepper

	// Verificação também aproveita o paralelismo para ser rápida.
	comparisonHash := argon2.IDKey([]byte(saltedPlain), salt, iterations, memory, parallelism, uint32(len(hash)))

	if subtle.ConstantTimeCompare(hash, comparisonHash) == 1 {
		return nil
	}

	return errors.New("senha incorreta")
}

// isCommon verifica se a senha está na lista de senhas mais comuns.
func isCommon(plain string) bool {
	lower := strings.ToLower(strings.TrimSpace(plain))
	_, found := commonPasswords[lower]
	return found
}

// commonPasswords contém as senhas mais frequentemente usadas.
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
