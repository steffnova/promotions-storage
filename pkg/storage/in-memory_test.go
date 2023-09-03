package storage

import (
	"errors"
	"fmt"
	"testing"

	"github.com/steffnova/promotions-storage/pkg/promotion"

	"pgregory.net/rapid"
)

func TestUnitInMemory_GetPromotion(t *testing.T) {
	properties := map[string]func(t *rapid.T){
		"NoPromotionFoundInEmptyStorage": func(t *rapid.T) {
			inMemory := InMemory()

			id := rapid.String().Draw(t, "id")
			if promotion := inMemory.GetPromotion(id); promotion != nil {
				t.Fatalf("Promotion must be nill when storage has no data")
			}
		},
		"PromotionFound": func(t *rapid.T) {
			inMemory := InMemory()

			promotions := promotion.PromotionsGen.Draw(t, "promotions")

			streamer := promotion.Streamer(func() (promotion.Stream, <-chan error) {
				errs := make(chan error)
				stream := make(chan promotion.Promotion, len(promotions))

				defer close(errs)
				defer close(stream)

				for _, promotion := range promotions {
					stream <- promotion
				}

				return promotion.Stream(stream), errs
			})

			if err := inMemory.LoadPromotions(streamer); err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}

			for _, promotion := range promotions {
				if p := inMemory.GetPromotion(promotion.ID); p == nil {
					t.Fatal("Expected to found all promotions")
				}
			}
		},
	}

	for name, property := range properties {
		t.Run(name, func(t *testing.T) { rapid.Check(t, property) })
	}
}

func TestUnitInMemory_NewData(t *testing.T) {
	properties := map[string]func(t *rapid.T){
		"NoError": func(t *rapid.T) {
			inMemory := InMemory()
			promotions := promotion.PromotionsGen.Draw(t, "promotions")

			streamer := promotion.Streamer(func() (promotion.Stream, <-chan error) {
				errs := make(chan error)
				stream := make(chan promotion.Promotion, len(promotions))

				go func() {
					defer close(errs)
					defer close(stream)

					for _, promotion := range promotions {
						stream <- promotion
					}
				}()

				return promotion.Stream(stream), errs
			})

			if err := inMemory.LoadPromotions(streamer); err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}
		},
		"StreamerError": func(t *rapid.T) {
			inMemory := InMemory()
			expecteError := rapid.Custom(func(t *rapid.T) error {
				return fmt.Errorf("Error code: %d", rapid.Int().Draw(t, "error-code"))
			}).Draw(t, "error")

			streamer := promotion.Streamer(func() (promotion.Stream, <-chan error) {
				errs := make(chan error, 1)
				stream := make(chan promotion.Promotion)

				go func() {
					defer close(errs)
					defer close(stream)

					errs <- expecteError
				}()

				return promotion.Stream(stream), errs
			})

			if err := inMemory.LoadPromotions(streamer); !errors.Is(err, expecteError) {
				t.Fatalf("Got error: %s, wanted err: %s", err, expecteError)
			}
		},
	}

	for name, property := range properties {
		t.Run(name, func(t *testing.T) { rapid.Check(t, property) })
	}
}
