package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/rs/cors"
	"github.com/zenazn/goji/web"

	"hobee-be/api"
	"hobee-be/pkg/log"
	"hobee-be/pkg/herrors"
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
		Addr: port,
		Handler: router(secret),
	}

	ctx := context.Background()

	// Start the server and log any errors it returns
	if err := srv.ListenAndServe(); err != nil {
		log.Error(ctx, herrors.New(fmt.Sprintf("error running server: %s", err.Error())))
	}
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/Users/nazar/n/src/hobee-be/build/index.html")
}

func router(secret string) *web.Mux{
	mux := web.New()

	// TODO: Setup logging, panic recovery and tracing on the top level, we want it everywhere?
	//mux.Post("/api/register", api.Register(secret))
	//mux.Get("/api/user", api.User(secret))
	//mux.Post("/api/login", api.Login(secret))

mux.Use(getCorsHandler())
	mux.Get("/test/login", api.TestLogin(secret))
	mux.Get("/test/logout", api.TestLogout(secret))

mux.Get("/api/identify", api.Identify(secret))
mux.Get("/api/got", api.GOT(secret))
	//mux.Get("/ws", controllers.WS(secret))

	compiledFELocation := "build/"
// TODO: Figure out relative path to /build so it works on heroku prod
// Check if the files I have in the folder and if it matches - serve it, otherwise default to index.html
mux.Handle("/got", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, compiledFELocation + "index.html")
})

	mux.Handle("/*", http.FileServer(http.Dir(compiledFELocation)))

	return mux
}

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
		AllowedMethods:   []string{"GET"},
	})

	return c.Handler
}
