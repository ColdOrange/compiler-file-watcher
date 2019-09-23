package watcher

import (
	"net/http"

	log "compiler-file-watcher/logging"
	"compiler-file-watcher/watcher/runner"
)

func handleUploadProtocol(w http.ResponseWriter, r *http.Request) {
	err := runner.NewProtocolRunner(w, r).Run()
	if err != nil {
		log.Errorf("handleUploadProtocol err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func handleUploadOssDesc(w http.ResponseWriter, r *http.Request) {
	err := runner.NewOssDescRunner(w, r).Run()
	if err != nil {
		log.Errorf("handleUploadOssDesc err: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func NewServer(addr string) *http.Server {
	handler := NewHandler()
	handler.Bind("/upload/protocol", handleUploadProtocol)
	handler.Bind("/upload/oss_desc", handleUploadOssDesc)

	return &http.Server{
		Addr:    addr,
		Handler: handler,
	}
}
