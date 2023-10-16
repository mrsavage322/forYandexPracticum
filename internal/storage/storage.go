package storage

type URLStorage interface {
	Get(key string) (string, bool)
	Set(key, value string)
}

type URLMapStorage struct {
	data map[string]string
}

func NewURLMapStorage() URLStorage {
	return &URLMapStorage{
		data: make(map[string]string),
	}
}

func (s *URLMapStorage) Get(key string) (string, bool) {
	value, ok := s.data[key]
	return value, ok
}

func (s *URLMapStorage) Set(key, value string) {
	s.data[key] = value
}
