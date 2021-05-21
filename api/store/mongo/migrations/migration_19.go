package migrations

import (
	"context"

	"github.com/sirupsen/logrus"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var migration_19 = migrate.Migration{
	Version: 19,
	Up: func(db *mongo.Database) error {
		logrus.Info("Applying migration 19 - Up")
		_, err := db.Collection("public_keys").Indexes().DropOne(context.TODO(), "fingerprint")
		return err
	},
	Down: func(db *mongo.Database) error {
		logrus.Info("Applying migration 19 - Down")
		mod := mongo.IndexModel{
			Keys:    bson.D{{"fingerprint", 1}},
			Options: options.Index().SetName("fingerprint").SetUnique(true),
		}
		_, err := db.Collection("public_keys").Indexes().CreateOne(context.TODO(), mod)
		return err
	},
}
