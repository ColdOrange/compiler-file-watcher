package main

import (
	"compiler-file-watcher/config"
	"compiler-file-watcher/watcher"
)

func main() {
	server := watcher.NewServer(config.WatcherConfig.Addr)
	server.ListenAndServe()
}
