package identity

import "context"

type PasswordHasher interface {
	Hash(string) (string, error)
}

type CreateCollaboratorInput struct {
	Name        string       `json:"name"`
	CPF         string       `json:"cpf"`
	Email       string       `json:"email"`
	Password    string       `json:"password"`
	Phone       string       `json:"phone"`
	Title       string       `json:"title"`
	ImageURL    string       `json:"image_url,omitempty"`
	Permissions []Permission `json:"permissions"`
}

type Service struct {
	repo   Repository
	hasher PasswordHasher
	newID  func() string
}

func NewService(repo Repository, hasher PasswordHasher, newID func() string) *Service {
	return &Service{repo: repo, hasher: hasher, newID: newID}
}
func (s *Service) CreateCollaborator(ctx context.Context, tenantID string, input CreateCollaboratorInput) (Collaborator, error) {
	if err := ValidatePassword(input.Password); err != nil {
		return Collaborator{}, err
	}
	hash, err := s.hasher.Hash(input.Password)
	if err != nil {
		return Collaborator{}, err
	}
	userID := s.newID()
	user, err := NewUser(userID, input.Name, input.CPF, input.Email, input.Phone, hash)
	if err != nil {
		return Collaborator{}, err
	}
	user.ImageURL = input.ImageURL
	membership, err := NewMembership(tenantID, userID, input.Title, input.Permissions)
	if err != nil {
		return Collaborator{}, err
	}
	collaborator := Collaborator{User: user, Membership: membership}
	if err := s.repo.CreateCollaborator(ctx, collaborator); err != nil {
		return Collaborator{}, err
	}
	return collaborator, nil
}

func (s *Service) ListCollaborators(ctx context.Context, tenantID, search string) ([]Collaborator, error) {
	return s.repo.ListCollaborators(ctx, tenantID, search)
}

func (s *Service) SetMembershipStatus(ctx context.Context, tenantID, userID string, status MembershipStatus) error {
	return s.repo.SetMembershipStatus(ctx, tenantID, userID, status)
}
