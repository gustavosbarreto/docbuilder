package models

import (
	"time"
)

type Namespace struct {
	Name         string             `json:"name"  validate:"required,hostname_rfc1123,excludes=."`
	Owner        string             `json:"owner"`
	TenantID     string             `json:"tenant_id" bson:"tenant_id,omitempty"`
	Members      []interface{}      `json:"members" bson:"members"`
	Settings     *NamespaceSettings `json:"settings"`
	Devices      int                `json:"devices" bson:",omitempty"`
	Sessions     int                `json:"sessions" bson:",omitempty"`
	MaxDevices   int                `json:"max_devices" bson:"max_devices"`
	DevicesCount int                `json:"devices_count" bson:"devices_count,omitempty"`
	Subscription *Subscription      `json:"subscription" bson:"subscription,omitempty"`
}

type Subscription struct {
	ID               string    `json:"id" bson:"id, omitempty"`
	CurrentPeriodEnd time.Time `json:"current_period_end" bson:"current_period_end, omitempty"`
	PriceID          string    `json:"price_id" bson:"price_id, omitempty"`
}

type NamespaceSettings struct {
	SessionRecord bool `json:"session_record" bson:"session_record,omitempty"`
}

type Member struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name,omitempty" bson:"-"`
}
