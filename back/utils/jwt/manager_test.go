package jwt

import (
	"gin-template/config"
	"gin-template/pkg/model/enum"
	"testing"
	"time"
)

// TestGenerateTokens tests the generation of tokens
func TestGenerateTokens(t *testing.T) {
	jwtConfig := config.JwtConfig{
		Secret:     "this_is_a_secret",
		Expiration: 24,
	}
	tokens, err := GenerateTokens(1, enum.SUPERADMIN, jwtConfig)
	if err != nil {
		t.Error(err)
	}

	if tokens.AccessToken == "" {
		t.Error("Access token is empty")
	}

	if tokens.RefreshToken == "" {
		t.Error("Refresh token is empty")
	}

	if tokens.ExpiresIn.IsZero() {
		t.Error("Expires in is empty")
	}

	if tokens.TokenID == "" {
		t.Error("Token ID is empty")
	}

	if tokens.AccessToken == tokens.RefreshToken {
		t.Error("Access token and refresh token are the same")
	}

	if time.Now().After(tokens.ExpiresIn) {
		t.Error("Expires in is in the past")
	}

	if tokens.ExpiresIn.Before(time.Now().Add(time.Duration(jwtConfig.Expiration)*time.Hour - time.Minute)) {
		t.Error("Expires in is too far in the future")
	}
}

// TestParseToken tests the parsing of a token
func TestParseToken(t *testing.T) {
	jwtConfig := config.JwtConfig{
		Secret:     "this_is_a_secret",
		Expiration: 24,
	}
	tokens, err := GenerateTokens(1, enum.SUPERADMIN, jwtConfig)
	if err != nil {
		t.Error(err)
	}

	claims, err := ParseToken(tokens.AccessToken, jwtConfig.Secret)
	if err != nil {
		t.Error(err)
	}

	if claims.UserId != 1 {
		t.Error("User ID is not 1")
	}

	if claims.Role != enum.SUPERADMIN {
		t.Error("Role is not SUPERADMIN")
	}

	if claims.ID != tokens.TokenID {
		t.Error("Token ID is not the same")
	}
}

// TestParseTokenWithInvalidToken tests the parsing of a token with an invalid token
func TestParseTokenWithInvalidToken(t *testing.T) {
	jwtConfig := config.JwtConfig{
		Secret:     "this_is_a_secret",
		Expiration: 24,
	}
	_, err := ParseToken("invalid_token", jwtConfig.Secret)
	if err == nil {
		t.Error("Expected error")
	}
}

// TestParseTokenWithInvalidSecret tests the parsing of a token with an invalid secret
func TestParseTokenWithInvalidSecret(t *testing.T) {
	jwtConfig := config.JwtConfig{
		Secret:     "this_is_a_secret",
		Expiration: 24,
	}
	tokens, err := GenerateTokens(1, enum.SUPERADMIN, jwtConfig)
	if err != nil {
		t.Error(err)
	}

	_, err = ParseToken(tokens.AccessToken, "invalid_secret")
	if err == nil {
		t.Error("Expected error")
	}
}

// TestRefreshTokenWithToken tests the refresh token with token
func TestRefreshTokenWithToken(t *testing.T) {
	jwtConfig := config.JwtConfig{
		Secret:     "this_is_a_secret",
		Expiration: 24,
	}
	tokens, err := GenerateTokens(1, enum.SUPERADMIN, jwtConfig)
	if err != nil {
		t.Error(err)
	}

	err = tokens.RefreshTokenWithToken(tokens.RefreshToken, jwtConfig)
	if err != nil {
		t.Error(err)
	}

	if tokens.AccessToken == "" {
		t.Error("Access token is empty")
	}

	if tokens.RefreshToken == "" {
		t.Error("Refresh token is empty")
	}

	if tokens.ExpiresIn.IsZero() {
		t.Error("Expires in is empty")
	}

	if tokens.TokenID == "" {
		t.Error("Token ID is empty")
	}

	if tokens.AccessToken == tokens.RefreshToken {
		t.Error("Access token and refresh token are the same")
	}

	if time.Now().After(tokens.ExpiresIn) {
		t.Error("Expires in is in the past")
	}

	if tokens.ExpiresIn.Before(time.Now().Add(time.Duration(jwtConfig.Expiration)*time.Hour - time.Minute)) {
		t.Error("Expires in is too far in the future")
	}
}

// BenchmarkGenerateTokens benchmarks the generation of tokens
func BenchmarkGenerateTokens(b *testing.B) {
	jwtConfig := config.JwtConfig{
		Secret:     "this_is_a_secret",
		Expiration: 24,
	}
	for i := 0; i < b.N; i++ {
		GenerateTokens(1, enum.SUPERADMIN, jwtConfig)
	}
}

// BenchmarkParseToken benchmarks the parsing of a token
func BenchmarkParseToken(b *testing.B) {
	jwtConfig := config.JwtConfig{
		Secret:     "this_is_a_secret",
		Expiration: 24,
	}
	tokens, _ := GenerateTokens(1, enum.SUPERADMIN, jwtConfig)
	for i := 0; i < b.N; i++ {
		ParseToken(tokens.AccessToken, jwtConfig.Secret)
	}
}

// BenchmarkRefreshTokenWithToken benchmarks the refresh token with token
func BenchmarkRefreshTokenWithToken(b *testing.B) {
	jwtConfig := config.JwtConfig{
		Secret:     "this_is_a_secret",
		Expiration: 24,
	}
	tokens, _ := GenerateTokens(1, enum.SUPERADMIN, jwtConfig)
	for i := 0; i < b.N; i++ {
		tokens.RefreshTokenWithToken(tokens.RefreshToken, jwtConfig)
	}
}
