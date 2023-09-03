package storage

import "github.com/steffnova/promotions-storage/pkg/promotion"

type mock struct {
	implLoadData     func(streamer promotion.Streamer) error
	implGetPromotion func(id string) *promotion.Promotion
}

func (m *mock) LoadPromotions(streamer promotion.Streamer) error {
	return m.implLoadData(streamer)
}

func (m *mock) GetPromotion(id string) *promotion.Promotion {
	return m.implGetPromotion(id)
}
