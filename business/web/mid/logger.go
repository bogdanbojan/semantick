package mid

import (
	"context"
	"net/http"
	"time"

	"github.com/bogdanbojan/semantick/foundation/web"
	"go.uber.org/zap"
)

// Logger ...
func Logger(log *zap.SugaredLogger) web.Middleware {

	m := func(handler web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			TraceID := "00000000000000000"
			statuscode := http.StatusOK
			now := time.Now()

			log.Infow("request started", "traceid", TraceID, "method", r.Method, "path", r.URL.Path,
				"remoteaddr", r.RemoteAddr)

			err := handler(ctx, w, r)

			log.Infow("request completed", "traceid", TraceID, "method", r.Method, "path", r.URL.Path,
				"remoteaddr", r.RemoteAddr, "statuscode", statuscode, "since", time.Since(now))

			return err
		}
		return h
	}
	return m
}
