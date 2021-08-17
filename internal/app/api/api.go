package api

import (
	"net/http"

	userhandler "dating/internal/app/api/handler/user"
	user "dating/internal/app/api/repositories/user"
	userService "dating/internal/app/api/services/user"

	matchhandler "dating/internal/app/api/handler/match"
	match "dating/internal/app/api/repositories/match"
	matchService "dating/internal/app/api/services/match"

	messagehandler "dating/internal/app/api/handler/message"
	message "dating/internal/app/api/repositories/message"
	messageService "dating/internal/app/api/services/message"

	mailhandler "dating/internal/app/api/handler/mail"
	mail "dating/internal/app/api/repositories/mail"
	mailService "dating/internal/app/api/services/mail"

	notificationhandler "dating/internal/app/api/handler/notification"
	notification "dating/internal/app/api/repositories/notification"
	notificationService "dating/internal/app/api/services/notification"

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

	middlewareFunc = func(http.HandlerFunc, *config.ErrorMessage) http.HandlerFunc
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
	patch  = http.MethodPatch
)

// Init init all handlers
func Init(conns *config.Configs, em config.ErrorMessage) (http.Handler, error) {
	logger := glog.New()

	var userRepo userService.Repository
	var matchRepo matchService.Repository

	var messageRepo messageService.Repository
	var mailRepo mailService.Repository
	var notificationRepo notificationService.Repository

	switch conns.Database.Type {
	case db.TypeMongoDB:
		s, err := config.Dial(&conns.Database.Mongo, logger)
		if err != nil {
			logger.Panicf("failed to dial to target server, err: %v", err)
		}
		userRepo = user.NewMongoRepository(s)
		matchRepo = match.NewMongoRepository(s)

		messageRepo = message.NewMongoRepository(s)
		mailRepo = mail.NewMongoRepository(s)
		notificationRepo = notification.NewMongoRepository(s)

	default:
		panic("database type not supported: " + conns.Database.Type)
	}

	userLogger := logger.WithField("package", "user")
	userSrv := userService.NewService(conns, &em, userRepo, mailRepo, userLogger)
	userHandler := userhandler.New(conns, &em, userSrv, userLogger)

	matchLogger := logger.WithField("package", "match")
	matchSrv := matchService.NewService(conns, &em, matchRepo, matchLogger)
	matchHandler := matchhandler.New(conns, &em, matchSrv, matchLogger)

	notificationLogger := logger.WithField("package", "notification")
	notificationSrv := notificationService.NewService(conns, &em, notificationRepo, notificationLogger)
	notificationHandler := notificationhandler.New(conns, &em, notificationSrv, notificationLogger)

	messageLogger := logger.WithField("package", "chat")
	messageSrv := messageService.NewService(conns, &em, messageRepo, messageLogger, notificationSrv)
	messageHandler := messagehandler.New(conns, &em, messageSrv, messageLogger)

	mailLogger := logger.WithField("package", "mail")
	mailSrv := mailService.NewService(conns, &em, mailRepo, mailLogger)
	mailHandler := mailhandler.New(conns, &em, mailSrv, mailLogger)

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
			path:        "/users",
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
		route{
			path:        "/users/{id:[a-z0-9-\\-]+}/disable",
			method:      patch,
			middlewares: []middlewareFunc{middleware.Auth},
			handler:     userHandler.DisableUsersByID,
		},
		route{
			path:        "/users/{id:[a-z0-9-\\-]+}/available",
			method:      get,
			middlewares: []middlewareFunc{middleware.Auth},
			handler:     userHandler.GetListUsersAvailable,
		},
		route{
			path:        "/users/{id:[a-z0-9-\\-]+}/matches",
			method:      get,
			middlewares: []middlewareFunc{middleware.Auth},
			handler:     userHandler.GetMatchedUsersByID,
		},
		route{
			path:        "/matches",
			method:      post,
			middlewares: []middlewareFunc{middleware.Auth},
			handler:     matchHandler.InsertMatch,
		},
		route{
			path:        "/matches",
			method:      delete,
			middlewares: []middlewareFunc{middleware.Auth},
			handler:     matchHandler.DeleteMatched,
		},
		route{
			path:        "/matches/{id:[a-z0-9-\\-]+}",
			method:      get,
			middlewares: []middlewareFunc{middleware.Auth},
			handler:     matchHandler.GetRoomsByUserId,
		},
		route{
			path:        "/messages/{id:[a-z0-9-\\-]+}",
			method:      get,
			middlewares: []middlewareFunc{middleware.Auth},
			handler:     messageHandler.GetMessagesByIdRoom,
		},

		route{
			path:        "/notification",
			method:      post,
			middlewares: []middlewareFunc{middleware.Auth},
			handler:     notificationHandler.AddDevice,
		},
		route{
			path:        "/notification",
			method:      delete,
			middlewares: []middlewareFunc{middleware.Auth},
			handler:     notificationHandler.RemoveDevice,
		},
		route{
			path:        "/notification/{id:[a-z0-9-\\-]+}",
			method:      get,
			middlewares: []middlewareFunc{middleware.Auth},
			handler:     notificationHandler.TestSend,
		},

		route{
			path:    "/ws",
			method:  get,
			handler: messageHandler.ServeWs,
		},
		// MailVerified
		route{
			path:    "/emails",
			method:  get,
			handler: mailHandler.SendMail,
		},
		route{
			path:    "/confirmation",
			method:  get,
			handler: mailHandler.MailVerified,
		},
	}

	loggingMW := middleware.Logging(logger.WithField("package", "middleware"))
	r := mux.NewRouter()
	r.PathPrefix("/swagger").Handler(http.StripPrefix("/swagger", http.FileServer(http.Dir("./swagger-ui/"))))
	r.Use(middleware.RequestID)
	r.Use(middleware.StatusResponseWriter)
	r.Use(loggingMW)
	r.Use(handlers.CompressHandler)

	for _, rt := range routes {
		h := rt.handler
		for _, mdw := range rt.middlewares {
			h = mdw(h, &em)
		}
		r.Path(rt.path).Methods(rt.method).HandlerFunc(h)
	}

	return r, nil
}
