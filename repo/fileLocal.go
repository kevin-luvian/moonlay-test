package repo

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/kennygrant/sanitize"
)

type FileLocalRepo struct {
	BasePath string
}

func NewFileRepo(basePath string) *FileLocalRepo {
	return &FileLocalRepo{
		BasePath: basePath,
	}
}

func (r *FileLocalRepo) GetFileContentType(fileID string) string {
	fileExt := filepath.Ext(fileID)
	switch fileExt {
	case ".txt":
		return "text/plain"
	case ".pdf":
		return "application/pdf"
	default:
		return "text/plain"
	}
}

func (r *FileLocalRepo) MatchExt(filename string, extensions []string) bool {
	filename = strings.TrimSpace(filename)
	filename = strings.ToLower(filename)
	fileExt := filepath.Ext(filename)

	for _, ext := range extensions {
		if ext == fileExt {
			return true
		}
	}

	return false
}

func (r *FileLocalRepo) EnsureBaseDir(filepath string) error {
	baseDir := path.Dir(filepath)
	info, err := os.Stat(baseDir)
	if err == nil && info.IsDir() {
		return nil
	}

	return os.MkdirAll(baseDir, 0755)
}

func (r *FileLocalRepo) SaveFile(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	fmt.Println("File Opened")

	defer src.Close()

	filePath := r.BasePath + "/" + file.Filename
	filePath = sanitize.Path(filePath)

	err = r.EnsureBaseDir(filePath)
	if err != nil {
		return "", err
	}

	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return filePath, nil
}

func (r *FileLocalRepo) DeleteFile(fileID string) error {
	if !strings.HasPrefix(fileID, r.BasePath) {
		return errors.New("invalid base path")
	}

	return os.Remove(fileID)
}

func (r *FileLocalRepo) StreamFile(fileID string) (io.Reader, error) {
	f, err := os.Open(fileID)
	return f, err
}
