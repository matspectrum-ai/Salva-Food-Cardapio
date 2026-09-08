package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/matspectrum-ai/salva-food/apps/api/internal/catalog"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/identity"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/onboarding"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/orders"
)

type IDGenerator func() string

type API struct {
	catalog           catalog.Repository
	orders            orders.Repository
	orderService      *orders.Service
	identityService   *identity.Service
	authService       *identity.AuthService
	onboardingService *onboarding.Service
	strictAuth        bool
	newID             IDGenerator
}

type catalogReader struct {
	store catalog.Repository
}

func (r catalogReader) LookupOrderProduct(ctx context.Context, tenantID, itemID string) (orders.CatalogProduct, error) {
	item, ok, err := r.store.GetItem(ctx, tenantID, itemID)
	if err != nil {
		return orders.CatalogProduct{}, err
	}
	if !ok {
		return orders.CatalogProduct{}, catalog.ErrItemNotFound
	}
	category, ok, err := r.store.GetCategory(ctx, tenantID, item.CategoryID)
	if err != nil {
		return orders.CatalogProduct{}, err
	}
	if !ok {
		return orders.CatalogProduct{}, catalog.ErrCategoryNotFound
	}
	return orders.CatalogProduct{
		ID:              item.ID,
		Name:            item.Name,
		PriceCents:      item.PriceCents,
		SoldOut:         item.SoldOut,
		CategorySoldOut: category.SoldOut,
	}, nil
}

func New(catalogStore catalog.Repository, orderStore orders.Repository, newID IDGenerator) *API {
	return newAPI(catalogStore, orderStore, nil, nil, nil, false, newID)
}

func NewAuthenticated(catalogStore catalog.Repository, orderStore orders.Repository, identityService *identity.Service, authService *identity.AuthService, newID IDGenerator) *API {
	return newAPI(catalogStore, orderStore, identityService, authService, nil, true, newID)
}

func NewAuthenticatedWithOnboarding(catalogStore catalog.Repository, orderStore orders.Repository, identityService *identity.Service, authService *identity.AuthService, onboardingService *onboarding.Service, newID IDGenerator) *API {
	return newAPI(catalogStore, orderStore, identityService, authService, onboardingService, true, newID)
}

func newAPI(catalogStore catalog.Repository, orderStore orders.Repository, identityService *identity.Service, authService *identity.AuthService, onboardingService *onboarding.Service, strictAuth bool, newID IDGenerator) *API {
	return &API{
		catalog: catalogStore, orders: orderStore,
		orderService:    orders.NewService(orderStore, catalogReader{store: catalogStore}, orders.IDGenerator(newID)),
		identityService: identityService, authService: authService, onboardingService: onboardingService, strictAuth: strictAuth, newID: newID,
	}
}
func (a *API) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/catalog/categories", a.createCategory)
	mux.HandleFunc("GET /api/v1/catalog/categories", a.listCategories)
	mux.HandleFunc("POST /api/v1/catalog/items", a.createItem)
	mux.HandleFunc("GET /api/v1/catalog/items", a.listItems)
	mux.HandleFunc("POST /api/v1/orders", a.createOrder)
	mux.HandleFunc("GET /api/v1/orders", a.listOrders)
	mux.HandleFunc("POST /api/v1/orders/{id}/transitions", a.transitionOrder)
	if a.onboardingService != nil {
		a.registerOnboardingRoutes(mux)
	}
	if a.authService != nil && a.identityService != nil {
		a.registerIdentityRoutes(mux)
	}
	return mux
}

type createCategoryRequest struct {
	Name      string `json:"name"`
	SortOrder int    `json:"sort_order"`
}

func (a *API) createCategory(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := a.tenantFromRequest(w, r)
	if !ok {
		return
	}
	var req createCategoryRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	category, err := catalog.NewCategory(tenantID, a.newID(), req.Name, req.SortOrder)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := a.catalog.CreateCategory(r.Context(), category); err != nil {
		writeError(w, http.StatusConflict, err)
		return
	}
	writeJSON(w, http.StatusCreated, category)
}

func (a *API) listCategories(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := a.tenantFromRequest(w, r)
	if !ok {
		return
	}
	categories, err := a.catalog.ListCategories(r.Context(), tenantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": categories})
}

type createItemRequest struct {
	CategoryID   string `json:"category_id"`
	Name         string `json:"name"`
	PriceCents   int64  `json:"price_cents"`
	SortOrder    int    `json:"sort_order"`
	SoldByWeight bool   `json:"sold_by_weight"`
}

func (a *API) createItem(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := a.tenantFromRequest(w, r)
	if !ok {
		return
	}
	var req createItemRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	category, ok, err := a.catalog.GetCategory(r.Context(), tenantID, req.CategoryID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, catalog.ErrCategoryNotFound)
		return
	}
	item, err := catalog.NewItem(tenantID, a.newID(), req.Name, category, req.PriceCents, req.SortOrder)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	item.SoldByWeight = req.SoldByWeight
	if err := a.catalog.CreateItem(r.Context(), item); err != nil {
		writeError(w, http.StatusConflict, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (a *API) listItems(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := a.tenantFromRequest(w, r)
	if !ok {
		return
	}
	items, err := a.catalog.ListItems(r.Context(), tenantID, r.URL.Query().Get("category_id"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": items})
}
func (a *API) createOrder(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := a.tenantFromRequest(w, r)
	if !ok {
		return
	}
	var input orders.CreateInput
	if !decodeJSON(w, r, &input) {
		return
	}
	order, replay, err := a.orderService.Create(r.Context(), tenantID, r.Header.Get("Idempotency-Key"), input)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, orders.ErrIdempotencyConflict) {
			status = http.StatusConflict
		}
		if errors.Is(err, catalog.ErrItemNotFound) || errors.Is(err, catalog.ErrCategoryNotFound) {
			status = http.StatusNotFound
		}
		writeError(w, status, err)
		return
	}
	if replay {
		w.Header().Set("Idempotent-Replay", "true")
	}
	writeJSON(w, http.StatusCreated, order)
}

func (a *API) listOrders(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := a.tenantFromRequest(w, r)
	if !ok {
		return
	}
	result, err := a.orders.List(r.Context(), tenantID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": result})
}

type transitionOrderRequest struct {
	Status orders.Status `json:"status"`
}

func (a *API) transitionOrder(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := a.tenantFromRequest(w, r)
	if !ok {
		return
	}
	var req transitionOrderRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	order, err := a.orders.Transition(r.Context(), tenantID, r.PathValue("id"), req.Status)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, orders.ErrOrderNotFound) {
			status = http.StatusNotFound
		}
		writeError(w, status, err)
		return
	}
	writeJSON(w, http.StatusOK, order)
}

func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) bool {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dst); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return false
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			err = errors.New("request body must contain a single JSON object")
		}
		writeError(w, http.StatusBadRequest, err)
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]any{
		"error": map[string]string{
			"message": err.Error(),
		},
	})
}
