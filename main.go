package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"tinyurl/pkg/config"
	"tinyurl/pkg/db"
	"tinyurl/pkg/log"
	"tinyurl/route"

	"github.com/flamego/cors"
	"github.com/flamego/flamego"
	lru "github.com/hashicorp/golang-lru"
	"github.com/ory/graceful"
	// "github.com/justinas/alice"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version        = "unknown"
	GitCommitLog   = "unknown"
	GitStatus      = "unknown"
	BuildTime      = "unknown"
	BuildGoVersion = "unknown"
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func StringifyMultiLine() string {
	return fmt.Sprintf("Version=%s\nGitCommitLog=%s\nBuildTime=%s\nGoVersion=%s\nruntime=%s/%s\n",
		Version, GitCommitLog, BuildTime, BuildGoVersion, runtime.GOOS, runtime.GOARCH)
}

func init() {
	flag.StringVar(&flagconf, "conf", "./configs/config.yml", "config path, eg: -conf config.yaml")
	db.InitDb()
}

func main() {
	fmt.Println(StringifyMultiLine())
	cfg, err := config.Parse(flagconf)
	if err != nil {
		panic(err)
	}
	logger, err := log.New(cfg.Log)
	if err != nil {
		panic(err)
	}
	l, err := lru.New(10240)
	if err != nil {
		panic(err)
	}
	f := flamego.Classic()
	f.Use(cors.CORS(
		cors.Options{
			Methods:          []string{"POST", "GET", "OPTIONS", "PUT", "DELETE", "UPDATE"},
			AllowDomain:      []string{"*"},
			AllowCredentials: true,
		},
	))
	f.Use(flamego.Static(flamego.StaticOptions{
		Directory: "./dist/",
	}))
	// f.Use(flamego.Static(flamego.StaticOptions{
	// 	Directory: "./dist/assets/",
	// }))
	f.Map(cfg, logger, db.DB, l)
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
