package clients

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/ice1n36/brain/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/config"
)

const (
	_mongoAtlasURITemplate = "mongodb+srv://%s:%s@cluster0.yspmb.mongodb.net/%s?retryWrites=true&w=majority"
)

type MobileDatabaseClient interface {
	WriteNetTraffic(ctx context.Context, m *models.MobileAppNetTraffic) error
}

type mongodb struct {
	username  string
	password  string
	database  string
	atlas_uri string
}

type mongoconfig struct {
	AtlasUsername string `yaml:"atlas_user"`
	AtlasPassword string `yaml:"atlas_pw"`
	AtlasDatabase string `yaml:"atlas_db"`
}

func NewMongoDatabaseClient(cfg config.Provider) (MobileDatabaseClient, error) {
	var mongocfg mongoconfig
	if err := cfg.Get("mongo").Populate(&mongocfg); err != nil {
		// fallback to using environment variables
		mongocfg.AtlasUsername = os.Getenv("ATLAS_USER")
		mongocfg.AtlasPassword = os.Getenv("ATLAS_PW")
		mongocfg.AtlasDatabase = os.Getenv("ATLAS_DB")
	}

	if len(mongocfg.AtlasUsername) == 0 {
		return nil, errors.New("Cannot retrieve atlas username from config or environment variable ATLAS_USER")
	}
	if len(mongocfg.AtlasPassword) == 0 {
		return nil, errors.New("Cannot retrieve atlas password from config or environment variable ATLAS_PW")
	}
	if len(mongocfg.AtlasDatabase) == 0 {
		return nil, errors.New("Cannot retrieve atlas database from config or environment variable ATLAS_DB")
	}

	atlas_uri := fmt.Sprintf(_mongoAtlasURITemplate, mongocfg.AtlasUsername, mongocfg.AtlasPassword, mongocfg.AtlasDatabase)
	return &mongodb{
		username:  mongocfg.AtlasUsername,
		password:  mongocfg.AtlasPassword,
		database:  mongocfg.AtlasDatabase,
		atlas_uri: atlas_uri,
	}, nil
}

func (m mongodb) WriteNetTraffic(ctx context.Context, mant *models.MobileAppNetTraffic) error {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.atlas_uri))
	if err != nil {
		return err
	}

	defer client.Disconnect(ctx)

	database := client.Database("brain_mobile_db")
	mobileNetTrafficCollection := database.Collection("mobile_net_traffic")

	// insert to collection
	insertResult, err := mobileNetTrafficCollection.InsertOne(ctx, mant)
	if err != nil {
		return err
	}
	fmt.Println(insertResult.InsertedID)
	return nil
}
