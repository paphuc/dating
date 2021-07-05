package middleware

import (
	"context"
	"dating/internal/app/config"
	"dating/internal/pkg/auth"
	"dating/internal/pkg/respond"
	"net/http"
)

func Auth(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenpath := auth.ExtractToken(r)
		if tokenpath == "" {
			respond.JSON(w, http.StatusUnauthorized, config.EM.Invalid_value.FailedAuthentication)
			return
		}
		claimMap, err := auth.IsAuthorized(tokenpath)

		if err != nil {
			respond.JSON(w, http.StatusUnauthorized, config.EM.Invalid_value.FailedAuthentication)
			return
		}
		ctx := context.WithValue(r.Context(), "props", claimMap)

		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}
