package main

import (
	"compiler-file-watcher/watcher"
	"github.com/spf13/viper"
)

var addr string

func init() {
	watcherConfig := viper.New()
	watcherConfig.SetConfigName("watcher")
	watcherConfig.AddConfigPath("./config/")
	_ = watcherConfig.ReadInConfig()

	addr = watcherConfig.GetString("addr")
}

func main() {
	server := watcher.NewServer(addr)
	server.ListenAndServe()
}
