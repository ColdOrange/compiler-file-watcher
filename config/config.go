package config

import (
	"fmt"
	"path"
	"strings"

	"github.com/spf13/viper"
)

var WatcherConfig struct {
	Addr        string // listen addr
	SourceDir   string // local project dir
	RemoteDir   string // remote project dir
	MakeCmdPath string // make cmd path
}

var LoggingConfig struct {
	Type    string // `text` or `json`
	Level   string // enable log level
	Console bool   // log to console
	File    string // or log to file
	Rotate  int    // rotate hours
}

func init() {
	initWatcherConfig()
	initLoggingConfig()
}

// Map file path from remote project dir to local project dir.
func GetLocalFilePath(remoteFilePath string) (string, error) {
	remoteDir := WatcherConfig.RemoteDir
	relativePath := strings.TrimPrefix(remoteFilePath, remoteDir)
	if relativePath == remoteFilePath {
		return "", fmt.Errorf("remoteFilePath: <%s> is not relative to remoteDir: <%s>", remoteFilePath, remoteDir)
	}

	localDir := WatcherConfig.SourceDir
	return path.Join(localDir, relativePath), nil
}

func initWatcherConfig() {
	watcherConfig := viper.New()
	watcherConfig.SetConfigName("watcher")
	watcherConfig.AddConfigPath("./config/")
	_ = watcherConfig.ReadInConfig()

	// addr
	addr := watcherConfig.GetString("addr")
	if addr == "" {
		panic("watcher.yml config err: `addr` not set")
	}

	// source dir
	sourceDir := watcherConfig.GetString("source_dir")
	if sourceDir == "" {
		panic("watcher.yml config err: `source_dir` not set")
	}

	// remote dir
	remoteDir := watcherConfig.GetString("remote_dir")
	if remoteDir == "" {
		panic("watcher.yml config err: `remote_dir` not set")
	}

	// make cmd path
	makeCmdPath := watcherConfig.GetString("make_cmd_path")
	if makeCmdPath == "" {
		panic("watcher.yml config err: `make_cmd_path` not set")
	}

	WatcherConfig = struct {
		Addr        string
		SourceDir   string
		RemoteDir   string
		MakeCmdPath string
	}{
		Addr:        addr,
		SourceDir:   sourceDir,
		RemoteDir:   remoteDir,
		MakeCmdPath: makeCmdPath,
	}
}

func initLoggingConfig() {
	loggingConfig := viper.New()
	loggingConfig.SetConfigName("logging")
	loggingConfig.AddConfigPath("./config/")
	_ = loggingConfig.ReadInConfig()

	// type
	logType := loggingConfig.GetString("type")
	if logType == "" {
		logType = "text"
	}

	// level
	level := loggingConfig.GetString("level")
	if level == "" {
		level = "info"
	} else {
		level = strings.ToLower(level)
	}

	// console
	console := loggingConfig.GetBool("console")

	// file
	file := loggingConfig.GetString("file")
	if file == "" {
		file = "compiler_file_watcher.log"
	}

	// rotate
	rotate := loggingConfig.GetInt("rotate")
	if rotate == 0 {
		rotate = 24
	}

	LoggingConfig = struct {
		Type    string
		Level   string
		Console bool
		File    string
		Rotate  int
	}{
		Type:    logType,
		Level:   level,
		Console: console,
		File:    file,
		Rotate:  rotate,
	}
}
