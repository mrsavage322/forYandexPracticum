package storage

type URLStorage interface {
	SetURL
	GetURL
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
}

type URLMapStorage struct {
	data map[string]string
}

func NewURLMapStorage() URLStorage {
	return &URLMapStorage{
		data: make(map[string]string),
	}
}
