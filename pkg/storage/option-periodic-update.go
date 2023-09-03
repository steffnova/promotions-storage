package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/steffnova/promotions-storage/pkg/promotion"
)

type optionPeriodicUpdate struct {
	Storage
	period  time.Duration
	context context.Context
}

func (option *optionPeriodicUpdate) LoadPromotions(streamer promotion.Streamer) error {
	if option.period < 0 {
		return fmt.Errorf("%w. Can't perform periodic update", ErrOptionPeriod)
	}

	// Ensure that promotions are loaded when the function is called.
	// If this is not done, first load would be performed when ticker
	// sends data to it's channel, which depending on the value of period
	// could leave storage empty for quite some time (in case period is
	// 30 min, 1Hr, etc...)
	if err := option.Storage.LoadPromotions(streamer); err != nil {
		return err
	}

	go func() {
		ticker := time.NewTicker(option.period)
		for running := true; running; {
			select {
			case <-option.context.Done():
				running = false
			case <-ticker.C:
				// Errors are not handled as they are in a separate go routine.
				// If it is required they could be sent somewhere for processing using channels.
				option.Storage.LoadPromotions(streamer)
			}
		}
	}()

	return nil
}

// OptionPeriodicUpdate returns new [Option] that ensures that LoadPromotions is being called periodically
// once invoked. The ctx parameter is used to signal an end of periodic update. The period parameter
// defines the period between two invocations.
func OptionPeriodicUpdate(ctx context.Context, period time.Duration) Option {
	return func(storage Storage) Storage {
		return &optionPeriodicUpdate{
			Storage: storage,
			context: ctx,
			period:  period,
		}
	}
}
