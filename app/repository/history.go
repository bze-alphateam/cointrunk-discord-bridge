package repository

import (
	"encoding/json"
	"fmt"
	"github.com/bze-alphateam/cointrunk-discord-bridge/app/entity"
	"os"
)

type History struct {
	FilePath string
	FileName string

	lastId int64
}

func NewHistory(fileName, filePath string) (*History, error) {
	if fileName == "" || filePath == "" {
		return nil, fmt.Errorf("invalid dependencies provided to history repository constructor")
	}

	return &History{
		FilePath: filePath,
		FileName: fileName,
	}, nil
}

func (h *History) GetLastInsertedID() int64 {
	if h.lastId > 0 {
		return h.lastId
	}

	stored := h.readHistoryFile()

	return stored.LastID
}

func (h *History) SaveLastInsertedID(id int64) error {
	h.lastId = id
	return h.writeHistoryInFile(id)
}

func (h *History) writeHistoryInFile(id int64) error {
	// Create a History struct with the last ID
	hist := entity.History{LastID: id}

	// Marshal struct to JSON
	jsonData, err := json.Marshal(hist)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Write JSON data to the file
	jsonFile, err := os.Create(h.getFileFullName()) // os.Create creates or truncates the file
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer jsonFile.Close()

	_, err = jsonFile.Write(jsonData)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (h *History) readHistoryFile() entity.History {
	jsonFile, err := os.Open(h.getFileFullName())
	if err != nil {
		return entity.History{LastID: 0}
	}
	defer jsonFile.Close()

	// Decode JSON data into struct
	var hist entity.History
	decoder := json.NewDecoder(jsonFile)
	err = decoder.Decode(&hist)
	if err != nil {
		fmt.Println(err)
		return entity.History{LastID: 0}
	}

	return hist
}

func (h *History) getFileFullName() string {
	return fmt.Sprintf("%s/%s", h.FilePath, h.FileName)
}
