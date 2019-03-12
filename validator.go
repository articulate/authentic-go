package authentic

import (
	"time"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type (
	Result struct {
		Claims  map[string]interface{}
		Valid   bool
		Expired bool
	}

	validator struct {
		CacheMaxAge  time.Duration
		ISSWhitelist []string
		keyManager   *keyManager
		clock        Clock
	}
)

// IsValid checks validity of token and ensures it is not expired
func (v *validator) IsValid(token string) bool {
	result := v.ValidateToken(token)
	return result.Valid && !result.Expired
}

// ValidateToken verifies JWT token, and returns result
func (v *validator) ValidateToken(token string) *Result {
	var (
		claims map[string]interface{}
		header jose.Header
		key    *jose.JSONWebKey
		result = &Result{}
	)
	parsedToken, err := jwt.ParseSigned(token)
	if err != nil {
		return result
	}
	// Note this is required to get the iss, the call to Claims later on validates the claims
	err = parsedToken.UnsafeClaimsWithoutVerification(&claims)
	if err != nil {
		return result
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
		return result
	}

	if err = parsedToken.Claims(key.Key, &claims); err != nil {
		return result
	}
	result.Valid = true
	result.Claims = claims
	result.Expired = v.IsExpired(result.Claims)

	return result
}

func (v *validator) IsExpired(claims map[string]interface{}) bool {
	if exp, ok := claims["exp"].(float64); ok {
		tm := time.Unix(int64(exp), 0)

		return v.clock.IsBeforeNow(tm)
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

func (v *validator) withClock(c Clock) Validator {
	v.clock = c
	return v
}
