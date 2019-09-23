package runner

import (
	"fmt"
	"net/http"
	"path"
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
	// save request file
	localFilePath, err := p.saveRequestFile()
	if err != nil {
		return fmt.Errorf("saveRequestFile err: %v", err)
	}

	// compile
	// TODO: config build dir
	buildDir := path.Join(config.WatcherConfig.SourceDir, "/build/libsrc/protocol")
	err = p.compile(buildDir)
	if err != nil {
		return fmt.Errorf("compile err: %v", err)
	}

	// upload compiled files
	baseFilePath := strings.TrimSuffix(localFilePath, ".proto")
	maybeGeneratedFiles := []string{
		baseFilePath + ".pb.h",
		baseFilePath + ".pb.cc",
		baseFilePath + ".pbp.h",
		baseFilePath + ".pbp.cc",
	}
	err = p.uploadCompiledFiles(maybeGeneratedFiles)
	if err != nil {
		return fmt.Errorf("upload compiled file err: %v", err)
	}

	return nil
}
