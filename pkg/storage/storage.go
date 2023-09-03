package storage

import "github.com/steffnova/promotions-storage/pkg/promotion"

// Storage is interface for storing Promotions and retrieving it.
type Storage interface {
	// LoadPromotions loads new promotions using a streamer. Once done old promotions are discarded
	LoadPromotions(streamer promotion.Streamer) error
	// GetStore retrieves store that can be used to access promotions stored within it
	GetPromotion(id string) *promotion.Promotion
}
