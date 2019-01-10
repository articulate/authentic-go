package authentic

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gopkg.in/square/go-jose.v2"
)

type (
	configResponse struct {
		jwkURI string `json:"jwks_uri"`
	}

	// keyManager manages the fetching and persisting of JWKs
	keyManager struct {
		cache  *cache
		client *http.Client
	}
)

const (
	wellKnown = "/.well-known/openid-configuration"
)

// NewKeyManager get key manager
func newKeyManager() *keyManager {
	return &keyManager{
		cache:  newCache(time.Hour * 10),
		client: &http.Client{},
	}
}

// Get key from OIDC config endpoint or cache if cache is not stale
// If we wanted better performance, we could spin off a thread to populate cache and simply return stale.
func (k *keyManager) Get(issuer, kid string) interface{} {
	if k.cache.KeyIsExpired(issuer, kid) {
		k.updateCache(issuer, kid)
	}

	return k.cache.GetKey(issuer, kid)
}

func (k *keyManager) updateCache(issuer, kid string) error {
	uri, err := url.Parse(strings.Trim(issuer, "/") + wellKnown)
	if err != nil {
		return err
	}
	jwksURI, err := k.fetchKeyURI(uri.String())
	if err != nil {
		return err
	}
	keyList, err := k.fetchKeys(jwksURI)
	if err != nil {
		return err
	}

	for _, key := range keyList {
		k.cache.SetKey(issuer, key.KeyID, &key)
	}

	return nil
}

func (k *keyManager) fetchKeyURI(uri string) (string, error) {
	var body map[string]interface{}
	resp, err := k.client.Get(uri)
	if err != nil {
		return "", errors.New("Failed to retrieve issuer config")
	}

	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", errors.New("Failed to parse the issuer config response body")
	}

	return body["jwks_uri"].(string), nil
}

func (k *keyManager) fetchKeys(uri string) ([]jose.JSONWebKey, error) {
	var keys jose.JSONWebKeySet
	resp, err := k.client.Get(uri)
	if err != nil {
		return nil, errors.New("Failed to retrieve JWKs")
	}

	if err = json.NewDecoder(resp.Body).Decode(&keys); err != nil {
		return nil, errors.New("Failed to parse JWKs from response body")
	}

	return keys.Keys, nil
}
