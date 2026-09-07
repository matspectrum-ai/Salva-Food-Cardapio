package identity

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"
)

var (
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrNoActiveMembership      = errors.New("no active membership")
	ErrTenantSelectionRequired = errors.New("tenant selection required")
	ErrTenantAccessDenied      = errors.New("tenant access denied")
	ErrSessionNotFound         = errors.New("session not found")
	ErrSessionExpired          = errors.New("session expired")
	ErrSessionRevoked          = errors.New("session revoked")
	ErrSessionAlreadyExists    = errors.New("session already exists")
)

type PasswordVerifier interface {
	Compare(hash, raw string) error
}

type Session struct {
	ID             string     `json:"id"`
	UserID         string     `json:"user_id"`
	ActiveTenantID string     `json:"active_tenant_id"`
	TokenHash      string     `json:"-"`
	ExpiresAt      time.Time  `json:"expires_at"`
	RevokedAt      *time.Time `json:"revoked_at,omitempty"`
}
type AuthRepository interface {
	FindUserByEmail(context.Context, string) (User, bool, error)
	GetUserByID(context.Context, string) (User, bool, error)
	GetMembership(context.Context, string, string) (Membership, bool, error)
	ListMembershipsByUser(context.Context, string) ([]Membership, error)
	CreateSession(context.Context, Session) error
	GetSessionByTokenHash(context.Context, string) (Session, bool, error)
	RevokeSession(context.Context, string) error
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	TenantID string `json:"tenant_id,omitempty"`
}

type Principal struct {
	User       User       `json:"user"`
	Membership Membership `json:"membership"`
}

type LoginResult struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Principal Principal `json:"principal"`
}

type AuthService struct {
	repo       AuthRepository
	verifier   PasswordVerifier
	newToken   func() (string, error)
	newID      func() string
	now        func() time.Time
	sessionTTL time.Duration
}

func NewAuthService(
	repo AuthRepository,
	verifier PasswordVerifier,
	newToken func() (string, error),
	newID func() string,
	now func() time.Time,
	sessionTTL time.Duration,
) *AuthService {
	return &AuthService{
		repo: repo, verifier: verifier, newToken: newToken,
		newID: newID, now: now, sessionTTL: sessionTTL,
	}
}

func (s *AuthService) Login(ctx context.Context, input LoginInput) (LoginResult, error) {
	user, ok, err := s.repo.FindUserByEmail(ctx, input.Email)
	if err != nil {
		return LoginResult{}, err
	}
	if !ok || s.verifier.Compare(user.PasswordHash, input.Password) != nil {
		return LoginResult{}, ErrInvalidCredentials
	}
	memberships, err := s.repo.ListMembershipsByUser(ctx, user.ID)
	if err != nil {
		return LoginResult{}, err
	}
	membership, err := selectLoginMembership(memberships, strings.TrimSpace(input.TenantID))
	if err != nil {
		return LoginResult{}, err
	}
	rawToken, err := s.newToken()
	if err != nil {
		return LoginResult{}, err
	}
	expiresAt := s.now().Add(s.sessionTTL)
	session := Session{
		ID: s.newID(), UserID: user.ID, ActiveTenantID: membership.TenantID,
		TokenHash: hashToken(rawToken), ExpiresAt: expiresAt,
	}
	if err := s.repo.CreateSession(ctx, session); err != nil {
		return LoginResult{}, err
	}
	return LoginResult{
		Token: rawToken, ExpiresAt: expiresAt,
		Principal: Principal{User: user, Membership: membership},
	}, nil
}

func selectLoginMembership(memberships []Membership, tenantID string) (Membership, error) {
	active := make([]Membership, 0, len(memberships))
	for _, membership := range memberships {
		if membership.Status == MembershipActive {
			active = append(active, membership)
		}
	}
	if tenantID != "" {
		for _, membership := range active {
			if membership.TenantID == tenantID {
				return membership, nil
			}
		}
		return Membership{}, ErrTenantAccessDenied
	}
	if len(active) == 0 {
		return Membership{}, ErrNoActiveMembership
	}
	if len(active) > 1 {
		return Membership{}, ErrTenantSelectionRequired
	}
	return active[0], nil
}

func (s *AuthService) Authenticate(ctx context.Context, rawToken string) (Principal, error) {
	rawToken = strings.TrimSpace(rawToken)
	if rawToken == "" {
		return Principal{}, ErrSessionNotFound
	}
	session, ok, err := s.repo.GetSessionByTokenHash(ctx, hashToken(rawToken))
	if err != nil {
		return Principal{}, err
	}
	if !ok {
		return Principal{}, ErrSessionNotFound
	}
	if session.RevokedAt != nil {
		return Principal{}, ErrSessionRevoked
	}
	if !s.now().Before(session.ExpiresAt) {
		return Principal{}, ErrSessionExpired
	}
	user, ok, err := s.repo.GetUserByID(ctx, session.UserID)
	if err != nil {
		return Principal{}, err
	}
	if !ok {
		return Principal{}, ErrSessionNotFound
	}
	membership, ok, err := s.repo.GetMembership(ctx, session.ActiveTenantID, session.UserID)
	if err != nil {
		return Principal{}, err
	}
	if !ok {
		return Principal{}, ErrSessionNotFound
	}
	return Principal{User: user, Membership: membership}, nil
}

func (s *AuthService) Logout(ctx context.Context, rawToken string) error {
	session, ok, err := s.repo.GetSessionByTokenHash(ctx, hashToken(strings.TrimSpace(rawToken)))
	if err != nil {
		return err
	}
	if !ok {
		return ErrSessionNotFound
	}
	if session.RevokedAt != nil {
		return nil
	}
	return s.repo.RevokeSession(ctx, session.ID)
}

func hashToken(rawToken string) string {
	sum := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(sum[:])
}
