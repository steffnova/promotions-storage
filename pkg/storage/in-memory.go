package storage

import (
	"fmt"
	"sync"

	"github.com/steffnova/promotions-storage/pkg/promotion"
)

type inMemory struct {
	newStoreLock sync.RWMutex
	promotions   map[string]promotion.Promotion
}

func (s *inMemory) LoadPromotions(streamer promotion.Streamer) error {
	promotions := map[string]promotion.Promotion{}

	stream, errs := streamer()
	for promotion := range stream {
		promotions[promotion.ID] = promotion
	}

	for err := range errs {
		if err != nil {
			return fmt.Errorf("%w: %w", ErrLoadPromotions, err)
		}
	}

	s.newStoreLock.Lock()
	defer s.newStoreLock.Unlock()

	s.promotions = promotions
	return nil
}

func (s *inMemory) GetPromotion(id string) *promotion.Promotion {
	s.newStoreLock.RLock()
	defer s.newStoreLock.RUnlock()
	promotion, ok := s.promotions[id]
	if !ok {
		return nil
	}

	return &promotion
}

// InMemory returns Storage that stores promotions data in memory
func InMemory() Storage {
	return &inMemory{}
}
