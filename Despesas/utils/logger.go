package utils

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	Logger  *log.Logger
	logFile *os.File
)

// InitLogger cria/usa um arquivo de log por dia: log_YYYY-MM-DD.txt
// e tenta sempre gravar na raiz do projeto (onde está o go.mod).
func InitLogger() error {
	// fecha arquivo anterior se InitLogger for chamado mais de uma vez
	if logFile != nil {
		_ = logFile.Close()
		logFile = nil
	}

	root, err := findProjectRoot()
	if err != nil {
		// fallback: diretório atual
		root, _ = os.Getwd()
	}

	date := time.Now().Format("2006-01-02")
	filename := "log_" + date + ".txt"
	fullPath := filepath.Join(root, filename)

	f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	logFile = f

	// logger com data/hora + arquivo:linha (ajuda muito a debugar)
	Logger = log.New(f, "", log.Ldate|log.Ltime|log.Lshortfile)

	Logger.Printf("logger iniciado em %s", fullPath)
	return nil
}

// findProjectRoot sobe as pastas até achar um go.mod (raiz do projeto).
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("go.mod não encontrado acima do diretório atual")
		}
		dir = parent
	}
}
