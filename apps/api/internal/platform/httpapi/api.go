package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/matspectrum-ai/salva-food/apps/api/internal/catalog"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/orders"
)

type IDGenerator func() string

type API struct {
	catalog      *catalog.MemoryStore
	orders       *orders.MemoryStore
	orderService *orders.Service
	newID        IDGenerator
}

type catalogReader struct {
	store *catalog.MemoryStore
}

func (r catalogReader) LookupOrderProduct(tenantID, itemID string) (orders.CatalogProduct, error) {
	item, ok := r.store.GetItem(tenantID, itemID)
	if !ok {
		return orders.CatalogProduct{}, catalog.ErrItemNotFound
	}
	category, ok := r.store.GetCategory(tenantID, item.CategoryID)
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

func New(catalogStore *catalog.MemoryStore, orderStore *orders.MemoryStore, newID IDGenerator) *API {
	return &API{
		catalog:      catalogStore,
		orders:       orderStore,
		orderService: orders.NewService(orderStore, catalogReader{store: catalogStore}, orders.IDGenerator(newID)),
		newID:        newID,
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
	return mux
}

type createCategoryRequest struct {
	Name      string `json:"name"`
	SortOrder int    `json:"sort_order"`
}

func (a *API) createCategory(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := tenantFromRequest(w, r)
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
	if err := a.catalog.CreateCategory(category); err != nil {
		writeError(w, http.StatusConflict, err)
		return
	}
	writeJSON(w, http.StatusCreated, category)
}

func (a *API) listCategories(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := tenantFromRequest(w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": a.catalog.ListCategories(tenantID)})
}

type createItemRequest struct {
	CategoryID   string `json:"category_id"`
	Name         string `json:"name"`
	PriceCents   int64  `json:"price_cents"`
	SortOrder    int    `json:"sort_order"`
	SoldByWeight bool   `json:"sold_by_weight"`
}

func (a *API) createItem(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := tenantFromRequest(w, r)
	if !ok {
		return
	}
	var req createItemRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	category, ok := a.catalog.GetCategory(tenantID, req.CategoryID)
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
	if err := a.catalog.CreateItem(item); err != nil {
		writeError(w, http.StatusConflict, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (a *API) listItems(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := tenantFromRequest(w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": a.catalog.ListItems(tenantID, r.URL.Query().Get("category_id"))})
}
func (a *API) createOrder(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := tenantFromRequest(w, r)
	if !ok {
		return
	}
	var input orders.CreateInput
	if !decodeJSON(w, r, &input) {
		return
	}
	order, replay, err := a.orderService.Create(tenantID, r.Header.Get("Idempotency-Key"), input)
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
	tenantID, ok := tenantFromRequest(w, r)
	if !ok {
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"data": a.orders.List(tenantID)})
}

type transitionOrderRequest struct {
	Status orders.Status `json:"status"`
}

func (a *API) transitionOrder(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := tenantFromRequest(w, r)
	if !ok {
		return
	}
	var req transitionOrderRequest
	if !decodeJSON(w, r, &req) {
		return
	}
	order, err := a.orders.Transition(tenantID, r.PathValue("id"), req.Status)
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

func tenantFromRequest(w http.ResponseWriter, r *http.Request) (string, bool) {
	tenantID := strings.TrimSpace(r.Header.Get("X-Tenant-ID"))
	if tenantID == "" {
		writeError(w, http.StatusBadRequest, errors.New("X-Tenant-ID header is required"))
		return "", false
	}
	return tenantID, true
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
