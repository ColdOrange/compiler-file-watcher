package runner

import (
	"fmt"
	"net/http"
	"strings"

	"compiler-file-watcher/config"
)

type TcaplusTableRunner struct {
	baseRunner
}

func NewTcaplusTableRunner(w http.ResponseWriter, r *http.Request) *TcaplusTableRunner {
	return &TcaplusTableRunner{baseRunner{w, r}}
}

func (p *TcaplusTableRunner) Run() error {
	// save request files
	localFilePaths, err := p.saveRequestFiles()
	if err != nil {
		return fmt.Errorf("saveRequestFile err: %v", err)
	}

	// compile
	err = p.compile(config.WatcherConfig.MakeCmdPath, []string{"gen_tcaplus"}, config.WatcherConfig.TcaplusTableConfig.BuildDir)
	if err != nil {
		return fmt.Errorf("compile err: %v", err)
	}

	// upload compiled files
	var maybeGeneratedFiles []string
	for _, localFilePath := range localFilePaths {
		baseFilePath := strings.Replace(strings.TrimSuffix(localFilePath, ".xml"),
			config.WatcherConfig.TcaplusTableConfig.SourceDir,
			config.WatcherConfig.TcaplusTableConfig.TargetDir, -1)
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
