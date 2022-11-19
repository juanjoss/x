package api

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/juanjoss/x/pkg/shutdown"
	"go.uber.org/zap"
)

type Config struct {
	Port            string        `mapstructure:"port"`
	Hostname        string        `mapstructure:"hostname"`
	ServerTimeout   time.Duration `mapstructure:"http-server-timeout"`
	ShutdownTimeout time.Duration `mapstructure:"server-shutdown-timeout"`
}

type server struct {
	router *mux.Router
	config *Config
	logger *zap.Logger
}

func NewServer(config *Config, logger *zap.Logger) *server {
	return &server{
		router: mux.NewRouter(),
		config: config,
		logger: logger,
	}
}

func (s *server) ListenAndServe() {
	s.router.HandleFunc("/", s.index).Methods(http.MethodGet)
	s.router.HandleFunc("/health", s.healthcheck).Methods(http.MethodGet)

	srv := &http.Server{
		Addr:         ":" + s.config.Port,
		WriteTimeout: s.config.ServerTimeout,
		ReadTimeout:  s.config.ServerTimeout,
		IdleTimeout:  2 * s.config.ServerTimeout,
		Handler:      s.router,
	}

	go func() {
		log.Printf("starting HTTP server at %s", srv.Addr)
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			s.logger.Fatal("unable to start server", zap.Error(err))
		}
	}()

	// graceful shutdown
	shutdown.Graceful(srv, s.config.ShutdownTimeout)
}
