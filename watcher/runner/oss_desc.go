package runner

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"compiler-file-watcher/config"
)

type OssDescRunner struct {
	baseRunner
}

func NewOssDescRunner(w http.ResponseWriter, r *http.Request) *OssDescRunner {
	return &OssDescRunner{baseRunner{w, r}}
}

func (p *OssDescRunner) Run() error {
	// save request files
	localFilePaths, err := p.saveRequestFiles()
	if err != nil {
		return fmt.Errorf("saveRequestFile err: %v", err)
	}

	// compile
	buildDir := path.Join(config.WatcherConfig.SourceDir, "/build/libsrc/log")
	err = p.compile(buildDir)
	if err != nil {
		return fmt.Errorf("compile err: %v", err)
	}

	// upload compiled files
	var maybeGeneratedFiles []string
	for _, localFilePath := range localFilePaths {
		baseFilePath := strings.TrimSuffix(localFilePath, ".xml")
		maybeGeneratedFiles = append(maybeGeneratedFiles, []string{
			baseFilePath + ".h",
			baseFilePath + ".cpp",
		}...)
	}
	err = p.uploadCompiledFiles(maybeGeneratedFiles)
	if err != nil {
		return fmt.Errorf("upload compiled file err: %v", err)
	}

	return nil
}
