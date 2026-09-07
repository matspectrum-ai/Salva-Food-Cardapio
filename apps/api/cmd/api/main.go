package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/matspectrum-ai/salva-food/apps/api/internal/catalog"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/orders"
	"github.com/matspectrum-ai/salva-food/apps/api/internal/platform/httpapi"
)

type healthResponse struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

func main() {
	catalogStore := catalog.NewMemoryStore()
	orderStore := orders.NewMemoryStore()
	api := httpapi.New(catalogStore, orderStore, newUUID)

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

func newUUID() string {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		panic(err)
	}
	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80
	buf := make([]byte, 36)
	hex.Encode(buf[0:8], bytes[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], bytes[4:6])
	hex.Encode(buf[14:18], bytes[6:8])
	buf[13] = '-'
	hex.Encode(buf[19:23], bytes[8:10])
	buf[18] = '-'
	hex.Encode(buf[24:36], bytes[10:16])
	buf[23] = '-'
	return string(buf)
}
