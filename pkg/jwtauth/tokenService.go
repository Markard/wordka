package jwtauth

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
	"time"
)

type Token struct {
	Sub int64
	Exp time.Time
	Iat time.Time
}

func NewToken(token *jwt.Token) (*Token, error) {
	subStr, err := token.Claims.GetSubject()
	if err != nil {
		return nil, err
	}
	sub, err := strconv.ParseInt(subStr, 10, 64)
	if err != nil {
		return nil, err
	}

	exp, err := token.Claims.GetExpirationTime()
	if err != nil {
		return nil, err
	}

	iat, err := token.Claims.GetIssuedAt()
	if err != nil {
		return nil, err
	}

	return &Token{
		Sub: sub,
		Iat: iat.Time,
		Exp: exp.Time,
	}, nil
}

type TokenService struct {
	privateKeyStr string
	publicKeyStr  string
}

func NewTokenService(privateKeyStr, publicKeyStr string) *TokenService {
	return &TokenService{
		privateKeyStr: privateKeyStr,
		publicKeyStr:  publicKeyStr,
	}
}

func (j *TokenService) CreateTokenStringWithES256(userId int64) (string, error) {
	privateKey, err := j.parseECDSAPrivateKeyStr()
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"sub": strconv.FormatInt(userId, 10),
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // expiration date
		"iat": time.Now().Unix(),                         // creation date
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *TokenService) VarifyTokenStringWithES256(tokenString string) (*Token, error) {
	publicKey, err := j.parseECDSAPublicKeyFromPEM()
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok || token.Method.Alg() != jwt.SigningMethodES256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	result, err := NewToken(token)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (j *TokenService) parseECDSAPrivateKeyStr() (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(j.privateKeyStr))
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return nil, errors.New("invalid PEM block for EC PRIVATE KEY")
	}

	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (j *TokenService) parseECDSAPublicKeyFromPEM() (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(j.publicKeyStr))
	if block == nil || block.Type != "EC PUBLIC KEY" {
		return nil, errors.New("invalid PEM block for EC PUBLIC KEY")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	ecdsaPub, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("invalid ECDSA public key")
	}

	return ecdsaPub, nil
}
