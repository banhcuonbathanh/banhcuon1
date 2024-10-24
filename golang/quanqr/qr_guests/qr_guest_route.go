package qr_guests

import (
	middleware "english-ai-full/ecomm-api"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	// middleware "english-ai-full/ecomm-api"
)

func RegisterGuestRoutes(r *chi.Mux, handler *GuestHandlerController) *chi.Mux {


	r.Get("/qr/guest/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("guest test is running"))
	})




	//----------------------------
	tokenMaker := handler.TokenMaker

    authMiddleware := middleware.NewAuthMiddleware(tokenMaker, nil)

	//----------------------------
	r.Route("/qr/guest", func(r chi.Router) {
		log.Print("golang/quanqr/qr_guests/qr_guest_route.go 12341234123")
		r.Use(authMiddleware.VerifyTableToken())
		r.Post("/login", handler.GuestLogin)
		r.Post("/refresh-token", handler.RefreshToken)

		// Protected routes (authentication required)
		r.Group(func(r chi.Router) {

			// log.Print("golang/quanqr/qr_guests/qr_guest_route.go")
			// r.Use(middleware.GetAuthMiddlewareFunc(tokenMaker))

			r.Post("/logout", handler.GuestLogout)
			// r.Post("/orders", handler.CreateOrders)
			// r.Get("/orders/{guestId}", handler.GetOrders)
		})
	})

	return r
}







// package middleware

// import (
//     "context"
//     "fmt"
//     "log"
//     "net/http"
//     "strings"
// )

// type AuthKey struct{}
// type TableAuthKey struct{}

// type Role string

// const (
//     RoleGuest    Role = "guest"
//     RoleEmployee Role = "employee"
//     RoleOwner    Role = "owner"
// )

// type AuthMiddleware struct {
//     tokenMaker *token.JWTMaker
//     logger     *log.Logger
// }

// func NewAuthMiddleware(tokenMaker *token.JWTMaker, logger *log.Logger) *AuthMiddleware {
//     if logger == nil {
//         logger = log.New(log.Writer(), "AUTH-MIDDLEWARE: ", log.LstdFlags)
//     }
//     return &AuthMiddleware{
//         tokenMaker: tokenMaker,
//         logger:     logger,
//     }
// }

// // VerifyAccessToken verifies the Bearer token in the Authorization header
// func (am *AuthMiddleware) VerifyAccessToken() func(http.Handler) http.Handler {
//     return func(next http.Handler) http.Handler {
//         return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//             am.logger.Println("Verifying access token")
//             claims, err := am.verifyClaimsFromAuthHeader(r)
//             if err != nil {
//                 am.logger.Printf("Error verifying access token: %v", err)
//                 http.Error(w, fmt.Sprintf("unauthorized: %v", err), http.StatusUnauthorized)
//                 return
//             }

//             am.logger.Printf("Access token verified successfully for user ID: %d, role: %s", claims.ID, claims.Role)
//             ctx := context.WithValue(r.Context(), AuthKey{}, claims)
//             next.ServeHTTP(w, r.WithContext(ctx))
//         })
//     }
// }

// // VerifyTableToken verifies the table token from X-Table-Token header
// func (am *AuthMiddleware) VerifyTableToken() func(http.Handler) http.Handler {
//     return func(next http.Handler) http.Handler {
//         return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//             am.logger.Println("Verifying table token")
//             tableToken := r.Header.Get("X-Table-Token")
//             if tableToken == "" {
//                 am.logger.Println("Table token missing")
//                 http.Error(w, "table token is required", http.StatusUnauthorized)
//                 return
//             }

//             claims, err := am.tokenMaker.VerifyShortToken(tableToken)
//             if err != nil {
//                 am.logger.Printf("Error verifying table token: %v", err)
//                 http.Error(w, "invalid table token", http.StatusUnauthorized)
//                 return
//             }

//             am.logger.Printf("Table token verified successfully for table ID: %d", claims.ID)
//             ctx := context.WithValue(r.Context(), TableAuthKey{}, claims)
//             next.ServeHTTP(w, r.WithContext(ctx))
//         })
//     }
// }

