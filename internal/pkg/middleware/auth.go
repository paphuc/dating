package middleware

import (
	"context"
	"dating/internal/pkg/auth"
	"dating/internal/pkg/respond"
	"net/http"

	"github.com/pkg/errors"
)

func Auth(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenpath := auth.ExtractToken(r)
		if tokenpath == "" {
			respond.Error(w, errors.New("You not login"), http.StatusUnauthorized)
			return
		}
		claimMap, err := auth.IsAuthorized(tokenpath)
		// fmt.Println(claimMap)
		if err != nil {
			respond.Error(w, err, http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "props", claimMap)

		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}
