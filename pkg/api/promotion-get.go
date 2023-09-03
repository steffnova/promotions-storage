package api

import (
	"net/http"

	"github.com/steffnova/promotions-storage/pkg/promotion"

	"github.com/gorilla/mux"
)

type GetPromotion func(string) *promotion.Promotion
type Encode func(any) ([]byte, error)

// PromotionGet is HTTP handler that retrives promotion specified by ID
func PromotionGET(getPromotion GetPromotion, encode Encode) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()
		id, ok := mux.Vars(request)["id"]
		if !ok {
			response.WriteHeader(http.StatusBadRequest)
			response.Write([]byte("missing id path parameter"))
			return
		}

		promotion := getPromotion(id)
		if promotion == nil {
			response.WriteHeader(http.StatusNotFound)
			return
		}

		data, err := encode(Promotion(*promotion))
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			return
		}

		response.WriteHeader(http.StatusOK)
		response.Write(data)
	}
}
