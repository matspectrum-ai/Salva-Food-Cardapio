package identity

import (
	"errors"
	"net/mail"
	"sort"
	"strings"
	"unicode"
)

var (
	ErrIDRequired          = errors.New("id is required")
	ErrTenantRequired      = errors.New("tenant id is required")
	ErrNameRequired        = errors.New("name is required")
	ErrInvalidCPF          = errors.New("invalid CPF")
	ErrInvalidEmail        = errors.New("invalid email")
	ErrPhoneRequired       = errors.New("phone is required")
	ErrPasswordHashMissing = errors.New("password hash is required")
	ErrTitleRequired       = errors.New("title is required")
	ErrPermissionsRequired = errors.New("at least one permission is required")
	ErrInvalidPermission   = errors.New("invalid permission")
	ErrInvalidStatus       = errors.New("invalid membership status")
	ErrWeakPassword        = errors.New("password does not meet policy")
)

type Permission string

const (
	PermissionReports          Permission = "reports"
	PermissionVirtualAssistant Permission = "virtual_assistant"
	PermissionSalesGrowth      Permission = "sales_growth"
	PermissionSubscription     Permission = "subscription"
	PermissionSettings         Permission = "settings"
	PermissionPOS              Permission = "pos"
	PermissionEstablishment    Permission = "establishment"
	PermissionCoupon           Permission = "coupon"
	PermissionSatisfaction     Permission = "satisfaction"
	PermissionMenuItems        Permission = "menu_items"
	PermissionDeliveryRegions  Permission = "delivery_regions"
	PermissionBuyMore          Permission = "buy_more"
	PermissionCashback         Permission = "cashback"
	PermissionSalesRecovery    Permission = "sales_recovery"
	PermissionHallWaiter       Permission = "hall_waiter"
	PermissionMenuEdit         Permission = "menu_edit"
	PermissionCashOpen         Permission = "cash_open"
	PermissionReadyOrderEdit   Permission = "ready_order_edit"
	PermissionFiscalIssue      Permission = "fiscal_issue"
	PermissionNFCEIssue        Permission = "nfce_issue"
	PermissionFinance          Permission = "finance"
	PermissionPurchasing       Permission = "purchasing"
	PermissionInventory        Permission = "inventory"
	PermissionDashboard        Permission = "dashboard"
	PermissionOnlinePayment    Permission = "online_payment"
)

var validPermissions = map[Permission]struct{}{
	PermissionReports: {}, PermissionVirtualAssistant: {}, PermissionSalesGrowth: {},
	PermissionSubscription: {}, PermissionSettings: {}, PermissionPOS: {},
	PermissionEstablishment: {}, PermissionCoupon: {}, PermissionSatisfaction: {},
	PermissionMenuItems: {}, PermissionDeliveryRegions: {}, PermissionBuyMore: {},
	PermissionCashback: {}, PermissionSalesRecovery: {}, PermissionHallWaiter: {},
	PermissionMenuEdit: {}, PermissionCashOpen: {}, PermissionReadyOrderEdit: {},
	PermissionFiscalIssue: {}, PermissionNFCEIssue: {}, PermissionFinance: {},
	PermissionPurchasing: {}, PermissionInventory: {}, PermissionDashboard: {},
	PermissionOnlinePayment: {},
}

func NormalizePermissionSet(input []Permission) ([]Permission, error) {
	if len(input) == 0 {
		return nil, ErrPermissionsRequired
	}
	seen := make(map[Permission]struct{}, len(input))
	for _, permission := range input {
		if _, ok := validPermissions[permission]; !ok {
			return nil, ErrInvalidPermission
		}
		seen[permission] = struct{}{}
	}
	result := make([]Permission, 0, len(seen))
	for permission := range seen {
		result = append(result, permission)
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return result, nil
}
func ValidatePassword(password string) error {
	if len([]rune(password)) < 8 {
		return ErrWeakPassword
	}
	var upper, lower, digit, symbol bool
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			upper = true
		case unicode.IsLower(r):
			lower = true
		case unicode.IsDigit(r):
			digit = true
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			symbol = true
		}
	}
	if !upper || !lower || !digit || !symbol {
		return ErrWeakPassword
	}
	return nil
}

type User struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	CPF          string `json:"cpf"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	ImageURL     string `json:"image_url,omitempty"`
	PasswordHash string `json:"-"`
}

func NewUser(id, name, cpf, email, phone, passwordHash string) (User, error) {
	if strings.TrimSpace(id) == "" {
		return User{}, ErrIDRequired
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return User{}, ErrNameRequired
	}
	cpf = normalizeDigits(cpf)
	if !validCPF(cpf) {
		return User{}, ErrInvalidCPF
	}
	email = strings.ToLower(strings.TrimSpace(email))
	parsed, err := mail.ParseAddress(email)
	if err != nil || parsed.Address != email || !strings.Contains(email, "@") {
		return User{}, ErrInvalidEmail
	}
	phone = strings.TrimSpace(phone)
	if phone == "" {
		return User{}, ErrPhoneRequired
	}
	if strings.TrimSpace(passwordHash) == "" {
		return User{}, ErrPasswordHashMissing
	}
	return User{ID: id, Name: name, CPF: cpf, Email: email, Phone: phone, PasswordHash: passwordHash}, nil
}

type MembershipStatus string

const (
	MembershipActive   MembershipStatus = "ACTIVE"
	MembershipInactive MembershipStatus = "INACTIVE"
)

type Membership struct {
	TenantID    string           `json:"tenant_id"`
	UserID      string           `json:"user_id"`
	Title       string           `json:"title"`
	Status      MembershipStatus `json:"status"`
	Permissions []Permission     `json:"permissions"`
}

func NewMembership(tenantID, userID, title string, permissions []Permission) (Membership, error) {
	if strings.TrimSpace(tenantID) == "" {
		return Membership{}, ErrTenantRequired
	}
	if strings.TrimSpace(userID) == "" {
		return Membership{}, ErrIDRequired
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return Membership{}, ErrTitleRequired
	}
	normalized, err := NormalizePermissionSet(permissions)
	if err != nil {
		return Membership{}, err
	}
	return Membership{
		TenantID: tenantID, UserID: userID, Title: title,
		Status: MembershipActive, Permissions: normalized,
	}, nil
}

func (m Membership) HasPermission(permission Permission) bool {
	for _, candidate := range m.Permissions {
		if candidate == permission {
			return true
		}
	}
	return false
}
func (m *Membership) SetStatus(status MembershipStatus) error {
	if status != MembershipActive && status != MembershipInactive {
		return ErrInvalidStatus
	}
	m.Status = status
	return nil
}

func normalizeDigits(value string) string {
	var b strings.Builder
	for _, r := range value {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func validCPF(cpf string) bool {
	if len(cpf) != 11 {
		return false
	}
	allSame := true
	for i := 1; i < len(cpf); i++ {
		if cpf[i] != cpf[0] {
			allSame = false
			break
		}
	}
	if allSame {
		return false
	}
	return cpfDigit(cpf[:9], 10) == int(cpf[9]-'0') && cpfDigit(cpf[:10], 11) == int(cpf[10]-'0')
}

func cpfDigit(prefix string, weight int) int {
	sum := 0
	for i := 0; i < len(prefix); i++ {
		sum += int(prefix[i]-'0') * (weight - i)
	}
	remainder := (sum * 10) % 11
	if remainder == 10 {
		return 0
	}
	return remainder
}
