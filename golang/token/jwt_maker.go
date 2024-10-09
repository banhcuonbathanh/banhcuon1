package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) *JWTMaker {
	return &JWTMaker{secretKey}
}

func (maker *JWTMaker) CreateToken(id int64, email string, role string, duration time.Duration) (string, *UserClaims, error) {
	claims, err := NewUserClaims(id, email, role, duration)
	if err != nil {
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(maker.secretKey))
	if err != nil {
		return "", nil, fmt.Errorf("error signing token: %w", err)
	}

	return tokenStr, claims, nil
}

func (maker *JWTMaker) VerifyToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// verify the signing method
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}

		return []byte(maker.secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token: %w", err)
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// func (maker *JWTMaker) VerifyToken(token string) (*UserClaims, error) {
//     keyFunc := func(token *jwt.Token) (interface{}, error) {
//         _, ok := token.Method.(*jwt.SigningMethodHMAC)
//         if !ok {
//             return nil, ErrInvalidToken
//         }
//         return []byte(maker.secretKey), nil
//     }

//     jwtToken, err := jwt.ParseWithClaims(token, &UserClaims{}, keyFunc)
//     if err != nil {
//         verr, ok := err.(*jwt.ValidationError)
//         if ok && errors.Is(verr.Inner, ErrExpiredToken) {
//             return nil, ErrExpiredToken
//         }
//         return nil, ErrInvalidToken
//     }

//     claims, ok := jwtToken.Claims.(*UserClaims)
//     if !ok {
//         return nil, ErrInvalidToken
//     }

//     return claims, nil
// }