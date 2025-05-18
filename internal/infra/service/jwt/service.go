package jwt

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

func NewToken(sub int64, iat time.Time, exp time.Time) *Token {
	return &Token{Sub: sub, Iat: iat, Exp: exp}
}

type Service struct {
	privateKeyStr string
	publicKeyStr  string
}

func NewService(privateKeyStr, publicKeyStr string) *Service {
	return &Service{
		privateKeyStr: privateKeyStr,
		publicKeyStr:  publicKeyStr,
	}
}

func (s Service) CreateTokenStringWithES256(userId int64) (string, error) {
	privateKey, err := s.parseECDSAPrivateKeyStr()
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

func (s Service) VerifyTokenStringWithES256(tokenString string) (*Token, error) {
	publicKey, err := s.parseECDSAPublicKeyFromPEM()
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

	return NewToken(sub, iat.Time, exp.Time), nil
}

func (s Service) parseECDSAPrivateKeyStr() (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(s.privateKeyStr))
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return nil, errors.New("invalid PEM block for EC PRIVATE KEY")
	}

	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func (s Service) parseECDSAPublicKeyFromPEM() (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(s.publicKeyStr))
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
