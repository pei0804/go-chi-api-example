package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/siddontang/go/log"
	"github.com/sirupsen/logrus"
)

// Server Server
type Server struct {
	router *chi.Mux
	logger *logrus.Logger
}

// New Server構造体のコンストラクタ
func New() *Server {
	return &Server{
		router: chi.NewRouter(),
		logger: logrus.New(),
	}
}

// Init 実行時にしたいこと
func (s *Server) Init(env string) {
	if env == "production" {
		s.logger.Formatter = &logrus.JSONFormatter{}
		s.logger.SetLevel(logrus.WarnLevel)
	} else {
		s.logger.Formatter = &logrus.TextFormatter{}
		s.logger.SetLevel(logrus.DebugLevel)
	}
}

// Middleware ミドルウェア
func (s *Server) Middleware() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.CloseNotify)
	s.router.Use(middleware.Timeout(time.Second * 60))
	s.router.Use(NewStructuredLogger(s.logger))
}

// Router ルーティング設定
func (s *Server) Router() {
	h := NewHandler(s.logger)
	s.router.Route("/api", func(api chi.Router) {
		api.Use(Auth("db connection"))
		api.Route("/members", func(members chi.Router) {
			members.Get("/{id}", h.Show)
			members.Get("/", h.List)
		})
	})
	s.router.Route("/api/auth", func(auth chi.Router) {
		auth.Get("/login", h.Login)
	})
}

// Auth 認証（dbはフェイク）
func Auth(db string) (fn func(http.Handler) http.Handler) {
	fn = func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.Header.Get("Auth")
			if token != "admin" {
				respondError(w, http.StatusUnauthorized, fmt.Errorf("利用権限がありません"))
				return
			}
			h.ServeHTTP(w, r)
		})
	}
	return
}

// NewStructuredLogger loggingに使うloggerを置き換える
func NewStructuredLogger(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{logger})
}

// StructuredLogger ロガー用
type StructuredLogger struct {
	Logger *logrus.Logger
}

// NewLogEntry ログ書き出し
func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{Logger: logrus.NewEntry(l.Logger)}
	logFields := logrus.Fields{}
	logFields["ts"] = time.Now().UTC().Format(time.RFC1123)
	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields["req_id"] = reqID
	}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	logFields["http_scheme"] = scheme
	logFields["http_proto"] = r.Proto
	logFields["http_method"] = r.Method

	logFields["remote_addr"] = r.RemoteAddr
	logFields["user_agent"] = r.UserAgent()

	logFields["uri"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)
	entry.Logger = entry.Logger.WithFields(logFields)
	entry.Logger.Infoln("request started")
	return entry
}

// StructuredLoggerEntry ログエントリ用
type StructuredLoggerEntry struct {
	Logger logrus.FieldLogger
}

// Panic 通常ログ
func (l *StructuredLoggerEntry) Write(status, bytes int, elapsed time.Duration) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"resp_status": status, "resp_bytes_length": bytes,
		"resp_elapsed_ms": float64(elapsed.Nanoseconds()) / 1000000.0,
	})

	l.Logger.Infoln("request complete")
}

// Panic パニック時用のログ
func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}

// Handler ハンドラ用
type Handler struct {
	logger *logrus.Logger
}

// NewHandler コンストラクタ
func NewHandler(logger *logrus.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

// Show endpoint
func (h *Handler) Show(w http.ResponseWriter, r *http.Request) {
	type json struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	res := json{ID: id, Name: fmt.Sprint("name_", id)}
	respondJSON(w, http.StatusOK, res)
}

// List endpoint
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	users := []struct {
		ID   int    `json:"id"`
		User string `json:"user"`
	}{
		{1, "hoge"},
		{2, "foo"},
		{3, "bar"},
	}
	respondJSON(w, http.StatusOK, users)
}

// Login endpoint
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token != "token" {
		respondError(w, http.StatusUnauthorized, fmt.Errorf("有効でないトークンです: %s", token))
		return
	}
	type json struct {
		Message string `json:"message"`
	}
	res := json{Message: "auth ok"}
	respondJSON(w, http.StatusOK, res)
}

// respondJSON レスポンスとして返すjsonを生成して、writerに書き込む
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "    ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

// respondError レスポンスとして返すエラーを生成する
func respondError(w http.ResponseWriter, code int, err error) {
	log.Error("err: ", err)
	if e, ok := err.(*HTTPError); ok {
		respondJSON(w, e.Code, e)
	} else if err != nil {
		he := HTTPError{
			Code:    code,
			Message: err.Error(),
		}
		respondJSON(w, code, he)
	}
}

// HTTPError エラー用
type HTTPError struct {
	Code    int
	Message string
}

func (he *HTTPError) Error() string {
	return fmt.Sprintf("code=%d, message=%v", he.Code, he.Message)
}

func main() {
	var (
		port = flag.String("port", "8080", "addr to bind")
		env  = flag.String("env", "develop", "実行環境 (production, staging, develop)")
	)
	flag.Parse()
	s := New()
	s.Init(*env)
	s.Middleware()
	s.Router()
	s.logger.Info("Starting app")
	http.ListenAndServe(fmt.Sprint(":", *port), s.router)
}
