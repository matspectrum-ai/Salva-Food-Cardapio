package identity

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	ErrUserAlreadyExists       = errors.New("user already exists")
	ErrMembershipAlreadyExists = errors.New("tenant membership already exists")
	ErrCollaboratorNotFound    = errors.New("collaborator not found")
)

type Collaborator struct {
	User       User       `json:"user"`
	Membership Membership `json:"membership"`
}

type Repository interface {
	CreateCollaborator(context.Context, Collaborator) error
	FindUserByEmail(context.Context, string) (User, bool, error)
	GetMembership(context.Context, string, string) (Membership, bool, error)
	ListCollaborators(context.Context, string, string) ([]Collaborator, error)
	SetMembershipStatus(context.Context, string, string, MembershipStatus) error
}
type MemoryStore struct {
	mu               sync.RWMutex
	users            map[string]User
	userByEmail      map[string]string
	memberships      map[string]map[string]Membership
	sessionsByToken  map[string]Session
	sessionTokenByID map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		users:            make(map[string]User),
		userByEmail:      make(map[string]string),
		memberships:      make(map[string]map[string]Membership),
		sessionsByToken:  make(map[string]Session),
		sessionTokenByID: make(map[string]string),
	}
}

func (s *MemoryStore) CreateCollaborator(_ context.Context, collaborator Collaborator) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.users[collaborator.User.ID]; exists {
		return ErrUserAlreadyExists
	}
	email := strings.ToLower(collaborator.User.Email)
	if _, exists := s.userByEmail[email]; exists {
		return ErrUserAlreadyExists
	}
	if s.memberships[collaborator.Membership.TenantID] == nil {
		s.memberships[collaborator.Membership.TenantID] = make(map[string]Membership)
	}
	if _, exists := s.memberships[collaborator.Membership.TenantID][collaborator.User.ID]; exists {
		return ErrMembershipAlreadyExists
	}
	s.users[collaborator.User.ID] = cloneUser(collaborator.User)
	s.userByEmail[email] = collaborator.User.ID
	s.memberships[collaborator.Membership.TenantID][collaborator.User.ID] = cloneMembership(collaborator.Membership)
	return nil
}
func (s *MemoryStore) FindUserByEmail(_ context.Context, email string) (User, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.userByEmail[strings.ToLower(strings.TrimSpace(email))]
	if !ok {
		return User{}, false, nil
	}
	return cloneUser(s.users[id]), true, nil
}

func (s *MemoryStore) GetMembership(_ context.Context, tenantID, userID string) (Membership, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	membership, ok := s.memberships[tenantID][userID]
	if !ok {
		return Membership{}, false, nil
	}
	return cloneMembership(membership), true, nil
}

func (s *MemoryStore) ListCollaborators(_ context.Context, tenantID, search string) ([]Collaborator, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	needle := strings.ToLower(strings.TrimSpace(search))
	result := make([]Collaborator, 0, len(s.memberships[tenantID]))
	for userID, membership := range s.memberships[tenantID] {
		user, ok := s.users[userID]
		if !ok {
			continue
		}
		if needle != "" && !matchesCollaborator(user, membership, needle) {
			continue
		}
		result = append(result, Collaborator{User: cloneUser(user), Membership: cloneMembership(membership)})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].User.Name == result[j].User.Name {
			return result[i].User.ID < result[j].User.ID
		}
		return result[i].User.Name < result[j].User.Name
	})
	return result, nil
}
func (s *MemoryStore) SetMembershipStatus(_ context.Context, tenantID, userID string, status MembershipStatus) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	membership, ok := s.memberships[tenantID][userID]
	if !ok {
		return ErrCollaboratorNotFound
	}
	if err := membership.SetStatus(status); err != nil {
		return err
	}
	s.memberships[tenantID][userID] = membership
	return nil
}

func matchesCollaborator(user User, membership Membership, needle string) bool {
	fields := []string{user.Name, user.Email, user.Phone, membership.Title}
	for _, field := range fields {
		if strings.Contains(strings.ToLower(field), needle) {
			return true
		}
	}
	return false
}

func cloneUser(user User) User {
	return user
}

func cloneMembership(membership Membership) Membership {
	clone := membership
	clone.Permissions = append([]Permission(nil), membership.Permissions...)
	return clone
}

func (s *MemoryStore) GetUserByID(_ context.Context, userID string) (User, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.users[userID]
	if !ok {
		return User{}, false, nil
	}
	return cloneUser(user), true, nil
}

func (s *MemoryStore) ListMembershipsByUser(_ context.Context, userID string) ([]Membership, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Membership, 0)
	for _, tenantMemberships := range s.memberships {
		if membership, ok := tenantMemberships[userID]; ok {
			result = append(result, cloneMembership(membership))
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].TenantID < result[j].TenantID })
	return result, nil
}

func (s *MemoryStore) CreateSession(_ context.Context, session Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.sessionTokenByID[session.ID]; exists {
		return ErrSessionAlreadyExists
	}
	if _, exists := s.sessionsByToken[session.TokenHash]; exists {
		return ErrSessionAlreadyExists
	}
	stored := cloneSession(session)
	s.sessionsByToken[session.TokenHash] = stored
	s.sessionTokenByID[session.ID] = session.TokenHash
	return nil
}

func (s *MemoryStore) GetSessionByTokenHash(_ context.Context, tokenHash string) (Session, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, ok := s.sessionsByToken[tokenHash]
	if !ok {
		return Session{}, false, nil
	}
	return cloneSession(session), true, nil
}

func (s *MemoryStore) RevokeSession(_ context.Context, sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	tokenHash, ok := s.sessionTokenByID[sessionID]
	if !ok {
		return ErrSessionNotFound
	}
	session := s.sessionsByToken[tokenHash]
	if session.RevokedAt != nil {
		return nil
	}
	now := time.Now().UTC()
	session.RevokedAt = &now
	s.sessionsByToken[tokenHash] = session
	return nil
}

func cloneSession(session Session) Session {
	clone := session
	if session.RevokedAt != nil {
		value := *session.RevokedAt
		clone.RevokedAt = &value
	}
	return clone
}

var _ Repository = (*MemoryStore)(nil)
var _ AuthRepository = (*MemoryStore)(nil)
