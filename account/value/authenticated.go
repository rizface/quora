package value

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Authenticated struct {
	Id       string  `json:"login"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Tokens   []Token `json:"tokens"`
}

type Claim struct {
	AccountId string
	Email     string
	Usernamae string
	jwt.RegisteredClaims
}

func getTokens(a Authenticated) ([]Token, error) {
	var (
		accessSecret  = []byte(os.Getenv("JWT_ACCESS_SECRET"))
		refreshSecret = []byte(os.Getenv("JWT_REFRESH_SECRET"))
		tokens        = []Token{}
		claim         = Claim{
			AccountId: a.Id,
			Email:     a.Email,
			Usernamae: a.Username,
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt: jwt.NewNumericDate(time.Now()),
				Issuer:   "quora",
			},
		}
		AccessType  = "access"
		RefreshType = "refresh"
	)

	var generateToken = func(tokenType string, secret []byte, expires time.Time) (string, error) {
		claim.ExpiresAt = jwt.NewNumericDate(time.Now().Add(24 * time.Hour))
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

		tokenString, err := token.SignedString(secret)
		if err != nil {
			return "", err
		}

		return tokenString, nil
	}

	token, err := generateToken(AccessType, accessSecret, time.Now().Add(24*time.Hour))
	if err != nil {
		return tokens, err
	}

	tokens = append(tokens, Token{
		Type:  AccessType,
		Value: token,
	})

	token, err = generateToken(RefreshType, refreshSecret, time.Now().Add(24*time.Hour*90))
	if err != nil {
		return tokens, err
	}

	tokens = append(tokens, Token{
		Type:  RefreshType,
		Value: token,
	})

	return tokens, nil
}

func NewAuthenticated(e AccountEntity) (Authenticated, error) {
	a := Authenticated{
		Id:       e.Id,
		Username: e.Username,
		Email:    e.Email,
	}

	tokens, err := getTokens(a)
	if err != nil {
		return Authenticated{}, err
	}

	a.Tokens = tokens

	return a, nil
}
