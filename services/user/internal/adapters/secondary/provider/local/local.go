package local

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/Yusufdot101/ripple/services/user/config"
	"github.com/Yusufdot101/ripple/services/user/internal/application/core/domain"
	"github.com/Yusufdot101/ripple/services/user/internal/ports"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const LocalProviderName = "local"

type LocalProvider struct {
	ProviderName string
	repo         ports.IdentityRepository
	Mailer       ports.Mailer
}

func NewLocalProvider(mailer ports.Mailer, repo ports.IdentityRepository) *LocalProvider {
	return &LocalProvider{
		repo:         repo,
		ProviderName: LocalProviderName,
		Mailer:       mailer,
	}
}

func (l *LocalProvider) Authenticate(ctx context.Context, credentials map[string]string) (*domain.UserIdentity, error) {
	method, email, name, password, err := parseCreds(credentials)
	if err != nil {
		return nil, err
	}

	identity, err := l.repo.FindIdentityByProviderAndSub(l.ProviderName, email)
	if err != nil && !errors.Is(err, domain.ErrRecordNotFound) {
		return nil, domain.ErrInvalidProviderInputs
	}

	if errors.Is(err, domain.ErrRecordNotFound) {
		if method != "signup" {
			return nil, domain.ErrInvalidProviderInputs
		}
		// send activation link to user
		if err := l.sendMail(name, email, password); err != nil {
			return nil, fmt.Errorf("error sending email: %w", err)
		}
		return nil, fmt.Errorf("%w: we sent an activation link to your inbox", domain.ErrUnverifiedAccount)
	}

	if method != "login" {
		return nil, domain.ErrInvalidProviderInputs
	}

	if identity.PasswordHash == nil {
		return nil, domain.ErrInvalidIdentity
	}

	err = bcrypt.CompareHashAndPassword(*identity.PasswordHash, []byte(password))
	if err != nil {
		return nil, domain.ErrInvalidProviderInputs
	}

	return identity, nil
}

var emailRX = regexp.MustCompile(
	"^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
)

func isValidEmail(email string) bool {
	return emailRX.MatchString(email)
}

func isValidPassword(password string) bool {
	return len(password) >= 8 && len(password) <= 72
}

func parseCreds(credentials map[string]string) (method, email, name, password string, err error) {
	method, hasMethod := credentials["method"]
	name, hasUsername := credentials["name"]
	email, hasEmail := credentials["email"]
	password, hasPassword := credentials["password"]

	if !hasMethod || method == "" || !hasEmail || email == "" || !hasPassword || password == "" {
		return "", "", "", "", domain.ErrInvalidProviderInputs
	}

	if method == "signup" && (!hasUsername || name == "") {
		return "", "", "", "", domain.ErrInvalidProviderInputs
	}

	if !isValidEmail(email) || !isValidPassword(password) {
		return "", "", "", "", domain.ErrInvalidProviderInputs
	}
	return method, email, name, password, nil
}

type EmailClaims struct {
	jwt.RegisteredClaims
	Name         string
	Email        string
	PasswordHash string
}

func (l *LocalProvider) sendMail(name, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	claims := EmailClaims{
		Name:         name,
		Email:        email,
		PasswordHash: string(hash),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    config.GetJWTIssuer(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.GetEmailVerificationTTL())),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(config.GetJWTSecret())
	if err != nil {
		return fmt.Errorf("sign verification token: %w", err)
	}

	go func() {
		data := struct{ Token string }{Token: tokenString}
		err := l.Mailer.Send(email, "activate_account.tmpl.html", data)
		log.Printf("error sending email: %v", err)
	}()

	return nil
}
