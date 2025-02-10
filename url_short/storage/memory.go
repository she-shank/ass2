package storage

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"

	"github.com/its-kos/assignment1/types"
)

const (
	dbFilePath = "db.json"
)

type MemoryStorage struct {
	dbLock *sync.Mutex
}

func NewMemoryStorage(dbLock *sync.Mutex) *MemoryStorage {
	return &MemoryStorage{
		dbLock: dbLock,
	}
}

func (m *MemoryStorage) CreateURL(newUrl types.URL) error {
	m.dbLock.Lock()
	defer m.dbLock.Unlock()

	file, err := os.OpenFile(dbFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	var data []*types.URL
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil && err != io.EOF {
		return err
	}

	data = append(data, &newUrl)

	file.Seek(0, 0)
	file.Truncate(0)
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(&data); err != nil {
		return err
	}
	return nil
}

func (m *MemoryStorage) GetURL(id string) (*types.URL, error) {

	m.dbLock.Lock()
	defer m.dbLock.Unlock()

	file, err := os.OpenFile(dbFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []*types.URL
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil && err != io.EOF {
		return nil, err
	}

	for _, url := range data {
		if url.ID == id {
			url.Hits++
			file.Seek(0, 0)
			file.Truncate(0)
			encoder := json.NewEncoder(file)
			if err := encoder.Encode(&data); err != nil {
				return nil, err
			}
			return url, nil
		}
	}
	return nil, errors.New("URL with id " + id + " not found")
}

func (m *MemoryStorage) UpdateURL(newUrl *types.URL) error {
	m.dbLock.Lock()
	defer m.dbLock.Unlock()

	file, err := os.OpenFile(dbFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	var data []*types.URL
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil && err != io.EOF {
		return err
	}

	for _, url := range data {
		if url.ID == newUrl.ID {
			url.ID = newUrl.ID
			url.URL = newUrl.URL
			url.CreatedAt = newUrl.CreatedAt
			url.Hits = newUrl.Hits
			url.TimeToLive = newUrl.TimeToLive
			break
		}
	}

	file.Seek(0, 0)
	file.Truncate(0)
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(&data); err != nil {
		return err
	}
	return nil
}

func (m *MemoryStorage) DeleteURL(id string) error {
	m.dbLock.Lock()
	defer m.dbLock.Unlock()

	file, err := os.OpenFile(dbFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	var data []*types.URL
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil && err != io.EOF {
		return err
	}

	var updatedData []*types.URL
	var deletedURL *types.URL

	for _, url := range data {
		if url.ID == id {
			deletedURL = url
		} else {
			updatedData = append(updatedData, url)
		}
	}

	file.Seek(0, 0)
	file.Truncate(0)
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(&updatedData); err != nil {
		return err
	}

	if deletedURL == nil {
		return errors.New("URL with id " + id + " not found")
	}
	return nil
}

func (m *MemoryStorage) GetAllURLs() ([]*types.URL, error) {
	m.dbLock.Lock()
	defer m.dbLock.Unlock()

	file, err := os.OpenFile(dbFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data []*types.URL
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil && err != io.EOF {
		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	} else {
		return data, nil
	}
}

func (m *MemoryStorage) DeleteAllURLs() error {
	m.dbLock.Lock()
	defer m.dbLock.Unlock()

	file, err := os.OpenFile(dbFilePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	file.Seek(0, 0)
	file.Truncate(0)
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(nil); err != nil {
		return err
	}
	return nil
}
