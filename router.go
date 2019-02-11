package main

import (
	"context"
	"fmt"
	"hobee-be/pkg/herrors"
	"net/http"

	"github.com/zenazn/goji/web"

	"hobee-be/api"
	"hobee-be/pkg/log"
)

type Server struct {}

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
		Handler: router(secret),
	}

	ctx := context.Background()

	// Start the server and log any errors it returns
	if err := srv.ListenAndServe(); err != nil {
		log.Error(ctx, herrors.New(fmt.Sprintf("error running server: %s", err.Error())))
	}
}

func router(secret string) *web.Mux{
	mux := web.New()

	// TODO: Setup logging, panic recovery and tracing on the top level, we want it everywhere?
	mux.Post("/api/register", api.Register(secret))
	mux.Get("/api/user", api.User(secret))
	mux.Post("/api/login", api.Login(secret))

	return mux
}
