package runner

import (
	"fmt"
	"net/http"
	"strings"

	"compiler-file-watcher/config"
)

type ProtocolRunner struct {
	baseRunner
}

func NewProtocolRunner(w http.ResponseWriter, r *http.Request) *ProtocolRunner {
	return &ProtocolRunner{baseRunner{w, r}}
}

func (p *ProtocolRunner) Run() error {
	// save request files
	localFilePaths, err := p.saveRequestFiles()
	if err != nil {
		return fmt.Errorf("saveRequestFile err: %v", err)
	}

	// compile
	err = p.compile(config.WatcherConfig.ProtocolConfig.BuildDir)
	if err != nil {
		return fmt.Errorf("compile err: %v", err)
	}

	// upload compiled files
	var maybeGeneratedFiles []string
	for _, localFilePath := range localFilePaths {
		baseFilePath := strings.Replace(strings.TrimSuffix(localFilePath, ".proto"),
			config.WatcherConfig.ProtocolConfig.SourceDir,
			config.WatcherConfig.ProtocolConfig.TargetDir, -1)
		maybeGeneratedFiles = append(maybeGeneratedFiles, []string{
			baseFilePath + ".pb.h",
			baseFilePath + ".pb.cc",
			baseFilePath + ".pbp.h",
			baseFilePath + ".pbp.cc",
		}...)
	}
	err = p.uploadCompiledFiles(maybeGeneratedFiles)
	if err != nil {
		return fmt.Errorf("upload compiled file err: %v", err)
	}

	return nil
}
