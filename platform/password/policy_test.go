package password

import (
	"testing"
)

func TestHashAndVerifyWithPepper(t *testing.T) {
	password := "senha-muito-segura-123"
	
	// Configura um pepper para o teste
	originalPepper := GlobalPepper
	GlobalPepper = "segredo-do-servidor-123"
	defer func() { GlobalPepper = originalPepper }()

	// Teste de Hash
	hash, err := Hash(password)
	if err != nil {
		t.Fatalf("Erro ao gerar hash: %v", err)
	}

	// Teste de Verificação Sucesso com o mesmo pepper
	err = Verify(password, hash)
	if err != nil {
		t.Errorf("Erro ao verificar senha correta com pepper correto: %v", err)
	}

	// Teste de Verificação Falha com pepper diferente
	GlobalPepper = "outro-segredo-errado"
	err = Verify(password, hash)
	if err == nil {
		t.Error("Deveria ter falhado pois o pepper mudou")
	}

	// Restaura e testa falha com senha errada
	GlobalPepper = "segredo-do-servidor-123"
	err = Verify("senha-errada-123", hash)
	if err == nil {
		t.Error("Deveria ter retornado erro para senha incorreta")
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{"Curta", "123", ErrTooShort},
		{"Longa", string(make([]byte, 129)), ErrTooLong},
		{"Comum", "administrator", ErrCommonPassword},
		{"Válida", "uma-senha-bem-longa-e-segura", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Validate(tt.password); err != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
