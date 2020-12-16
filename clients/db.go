package clients

import (
	"context"
	"fmt"
	"os"

	"github.com/ice1n36/brain/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func NewMongoDatabaseClient() MobileDatabaseClient {
	// TODO: add config support here (so docker deployment can be enabled)
	username := os.Getenv("ATLAS_USER")
	password := os.Getenv("ATLAS_PW")
	database := os.Getenv("ATLAS_DB")
	atlas_uri := fmt.Sprintf(_mongoAtlasURITemplate, username, password, database)
	return &mongodb{
		username:  username,
		password:  password,
		database:  database,
		atlas_uri: atlas_uri,
	}
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
