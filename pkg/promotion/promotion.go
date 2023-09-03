package promotion

import (
	"fmt"
	"time"

	"pgregory.net/rapid"
)

// Promotion is a structure that represents promotion data
type Promotion struct {
	ID             string // ID promotion ID
	Price          string // Price is price is expressed as string but can be parsed to float64
	ExparationDate string // Exparation data is a string and can be parsed to time.Time
}

var PromotionGen = rapid.Custom(func(t *rapid.T) Promotion {
	return Promotion{
		ID: rapid.Map(rapid.Int(), func(n int) string {
			return fmt.Sprintf("%d", n)
		}).Draw(t, "id"),
		Price:          fmt.Sprintf("%v", rapid.Float64().Draw(t, "price")),
		ExparationDate: time.Now().String(),
	}
})

var PromotionsGen = rapid.Custom(func(t *rapid.T) []Promotion {
	n := rapid.IntRange(1, 100).Draw(t, "n")
	promotions := make([]Promotion, n)
	for index := range promotions {
		promotions[index] = PromotionGen.Draw(t, "promotion")
	}
	return promotions
})
