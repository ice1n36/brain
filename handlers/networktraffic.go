package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ice1n36/brain/clients"
	"github.com/ice1n36/brain/models"
	"go.uber.org/config"
	"net/http"
)

func WriteNetworkTrafficHandler(w http.ResponseWriter, r *http.Request) {
	configDir := os.Getenv("CONFIG_DIR")
	if len(configDir) == 0 {
		http.Error(w, "must pass in CONFIG_DIR", http.StatusInternalServerError)
		return
	}

	cfg, err := config.NewYAML(config.File(filepath.Join(configDir, "secrets.yaml")))
	if err != nil {
		http.Error(w, "bad config", http.StatusInternalServerError)
		return
	}
	mongoClient, err := clients.NewMongoDatabaseClient(cfg)
	if err != nil {
		fmt.Printf("mongo client failed %v\n", err)
		http.Error(w, "could not create mongo client", http.StatusInternalServerError)
		return
	}
	var reqBody models.MobileAppNetTrafficRequest
	err = json.NewDecoder(r.Body).Decode(&reqBody)
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
