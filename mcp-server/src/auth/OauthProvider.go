package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-jose/go-jose/v4"
	"github.com/rizvn/panics"
)

type OauthProvider struct {
	JwksUri string `json:"jwks_uri"`

	IssuerUrl       string
	OpenIDConfigUrl string

	jwks       map[string]jose.JSONWebKey
	httpClient *http.Client
}

func (r *OauthProvider) Init() {
	if r.IssuerUrl == "" {
		panics.OnError(fmt.Errorf("OauthProvider Init: IssuerUrl URL is required"), "")
	}

	// Set OpenIDConfigUrl if not provided
	if r.OpenIDConfigUrl == "" {
		r.OpenIDConfigUrl = r.IssuerUrl + "/.well-known/openid-configuration"
	}

	resp, err := http.Get(r.OpenIDConfigUrl)
	panics.OnError(err, "failed to fetch OIDC oauthProvider document")
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	panics.OnError(err, "failed to read OIDC oauthProvider document")

	err = json.Unmarshal(body, r)
	panics.OnError(err, "failed to unmarshal OIDC oauthProvider document")

	// Fetch JWKS
	jwks, err := r.fetchJwks()
	panics.OnError(err, "failed to fetch JWKS from provider")

	r.jwks = jwks
	r.httpClient = &http.Client{}

}

func (r *OauthProvider) GetKey(kid string) jose.JSONWebKey {
	jwk, ok := r.jwks[kid]
	if !ok {
		// Try to refresh JWKS
		jwks, err := r.fetchJwks()
		panics.OnError(err, "failed to refresh JWKS from provider")

		// Update local jwks
		r.jwks = jwks

		// Try to get the key again
		jwk, ok = r.jwks[kid]
		panics.OnFalse(ok, "JWK with kid %s not found in provider JWKS", kid)
	}
	return jwk
}

func (r *OauthProvider) fetchJwks() (map[string]jose.JSONWebKey, error) {
	resp, err := http.Get(r.JwksUri)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jwks jose.JSONWebKeySet
	err = json.Unmarshal(body, &jwks)
	if err != nil {
		return nil, err
	}

	keyMap := make(map[string]jose.JSONWebKey)
	for _, key := range jwks.Keys {
		keyMap[key.KeyID] = key
	}
	return keyMap, nil
}
