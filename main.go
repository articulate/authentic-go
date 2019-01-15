package authentic

import (
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	cacheMaxAgeHours = 10
)

type (
	// MiddlewareCreator helpers that create middleware to validate JWTs on incoming requests
	MiddlewareCreator interface {
		// CreateGinMiddleware specifically for the Gin framework
		CreateGinMiddleware() gin.HandlerFunc
		OnSuccess(gin.HandlerFunc) MiddlewareCreator
		OnFailure(gin.HandlerFunc) MiddlewareCreator
	}

	// Validator interface for validating a JWT
	Validator interface {
		// IsValid check the validity of a JWT token
		IsValid(string) bool
		WithWhitelist(...string) Validator
		WithCacheMaxAge(time.Duration) Validator
	}
)

// NewValidator configures and creates a JWT validator
func NewValidator() Validator {
	// Setup sensible defaults prior to exposing validator which has `WithX` helper functions
	return &validator{
		CacheMaxAge:  time.Hour * time.Duration(cacheMaxAgeHours),
		ISSWhitelist: strings.Split(os.Getenv("ISS_WHITELIST"), "|"),
		keyManager:   newKeyManager(),
	}
}

// NewMiddlewareCreator creates a new Middleware creation helper
func NewMiddlewareCreator() MiddlewareCreator {
	return &middlewareCreator{
		Validator: NewValidator(),
	}
}
