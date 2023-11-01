package app

import (
	"encoding/json"
	"io"
	"os"
	"strconv"
)

const DefaultFilePath = "/tmp/short-url-db.json"

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

func NewURLMapStorage() (URLStorage, error) {
	filename := DefaultFilePath
	if FilePATH != "" {
		filename = FilePATH
	}
	data, err := loadDataFromFile(filename)
	if err != nil {
		return nil, err
	}
	return &URLMapStorage{
		data:     data,
		filename: filename,
	}, nil
}

func (s *URLMapStorage) SaveToFile() error {
	file, err := os.OpenFile(s.filename, os.O_RDWR|os.O_CREATE, 0666)
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

func loadDataFromFile(filename string) (map[string]string, error) {
	data := make(map[string]string)
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return data, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(content, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}
