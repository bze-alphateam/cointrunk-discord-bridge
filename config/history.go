package config

import (
	"fmt"
	"os"
)

const (
	defaultHistoryFilePath = "./"
	defaultHistoryFileName = "history.json"
)

type History struct {
	FilePath string `yaml:"file_path"`
	FileName string `yaml:"file_name"`
}

func NewHistoryConfig() History {
	filePath := os.Getenv("HISTORY_FILE_PATH")
	if filePath == "" {
		filePath = defaultHistoryFilePath
	}

	return History{
		FilePath: filePath,
		FileName: defaultHistoryFileName,
	}
}

func (h History) Validate() error {
	if len(h.FileName) == 0 {
		return fmt.Errorf("invalid hitory file name")
	}

	if len(h.FilePath) == 0 {
		return fmt.Errorf("invalid hitory file path")
	}

	return nil
}
