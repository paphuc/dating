package middleware

import (
	"net/http"

	"dating/internal/app/config"
	"dating/internal/pkg/auth"
	"dating/internal/pkg/glog"
	"dating/internal/pkg/respond"
)

func Auth(h http.HandlerFunc, em *config.ErrorMessage) http.HandlerFunc {
	logger := glog.New().WithField("package", "middleware")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenpath := auth.ExtractToken(r)
		if tokenpath == "" {
			logger.Errorc(r.Context(), "The request does not contain token")
			respond.JSON(w, http.StatusUnauthorized, &em.InvalidValue.FailedAuthentication)
			return
		}
		_, err := auth.IsAuthorized(tokenpath)

		if err != nil {
			logger.Errorc(r.Context(), "Not authorized, error: %v", err)
			respond.JSON(w, http.StatusUnauthorized, &em.InvalidValue.FailedAuthentication)
			return
		}

		h.ServeHTTP(w, r)
	})
}
