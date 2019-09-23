package runner

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"compiler-file-watcher/config"
	log "compiler-file-watcher/logging"
)

type baseRunner struct {
	w http.ResponseWriter
	r *http.Request
}

// Retrieve uploaded file from http request and save to local project directory.
// Return local file path on success, else return an error.
func (b *baseRunner) saveRequestFile() (string, error) {
	// get file path
	filePath := b.r.FormValue("filepath")
	if filePath == "" {
		return "", errors.New("filepath is empty")
	}

	// get mapped local file path
	localFilePath, err := config.GetLocalFilePath(filePath)
	if err != nil {
		return "", err
	}

	// get file content
	formFile, _, err := b.r.FormFile("file")
	if err != nil {
		return "", err
	}
	defer formFile.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, formFile)
	if err != nil {
		return "", err
	}

	// open local file
	file, err := os.OpenFile(localFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// save to local file
	_, err = file.Write(buf.Bytes())
	if err != nil {
		return "", err
	}

	return localFilePath, nil
}

// Run `make` command to compile and generate files to response.
func (b *baseRunner) compile(buildDir string) error {
	cmd := &exec.Cmd{
		Path: config.WatcherConfig.MakeCmdPath,
		Dir:  buildDir,
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(out))
	}

	return nil
}

// Upload compiled files in representation of multipart form data.
// Each file is composed by a <file_path, file_content> multipart pair.
func (b *baseRunner) uploadCompiledFiles(fileList []string) error {
	// multipart writer
	mw := multipart.NewWriter(b.w)
	b.w.Header().Set("Content-Type", mw.FormDataContentType())

	// upload each file as multipart form data
	for i, filePath := range fileList {
		if _, err := os.Stat(filePath); err == nil {
			err := uploadMultipartFile(i, filePath, mw)
			if err != nil {
				return err
			}
			log.Infof("upload compiled file: <%s> success", filePath)
		}
	}

	// must close multipart writer!
	if err := mw.Close(); err != nil {
		return err
	}

	return nil
}

func uploadMultipartFile(index int, filePath string, mw *multipart.Writer) error {
	// file path
	fw, err := mw.CreateFormField(fmt.Sprintf("filepath%d", index))
	if err != nil {
		return err
	}
	_, err = fw.Write([]byte(filePath))
	if err != nil {
		return err
	}

	// file content
	fw, err = mw.CreateFormFile(fmt.Sprintf("file%d", index), filepath.Base(filePath))
	if err != nil {
		return err
	}

	// read file
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0777)
	if err != nil {
		return err
	}
	defer file.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		return err
	}

	// write
	_, err = fw.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}
