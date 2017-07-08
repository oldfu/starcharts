package main

import (
	"net/http"
	"os"
	"time"

	"github.com/apex/httplog"
	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/caarlos0/starchart/internal/cache"
	"github.com/caarlos0/starchart/internal/config"
	"github.com/caarlos0/starchart/internal/controller"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func init() {
	log.SetHandler(text.New(os.Stderr))
}

func main() {
	var config = config.Get()
	var ctx = log.WithField("port", config.Port)
	var cache = cache.New(config.RedisURL)
	defer cache.Close()

	var r = mux.NewRouter()
	r.Path("/").
		Methods(http.MethodGet).
		HandlerFunc(controller.Index())
	r.Path("/{owner}/{repo}.svg").
		Methods(http.MethodGet).
		HandlerFunc(controller.GetRepoChart(config, cache))
	r.Path("/{owner}/{repo}").
		Methods(http.MethodGet).
		HandlerFunc(controller.GetRepo())

	var srv = &http.Server{
		Handler:      httplog.New(handlers.CompressHandler(r)),
		Addr:         "0.0.0.0:" + config.Port,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}
	ctx.Info("starting up...")
	ctx.WithError(srv.ListenAndServe()).Error("failed to start up server")
}
