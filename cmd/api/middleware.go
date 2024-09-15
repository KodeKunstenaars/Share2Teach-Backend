package main

import (
	"net/http"
)

func (app *application) enableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, X-CSRF-Token, Authorization")
			return
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

//func (app *application) authRequired(next http.Handler, requiredRole string) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		token, claims, err := app.auth.GetTokenFromHeaderAndVerify(w, r)
//		if err != nil {
//			http.Error(w, "Unauthorized", http.StatusUnauthorized)
//			return
//		}
//
//		if !token.Valid {
//			http.Error(w, "Invalid token", http.StatusUnauthorized)
//			return
//		}
//
//		userRole := claims.Role
//
//		if requiredRole != "" && userRole != requiredRole {
//			http.Error(w, "Forbidden - insufficient permissions", http.StatusForbidden)
//			return
//		}
//
//		next.ServeHTTP(w, r)
//	})
//}

func (app *application) authRequired(next http.Handler, allowedRoles ...string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// extract token and claims
		token, claims, err := app.auth.GetTokenFromHeaderAndVerify(w, r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		userRole := claims.Role

		// Check if the user's role is in the allowedRoles list
		if len(allowedRoles) > 0 {
			allowed := false
			for _, role := range allowedRoles {
				if userRole == role {
					allowed = true
					break
				}
			}
			if !allowed {
				http.Error(w, "Forbidden - insufficient permissions", http.StatusForbidden)
				return
			}
		}

		// Proceed to the next handler if role is authorized
		next.ServeHTTP(w, r)
	})
}
