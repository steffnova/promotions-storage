package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/steffnova/promotions-storage/pkg/promotion"

	"github.com/gorilla/mux"
	"pgregory.net/rapid"
)

func TestUnitPromotionGet(t *testing.T) {
	properties := map[string]func(*rapid.T){
		"StatusBadRequest": func(t *rapid.T) {
			request := httptest.NewRequest("GET", "/properties/{id}", nil)
			recorder := httptest.NewRecorder()

			getPromotion := GetPromotion(func(s string) *promotion.Promotion {
				promotion := promotion.PromotionGen.Draw(t, "promotion")
				return &promotion
			})

			handler := PromotionGET(getPromotion, json.Marshal)
			handler.ServeHTTP(recorder, request)

			response := recorder.Result()
			if response.StatusCode != http.StatusBadRequest {
				t.Fatalf("\nInvalid status code: %d\nExpected: status code: %d", response.StatusCode, http.StatusBadRequest)
			}
		},
		"StatusInternalServerError": func(t *rapid.T) {
			id := rapid.Map(rapid.Int(), func(n int) string {
				return fmt.Sprintf("%d", n)
			}).Draw(t, "id")

			getPromotion := GetPromotion(func(s string) *promotion.Promotion {
				promotion := promotion.PromotionGen.Draw(t, "promotion")
				return &promotion
			})

			encode := Encode(func(any) ([]byte, error) {
				return nil, rapid.Map(rapid.Int(), func(n int) error {
					return fmt.Errorf("Error code: %d", n)
				}).Draw(t, "encodeError")
			})

			recorder := httptest.NewRecorder()
			multiplexer := mux.NewRouter()
			multiplexer.HandleFunc("/promotions/{id}", PromotionGET(getPromotion, encode)).Methods(http.MethodGet)

			request := httptest.NewRequest("GET", fmt.Sprintf("/promotions/%s", id), nil)
			multiplexer.ServeHTTP(recorder, request)

			response := recorder.Result()
			if response.StatusCode != http.StatusInternalServerError {
				t.Fatalf("\nInvalid status code: %d\nExpected: status code: %d", response.StatusCode, http.StatusInternalServerError)
			}
		},
		"StatusNotFound": func(t *rapid.T) {
			id := rapid.Map(rapid.Int(), func(n int) string {
				return fmt.Sprintf("%d", n)
			}).Draw(t, "id")

			getPromotion := GetPromotion(func(s string) *promotion.Promotion {
				return nil
			})

			encode := Encode(func(any) ([]byte, error) {
				return nil, nil
			})

			recorder := httptest.NewRecorder()
			multiplexer := mux.NewRouter()
			multiplexer.HandleFunc("/promotions/{id}", PromotionGET(getPromotion, encode)).Methods(http.MethodGet)

			request := httptest.NewRequest("GET", fmt.Sprintf("/promotions/%s", id), nil)
			multiplexer.ServeHTTP(recorder, request)

			response := recorder.Result()
			if response.StatusCode != http.StatusNotFound {
				t.Fatalf("\nInvalid status code: %d\nExpected: status code: %d", response.StatusCode, http.StatusNotFound)
			}
		},
		"StatusOK": func(t *rapid.T) {
			promo := promotion.PromotionGen.Draw(t, "promotion")
			encoded, err := json.Marshal(promo)
			if err != nil {
				t.Fatalf("failed to encode promotion: %s", err)
			}

			getPromotion := GetPromotion(func(id string) *promotion.Promotion {
				if id != promo.ID {
					t.Fatalf("\nGetPromotion: invalid id %s.\nExpected: %s", id, promo.ID)
				}
				return &promo
			})

			encode := Encode(func(input any) ([]byte, error) {
				if input != Promotion(promo) {
					t.Fatalf("\nEncode: invalid input: %#v.\nExpected: %#v", input, Promotion(promo))
				}
				return encoded, nil
			})

			recorder := httptest.NewRecorder()
			multiplexer := mux.NewRouter()
			multiplexer.HandleFunc("/promotions/{id}", PromotionGET(getPromotion, encode)).Methods(http.MethodGet)

			request := httptest.NewRequest("GET", fmt.Sprintf("/promotions/%s", promo.ID), nil)
			multiplexer.ServeHTTP(recorder, request)

			response := recorder.Result()
			if response.StatusCode != http.StatusOK {
				t.Fatalf("\nInvalid status code: %d\nExpected: status code: %d", response.StatusCode, http.StatusOK)
			}

			data, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %s", err)
			}

			if !reflect.DeepEqual(data, encoded) {
				t.Fatalf("invalid reponse data")
			}
		},
	}

	for name, property := range properties {
		t.Run(name, func(t *testing.T) { rapid.Check(t, property) })
	}

}
