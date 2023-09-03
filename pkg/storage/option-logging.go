package storage

import (
	"fmt"
	"io"

	"github.com/steffnova/promotions-storage/pkg/promotion"
)

type optionLogging struct {
	Storage
	writer io.Writer
}

func (option *optionLogging) LoadPromotions(streamer promotion.Streamer) error {
	if err := option.Storage.LoadPromotions(streamer); err != nil {
		fmt.Fprintf(option.writer, "Error while updating storage data: %s\n", err)
		return err
	}

	fmt.Fprintf(option.writer, "Successfully updated storage data\n")
	return nil
}

func (option *optionLogging) GetPromotion(id string) *promotion.Promotion {
	promotion := option.Storage.GetPromotion(id)
	if promotion == nil {
		fmt.Fprintf(option.writer, "No promotion found for ID: %s\n", id)
	} else {
		fmt.Fprintf(option.writer, "Found promotion for ID: %s. Promotion: %+v\n", id, *promotion)
	}
	return promotion
}

// OptionLogging returns new [Option] that adds logging to [Storage] methods.
// The writer parameter specifies writer that is used for logging.
func OptionLogging(writer io.Writer, enable bool) Option {
	return func(storage Storage) Storage {
		if !enable {
			return storage
		}
		return &optionLogging{
			Storage: storage,
			writer:  writer,
		}
	}
}
