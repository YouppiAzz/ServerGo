package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type JWTClaims struct {
	UserID int64 `json:"user_id"`
	Exp    int64 `json:"exp"`
	Iat    int64 `json:"iat"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(userID int64, secret string) (string, error) {
	header := map[string]interface{}{
		"alg": "HS256",
		"typ": "JWT",
	}

	now := time.Now()
	claims := JWTClaims{
		UserID: userID,
		Exp:    now.Add(time.Hour * 24).Unix(),
		Iat:    now.Unix(),
	}

	headerJSON, _ := json.Marshal(header)
	claimsJSON, _ := json.Marshal(claims)

	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)
	claimsB64 := base64.RawURLEncoding.EncodeToString(claimsJSON)

	message := headerB64 + "." + claimsB64

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	signature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	return message + "." + signature, nil
}

func ValidateToken(tokenString, secret string) (int64, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid token format")
	}

	message := parts[0] + "." + parts[1]
	signature := parts[2]

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(message))
	expectedSignature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return 0, fmt.Errorf("invalid signature")
	}

	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid claims encoding")
	}

	var claims JWTClaims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return 0, fmt.Errorf("invalid claims format")
	}

	if time.Now().Unix() > claims.Exp {
		return 0, fmt.Errorf("token expired")
	}

	return claims.UserID, nil
}

// RefreshToken generates a new token with extended expiration
func RefreshToken(tokenString, secret string) (string, error) {
	userID, err := ValidateToken(tokenString, secret)
	if err != nil {
		return "", err
	}

	return GenerateToken(userID, secret)
}
