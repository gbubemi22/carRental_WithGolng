package controller

import (
	"context"
	"fmt"
	database "github.com/gbubemi22/go-rentals/database"
	"github.com/gbubemi22/go-rentals/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"strconv"
	"time"
)

//var userCollection *mongo.Collection

func init() {
	userCollection = database.OpenCollection(database.Client, "user")
}

func GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))
		if err != nil || startIndex < 0 {
			startIndex = 0
		}

		matchStage := bson.D{{"$match", bson.D{{}}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{
					{"$cond", bson.A{
						bson.D{{"$isArray", "$data"}},
						bson.D{{"$slice", bson.A{"$data", startIndex, recordPerPage}}},
						bson.A{},
					}},
				}},
			}},
		}
		fmt.Println("startIndex:", startIndex)
		fmt.Println("recordPerPage:", recordPerPage)

		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, projectStage,
		})
		defer cancel()

		if err != nil {
			log.Println("Error occurred while listing user items:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list user items"})
			return
		}

		var allUsers []bson.M
		if err = result.All(ctx, &allUsers); err != nil {
			log.Fatal(err)
		}

		if len(allUsers) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "No cars found"})
			return
		}
		c.JSON(http.StatusOK, allUsers[0])

	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)

	}
}
