package app

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
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
	Get(key string) (string, error)
	GetReverse(ctx context.Context, key, userID string) (string, error)
	GetDBNoCookie(key string) (string, error)
	GetDB(ctx context.Context, key, userID string) (string, error)
	GetDBAll(ctx context.Context, userID string) (map[string]string, error)
}

func (s *URLMapStorage) Get(key string) (string, error) {
	value, ok := s.data[key]
	if !ok {
		ok = false
	}
	return value, nil
}

type SetURL interface {
	Set(key, value string) error
	SetDB(ctx context.Context, key, value, userID string) error
	DeleteDBPrepare(ctx context.Context, key, userID string) error
	DeleteDBFinally(ctx context.Context, key, userID string) error
}

func (s *URLMapStorage) Set(key, value string) error {
	s.data[key] = value
	if Cfg.FilePATH != "" {
		s.SaveToFile()
	}
	return nil
}

type URLMapStorage struct {
	data     map[string]string
	filename string
}

func NewURLMapStorage() URLStorage {
	if Cfg.FilePATH != "" {
		filename = Cfg.FilePATH
	}
	loadDataFromFile(filename)
	return &URLMapStorage{
		data:     data,
		filename: filename,
	}
}

func NewURLDBStorage(connString string) URLStorage {
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		return nil
	}
	//defer pool.Close()

	urlStorage := &URLDBStorage{
		pool:  pool,
		error: err,
	}

	if err := urlStorage.CreateTable(); err != nil {
		return nil
	}
	return urlStorage
}

type URLDBStorage struct {
	pool *pgxpool.Pool
	error
}

func (s *URLDBStorage) GetDB(ctx context.Context, key, userID string) (string, error) {
	var originalURL string
	err := s.pool.QueryRow(ctx, "SELECT original_url FROM url_storage WHERE short_url = $1 AND user_id = $2", key, userID).Scan(&originalURL)
	if err != nil {
		return "", err
	}
	return originalURL, err
}

func (s *URLDBStorage) GetReverse(ctx context.Context, key, userID string) (string, error) {
	var originalURL string
	err := s.pool.QueryRow(ctx, "SELECT short_url FROM url_storage WHERE original_url = $1 AND user_id = $2", key, userID).Scan(&originalURL)
	if err != nil {
		return "", err
	}
	return originalURL, err
}

func (s *URLDBStorage) GetDBNoCookie(key string) (string, error) {
	var originalURL string
	err := s.pool.QueryRow(context.Background(), "SELECT original_url FROM url_storage WHERE short_url = $1 and is_deleted = false", key).Scan(&originalURL)
	if err != nil {
		return "", err
	}
	return originalURL, err
}

func (s *URLDBStorage) GetDBAll(ctx context.Context, userID string) (map[string]string, error) {
	rows, err := s.pool.Query(ctx, "SELECT original_url, short_url FROM url_storage WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	urlMap := make(map[string]string)

	for rows.Next() {
		var shortURL, originalURL string
		err := rows.Scan(&originalURL, &shortURL)
		if err != nil {
			return nil, err
		}
		urlMap[shortURL] = originalURL
	}
	return urlMap, nil
}

func (s *URLDBStorage) DeleteDBFinally(ctx context.Context, key, userID string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		sugar.Info("Error beginning transaction:", err)
		return err
	}
	_, err = tx.Exec(ctx, `
		DELETE FROM url_storage
		WHERE short_url = $1 and user_id = $2 and is_deleted = true
	`, key, userID)

	if err != nil {
		tx.Rollback(ctx)
		sugar.Info("Error rolling back transaction:", err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		sugar.Info("Error committing transaction:", err)
		return err
	}
	return nil
}

func (s *URLDBStorage) DeleteDBPrepare(ctx context.Context, key, userID string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		sugar.Info("Error beginning transaction:", err)
		return err
	}
	_, err = tx.Exec(ctx, `
		UPDATE url_storage
		SET is_deleted = true
		WHERE short_url = $1 and user_id = $2
	`, key, userID)

	if err != nil {
		tx.Rollback(ctx)
		sugar.Info("Error rolling back transaction:", err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		sugar.Info("Error committing transaction:", err)
		return err
	}
	return nil
}

func (s *URLDBStorage) SetDB(ctx context.Context, key, value, userID string) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		sugar.Info("Error beginning transaction:", err)
		return err
	}
	_, err = tx.Exec(ctx, `
		INSERT INTO url_storage (short_url, original_url, user_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (original_url)
		DO UPDATE SET uuid = 1 
	`, key, value, userID)

	if err != nil {
		tx.Rollback(ctx)
		sugar.Info("Error rolling back transaction:", err)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		sugar.Info("Error committing transaction:", err)
		return err
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

func (s *URLDBStorage) SaveToFile() error {
	return nil
}

func (s *URLDBStorage) CreateTable() error {
	_, err := s.pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS url_storage (
            uuid SERIAL PRIMARY KEY,
            short_url VARCHAR UNIQUE NOT NULL,
            original_url VARCHAR UNIQUE NOT NULL,
            user_id VARCHAR,
            is_deleted BOOL DEFAULT false                                
        );
    `)

	s.pool.Exec(context.Background(), `
        INSERT INTO url_storage (short_url, original_url, user_id) 
        VALUES ('first_short_url', 'first_original_url', 'first_user_id');
    `)
	return err
}

func (s *URLMapStorage) GetReverse(ctx context.Context, key, userID string) (string, error) {
	return "", nil
}

func (s *URLMapStorage) GetDB(ctx context.Context, key, userID string) (string, error) {
	return "", nil
}

func (s *URLDBStorage) Get(key string) (string, error) {
	return "", nil
}

func (s *URLMapStorage) SetDB(ctx context.Context, key, value, userID string) error {
	return nil
}

func (s *URLDBStorage) Set(key, value string) error {
	return nil
}

func (s *URLMapStorage) GetDBAll(ctx context.Context, userID string) (map[string]string, error) {
	return nil, nil
}

func (s *URLMapStorage) GetDBNoCookie(key string) (string, error) {
	return "", nil
}

func (s *URLMapStorage) DeleteDBPrepare(ctx context.Context, key, userID string) error {
	return nil
}

func (s *URLMapStorage) DeleteDBFinally(ctx context.Context, key, userID string) error {
	return nil
}
