package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"

	"github.com/matspectrum-ai/salva-food/apps/api/internal/catalog"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/orders"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/platform/httpapi"
	postgresstore "github.com/matspectrum-ai/salva-food/apps/api/internal/platform/postgres"
)

type healthResponse struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

func main() {
	ctx := context.Background()
	var catalogRepo catalog.Repository
	var orderRepo orders.Repository
	var closeStore func()

	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		store, err := postgresstore.Open(ctx, databaseURL)
		if err != nil {
			log.Fatalf("open postgres store: %v", err)
		}
		catalogRepo = store
		orderRepo = store
		closeStore = store.Close
		log.Printf("storage backend: postgres")
	} else {
		catalogRepo = catalog.NewMemoryStore()
		orderRepo = orders.NewMemoryStore()
		closeStore = func() {}
		log.Printf("storage backend: memory")
	}
	defer closeStore()

	api := httpapi.New(catalogRepo, orderRepo, uuid.NewString)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.Handle("/", api.Handler())

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("salva-food-api listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(healthResponse{Service: "salva-food-api", Status: "ok"})
}
