package config

import "github.com/golang-jwt/jwt/v5"

var JWT_KEY = []byte("6202834921c0b3d34c59763bd6def626ace1ce6958f98c03c42ac0b4bad65dcc")

type JWTClaim struct {
	ID       uint
	Username string
	jwt.RegisteredClaims
}
