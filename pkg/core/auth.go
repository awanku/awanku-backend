package core

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"
)

func BuildOauthToken(secretKey []byte, length int) (*OauthToken, error) {
	accessToken, err := generateHash(length)
	if err != nil {
		return nil, err
	}
	refreshToken, err := generateHash(length)
	if err != nil {
		return nil, err
	}

	accessTokenHash, err := HashHMAC(secretKey, []byte(accessToken))
	if err != nil {
		return nil, err
	}

	refreshTokenHash, err := HashHMAC(secretKey, []byte(refreshToken))
	if err != nil {
		return nil, err
	}

	token := &OauthToken{
		AccessToken:      accessToken,
		AccessTokenHash:  accessTokenHash,
		RefreshToken:     refreshToken,
		RefreshTokenHash: refreshTokenHash,
		ExpiresAt:        time.Now().Add(60 * time.Minute),
	}
	return token, nil
}

func BuildOauthAuthorizationCode(length int) (string, error) {
	buff, err := generateHash(length)
	if err != nil {
		return "", nil
	}
	code := base64.URLEncoding.EncodeToString(buff)
	return code, nil
}

func ValidateHMAC(secretKey, plain, hashed []byte) (bool, error) {
	computed, err := HashHMAC(secretKey, plain)
	if err != nil {
		return false, err
	}
	return hmac.Equal(computed, hashed), nil
}

func HashHMAC(secretKey, plain []byte) ([]byte, error) {
	mac := hmac.New(sha256.New, secretKey)
	_, err := mac.Write(plain)
	if err != nil {
		return nil, err
	}
	return mac.Sum(nil), nil
}

func generateHash(length int) ([]byte, error) {
	buff := make([]byte, length)
	_, err := rand.Read(buff)
	if err != nil {
		return nil, err
	}
	return buff, nil
}
