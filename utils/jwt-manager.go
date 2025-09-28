package utils

import (
	"crypto/rsa"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtManager struct {
	RsaPrivateKey      *rsa.PrivateKey
	RsaPublicKey       *rsa.PublicKey
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

func NewJwtManager() *JwtManager {
	var jwtManager JwtManager
	privateKey := getKey("private.pem")
	publicKey := getKey("public.pem")

	rsaPrivKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		log.Fatalf("failed to create private key: %v", err)
	}

	rsaPubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		log.Fatalf("failed to create private key: %v", err)
	}

	jwtManager.RsaPrivateKey = rsaPrivKey
	jwtManager.RsaPublicKey = rsaPubKey
	jwtManager.AccessTokenExpiry = (15 * time.Minute) - 1  /// access token expiry in ms slightly less than expiry time
	jwtManager.RefreshTokenExpiry = (60 * time.Minute) - 1 /// referes token expiry in ms
	return &jwtManager
}

func getKey(fileName string) string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)

	}
	keyPath := filepath.Dir(ex)
	filePath := filepath.Join(keyPath, "keys", fileName)
	key, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	return string(key)
}

func (j *JwtManager) CreateAccessToken(username string, iss time.Time) (string, error) {

	claims := jwt.MapClaims{
		"sub": username,
		"iss": "http://example-xyz.com",
		"aud": "clients",
		"iat": iss.Unix(),
		"nbf": iss.Unix(),
		"exp": iss.Add(j.AccessTokenExpiry).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "v1"

	return token.SignedString(j.RsaPrivateKey)
}

func (j *JwtManager) CreateRefreshToken(username string, iss time.Time) (string, error) {

	claims := jwt.MapClaims{
		"typ": "refresh",
		"sub": username,
		"iss": "http://example-xyz.com",
		"aud": "clients",
		"iat": iss.Unix(),
		"nbf": iss.Unix(),
		"exp": iss.Add(j.RefreshTokenExpiry).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = "v1"

	return token.SignedString(j.RsaPrivateKey)
}
