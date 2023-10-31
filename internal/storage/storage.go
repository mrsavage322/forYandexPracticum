package storage

import (
	"encoding/json"
	"io"
	_ "io/ioutil"
	"os"
	"strconv"
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
	file, err := os.OpenFile(s.filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	i := 1
	for key, value := range s.data {
		urlData := URLData{
			UUID:        strconv.Itoa(i), // Преобразуем целое число в строку
			ShortURL:    key,
			OriginalURL: value, // Ваш код устанавливает OriginalURL равным ShortURL
		}
		err := encoder.Encode(urlData)
		if err != nil {
			return err
		}
		i++
	}

	return nil
}

func loadDataFromFile(filename string, data *map[string]string) {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return
	}

	err = json.Unmarshal(content, data)
	if err != nil {
		return
	}
}
