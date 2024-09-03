package app

import (
	"sync"

	"github.com/e1m0re/passman/internal/model"
)

type Store struct {
	itemsList map[string]*model.DatumInfo

	sync.RWMutex
}

func (s *Store) GetItemById(id string) *model.DatumInfo {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	v, ok := s.itemsList[id]
	if !ok {
		return nil
	}

	return v
}

func (s *Store) AddItem(datum *model.DatumInfo) {
	s.itemsList[datum.File] = datum
}

func (s *Store) UpdateList(list []*model.DatumInfo) {
	for _, item := range list {
		s.itemsList[item.File] = item
	}
}

func (s *Store) GetItems() []*model.DatumInfo {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	idx := 0
	result := make([]*model.DatumInfo, len(s.itemsList))
	for _, item := range s.itemsList {
		result[idx] = item
		idx++
	}

	return result
}

func NewStore() *Store {
	return &Store{
		itemsList: make(map[string]*model.DatumInfo),
	}
}
