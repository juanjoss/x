package shutdown

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Graceful(server *http.Server, timeout time.Duration) {
	stop := make(chan os.Signal, 1)
	defer close(stop)

	signal.Notify(stop, shutdownSignals...)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}
