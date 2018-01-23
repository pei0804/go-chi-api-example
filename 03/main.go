package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/ascarter/requestid"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/docgen"
	"github.com/jmoiron/sqlx"
)

// Server Server
type Server struct {
	router *chi.Mux
	db     *sqlx.DB
}

// New Server構造体のコンストラクタ
func New() *Server {
	return &Server{
		router: chi.NewRouter(),
	}
}

// Init 実行時にしたいこと
func (s *Server) Init(env string) {
	cs, err := NewConfigsFromFile(filepath.Join("db.yml"))
	if err != nil {
		log.Fatalf("cannot open database configuration. please, $make config-set. faild: %s", err)
	}
	s.db, err = cs.Open(env)
	if err != nil {
		log.Fatalf("database initialization failed: %s", err)
	}
	if s.db.Ping() != nil {
		log.Fatalf("database ping failed: %s", s.db.Ping())
	}
}

// Middleware ミドルウェア
func (s *Server) Middleware(env string) {
	s.router.Use(CorsConfig[env].Handler)
	s.router.Use(requestid.RequestIDHandler)
	s.router.Use(middleware.CloseNotify)
	s.router.Use(loggingMiddleware)
	s.router.Use(middleware.Timeout(time.Second * 60))
}

// Router ルーティング設定
func (s *Server) Router() {
	c := NewController()
	s.router.Route("/api", func(api chi.Router) {
		api.Use(Auth("db connection"))
		api.Route("/members", func(members chi.Router) {
			members.Get("/{id}", handler(c.Show).ServeHTTP)
			members.Get("/", handler(c.List).ServeHTTP)
		})
	})
	s.router.Route("/api/auth", func(auth chi.Router) {
		auth.Get("/login", handler(c.Login).ServeHTTP)
	})
}

func main() {
	var (
		port   = flag.String("port", "8080", "addr to bind")
		env    = flag.String("env", "develop", "実行環境 (production, staging, develop)")
		gendoc = flag.Bool("gendoc", true, "ドキュメント自動生成")
	)
	flag.Parse()
	s := New()
	s.Init(*env)
	s.Middleware(*env)
	s.Router()
	logcfg := DevLogConfig([]string{"api"})
	logcfg.Format = "%{time} [%{module}:%{level:.4s}]%{message}[%{shortfile}]"
	logcfg.Colors = true
	logcfg.InitLogger()
	logger := MustGetLogger("api")
	logger.Error("aaa")
	logger.Info("a")
	if *gendoc {
		doc := docgen.MarkdownRoutesDoc(s.router, docgen.MarkdownOpts{
			ProjectPath: "github.com/pei0804/go-chi-api-example",
			Intro:       "generated docs.",
		})
		file, err := os.Create("doc/doc.md")
		if err != nil {
			log.Printf("err: %v", err)
		}
		defer file.Close()
		file.Write(([]byte)(doc))
	}
	http.ListenAndServe(fmt.Sprint(":", *port), s.router)
}
