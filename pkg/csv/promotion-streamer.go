package csv

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"

	"github.com/steffnova/promotions-storage/pkg/promotion"
)

// PromotionStreamer returns [promotion.Streamer] that streams
// read CSV data (using reader returned by getReader parameter) into
// it's stream. Error channal returned by streamer can be used to track
// errors. Stream is closed if an error occurs during read.
func PromotionStreamer(reader Reader) promotion.Streamer {
	return func() (promotion.Stream, <-chan error) {
		out := make(chan promotion.Promotion)
		errs := make(chan error, 1)

		reader, err := reader()
		if err != nil {
			errs <- fmt.Errorf("%w: %w", ErrStreamer, err)
			close(errs)
			close(out)
			return out, errs
		}

		go func() {
			defer close(out)
			defer close(errs)

			csvReader := csv.NewReader(reader)

			counter := 0
			for {
				records, err := csvReader.Read()
				if errors.Is(err, io.EOF) {
					break
				}
				if err != nil {
					errs <- fmt.Errorf("%w: %d. %w. Streaming stops", ErrStreaming, counter, err)
					break
				}

				out <- promotion.Promotion{
					ID:             records[0],
					Price:          records[1],
					ExparationDate: records[2],
				}
				counter++
			}
		}()

		return out, errs
	}
}
