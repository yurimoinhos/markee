package infrastructure

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"github.com/aggi-tech/aggipay/modules/auth/domain"
	"github.com/aggi-tech/aggipay/platform/config"
	gooidc "github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

// ── JWT Provider ─────────────────────────────────────────────────────────────

type Claims struct {
	UserID      string   `json:"user_id"`
	Email       string   `json:"email"`
	Role        string   `json:"role"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

type JWTProvider struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

func NewJWTProvider(cfg config.JWTConfig) *JWTProvider {
	return &JWTProvider{
		secret: []byte(cfg.Secret),
		issuer: cfg.Issuer,
		ttl:    cfg.TTL,
	}
}

func (p *JWTProvider) GenerateToken(user domain.User) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:      user.ID,
		Email:       user.Email,
		Role:        user.Role,
		Roles:       user.Roles,
		Permissions: user.Permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    p.issuer,
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(p.ttl)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(p.secret)
}

func (p *JWTProvider) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("método de criptografia inválido")
		}
		return p.secret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("o token informado é inválido")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("as permissões do token são inválidas para essa requisição")
	}
	if claims.UserID == "" || claims.Email == "" || claims.Role == "" {
		return nil, errors.New("claims obrigatórias ausentes")
	}

	return claims, nil
}

// ── OIDC Provider ─────────────────────────────────────────────────────────────

// OIDCUserInfo contém os dados extraídos do ID token após autenticação bem-sucedida.
type OIDCUserInfo struct {
	Sub       string
	Email     string
	FirstName string
	LastName  string
	Provider  string
}

// OIDCAuthParams contém os parâmetros gerados para iniciar o fluxo OIDC.
type OIDCAuthParams struct {
	AuthURL      string `json:"auth_url"`
	State        string `json:"state"`
	CodeVerifier string `json:"code_verifier"`
}

type stateClaims struct {
	Nonce string `json:"nonce"`
	jwt.RegisteredClaims
}

// OIDCProvider encapsula o fluxo OAuth2/OIDC com PKCE para um provider específico.
type OIDCProvider struct {
	provider    *gooidc.Provider
	oauthConfig oauth2.Config
	jwtSecret   []byte
	providerID  string
}

// NewGoogleOIDCProvider cria um OIDCProvider configurado para o Google.
func NewGoogleOIDCProvider(ctx context.Context, cfg config.GoogleOIDCConfig, jwtSecret string) (*OIDCProvider, error) {
	provider, err := gooidc.NewProvider(ctx, cfg.Issuer)
	if err != nil {
		return nil, err
	}

	oauthConfig := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{gooidc.ScopeOpenID, "email", "profile"},
	}

	return &OIDCProvider{
		provider:    provider,
		oauthConfig: oauthConfig,
		jwtSecret:   []byte(jwtSecret),
		providerID:  "google",
	}, nil
}

// AuthURL gera os parâmetros do fluxo OIDC com PKCE.
func (p *OIDCProvider) AuthURL() (OIDCAuthParams, error) {
	codeVerifier, err := generateRandomBase64(32)
	if err != nil {
		return OIDCAuthParams{}, err
	}

	codeChallenge := computeS256Challenge(codeVerifier)

	nonce, err := generateRandomBase64(16)
	if err != nil {
		return OIDCAuthParams{}, err
	}

	stateJWT, err := p.signState(nonce)
	if err != nil {
		return OIDCAuthParams{}, err
	}

	authURL := p.oauthConfig.AuthCodeURL(
		stateJWT,
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
		oauth2.SetAuthURLParam("nonce", nonce),
	)

	return OIDCAuthParams{AuthURL: authURL, State: stateJWT, CodeVerifier: codeVerifier}, nil
}

// Exchange troca o authorization code pelo OIDCUserInfo.
func (p *OIDCProvider) Exchange(ctx context.Context, code, stateJWT, codeVerifier string) (*OIDCUserInfo, error) {
	if err := p.validateState(stateJWT); err != nil {
		return nil, errors.New("state inválido ou expirado")
	}

	token, err := p.oauthConfig.Exchange(ctx, code,
		oauth2.SetAuthURLParam("code_verifier", codeVerifier),
	)
	if err != nil {
		return nil, errors.New("falha ao trocar authorization code")
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("id_token ausente na resposta do provider")
	}

	verifier := p.provider.Verifier(&gooidc.Config{ClientID: p.oauthConfig.ClientID})
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, errors.New("id_token inválido")
	}

	var c struct {
		Sub        string `json:"sub"`
		Email      string `json:"email"`
		GivenName  string `json:"given_name"`
		FamilyName string `json:"family_name"`
	}
	if err := idToken.Claims(&c); err != nil {
		return nil, errors.New("falha ao decodificar claims do id_token")
	}

	return &OIDCUserInfo{
		Sub:       c.Sub,
		Email:     c.Email,
		FirstName: c.GivenName,
		LastName:  c.FamilyName,
		Provider:  p.providerID,
	}, nil
}

func (p *OIDCProvider) signState(nonce string) (string, error) {
	now := time.Now()
	claims := stateClaims{
		Nonce: nonce,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(5 * time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(p.jwtSecret)
}

func (p *OIDCProvider) validateState(stateJWT string) error {
	_, err := jwt.ParseWithClaims(stateJWT, &stateClaims{}, func(t *jwt.Token) (any, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("algoritmo inválido")
		}
		return p.jwtSecret, nil
	})
	return err
}

func generateRandomBase64(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func computeS256Challenge(verifier string) string {
	h := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(h[:])
}
