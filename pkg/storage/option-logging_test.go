package storage

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/steffnova/promotions-storage/pkg/promotion"

	"pgregory.net/rapid"
)

func TestUnitOptionLogging_LoadPromotions(t *testing.T) {
	properties := map[string]func(t *rapid.T){
		"PassThrough": func(t *rapid.T) {
			err := rapid.Custom(func(t *rapid.T) error {
				isError := rapid.Bool().Draw(t, "isError")
				if isError {
					return rapid.Map(rapid.Int(), func(n int) error {
						return fmt.Errorf("Error code: %d", n)
					}).Draw(t, "errorCode")
				}
				return nil
			}).Draw(t, "error")

			enableLogging := rapid.Bool().Draw(t, "enableLog")

			mockStorage := mock{
				implLoadData: func(streamer promotion.Streamer) error {
					return err
				},
			}

			storage := WithOptions(&mockStorage, OptionLogging(bytes.NewBufferString(""), enableLogging))

			// Logging option should act as a passthrough (it shouldn't affect output)
			err1 := mockStorage.LoadPromotions(nil)
			err2 := storage.LoadPromotions(nil)
			if !errors.Is(err1, err2) {
				t.Fatalf("Expecting same output.\nErr1: %v\nErr2: %v\n", err1, err2)
			}
		},
	}

	for name, property := range properties {
		t.Run(name, func(t *testing.T) { rapid.Check(t, property) })
	}
}

func TestUnitOptionLogging_GetPromotion(t *testing.T) {
	properties := map[string]func(t *rapid.T){
		"PassThrough": func(t *rapid.T) {
			result := rapid.Custom(func(t *rapid.T) *promotion.Promotion {
				found := rapid.Bool().Draw(t, "found")
				if !found {
					return nil
				}
				result := promotion.PromotionGen.Draw(t, "promotion")
				return &result
			}).Draw(t, "result")

			mockStorage := mock{
				implGetPromotion: func(id string) *promotion.Promotion {
					return result
				},
			}

			enableLogging := rapid.Bool().Draw(t, "enableLog")

			storage := WithOptions(&mockStorage, OptionLogging(bytes.NewBufferString(""), enableLogging))

			// Logging option should act as a passthrough (it shouldn't affect output)
			prom1 := mockStorage.GetPromotion("")
			prom2 := storage.GetPromotion("")
			if !reflect.DeepEqual(prom1, prom2) {
				t.Fatalf("Expecting same output.\nPromotion1: %v\nPromotion2: %v\n", prom1, prom2)
			}
		},
	}

	for name, property := range properties {
		t.Run(name, func(t *testing.T) { rapid.Check(t, property) })
	}
}
