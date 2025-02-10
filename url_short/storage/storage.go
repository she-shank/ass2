package storage

import "github.com/its-kos/assignment1/types"

type Storage interface {
	CreateURL(url types.URL) error
	GetURL(id string) (*types.URL, error)
	UpdateURL(url *types.URL) error
	DeleteURL(id string) error
	GetAllURLs() ([]*types.URL, error)
	DeleteAllURLs() error
}