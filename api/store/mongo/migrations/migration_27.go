package migrations

import (
	"context"

	"github.com/shellhub-io/shellhub/pkg/models"
	"github.com/sirupsen/logrus"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var migration27 = migrate.Migration{
	Version:     27,
	Description: "Create a new field on namespaces to store the API tokens",
	Up: func(db *mongo.Database) error {
		logrus.Info("Applying migration 27 - Up")
		_, err := db.Collection("namespaces").UpdateMany(context.TODO(), bson.M{}, bson.M{"$set": bson.M{"api_tokens": []models.Token{}}})

		return err
	},
	Down: func(db *mongo.Database) error {
		logrus.Info("Applying migration 27 - Down")

		return nil
	},
}
