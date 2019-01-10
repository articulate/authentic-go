package authentic

import (
	"time"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type (
	validator struct {
		CacheMaxAge  time.Duration
		ISSWhitelist []string
		keyManager   *keyManager
	}
)

// IsValid verifies JWT token
func (v *validator) IsValid(token string) bool {
	var (
		claims map[string]interface{}
		header jose.Header
		key    *jose.JSONWebKey
	)
	parsedToken, err := jwt.ParseSigned(token)
	if err != nil {
		return false
	}
	err = parsedToken.UnsafeClaimsWithoutVerification(&claims)
	if err != nil {
		return false
	}

	for _, h := range parsedToken.Headers {
		jwk := v.keyManager.Get(claims["iss"].(string), h.KeyID)
		if jwk != nil {
			key = jwk.(*jose.JSONWebKey)
			header = h
			continue
		}
	}

	// If JWT has no corresponding JWK, JWK is invalid, or JWT encryption algorithm does not match the JWKs
	if key == nil || !key.Valid() || key.Algorithm != header.Algorithm {
		return false
	}

	if err = parsedToken.Claims(key.Key, &claims); err != nil {
		return false
	}

	return true
}

// WithWhitelist set ISS whitelist
func (v *validator) WithWhitelist(whitelist ...string) Validator {
	v.ISSWhitelist = whitelist
	return v
}

// WithCacheMaxAge set cache max age
func (v *validator) WithCacheMaxAge(c time.Duration) Validator {
	v.CacheMaxAge = c
	return v
}
