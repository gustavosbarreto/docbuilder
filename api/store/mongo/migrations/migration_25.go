package migrations

import (
	"context"

	"github.com/sirupsen/logrus"
	migrate "github.com/xakep666/mongo-migrate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var migration_25 = migrate.Migration{
	Version:     25,
	Description: "remove devices with no namespaces related",
	Up: func(db *mongo.Database) error {
		logrus.Info("Applying migration 25 - Up")
		query := []bson.M{
			{
				"$lookup": bson.M{
					"from":         "namespaces",
					"localField":   "tenant_id",
					"foreignField": "tenant_id",
					"as":           "namespace",
				},
			},
			{
				"$addFields": bson.M{
					"namespace": bson.M{"$anyElementTrue": []interface{}{"$namespace"}},
				},
			},

			{
				"$match": bson.M{
					"namespace": bson.M{"$eq": true},
				},
			},

			{
				"$unset": "namespace",
			},

			{
				"$out": "devices",
			},
		}

		_, err := db.Collection("devices").Aggregate(context.TODO(), query)
		return err
	},
	Down: func(db *mongo.Database) error {
		logrus.Info("Applying migration 25 - Down")
		return nil
	},
}
