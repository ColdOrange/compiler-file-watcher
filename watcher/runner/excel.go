package runner

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"time"

	"compiler-file-watcher/config"
)

type ExcelRunner struct {
	baseRunner
}

func NewExcelRunner(w http.ResponseWriter, r *http.Request) *ExcelRunner {
	return &ExcelRunner{baseRunner{w, r}}
}

func (p *ExcelRunner) Run() error {
	// save request files
	_, err := p.saveRequestFiles()
	if err != nil {
		return fmt.Errorf("saveRequestFile err: %v", err)
	}

	// compile
	timestamp := time.Now()
	err = p.compile(config.WatcherConfig.MakeCmdPath, nil, config.WatcherConfig.ExcelConfig.BuildDir)
	if err != nil {
		return fmt.Errorf("compile err: %v", err)
	}

	// upload new generated files
	var newGeneratedFiles []string
	allExcelConfigFiles, err := ioutil.ReadDir(config.WatcherConfig.ExcelConfig.TargetDir)
	if err != nil {
		return fmt.Errorf("ioutil.ReadDir err: %v", err)
	}
	for _, file := range allExcelConfigFiles {
		if !file.IsDir() && file.ModTime().After(timestamp) {
			newGeneratedFiles = append(newGeneratedFiles,
				path.Join(config.WatcherConfig.ExcelConfig.TargetDir, file.Name()))
		}
	}

	err = p.uploadCompiledFiles(newGeneratedFiles)
	if err != nil {
		return fmt.Errorf("upload compiled file err: %v", err)
	}

	return nil
}
