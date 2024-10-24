package middleware

import (
	"context"
	"english-ai-full/token"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type AuthKey struct{}

// --------------
type TableAuthKey struct{} // new add

type AuthMiddleware struct {
    tokenMaker *token.JWTMaker
    logger     *log.Logger
}
func NewAuthMiddleware(tokenMaker *token.JWTMaker, logger *log.Logger) *AuthMiddleware {
    if logger == nil {
        logger = log.New(log.Writer(), "AUTH-MIDDLEWARE: ", log.LstdFlags)
    }
    return &AuthMiddleware{
        tokenMaker: tokenMaker,
        logger:     logger,
    }
}
// VerifyAccessToken verifies the Bearer token in the Authorization header
func (am *AuthMiddleware) VerifyAccessToken() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            am.logger.Println("Verifying access token")
            claims, err := am.verifyClaimsFromAuthHeader(r)
            if err != nil {
                am.logger.Printf("Error verifying access token: %v", err)
                http.Error(w, fmt.Sprintf("unauthorized: %v", err), http.StatusUnauthorized)
                return
            }

            am.logger.Printf("Access token verified successfully for user ID: %d, role: %s", claims.ID, claims.Role)
            ctx := context.WithValue(r.Context(), AuthKey{}, claims)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
func (am *AuthMiddleware) VerifyTableToken() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            am.logger.Println("Verifying table token")
            tableToken := r.Header.Get("X-Table-Token")
            if tableToken == "" {
                am.logger.Println("Table token missing")
                http.Error(w, "table token is required", http.StatusUnauthorized)
                return
            }

            claims, err := am.tokenMaker.VerifyShortToken(tableToken)
            if err != nil {
                am.logger.Printf("Error verifying table token: %v", err)
                http.Error(w, "invalid table token", http.StatusUnauthorized)
                return
            }

            am.logger.Printf("Table token verified successfully for table ID: %d", claims.ID)
            ctx := context.WithValue(r.Context(), TableAuthKey{}, claims)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
// RequireRoles middleware checks if the user has any of the required roles
func (am *AuthMiddleware) RequireRoles(allowedRoles ...Role) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            claims, err := am.verifyClaimsFromAuthHeader(r)
            if err != nil {
                am.logger.Printf("Error verifying token in role check: %v", err)
                http.Error(w, fmt.Sprintf("unauthorized: %v", err), http.StatusUnauthorized)
                return
            }

            if !isRoleAllowed(Role(claims.Role), allowedRoles) {
                am.logger.Printf("User role %s not allowed. Required roles: %v", claims.Role, allowedRoles)
                http.Error(w, "forbidden: insufficient permissions", http.StatusForbidden)
                return
            }

            am.logger.Printf("Role verification successful for user ID: %d, role: %s", claims.ID, claims.Role)
            ctx := context.WithValue(r.Context(), AuthKey{}, claims)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
// Helper functions
func (am *AuthMiddleware) verifyClaimsFromAuthHeader(r *http.Request) (*token.UserClaims, error) {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        return nil, fmt.Errorf("authorization header is missing")
    }

    fields := strings.Fields(authHeader)
    if len(fields) != 2 || fields[0] != "Bearer" {
        return nil, fmt.Errorf("invalid authorization header format")
    }

    token := fields[1]
    claims, err := am.tokenMaker.VerifyToken(token)
    if err != nil {
        return nil, fmt.Errorf("invalid token: %w", err)
    }

    return claims, nil
}





// Helper functions to extract claims from context
func GetAccessClaims(ctx context.Context) (*token.UserClaims, error) {
    claims, ok := ctx.Value(AuthKey{}).(*token.UserClaims)
    if !ok {
        return nil, fmt.Errorf("access claims not found in context")
    }
    return claims, nil
}

func GetTableClaims(ctx context.Context) (*token.UserClaims, error) {
    claims, ok := ctx.Value(TableAuthKey{}).(*token.UserClaims)
    if !ok {
        return nil, fmt.Errorf("table claims not found in context")
    }
    return claims, nil
}

//------------
type Role string

const (
    RoleGuest    Role = "guest"
    RoleEmployee Role = "employee"
    RoleOwner    Role = "owner"
)

func GetAuthMiddlewareFunc(tokenMaker *token.JWTMaker) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            log.Println("Auth middleware started")
            claims, err := verifyClaimsFromAuthHeader(r, tokenMaker)
            if err != nil {
                log.Printf("Error verifying token: %v", err)
                http.Error(w, fmt.Sprintf("error verifying token: %v", err), http.StatusUnauthorized)
                return
            }

            log.Printf("Claims verified successfully: %+v", claims)
            ctx := context.WithValue(r.Context(), AuthKey{}, claims)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func GetRoleMiddlewareFunc(tokenMaker *token.JWTMaker, allowedRoles ...Role) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            claims, err := verifyClaimsFromAuthHeader(r, tokenMaker)
            if err != nil {
                http.Error(w, fmt.Sprintf("error verifying token: %v", err), http.StatusUnauthorized)
                return
            }

            if !isRoleAllowed(Role(claims.Role), allowedRoles) {
                http.Error(w, "user does not have the required role", http.StatusForbidden)
                return
            }

            ctx := context.WithValue(r.Context(), AuthKey{}, claims)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func isRoleAllowed(userRole Role, allowedRoles []Role) bool {
    for _, role := range allowedRoles {
        if userRole == role {
            return true
        }
    }
    return false
}

func verifyClaimsFromAuthHeader(r *http.Request, tokenMaker *token.JWTMaker) (*token.UserClaims, error) {
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        return nil, fmt.Errorf("authorization header is missing")
    }

    fields := strings.Fields(authHeader)
    if len(fields) != 2 || fields[0] != "Bearer" {
        return nil, fmt.Errorf("invalid authorization header")
    }

    token := fields[1]
    log.Printf("Token received: %s", token) // Log the token (be careful with this in production)
    claims, err := tokenMaker.VerifyToken(token)
    if err != nil {
        log.Printf("Error verifying token: %v", err)
        return nil, fmt.Errorf("invalid token: %w", err)
    }

    log.Printf("Claims verified: %+v", claims)
    return claims, nil
}


// // Middleware that allows all authenticated users
// authMiddleware := GetAuthMiddlewareFunc(tokenMaker)

// // Middleware that only allows employees and owners
// employeeAndOwnerMiddleware := GetRoleMiddlewareFunc(tokenMaker, RoleEmployee, RoleOwner)

// // Middleware that only allows owners
// ownerOnlyMiddleware := GetRoleMiddlewareFunc(tokenMaker, RoleOwner)