package zapextra

import (
	"net/http"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggingHandler struct {
	log     *zap.Logger
	handler http.Handler
	fields  []zapcore.Field
}

func (lh loggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s := &responseSizer{w: w}
	t := time.Now()
	lh.handler.ServeHTTP(s, r)

	d := time.Since(t)
	fs := append(
		lh.fields,
		zap.String("host", r.Host),
		zap.String("remote_addr", getHTTPHostname(r.RemoteAddr)),
		zap.String("username", "-"),
		zap.String("method", r.Method),
		zap.String("path", r.RequestURI),
		zap.Int("status", s.code),
		zap.Uint64("size", s.size),
		zap.Duration("duration_human", d),
		zap.Int64("duration_ns", d.Nanoseconds()),
	)
	lh.log.Info(
		"Request",
		fs...,
	)
}

// LoggingHandler enables handling
func LoggingHandler(l *zap.Logger, h http.Handler, fs ...zapcore.Field) http.Handler {
	return loggingHandler{l, h, fs}
}
