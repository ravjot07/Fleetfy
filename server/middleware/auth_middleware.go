package middleware

import (
    "net/http"
    "strings"
)

func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            role := r.Header.Get("Role") // In real-world apps, this should be from JWT or session

            for _, allowedRole := range allowedRoles {
                if strings.EqualFold(role, allowedRole) {
                    next.ServeHTTP(w, r)
                    return
                }
            }

            http.Error(w, "Forbidden", http.StatusForbidden)
        })
    }
}
