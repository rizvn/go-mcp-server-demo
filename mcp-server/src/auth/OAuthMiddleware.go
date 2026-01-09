package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/modelcontextprotocol/go-sdk/oauthex"
)

type OAuthMiddleware struct {
	provider       *OauthProvider
	TargetAudience string
	IssuerUrl      string
	Scope          string
}

func (r *OAuthMiddleware) Init() {
	r.provider = &OauthProvider{
		IssuerUrl: r.IssuerUrl,
	}
	r.provider.Init()
}

func (r *OAuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {

		// Handle OAuth 2.1 protected resource metadata endpoint
		// no authorization required
		if rq.URL.Path == "/.well-known/oauth-protected-resource" {
			r.HandleProtectedResourceMetadata(w, rq)
			return
		}

		// Check Authorization header
		authHeader := rq.Header.Get("Authorization")
		if authHeader == "" {
			r.sendUnauthorized(w, rq)
			return
		}

		// Extract Bearer token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			r.sendUnauthorized(w, rq)
			return
		}

		// Validate JWT token using JWKS with algorithm validation
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure the token's signing method is RSA
			jwk := r.provider.GetKey(token.Header["kid"].(string))
			return jwk.Key, nil
		}, jwt.WithValidMethods([]string{"RS256"}))

		if err != nil {
			slog.Error(fmt.Sprintf("Failed to parse token: %v", err))
			r.sendUnauthorized(w, rq)
			return
		}

		if !token.Valid {
			slog.Error(fmt.Sprintf("Invalid token"))
			r.sendUnauthorized(w, rq)
			return
		}

		// Get claims for validation
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			slog.Error(fmt.Sprintf("Invalid claims type"))
			r.sendUnauthorized(w, rq)
			return
		}

		// Validate audience (MUST): Verify this resource server is in the audience
		if !r.validateAudience(claims) {
			slog.Error(fmt.Sprintf("Invalid audience"))
			r.sendUnauthorized(w, rq)
			return
		}

		// Validate issuer (MUST): Verify token is issued by expected authorization server
		if !r.validateIssuer(claims) {
			slog.Error(fmt.Sprintf("Invalid issuer"))
			r.sendUnauthorized(w, rq)
			return
		}

		// Validate expiration (MUST): Ensure token is not expired
		// Note: jwt.Parse already validates exp by default, but we explicitly check here for clarity
		if !r.validateExpiration(claims) {
			slog.Error(fmt.Sprintf("Token expired"))
			r.sendUnauthorized(w, rq)
			return
		}

		// Validate scope: Verify token has required scopes (optional, depends on your requirements)
		if !r.validateScope(claims) {
			slog.Error(fmt.Sprintf("Insufficient scope"))
			r.sendUnauthorized(w, rq)
			return
		}

		// store user info in context for downstream handlers
		if email, ok := claims["email"].(string); ok {
			rq = rq.WithContext(context.WithValue(rq.Context(), "user_email", email))
		}

		// Authorization successful - proceed to next handler
		next.ServeHTTP(w, rq)
	})
}

// validateAudience validates that the token's audience matches this resource server
func (r *OAuthMiddleware) validateAudience(claims jwt.MapClaims) bool {
	aud, ok := claims["aud"]
	if !ok {
		return false
	}

	// aud can be a string or array of strings
	switch v := aud.(type) {
	case string:
		return v == r.TargetAudience
	case []interface{}:
		for _, a := range v {
			if audStr, ok := a.(string); ok && audStr == r.TargetAudience {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// validateIssuer validates that the token's issuer matches the expected authorization server
func (r *OAuthMiddleware) validateIssuer(claims jwt.MapClaims) bool {
	iss, ok := claims["iss"].(string)
	if !ok {
		return false
	}
	return iss == r.provider.IssuerUrl
}

// validateExpiration validates that the token has not expired
func (r *OAuthMiddleware) validateExpiration(claims jwt.MapClaims) bool {
	exp, ok := claims["exp"].(float64)
	if !ok {
		return false
	}
	// Allow 60 seconds of clock skew
	return time.Now().Unix() < int64(exp)+60
}

// validateScope validates that the token has required scopes
func (r *OAuthMiddleware) validateScope(claims jwt.MapClaims) bool {
	scope, ok := claims["scope"].(string)
	if !ok {
		return false
	}
	// Scope is a space-separated string (OAuth 2.0 standard)
	// Check if "mcp:tools" is present
	for _, s := range strings.Split(scope, " ") {
		if s == "mcp:tools" {
			return true
		}
	}
	return false
}

// sendUnauthorized sends a 401 response with WWW-Authenticate header
func (r *OAuthMiddleware) sendUnauthorized(w http.ResponseWriter, rq *http.Request) {
	metadataURL := r.TargetAudience + "/.well-known/oauth-protected-resource"
	// tell client where to get resource metadata to authenticate
	w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Bearer resource_metadata="%s", scope="openid profile email"`, metadataURL))
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}

// OAuth 2.1 metadata endpoint (no authorization required)
// clients will use this to discover the resource server's authorization server and scopes
func (r *OAuthMiddleware) HandleProtectedResourceMetadata(w http.ResponseWriter, rq *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if rq.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// the resource identifier (what resource this server protects),
	// which OAuth scopes the resource supports,
	// which authorization servers (issuer URLs) are authoritative for access tokens for this resourc
	metadata := oauthex.ProtectedResourceMetadata{
		Resource:             r.TargetAudience,
		ScopesSupported:      []string{r.Scope},
		AuthorizationServers: []string{r.IssuerUrl},
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(metadata)
}
