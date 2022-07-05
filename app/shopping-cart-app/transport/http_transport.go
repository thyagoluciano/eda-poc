package transport

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"

	"br.com.thyagoluciano.poc/app/shopping-cart-app/endpoint"
	"br.com.thyagoluciano.poc/app/shopping-cart-app/endpoint/api"
)

func NewHttpServer(ctx context.Context, endpoints endpoint.Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware)

	r.Methods("POST").Path("/cart/checked-out").Handler(httptransport.NewServer(
		endpoints.PurchaseOrderCheckedOut,
		decodeCustomerOrderRequest,
		encodeResponse,
	))

	return r
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func decodeCustomerOrderRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req api.PurchaseOrderCheckedOutRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	return req, nil
}
