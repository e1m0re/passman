package store

import (
	"context"

	"github.com/e1m0re/passman/internal/model"
	"github.com/e1m0re/passman/internal/repository"
)

//go:generate go run github.com/vektra/mockery/v2@v2.44.2 --name=StoreService
type StoreService interface {
	// AddItem creates new datum item.
	AddItem(ctx context.Context, datumInfo model.DatumInfo) (*model.DatumItem, error)
}

type storeService struct {
	datumRepository repository.DatumRepository
}

// AddItem creates new datum item.
func (s storeService) AddItem(ctx context.Context, datumInfo model.DatumInfo) (*model.DatumItem, error) {
	return s.datumRepository.AddItem(ctx, datumInfo)
}

var _ StoreService = (*storeService)(nil)

// NewStoreService initiates new instance of StoreService.
func NewStoreService(datumRepository repository.DatumRepository) StoreService {
	return &storeService{
		datumRepository: datumRepository,
	}
}
