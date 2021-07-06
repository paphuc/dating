package api

import (
	"net/http"

	userhandler "dating/internal/app/api/handler/user"
	user "dating/internal/app/api/repositories/user"
	userService "dating/internal/app/api/services/user"
	"dating/internal/app/config"

	"dating/internal/app/db"
	"dating/internal/pkg/glog"
	"dating/internal/pkg/health"
	"dating/internal/pkg/middleware"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type (
	// InfraConns holds infrastructure services connections like MongoDB, Redis, Kafka,...
	InfraConns struct {
		Databases db.Connections
	}

	middlewareFunc = func(http.HandlerFunc) http.HandlerFunc
	route          struct {
		path        string
		method      string
		handler     http.HandlerFunc
		middlewares []middlewareFunc
	}
)

const (
	get    = http.MethodGet
	post   = http.MethodPost
	put    = http.MethodPut
	delete = http.MethodDelete
)

// Init init all handlers
func Init(conns *config.Configs, em config.ErrorMessage) (http.Handler, error) {
	logger := glog.New()

	var userRepo userService.Repository
	switch conns.Database.Type {
	case db.TypeMongoDB:
		s, err := config.Dial(&conns.Database.Mongo, logger)
		if err != nil {
			logger.Panicf("failed to dial to target server, err: %v", err)
		}
		userRepo = user.NewMongoRepository(s)
	default:
		panic("database type not supported: " + conns.Database.Type)
	}

	userLogger := logger.WithField("package", "user")
	userSrv := userService.NewService(conns, &em, userRepo, userLogger)
	userHandler := userhandler.New(conns, &em, userSrv, userLogger)

	routes := []route{
		// infra
		route{
			path:    "/readiness",
			method:  get,
			handler: health.Readiness().ServeHTTP,
		},
		// services
		route{
			path:    "/signup",
			method:  post,
			handler: userHandler.SignUp,
		},
		route{
			path:    "/login",
			method:  post,
			handler: userHandler.Login,
		},
		route{
			path:        "/users/{id:[a-z0-9-\\-]+}",
			method:      get,
			middlewares: []middlewareFunc{middleware.Auth},
			handler:     userHandler.GetUserByID,
		},
		route{
			path:        "/users/update",
			method:      put,
			middlewares: []middlewareFunc{middleware.Auth},
			handler:     userHandler.UpdateUserByID,
		},
		route{
			path:        "/users",
			method:      get,
			middlewares: []middlewareFunc{middleware.Auth},
			handler:     userHandler.GetListUsers,
		},
	}

	loggingMW := middleware.Logging(logger.WithField("package", "middleware"))
	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r.Use(middleware.RequestID)
	r.Use(middleware.StatusResponseWriter)
	r.Use(loggingMW)
	r.Use(handlers.CompressHandler)

	for _, rt := range routes {
		h := rt.handler
		for _, mdw := range rt.middlewares {
			h = mdw(h)
		}
		r.Path(rt.path).Methods(rt.method).HandlerFunc(h)
	}

	return r, nil
}
