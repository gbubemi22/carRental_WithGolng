package controller

import (
	"context"
	"fmt"
	database "github.com/gbubemi22/go-rentals/database"
	"github.com/gbubemi22/go-rentals/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)


   func GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		userCollection := database.OpenCollection(database.Client, "user")

		result, err := userCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get users"})
			return
		}


		var allUsers []bson.M

		if err = result.All(ctx, &allUsers); err != nil {
			log.Fatal(err)
		}

		if len(allUsers) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Users not found "})
			return
		}
		c.JSON(http.StatusOK, allUsers[0])
		fmt.Println(allUsers)

	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)


		var user models.User
		
		userCollection := database.OpenCollection(database.Client, "user")
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)

	}
}
