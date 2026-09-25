package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/matspectrum-ai/salva-food/apps/api/internal/catalog"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/customers"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/identity"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/onboarding"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/orders"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/platform/httpapi"
	postgresstore "github.com/matspectrum-ai/salva-food/apps/api/internal/platform/postgres"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/realtime"
)

type healthResponse struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

func main() {
	ctx := context.Background()
	var catalogRepo catalog.Repository
	var orderRepo orders.Repository
	var customerRepo customers.Repository
	var identityRepo identity.Repository
	var authRepo identity.AuthRepository
	var onboardingRepo onboarding.Repository
	var closeStore func()
	var eventBus realtime.Bus
	var publishDirect bool

	if databaseURL := os.Getenv("DATABASE_URL"); databaseURL != "" {
		store, err := postgresstore.Open(ctx, databaseURL)
		if err != nil {
			log.Fatalf("open postgres store: %v", err)
		}
		catalogRepo = store
		orderRepo = store
		identityRepo = store
		authRepo = store
		onboardingRepo = store
		customerRepo = store
		redisURL := os.Getenv("REDIS_URL")
		if redisURL == "" {
			log.Fatal("REDIS_URL is required when DATABASE_URL is configured")
		}
		redisOptions, err := redis.ParseURL(redisURL)
		if err != nil {
			log.Fatalf("parse REDIS_URL: %v", err)
		}
		redisClient := redis.NewClient(redisOptions)
		if err := redisClient.Ping(ctx).Err(); err != nil {
			_ = redisClient.Close()
			store.Close()
			log.Fatalf("ping redis: %v", err)
		}
		eventBus = realtime.NewRedisBus(redisClient, "salva-food")
		publisher := postgresstore.NewOutboxPublisher(store, eventBus)
		go publisher.Run(ctx)
		closeStore = func() { _ = redisClient.Close(); store.Close() }
		log.Printf("storage backend: postgres")
	} else {
		catalogRepo = catalog.NewMemoryStore()
		orderRepo = orders.NewMemoryStore()
		customerRepo = customers.NewMemoryStore()
		identityStore := identity.NewMemoryStore()
		identityRepo = identityStore
		authRepo = identityStore
		onboardingRepo = onboarding.NewMemoryRepository(identityStore)
		closeStore = func() {}
		eventBus = realtime.NewMemoryBus()
		publishDirect = true
		log.Printf("storage backend: memory")
	}
	defer closeStore()

	passwords := identity.NewBcryptPasswordManager(0)
	identityService := identity.NewService(identityRepo, passwords, uuid.NewString)
	onboardingService := onboarding.NewService(onboardingRepo, passwords, uuid.NewString)
	authService := identity.NewAuthService(
		authRepo, passwords, newOpaqueToken, uuid.NewString, time.Now, sessionTTL(),
	)
	api := httpapi.NewAuthenticatedWithOnboardingAndRealtime(catalogRepo, orderRepo, identityService, authService, onboardingService, uuid.NewString, customerRepo, eventBus, publishDirect)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)
	mux.Handle("/", api.Handler())

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":18080"
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

func newOpaqueToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func sessionTTL() time.Duration {
	value := os.Getenv("SESSION_TTL")
	if value == "" {
		return 12 * time.Hour
	}
	duration, err := time.ParseDuration(value)
	if err != nil || duration <= 0 {
		log.Fatalf("invalid SESSION_TTL %q", value)
	}
	return duration
}
