package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	t.Run("Successfully hashes password", func(t *testing.T) {
		password := "mySecurePassword123"
		hash, err := HashPassword(password)

		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.NotEqual(t, password, hash) // Hash should be different from original
	})

	t.Run("Same password produces different hashes", func(t *testing.T) {
		password := "testPassword"
		hash1, err1 := HashPassword(password)
		hash2, err2 := HashPassword(password)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, hash1, hash2) // bcrypt uses random salt, so hashes differ
	})

	t.Run("Empty password can be hashed", func(t *testing.T) {
		hash, err := HashPassword("")

		assert.NoError(t, err)
		assert.NotEmpty(t, hash)
	})
}

func TestCheckPasswordHash(t *testing.T) {
	t.Run("Correct password returns true", func(t *testing.T) {
		password := "correctPassword123"
		hash, _ := HashPassword(password)

		result := CheckPasswordHash(password, hash)
		assert.True(t, result)
	})

	t.Run("Incorrect password returns false", func(t *testing.T) {
		password := "correctPassword123"
		hash, _ := HashPassword(password)

		result := CheckPasswordHash("wrongPassword", hash)
		assert.False(t, result)
	})

	t.Run("Empty password with valid hash returns false", func(t *testing.T) {
		password := "somePassword"
		hash, _ := HashPassword(password)

		result := CheckPasswordHash("", hash)
		assert.False(t, result)
	})

	t.Run("Invalid hash format returns false", func(t *testing.T) {
		result := CheckPasswordHash("password", "not-a-valid-bcrypt-hash")
		assert.False(t, result)
	})

	t.Run("Empty hash returns false", func(t *testing.T) {
		result := CheckPasswordHash("password", "")
		assert.False(t, result)
	})
}

func TestGenerateJWT(t *testing.T) {
	// Save original env and restore after tests
	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)

	t.Run("Generates valid token with custom secret", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret-key")

		token, err := GenerateJWT(1, RoleStudent)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		// JWT tokens have 3 parts separated by dots
		assert.Contains(t, token, ".")
	})

	t.Run("Generates valid token with fallback secret", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "")

		token, err := GenerateJWT(1, RoleStudent)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("Different users produce different tokens", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret")

		token1, _ := GenerateJWT(1, RoleStudent)
		token2, _ := GenerateJWT(2, RoleStudent)

		assert.NotEqual(t, token1, token2)
	})

	t.Run("Different roles produce different tokens", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret")

		token1, _ := GenerateJWT(1, RoleStudent)
		token2, _ := GenerateJWT(1, RoleProfessor)

		assert.NotEqual(t, token1, token2)
	})

	t.Run("Token contains valid claims", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret")

		userID := uint(42)
		role := RoleProfessor

		token, err := GenerateJWT(userID, role)
		assert.NoError(t, err)

		// Validate the token and check claims
		claims, err := ValidateJWT(token)
		assert.NoError(t, err)
		assert.Equal(t, userID, claims.UserID)
		assert.Equal(t, role, claims.Role)
	})
}

func TestValidateJWT(t *testing.T) {
	// Save original env and restore after tests
	originalSecret := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", originalSecret)

	t.Run("Validates correct token", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret")

		token, _ := GenerateJWT(1, RoleStudent)
		claims, err := ValidateJWT(token)

		assert.NoError(t, err)
		assert.NotNil(t, claims)
		assert.Equal(t, uint(1), claims.UserID)
		assert.Equal(t, RoleStudent, claims.Role)
	})

	t.Run("Rejects invalid token format", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret")

		claims, err := ValidateJWT("not-a-valid-token")

		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("Rejects token with wrong secret", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "secret1")
		token, _ := GenerateJWT(1, RoleStudent)

		os.Setenv("JWT_SECRET", "secret2")
		claims, err := ValidateJWT(token)

		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("Rejects empty token", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret")

		claims, err := ValidateJWT("")

		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("Token expiration is set correctly", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret")

		token, _ := GenerateJWT(1, RoleStudent)
		claims, err := ValidateJWT(token)

		assert.NoError(t, err)
		assert.NotNil(t, claims.ExpiresAt)

		// Token should expire approximately 24 hours from now
		expectedExpiry := time.Now().Add(24 * time.Hour)
		actualExpiry := claims.ExpiresAt.Time

		// Allow 5 seconds tolerance
		diff := actualExpiry.Sub(expectedExpiry)
		assert.Less(t, diff.Abs(), 5*time.Second)
	})

	t.Run("Token issued at is set correctly", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret")

		beforeGeneration := time.Now().Add(-1 * time.Second) // Add buffer for timing
		token, _ := GenerateJWT(1, RoleStudent)
		afterGeneration := time.Now().Add(1 * time.Second) // Add buffer for timing

		claims, err := ValidateJWT(token)

		assert.NoError(t, err)
		assert.NotNil(t, claims.IssuedAt)

		issuedAt := claims.IssuedAt.Time
		assert.True(t, issuedAt.After(beforeGeneration) || issuedAt.Equal(beforeGeneration),
			"IssuedAt should be after or equal to beforeGeneration")
		assert.True(t, issuedAt.Before(afterGeneration) || issuedAt.Equal(afterGeneration),
			"IssuedAt should be before or equal to afterGeneration")
	})
}

func TestClaimsStructure(t *testing.T) {
	t.Run("Claims contains all required fields", func(t *testing.T) {
		os.Setenv("JWT_SECRET", "test-secret")

		token, _ := GenerateJWT(123, RoleAdmin)
		claims, err := ValidateJWT(token)

		assert.NoError(t, err)
		assert.Equal(t, uint(123), claims.UserID)
		assert.Equal(t, RoleAdmin, claims.Role)
		assert.NotNil(t, claims.ExpiresAt)
		assert.NotNil(t, claims.IssuedAt)
	})
}
