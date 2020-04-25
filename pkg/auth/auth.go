package auth

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/kkeuning/go-api-example/pkg/services"
	"golang.org/x/crypto/bcrypt"
)

// APIKeyMiddleware is a simplified shared secret authentication for api requests
type APIKeyMiddleware struct {
	apiKey []byte
	Path   string
}

// InitializeKey will load a hash that could be safely persisted instead of the actual key
func (akm *APIKeyMiddleware) InitializeKey(key string) {
	k := []byte(key)
	bCryptKey, err := bcrypt.GenerateFromPassword(k, bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Key initialization error.")
	}
	akm.apiKey = bCryptKey
}

// KeyIsValid checks that userAPIKey matches userAPIKey hash
func (akm *APIKeyMiddleware) KeyIsValid(userAPIKey string) bool {
	if e := bcrypt.CompareHashAndPassword(akm.apiKey, []byte(userAPIKey)); e != nil {
		log.Println("Key Mismatch")
		return false
	}
	return true
}

// Middleware function, which will be called for each api request
func (akm *APIKeyMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if akm.Path == "" || strings.Contains(r.URL.Path, akm.Path) {
			log.Printf("Authorization required for %s", r.URL.Path)
			key := r.Header.Get("Authorization")
			if akm.KeyIsValid(key) {
				log.Printf("User is authorized.")
				ctx := context.WithValue(r.Context(), services.AuthorizationKey, key)
				// Pass down the request to the next middleware (or final handler)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			// Write an error and return to stop the handler chain
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		// Auth not required for non-api endpoints
		// Just pass the request to the next middleware (or final handler)
		next.ServeHTTP(w, r)
	})
}
