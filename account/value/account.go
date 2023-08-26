package value

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AccountPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AccountEntity struct {
	Id             string    `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Password       string    `json:"-"`
	EmailConfirmed bool      `json:"emailConfirmed"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	VerifiedAt     time.Time `json:"verifiedAt"`
}

func NewAccountEntity(p AccountPayload) AccountEntity {
	return AccountEntity{
		Id:             uuid.NewString(),
		Username:       p.Username,
		Email:          p.Email,
		Password:       p.Password,
		EmailConfirmed: false,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

func (a AccountEntity) Validate() error {
	return validation.Errors{
		"username": validation.Validate(a.Username, validation.Required),
		"email":    validation.Validate(a.Email, validation.Required, is.Email),
		"password": validation.Validate(a.Password, validation.Required, validation.Length(8, 0)),
	}.Filter()
}

func (a AccountEntity) GetPasswordHash() (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
