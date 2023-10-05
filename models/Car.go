package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)


type Car struct {
	ID          primitive.ObjectID  `bson:"_id" json:"_id"`
	User_ID *primitive.ObjectID  `bson:"user_id" json:"user_id"`

	CarName     *string             `bson:"carname" json:"carname" validate:"required"`
	CarModel    *string             `bson:"carmodel" json:"carmodel" validate:"required"`
	CarBrand    *string             `bson:"carbrand" json:"carbrand" validate:"required"`
	CarYear     string              `bson:"carYear" json:"carYear"`
	RentalPrice *int                `bson:"rentalPrice" json:"rentalPrice"`
	IsAvailable bool                `bson:"isAvailable" json:"isAvailable"`
	RentedBy    *primitive.ObjectID `bson:"rentedBy" json:"rentedBy"`
	Image       []string            `bson:"image" json:"image"`
	CreatedAt   time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time           `bson:"updatedAt" json:"updatedAt"`
	CarID       string              `bson:"car_id" json:"car_id"`
}


// type Car struct {
// 	ID          primitive.ObjectID  `bson:"_id"`
// 	User_ID     primitive.ObjectID  `json:"user_id" binding:"required"`
// 	CarName     *string             `json:"carname" validate: "carname required`
// 	CarModel    *string             `json:"carmodel" validate: " required"`
// 	CarBrand    *string             `json:"carbrand" validate: ", required"`
// 	CarYear     string              `bson:"carYear" json:"carYear"`
// 	RentalPrice *int                `json:"rentalPrice"`
// 	IsAvailable bool                `json:"isAvailable"`
// 	RentedBy    *primitive.ObjectID `json:"rentedBy,"`
// 	Image       []string            `json:"image"`
// 	CreatedAt   time.Time           `json:"createdAt"`
// 	UpdatedAt   time.Time           `json:"updatedAt"`
// 	CarID       string              `json:"car_id"`
// }
