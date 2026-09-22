package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yoannduc/go-project-template/internal/domain"
	"github.com/yoannduc/go-project-template/internal/dtos"
	examplehdl "github.com/yoannduc/go-project-template/internal/handlers/example"
	examplerepo "github.com/yoannduc/go-project-template/internal/repositories/example"
	examplesrv "github.com/yoannduc/go-project-template/internal/services/example"
	"github.com/yoannduc/go-project-template/pkg/env"
	"github.com/yoannduc/go-project-template/pkg/logger"
	"github.com/yoannduc/go-project-template/pkg/mapper"
	"github.com/yoannduc/go-project-template/pkg/memorydb"
	httpmw "github.com/yoannduc/go-project-template/pkg/middlewares/http"
)

const (
	httpAddrEnvVar = "GO_HTTP_ADDRESS"
)

var (
	httpAddr = ":8080"
)

func main() {
	mux := http.NewServeMux()

	var h http.Handler = mux
	h = httpmw.SuccessJSONContentType(h)
	h = httpmw.LoadUserIP(h)
	h = httpmw.LoadTraceID(h)
	h = httpmw.Timeout(1 * time.Second)(h)
	h = httpmw.LogResponse(h)
	// Add cors only on dev env
	if env.IsDev() {
		h = httpmw.DevCORS(h)
	}

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain;charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Example microservice"))
	})

	examplehdl.New(
		examplesrv.New(
			examplerepo.NewMemoryRepository(memorydb.NewMemoryDB[domain.Example]()),
			mapper.New[domain.Example, dtos.Example](),
		),
	).LoadRoutes(mux)

	// Override default httpAddr if env var is present
	if a := os.Getenv(httpAddrEnvVar); a != "" {
		httpAddr = a
	}

	// Explicitly create http server & set timeout
	// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	// https://adam-p.ca/blog/2022/01/golang-http-server-timeouts/
	s := &http.Server{
		Addr:    httpAddr,
		Handler: h,
		// ReadTimeout is set to 60s to allow files import
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		// IdleTimeout is the maximum amount of time to wait for the next
		// request when keep-alives are enabled.
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	// https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/notify-with-context/server.go
	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Get().LogAttrs(context.Background(), logger.LevelFatal, "Server start error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	logger.Get().Info("Server started on address " + httpAddr)

	// Create context that listens for the interrupt signal from the OS.
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of
	// shutdown.
	stop()
	logger.Get().Info("Shutting down gracefully, press Ctrl+C again to force")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		logger.Get().LogAttrs(context.Background(), logger.LevelFatal, "Server forced to shutdown", slog.Any("error", err))
		os.Exit(1)
	}

	logger.Get().Info("Server exiting")
}
