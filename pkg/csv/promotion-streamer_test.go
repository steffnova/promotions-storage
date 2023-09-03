package csv

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/steffnova/promotions-storage/pkg/promotion"

	"pgregory.net/rapid"
)

func TestUnitPromotionStreamer(t *testing.T) {
	properties := map[string]func(*rapid.T){
		"ErrStreamer": func(t *rapid.T) {
			reader := Reader(func() (io.Reader, error) {
				err := rapid.Map(rapid.Int(), func(n int) error {
					return fmt.Errorf("Error code: %d", n)
				}).Draw(t, "err")
				return nil, err
			})

			_, errs := PromotionStreamer(reader)()

			for err := range errs {
				switch {
				case err == nil:
					t.Fatalf("Expected an error: %s", ErrStreamer)
				case errors.Is(err, ErrStreamer):
					return
				default:
					t.Fatalf("Got error: %v, want error: %v", err, ErrStreamer)

				}
			}
		},
		"ErrStreaming": func(t *rapid.T) {
			reader := Reader(func() (io.Reader, error) {
				err := rapid.Map(rapid.Int(), func(n int) error {
					return fmt.Errorf("Error code: %d", n)
				}).Draw(t, "err")
				mock := mockReader{
					implRead: func(p []byte) (int, error) {
						return 0, err
					},
				}
				return &mock, nil
			})

			_, errs := PromotionStreamer(reader)()

			for err := range errs {
				switch {
				case err == nil:
					t.Fatalf("Expected an error: %s", ErrStreaming)
				case errors.Is(err, ErrStreaming):
					return
				default:
					t.Fatalf("Got error: %v, want error: %v", err, ErrStreaming)

				}
			}
		},
		"ProcessedCSVPromotions": func(t *rapid.T) {
			promotions1 := promotion.PromotionsGen.Draw(t, "promotions")
			buffer := bytes.NewBufferString("")
			writer := csv.NewWriter(buffer)
			writer.Write([]string{})
			for _, promotion := range promotions1 {
				writer.Write([]string{promotion.ID, promotion.Price, promotion.ExparationDate})
			}
			writer.Flush()

			reader := Reader(func() (io.Reader, error) {
				return bytes.NewBufferString(buffer.String()), nil
			})

			promotions, errs := PromotionStreamer(reader)()
			promotions2 := []promotion.Promotion{}
			for promotion := range promotions {
				t.Log("prmotions")
				promotions2 = append(promotions2, promotion)
			}

			for err := range errs {
				if err != nil {
					t.Fatalf("Unexpected error: %s", err)
				}
			}

			if !reflect.DeepEqual(promotions1, promotions2) {
				t.Fatalf("Processed promotions should match original ones")
			}
		},
	}

	for name, property := range properties {
		t.Run(name, func(t *testing.T) { rapid.Check(t, property) })
	}
}
