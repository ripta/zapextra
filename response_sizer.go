package zapextra

import (
	"net/http"
	"sync/atomic"
)

// Ensure responseSizer is always a ResponseWriter at compile time
var _ http.ResponseWriter = &responseSizer{}

type responseSizer struct {
	w    http.ResponseWriter
	code int
	size uint64
}

func (s *responseSizer) Header() http.Header {
	return s.w.Header()
}

func (s *responseSizer) Write(b []byte) (int, error) {
	n, err := s.w.Write(b)
	atomic.AddUint64(&s.size, uint64(n))
	return n, err
}

func (s *responseSizer) WriteHeader(code int) {
	s.w.WriteHeader(code)
	s.code = code
}
