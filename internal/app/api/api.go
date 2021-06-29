package api

import (
	"net/http"

	memberhandler "dating/internal/app/api/handler/member"
	"dating/internal/app/db"
	"dating/internal/app/member"
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
func Init(conns *InfraConns) (http.Handler, error) {
	logger := glog.New()

	var memberRepo member.Repository
	switch conns.Databases.Type {
	case db.TypeMongoDB:
		memberRepo = member.NewMongoRepository(conns.Databases.MongoDB)
	default:
		panic("database type not supported: " + conns.Databases.Type)
	}

	memberLogger := logger.WithField("package", "member")
	memberSrv := member.NewService(memberRepo, memberLogger)
	memberHandler := memberhandler.New(memberSrv, memberLogger)

	routes := []route{
		// infra
		route{
			path:    "/readiness",
			method:  get,
			handler: health.Readiness().ServeHTTP,
		},
		// services
		route{
			path:    "/api/v1/member/{id:[a-z0-9-\\-]+}",
			method:  get,
			handler: memberHandler.Get,
		},
	}

	loggingMW := middleware.Logging(logger.WithField("package", "middleware"))
	r := mux.NewRouter()
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

// Close close all underlying connections
func (c *InfraConns) Close() {
	c.Databases.Close()
}
