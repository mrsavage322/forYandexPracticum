package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5"
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

func NewURLDBStorage(connString string) URLStorage {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return nil
	}

	urlStorage := &URLDBStorage{
		conn: conn,
	}

	if err := urlStorage.CreateTable(); err != nil {
		return nil
	}
	return urlStorage
}

type URLDBStorage struct {
	conn *pgx.Conn
}

func (s *URLDBStorage) Get(key string) (string, bool) {
	var originalURL string
	err := s.conn.QueryRow(context.Background(), "SELECT original_url FROM url_storage WHERE short_url = $1", key).Scan(&originalURL)
	if err != nil {
		return "", false
	}
	return originalURL, true
}

func (s *URLDBStorage) Set(key, value string) {
	_, err := s.conn.Exec(context.Background(), "INSERT INTO url_storage (short_url, original_url) VALUES ($1, $2)", key, value)
	if err != nil {
		fmt.Println("Error inserting into database:", err)
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

func (s *URLDBStorage) SaveToFile() error {
	return nil
}

func (s *URLDBStorage) CloseDB() {
	s.conn.Close(context.Background())
}

func (s *URLDBStorage) CreateTable() error {
	_, err := s.conn.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS url_storage (
            uuid SERIAL PRIMARY KEY,
            short_url VARCHAR UNIQUE NOT NULL,
            original_url VARCHAR NOT NULL
        );
    `)
	return err
}
