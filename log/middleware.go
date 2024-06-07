package log

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

const (
	RequestBodyMaxSize  = 64 * 1024 // 64KB
	ResponseBodyMaxSize = 64 * 1024 // 64KB
)

func New(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			path := r.URL.Path
			query := r.URL.RawQuery

			br := newBodyReader(r.Body, RequestBodyMaxSize)
			r.Body = br

			// dump response body
			bw := newBodyWriter(w, ResponseBodyMaxSize)

			defer func() {
				status := bw.Status()
				method := r.Method
				host := r.Host
				end := time.Now()
				latency := end.Sub(start)
				//userAgent := r.UserAgent()
				ip := r.RemoteAddr
				referer := r.Referer()

				baseAttributes := []slog.Attr{}

				requestAttributes := []slog.Attr{
					slog.Time("time", start),
					slog.String("method", method),
					slog.String("host", host),
					slog.String("path", path),
					slog.String("query", query),
					slog.String("ip", ip),
					slog.String("referer", referer),
				}

				responseAttributes := []slog.Attr{
					slog.Time("time", end),
					slog.Duration("latency", latency),
					slog.Int("status", status),
				}

				//requestAttributes = append(requestAttributes, slog.String("user-agent", userAgent))

				attributes := append(
					[]slog.Attr{
						{
							Key:   "request",
							Value: slog.GroupValue(requestAttributes...),
						},
						{
							Key:   "response",
							Value: slog.GroupValue(responseAttributes...),
						},
					},
					baseAttributes...,
				)

				level := slog.LevelInfo
				if status >= http.StatusInternalServerError {
					level = slog.LevelError
				} else if status >= http.StatusBadRequest && status < http.StatusInternalServerError {
					level = slog.LevelWarn
				}

				logger.LogAttrs(r.Context(), level, strconv.Itoa(status)+": "+http.StatusText(status), attributes...)
			}()

			next.ServeHTTP(bw, r)
		})
	}
}