// // RequireRoles middleware checks if the user has any of the required roles
// func (am *AuthMiddleware) RequireRoles(allowedRoles ...Role) func(http.Handler) http.Handler {
//     return func(next http.Handler) http.Handler {
//         return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//             claims, err := am.verifyClaimsFromAuthHeader(r)
//             if err != nil {
//                 am.logger.Printf("Error verifying token in role check: %v", err)
//                 http.Error(w, fmt.Sprintf("unauthorized: %v", err), http.StatusUnauthorized)
//                 return
//             }

//             if !isRoleAllowed(Role(claims.Role), allowedRoles) {
//                 am.logger.Printf("User role %s not allowed. Required roles: %v", claims.Role, allowedRoles)
//                 http.Error(w, "forbidden: insufficient permissions", http.StatusForbidden)
//                 return
//             }

//             am.logger.Printf("Role verification successful for user ID: %d, role: %s", claims.ID, claims.Role)
//             ctx := context.WithValue(r.Context(), AuthKey{}, claims)
//             next.ServeHTTP(w, r.WithContext(ctx))
//         })
//     }
// }

// // Helper functions
// func (am *AuthMiddleware) verifyClaimsFromAuthHeader(r *http.Request) (*token.UserClaims, error) {
//     authHeader := r.Header.Get("Authorization")
//     if authHeader == "" {
//         return nil, fmt.Errorf("authorization header is missing")
//     }

//     fields := strings.Fields(authHeader)
//     if len(fields) != 2 || fields[0] != "Bearer" {
//         return nil, fmt.Errorf("invalid authorization header format")
//     }

//     token := fields[1]
//     claims, err := am.tokenMaker.VerifyToken(token)
//     if err != nil {
//         return nil, fmt.Errorf("invalid token: %w", err)
//     }

//     return claims, nil
// }

// func isRoleAllowed(userRole Role, allowedRoles []Role) bool {
//     for _, role := range allowedRoles {
//         if userRole == role {
//             return true
//         }
//     }
//     return false
// }

// // Helper functions to extract claims from context
// func GetAccessClaims(ctx context.Context) (*token.UserClaims, error) {
//     claims, ok := ctx.Value(AuthKey{}).(*token.UserClaims)
//     if !ok {
//         return nil, fmt.Errorf("access claims not found in context")
//     }
//     return claims, nil
// }

// func GetTableClaims(ctx context.Context) (*token.UserClaims, error) {
//     claims, ok := ctx.Value(TableAuthKey{}).(*token.UserClaims)
//     if !ok {
//         return nil, fmt.Errorf("table claims not found in context")
//     }
//     return claims, nil
// }

// // Example usage in routes:
// func ExampleRouteSetup(r *chi.Mux, tokenMaker *token.JWTMaker) {
//     authMiddleware := NewAuthMiddleware(tokenMaker, nil)

//     // Public routes
//     r.Get("/health", healthCheck)

//     // Guest routes with table token verification
//     r.Group(func(r chi.Router) {
//         r.Use(authMiddleware.VerifyTableToken())
//         r.Post("/guest/login", guestLogin)
//     })

//     // Protected guest routes
//     r.Group(func(r chi.Router) {
//         r.Use(authMiddleware.VerifyAccessToken())
//         r.Use(authMiddleware.RequireRoles(RoleGuest))
//         r.Post("/orders", createOrder)
//         r.Get("/orders/{id}", getOrder)
//     })

//     // Employee routes
//     r.Group(func(r chi.Router) {
//         r.Use(authMiddleware.VerifyAccessToken())
//         r.Use(authMiddleware.RequireRoles(RoleEmployee, RoleOwner))
//         r.Get("/tables", getTables)
//         r.Post("/tables", createTable)
//     })

//     // Owner-only routes
//     r.Group(func(r chi.Router) {
//         r.Use(authMiddleware.VerifyAccessToken())
//         r.Use(authMiddleware.RequireRoles(RoleOwner))
//         r.Post("/employees", createEmployee)
//         r.Delete("/employees/{id}", deleteEmployee)
//     })
// }

// // Example handler showing how to access claims
// func ExampleHandler(w http.ResponseWriter, r *http.Request) {
//     // Get access claims
//     claims, err := GetAccessClaims(r.Context())
//     if err != nil {
//         http.Error(w, "unauthorized", http.StatusUnauthorized)
//         return
//     }

//     // Use claims
//     fmt.Printf("Request from user ID: %d, role: %s\n", claims.ID, claims.Role)
// }