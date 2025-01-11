package aggregator

import (
	"encoding/json"
	"log"
	"net/http"
	"reflect"
	"storya-gateway-backend/internal/client"
	"storya-gateway-backend/internal/config"
	"storya-gateway-backend/internal/pb/github.com/webbsalad/storya-content-backend/content"
	"storya-gateway-backend/internal/pb/github.com/webbsalad/storya-passport-backend/passport"
)

// only test complex request handlers
func MixedHandler(cfg config.Config, grpcClients *client.GRPCClients) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		getRandResp, err := grpcClients.ContentClient.GetRand(r.Context(), &content.GetRandRequest{ContentType: 0, Count: 1})
		if err != nil {
			log.Printf("get rand: %v", err)
			return
		}

		items := getRandResp.Items

		for _, item := range items {
			getItemResp, err := grpcClients.ContentClient.Get(r.Context(), &content.GetItemRequest{ItemId: item.Id})
			if err != nil {
				log.Printf("get item: %v", err)
				return
			}

			if !reflect.DeepEqual(item, getItemResp) {
				log.Printf("wrong item")
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode("success"); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}

	}
}

// only test complex mixed services request handlers
func MixedClientsHandler(cfg config.Config, grpcClients *client.GRPCClients) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var checkRequest struct {
			Token string `json:"token"`
		}

		err := json.NewDecoder(r.Body).Decode(&checkRequest)
		if err != nil {
			return
		}

		checkResp, err := grpcClients.PassportClient.CheckToken(r.Context(), &passport.CheckTokenRequest{Token: checkRequest.Token})
		if err != nil {
			log.Printf("check token: %v", err)
			return
		}

		getUserResp, err := grpcClients.PassportClient.GetUser(r.Context(), &passport.GetUserRequest{UserId: checkResp.UserId})
		if err != nil {
			log.Printf("get user: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(getUserResp); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
