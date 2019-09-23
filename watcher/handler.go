package watcher

import (
	"net/http"
	"regexp"
	"time"

	log "compiler-file-watcher/logging"
)

type Handler struct {
	mux map[string]func(http.ResponseWriter, *http.Request)
	reg map[string]*regexp.Regexp
}

func NewHandler() *Handler {
	Handler := &Handler{}
	Handler.mux = make(map[string]func(http.ResponseWriter, *http.Request))
	Handler.reg = make(map[string]*regexp.Regexp)
	return Handler
}

func (Handler *Handler) Bind(pattern string, handler func(w http.ResponseWriter, r *http.Request)) {
	Handler.mux[pattern] = handler
	Handler.reg[pattern] = regexp.MustCompile(pattern)
}

func (Handler *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	zero := time.Now()
	defer func() {
		duration := time.Since(zero)
		log.Debugf("%s %s [%.3fms]", r.Method, r.URL, duration.Seconds()*1000)
	}()

	for pattern, regex := range Handler.reg {
		if regex.MatchString(r.URL.Path) {
			Handler.mux[pattern](w, r)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}
