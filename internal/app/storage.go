package app

import (
	"context"
	"encoding/json"
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

type URLDatabaseStorage interface {
	Set(key, value string) error
	Get(key string) (string, error)
}

type URLDatabase struct {
	conn *pgx.Conn
}

func NewURLDatabase() (URLDatabaseStorage, error) {
	conn, err := pgx.Connect(context.Background(), DatabaseAddr)
	if err != nil {
		return nil, err
	}

	db := &URLDatabase{conn: conn}
	err = db.CreateTableIfNotExists()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (db *URLDatabase) Set(key, value string) error {
	tx, err := db.conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), "INSERT INTO url_storage (short_url, original_url) VALUES ($1, $2)", key, value)
	if err != nil {
		return err
	}
	return tx.Commit(context.Background())
}

func (db *URLDatabase) Get(key string) (string, error) {
	var originalURL string

	tx, err := db.conn.Begin(context.Background())
	if err != nil {
		return "", err
	}
	defer tx.Rollback(context.Background())

	err = tx.QueryRow(context.Background(), "SELECT original_url FROM url_storage WHERE short_url = $1", key).Scan(&originalURL)
	if err != nil {
		return "", err
	}
	if err := tx.Commit(context.Background()); err != nil {
		return "", err
	}
	return originalURL, nil
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

func NewDBMapStorage() URLDatabaseStorage {
	if DatabaseAddr != "" {
		db, err := NewURLDatabase()
		if err != nil {
			panic(err)
		}
		return db
	}
	return nil
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

func (db *URLDatabase) CreateTableIfNotExists() error {
	_, err := db.conn.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS url_storage (
            uuid SERIAL PRIMARY KEY,
            short_url VARCHAR NOT NULL,
            original_url TEXT NOT NULL
        );
    `)
	return err
}
