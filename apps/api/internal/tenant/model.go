package tenant

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrIDRequired      = errors.New("id is required")
	ErrTenantRequired  = errors.New("tenant id is required")
	ErrNameRequired    = errors.New("name is required")
	ErrInvalidTimezone = errors.New("invalid timezone")
)

type Tenant struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Establishment struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
	Name     string `json:"name"`
	Timezone string `json:"timezone"`
}

func New(id, name string) (Tenant, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Tenant{}, ErrIDRequired
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return Tenant{}, ErrNameRequired
	}
	return Tenant{ID: id, Name: name}, nil
}

func NewEstablishment(id, tenantID, name, timezone string) (Establishment, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Establishment{}, ErrIDRequired
	}
	tenantID = strings.TrimSpace(tenantID)
	if tenantID == "" {
		return Establishment{}, ErrTenantRequired
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return Establishment{}, ErrNameRequired
	}
	timezone = strings.TrimSpace(timezone)
	if timezone == "" {
		return Establishment{}, ErrInvalidTimezone
	}
	if _, err := time.LoadLocation(timezone); err != nil {
		return Establishment{}, ErrInvalidTimezone
	}
	return Establishment{
		ID:       id,
		TenantID: tenantID,
		Name:     name,
		Timezone: timezone,
	}, nil
}
