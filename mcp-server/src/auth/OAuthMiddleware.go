package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type OAuthMiddleware struct {
	provider          *OauthProvider
	TargetAudienceUrl string
	IssuerUrl         string
}

func (r *OAuthMiddleware) Init() {
	r.provider = &OauthProvider{
		IssuerUrl: r.IssuerUrl,
	}
	r.provider.Init()
}

func (r *OAuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, rq *http.Request) {
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
			log.Printf("Failed to parse token: %v", err)
			r.sendUnauthorized(w, rq)
			return
		}

		if !token.Valid {
			log.Printf("Invalid token")
			r.sendUnauthorized(w, rq)
			return
		}

		// Get claims for validation
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Printf("Invalid claims type")
			r.sendUnauthorized(w, rq)
			return
		}

		// Debug: Dump JWT access token before validation
		log.Printf("=== JWT Access Token Debug ===")
		log.Printf("Raw Token: %s", tokenString)
		claimsJSON, _ := json.MarshalIndent(claims, "", "  ")
		log.Printf("Claims: %s", string(claimsJSON))
		log.Printf("===============================")

		// Validate audience (MUST): Verify this resource server is in the audience
		if !r.validateAudience(claims) {
			log.Printf("Invalid audience")
			r.sendUnauthorized(w, rq)
			return
		}

		// Validate issuer (MUST): Verify token is issued by expected authorization server
		if !r.validateIssuer(claims) {
			log.Printf("Invalid issuer")
			r.sendUnauthorized(w, rq)
			return
		}

		// Validate expiration (MUST): Ensure token is not expired
		// Note: jwt.Parse already validates exp by default, but we explicitly check here for clarity
		if !r.validateExpiration(claims) {
			log.Printf("Token expired")
			r.sendUnauthorized(w, rq)
			return
		}

		// Validate scope: Verify token has required scopes (optional, depends on your requirements)
		if !r.validateScope(claims) {
			log.Printf("Insufficient scope")
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
		return v == r.TargetAudienceUrl
	case []interface{}:
		for _, a := range v {
			if audStr, ok := a.(string); ok && audStr == r.TargetAudienceUrl {
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
	metadataURL := r.TargetAudienceUrl + "/.well-known/oauth-protected-resource"
	// tell client where to get resource metadata to authenticate
	w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Bearer resource_metadata="%s", scope="openid profile email"`, metadataURL))
	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}
