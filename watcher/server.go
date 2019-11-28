package watcher

import (
	"fmt"
	"net/http"

	log "compiler-file-watcher/logging"
	"compiler-file-watcher/watcher/runner"
)

func handleUploadProtocol(w http.ResponseWriter, r *http.Request) {
	err := runner.NewProtocolRunner(w, r).Run()
	if err != nil {
		errMsg := fmt.Sprintf("handleUploadProtocol err: %v", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		log.Error(errMsg)
	}
}

func handleUploadOssDesc(w http.ResponseWriter, r *http.Request) {
	err := runner.NewOssDescRunner(w, r).Run()
	if err != nil {
		errMsg := fmt.Sprintf("handleUploadOssDesc err: %v", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		log.Error(errMsg)
	}
}

func handleUploadExcel(w http.ResponseWriter, r *http.Request) {
	err := runner.NewExcelRunner(w, r).Run()
	if err != nil {
		errMsg := fmt.Sprintf("handleUploadExcel err: %v", err)
		http.Error(w, errMsg, http.StatusInternalServerError)
		log.Error(errMsg)
	}
}

func NewServer(addr string) *http.Server {
	handler := NewHandler()
	handler.Bind("/upload/protocol", handleUploadProtocol)
	handler.Bind("/upload/oss_desc", handleUploadOssDesc)
	handler.Bind("/upload/excel", handleUploadExcel)

	return &http.Server{
		Addr:    addr,
		Handler: handler,
	}
}
