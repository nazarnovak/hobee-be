package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/cors"
	"github.com/zenazn/goji/web"

	"github.com/nazarnovak/hobee-be/api"
	"github.com/nazarnovak/hobee-be/pkg/herrors"
	"github.com/nazarnovak/hobee-be/pkg/log"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

type handlerFunc func(w http.ResponseWriter, r *http.Request) error

type hobeeHandler struct {
	handler handlerFunc
}

func (hh *hobeeHandler) ServeHTTPC(w http.ResponseWriter, r *http.Request) {
	if err := hh.handler(w, r); err != nil {
		// do something about the error
	}

	// if no error show 200
}

type httpServer interface {
	ServeHTTPC(http.ResponseWriter, *http.Request)
}

func HobeeHandler(hf handlerFunc) httpServer {
	return &hobeeHandler{handler: hf}
}

func (s *Server) Start(secret, port string) {
	srv := &http.Server{
		Addr:    port,
		Handler: router(secret),
	}

	ctx := context.Background()

	// Start the server and log any errors it returns
	if err := srv.ListenAndServe(); err != nil {
		log.Error(ctx, herrors.New(fmt.Sprintf("error running server: %s", err.Error())))
	}
}

func router(secret string) *web.Mux {
	mux := web.New()

	// TODO: Setup logging, panic recovery and tracing on the top level, we want it everywhere?
	//mux.Post("/api/register", api.Register(secret))
	//mux.Get("/api/user", api.User(secret))
	//mux.Post("/api/login", api.Login(secret))

	mux.Use(getCorsHandler())
	//mux.Get("/test/login", api.TestLogin(secret))
	//mux.Get("/test/logout", api.TestLogout(secret))

	mux.Get("/api/identify", api.Identify(secret))
	mux.Get("/api/chat", api.Chat(secret))
	mux.Get("/api/messages", api.Messages(secret))
	mux.Get("/api/online", api.Online(secret))
	mux.Get("/api/result", api.Result(secret))
	mux.Get("/api/history", api.History(secret))
	mux.Post("/api/contact", api.Contact(secret))
	mux.Post("/api/feedback", api.Feedback(secret))

	//mux.Get("/api/test", Test())

	return mux
}

//func Test() func(w http.ResponseWriter, r *http.Request) {
//	return func(w http.ResponseWriter, r *http.Request) {
//		fmt.Println("Test hello!")
//		w.WriteHeader(200)
//		w.Write([]byte("Hello!"))
//	}
//}

func getCorsHandler() func(http.Handler) http.Handler {
	allowedOrigins := []string{}
// TODO: Add mode dev + mode prod here to separate sites
	allowedOrigins = append(allowedOrigins, "http://localhost:3000")
	// External IP
	allowedOrigins = append(allowedOrigins, "http://84.219.232.19:3000")

	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedHeaders:   []string{"Accept", "Authorization", "Cache-Control", "Content-Type", "Origin", "User-Agent", "Viewport", "X-Requested-With"},
		MaxAge:           1728000,
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST"},
	})

	return c.Handler
}
