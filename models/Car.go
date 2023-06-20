package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Car struct {
	ID          primitive.ObjectID  `bson:"_id"`
	AgentID     primitive.ObjectID  `bson:"agentId" json:"agentId"`
	CarName     *string             `bson:"carname" json:"carname"`
	CarModel    *string             `bson:"carmodel" json:"carmodel"`
	CarBrand    *string             `bson:"carbrand" json:"carbrand"`
	CarYear     *int                `bson:"carYear" json:"carYear"`
	RentalPrice *int                `bson:"rentalPrice" json:"rentalPrice"`
	IsAvailable bool                `bson:"isAvailable" json:"isAvailable"`
	RentedBy    *primitive.ObjectID `bson:"rentedBy,omitempty" json:"rentedBy,omitempty"`
	Image       []string           `bson:"image,omitempty" json:"image,omitempty"`
	CreatedAt   time.Time           `bson:"createdAt,omitempty" json:"createdAt,omitempty"`
	UpdatedAt   time.Time           `bson:"updatedAt,omitempty" json:"updatedAt,omitempty"`
	CarID      string             `json:"car_id"`
}
