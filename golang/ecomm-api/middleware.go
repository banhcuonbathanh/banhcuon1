package middleware

import (
	"context"
	"english-ai-full/token"
	"fmt"
	"net/http"
	"strings"
)

type AuthKey struct{}

type Role string

const (
    RoleGuest    Role = "guest"
    RoleEmployee Role = "employee"
    RoleOwner    Role = "owner"
)

func GetAuthMiddlewareFunc(tokenMaker *token.JWTMaker) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            claims, err := verifyClaimsFromAuthHeader(r, tokenMaker)
            if err != nil {
                http.Error(w, fmt.Sprintf("error verifying token: %v", err), http.StatusUnauthorized)
                return
            }

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
    claims, err := tokenMaker.VerifyToken(token)
    if err != nil {
        return nil, fmt.Errorf("invalid token: %w", err)
    }

    return claims, nil
}


// // Middleware that allows all authenticated users
// authMiddleware := GetAuthMiddlewareFunc(tokenMaker)

// // Middleware that only allows employees and owners
// employeeAndOwnerMiddleware := GetRoleMiddlewareFunc(tokenMaker, RoleEmployee, RoleOwner)

// // Middleware that only allows owners
// ownerOnlyMiddleware := GetRoleMiddlewareFunc(tokenMaker, RoleOwner)