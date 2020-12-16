package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MobileAppNetTraffic struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	AppID      string             `bson:"app_id,omitempty"`
	AppVersion string             `bson:"app_version,omitempty"`
	SessionID  string             `bson:"session_id,omitempty"`
	Host       string             `bson:"host,omitempty"`
	Port       string             `bson:"port,omitempty"`
}
