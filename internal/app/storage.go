package app

import (
	"encoding/json"
	"io"
	"os"
	"strconv"
)

const DefaultFilePath = "/tmp/short-url-db.json"

var filename, data = DefaultFilePath, make(map[string]string)

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
	if FilePATH != "" {
		s.SaveToFile()
	}
}

type URLMapStorage struct {
	data     map[string]string
	filename string
}

func NewURLMapStorage() URLStorage {
	if FilePATH != "" {
		filename = FilePATH
	}
	loadDataFromFile(filename)
	return &URLMapStorage{
		data:     data,
		filename: filename,
	}
}

func (s *URLMapStorage) SaveToFile() error {
	file, err := os.OpenFile(s.filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	i := 1
	for key, value := range s.data {
		urlData := URLData{
			UUID:        strconv.Itoa(i),
			ShortURL:    key,
			OriginalURL: value,
		}
		err := encoder.Encode(urlData)
		if err != nil {
			return err
		}
		i++
	}

	return nil
}

func loadDataFromFile(filename string) map[string]string {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		return nil
	}
	return nil
}
