package local

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"

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
	mailer       ports.Mailer
	tsvc         ports.TokenService
}

func NewLocalProvider(mailer ports.Mailer, repo ports.IdentityRepository, tsvc ports.TokenService) *LocalProvider {
	return &LocalProvider{
		repo:         repo,
		ProviderName: LocalProviderName,
		mailer:       mailer,
		tsvc:         tsvc,
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
		return nil, l.handleSignup(method, name, email, password)
	}

	err = l.handleLogin(method, identity, password)
	if err != nil {
		return nil, err
	}

	return identity, nil
}

func (l *LocalProvider) handleSignup(method, name, email, password string) error {
	if method != "signup" {
		return domain.ErrInvalidProviderInputs
	}
	var user *domain.User
	var err error
	var identity *domain.UserIdentity
	err = l.repo.WithTx(func(repo ports.Repository) error {
		// create user
		user, err = repo.FindUserByEmail(email)
		if err != nil && !errors.Is(err, domain.ErrRecordNotFound) {
			return err
		}

		// save
		if errors.Is(err, domain.ErrRecordNotFound) {
			// create entry
			user = &domain.User{
				Email: email,
				Name:  name,
			}
			err = repo.InsertUser(user)
			if err != nil {
				return err
			}
		}

		// create identity; with verified = false
		passwordHash := []byte(password)
		identity = domain.NewIdentity(l.ProviderName, email)
		// link
		identity.UserID = user.ID
		identity.PasswordHash = &passwordHash
		// save
		err = repo.InsertIdentity(identity)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	// create activation toke
	token, err := l.tsvc.New(domain.UUID, domain.ACTIVATE, user.ID)
	if err != nil {
		return err
	}
	// save activation token
	if err = l.tsvc.Save(token); err != nil {
		return err
	}

	// send activation link to user
	if err := l.sendMail(email, token.TokenString, identity.ID); err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}
	return fmt.Errorf("%w: we are sending an activation link to your inbox", domain.ErrUnverifiedAccount)
}

func (l *LocalProvider) handleLogin(method string, identity *domain.UserIdentity, password string) error {
	if method != "login" {
		return domain.ErrInvalidProviderInputs
	}

	if identity.PasswordHash == nil {
		return domain.ErrInvalidIdentity
	}

	err := bcrypt.CompareHashAndPassword(*identity.PasswordHash, []byte(password))
	if err != nil {
		return domain.ErrInvalidProviderInputs
	}
	return nil
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

func (l *LocalProvider) sendMail(email, token string, identityID uint) error {
	go func() {
		url := fmt.Sprintf("%s/auth/verify?token=%s&identity=%d", config.GetServiceURL(), token, identityID)
		data := struct {
			ActivationURL string
		}{
			ActivationURL: url,
		}
		err := l.mailer.Send(email, "activate_account.tmpl.html", data)
		if err != nil {
			log.Printf("error sending activation email to %s: %v", email, err)
		}
	}()

	return nil
}
