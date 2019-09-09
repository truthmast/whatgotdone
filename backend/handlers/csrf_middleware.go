package handlers

import (
	"log"
	"os"

	"github.com/gorilla/csrf"
)

func newCsrfMiddleware() httpMiddlewareHandler {
	csrfSeed := os.Getenv("CSRF_SECRET_SEED")
	if csrfSeed == "" {
		log.Fatalf("CSRF_SECRET_SEED environment variable must be set")
	}
	return csrf.Protect(
		[]byte(csrfSeed),
		// The v2 suffix is just to prevent clients from re-using a previous version of the cookie.
		csrf.CookieName("csrf_base_v2"),
		csrf.Path("/"),
		csrf.Secure(false))
}
