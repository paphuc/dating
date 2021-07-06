package middleware

import (
	"context"
	"net/http"

	"dating/internal/pkg/auth"
	"dating/internal/pkg/glog"
	"dating/internal/pkg/respond"
)

func Auth(h http.HandlerFunc) http.HandlerFunc {
	logger := glog.New().WithField("package", "middleware")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenpath := auth.ExtractToken(r)
		if tokenpath == "" {
			logger.Infof("The request does not contain token")
			respond.JSON(w, http.StatusUnauthorized, "Failed to authentication user. (IVFA)")
			return
		}
		claimMap, err := auth.IsAuthorized(tokenpath)
		if err != nil {
			logger.Infof("Not authorized, error: ", err)
			respond.JSON(w, http.StatusUnauthorized, "Failed to authentication user. (IVFA)")
			return
		}
		ctx := context.WithValue(r.Context(), "props", claimMap)

		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}
