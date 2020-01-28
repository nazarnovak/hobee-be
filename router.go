package main

import (
	"context"
	"fmt"
	"github.com/zenazn/goji/web"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

	//mux.Use(getCorsHandler())
	//	mux.Get("/test/login", api.TestLogin(secret))
	//	mux.Get("/test/logout", api.TestLogout(secret))

	mux.Get("/api/identify", api.Identify(secret))
	mux.Get("/api/chat", api.Chat(secret))
	mux.Get("/api/messages", api.Messages(secret))
	mux.Get("/api/result", api.Result(secret))
	mux.Get("/api/history", api.History(secret))
	mux.Post("/api/contact", api.Contact(secret))

	compiledFEFolder := "build"

	// TODO: Temporary disable for deploys to work. Works locally tho
	//staticFiles, err := getStaticFiles(compiledFEFolder)
	//if err != nil {
	//	log.Critical(context.Background(), err)
	//}
	//
	//for _, staticFile := range staticFiles {
	//	mux.Handle(staticFile, http.FileServer(http.Dir(compiledFEFolder)))
	//}
	//
	//mux.Handle("/*", func(w http.ResponseWriter, r *http.Request) {
	//	http.ServeFile(w, r, fmt.Sprintf("%s/%s", compiledFEFolder, "index.html"))
	//})
	mux.Handle("/*", http.FileServer(http.Dir(compiledFEFolder)))

	return mux
}

func getStaticFiles(root string) ([]string, error) {
	files := []string{}

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return herrors.Wrap(err)
		}

		path = strings.TrimLeft(path, "build")

		if path == "/.DS_Store" || path == "/index.html" || path == "/static" || path == "/static/css" ||
			path == "/static/js" || path == "/static/media" {
			return nil
		}

		files = append(files, path)

		return nil
	}

	err := filepath.Walk(root, walkFn)
	if err != nil {
		return nil, herrors.Wrap(err)
	}

	return files, nil
}

//func getCorsHandler() func(http.Handler) http.Handler {
//	allowedOrigins := []string{}
//// TODO: Add mode dev + mode prod here to separate sites
//	allowedOrigins = append(allowedOrigins, "http://localhost:3000")
//// External IP
//allowedOrigins = append(allowedOrigins, "http://84.219.232.19:3000")
//
//	c := cors.New(cors.Options{
//		AllowedOrigins:   allowedOrigins,
//		AllowedHeaders:   []string{"Accept", "Authorization", "Cache-Control", "Content-Type", "Origin", "User-Agent", "Viewport", "X-Requested-With"},
//		MaxAge:           1728000,
//		AllowCredentials: true,
//		AllowedMethods:   []string{"GET"},
//	})
//
//	return c.Handler
//}
