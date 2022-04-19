package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"tinyurl/pkg/config"
	"tinyurl/pkg/db"
	"tinyurl/pkg/log"
	"tinyurl/route"

	"github.com/flamego/flamego"
	"github.com/ory/graceful"
	// "github.com/justinas/alice"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "./configs/config.yml", "config path, eg: -conf config.yaml")
}

func main() {
	cfg, err := config.Parse(flagconf)
	if err != nil {
		panic(err)
	}
	logger, err := log.New(cfg.Log)
	if err != nil {
		panic(err)
	}
	f := flamego.Classic()
	f.Map(cfg, logger, db.DB)
	// f.Use(auth.Basic("admin", "1111"))
	route.Route(f)
	server := graceful.WithDefaults(
		&http.Server{
			Addr:    "0.0.0.0:2830",
			Handler: f,
		},
	)

	if err := graceful.Graceful(server.ListenAndServe, server.Shutdown); err != nil {
		fmt.Println("main: Failed to gracefully shutdown")
	}
	fmt.Println("main: Server was shutdown gracefully")
}
