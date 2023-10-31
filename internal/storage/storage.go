package storage

import (
	"encoding/json"
	"os"
)

type URLStorage interface {
	SetURL
	GetURL
	SaveToFile() error
}

type URLData struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type GetURL interface {
	Get(key string) (string, bool)
}

func (s *URLMapStorage) Get(key string) (string, bool) {
	value, ok := s.data[key]
	return value, ok
}

type SetURL interface {
	Set(key, value string)
}

func (s *URLMapStorage) Set(key, value string) {
	s.data[key] = value
	s.SaveToFile()
}

type URLMapStorage struct {
	data     map[string]string
	filename string
}

func NewURLMapStorage() URLStorage {
	filename := "/tmp/short-url-db.json"
	data := make(map[string]string)
	loadDataFromFile(filename, &data)
	return &URLMapStorage{
		data:     data,
		filename: filename,
	}
}

func (s *URLMapStorage) SaveToFile() error {
	data, err := json.Marshal(s.data)
	if err != nil {
		return err
	}
	err = os.WriteFile(s.filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func loadDataFromFile(filename string, data *map[string]string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return
	}
	err = json.Unmarshal(content, data)
	if err != nil {
		return
	}
}
