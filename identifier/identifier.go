package identifier

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rizface/quora/stdres"
)

// copy of claim struct account/value/authenticated.go
type (
	Claim struct {
		AccountId string `json:"accountId"`
		Email     string `json:"email"`
		Username  string `json:"username"`
		jwt.RegisteredClaims
	}
	ClaimKeyword string
)

const (
	ClaimKey ClaimKeyword = "claim"
)

func validateTokenForm(splittedToken []string) error {
	if len(splittedToken) != 2 {
		return errors.New("token has invalid segment")
	}

	if splittedToken[0] != "Bearer" {
		return errors.New("token has invalid scheme")
	}

	return nil
}

func getClaim(tokenString string) (*Claim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claim{}, func(t *jwt.Token) (interface{}, error) {
		// this line will be executed when token is not match with the secret
		return []byte(os.Getenv("JWT_ACCESS_SECRET")), nil
	})
	if err != nil {
		return &Claim{}, err
	}

	claim, ok := token.Claims.(*Claim)
	if ok && token.Valid {
		// this line will be executed when token is match with secret
		// but the claim is invalid, for example token had expired
		return claim, nil
	}

	return &Claim{}, errors.New("invalid token")
}

func Identifier(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		splittedToken := strings.Split(r.Header.Get("Authorization"), " ")

		if err := validateTokenForm(splittedToken); err != nil {
			stdres.Writer(w, stdres.Response{
				Code: http.StatusUnauthorized,
				Info: err.Error(),
			})

			return
		}

		var (
			token      = splittedToken[1]
			claim, err = getClaim(token)
		)
		if err != nil {
			stdres.Writer(w, stdres.Response{
				Code: http.StatusUnauthorized,
				Info: err.Error(),
			})

			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ClaimKey, claim)))
	})
}

func GetFromContext(ctx context.Context) (*Claim, error) {
	claim := ctx.Value(ClaimKey)
	if claim == nil {
		return &Claim{}, errors.New("claim not found")
	}

	return claim.(*Claim), nil
}
