package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ice1n36/brain/clients"
	"github.com/ice1n36/brain/models"
	"net/http"
)

func WriteNetworkTrafficHandler(w http.ResponseWriter, r *http.Request) {
	mongoClient := clients.NewMongoDatabaseClient()
	var reqBody models.MobileAppNetTrafficRequest
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = mongoClient.WriteNetTraffic(context.TODO(),
		&models.MobileAppNetTraffic{
			AppID:      reqBody.AppID,
			AppVersion: reqBody.AppVersion,
			SessionID:  reqBody.SessionID,
			Host:       reqBody.Host,
			Port:       reqBody.Port,
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		fmt.Fprintf(w, "OK")
	}
}
