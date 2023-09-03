package storage

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sync"
	"testing"
	"time"

	"github.com/steffnova/promotions-storage/pkg/promotion"

	"pgregory.net/rapid"
)

func TestUnitAutoUpdate_LoadPromotions(t *testing.T) {
	properties := map[string]func(t *rapid.T){
		"ErrorPeriod": func(t *rapid.T) {
			mockStorage := mock{
				implLoadData: func(streamer promotion.Streamer) error {
					return nil
				},
			}

			period := time.Millisecond * time.Duration(rapid.IntRange(math.MinInt/int(time.Millisecond), -1).Draw(t, "period"))

			storage := WithOptions(&mockStorage, OptionPeriodicUpdate(context.Background(), period))
			if err := storage.LoadPromotions(nil); !errors.Is(err, ErrOptionPeriod) {
				t.Fatalf("Got error: %s, want error: %s", err, ErrOptionPeriod)
			}
		},
		"ErrorLoadPromotions": func(t *rapid.T) {
			expectedError := rapid.Map(rapid.Int(), func(n int) error {
				return fmt.Errorf("Error code: %d", n)
			}).Draw(t, "error")
			mockStorage := mock{
				implLoadData: func(streamer promotion.Streamer) error {
					return expectedError
				},
			}

			period := time.Millisecond * time.Duration(time.Duration(rapid.IntRange(1, 5).Draw(t, "period")))

			storage := WithOptions(&mockStorage, OptionPeriodicUpdate(context.Background(), period))
			if err := storage.LoadPromotions(nil); !errors.Is(err, expectedError) {
				t.Fatalf("Got error: %s, want error: %s", err, expectedError)
			}
		},
		"StopAfterNUpdates": func(t *rapid.T) {
			ctx, cancel := context.WithCancel(context.Background())
			counter := rapid.IntRange(0, 10).Draw(t, "counter")
			wg := sync.WaitGroup{}
			wg.Add(counter)
			mockStorage := mock{
				implLoadData: func() func(streamer promotion.Streamer) error {
					return func(streamer promotion.Streamer) error {
						if counter == 0 {
							cancel()
						} else {
							wg.Done()
							counter--
						}
						return nil
					}
				}(),
			}

			period := time.Millisecond * time.Duration(rapid.IntRange(1, 5).Draw(t, "period"))
			storage := WithOptions(&mockStorage, OptionPeriodicUpdate(ctx, period))

			storage.LoadPromotions(func() (promotion.Stream, <-chan error) {
				stream := make(chan promotion.Promotion)
				errs := make(chan error)
				defer close(stream)
				defer close(errs)
				return stream, errs
			})

			wg.Wait()
		},
	}

	for name, property := range properties {
		t.Run(name, func(t *testing.T) { rapid.Check(t, property) })
	}
}
