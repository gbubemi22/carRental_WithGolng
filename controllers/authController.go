package controller

import (
	"context"
	//"fmt"
	database "github.com/gbubemi22/go-rentals/database"
	"github.com/gbubemi22/go-rentals/models"
	Utils "github.com/gbubemi22/go-rentals/utils"
	"github.com/dgrijalva/jwt-go"
	//"github.com/joho/godotenv"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
	"os"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

var validate = validator.New()

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User

		// convert the JSON data coming from Postman to something that Golang understands
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// validate the data based on the user struct
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Create a channel to receive the result of the asynchronous tasks
		resultCh := make(chan interface{})

		// Perform the email and phone number existence checks concurrently
		go func() {
			defer close(resultCh)

			// Check if the email has already been used by another user
			emailCountCh := make(chan int64)
			go func() {
				defer close(emailCountCh)
				count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
				if err != nil {
					log.Panic(err)
					emailCountCh <- -1 // Indicate an error
					return
				}
				emailCountCh <- count
			}()

			// Check if the phone number has already been used by another user
			phoneCountCh := make(chan int64)
			go func() {
				defer close(phoneCountCh)
				count, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
				if err != nil {
					log.Panic(err)
					phoneCountCh <- -1 // Indicate an error
					return
				}
				phoneCountCh <- count
			}()

			// Wait for both emailCountCh and phoneCountCh to receive values
			emailCount, emailCountOk := <-emailCountCh
			phoneCount, phoneCountOk := <-phoneCountCh

			// Check for any errors
			if emailCount < 0 || phoneCount < 0 || !emailCountOk || !phoneCountOk {
				resultCh <- "Error occurred while checking for email or phone number"
				return
			}

			if emailCount > 0 || phoneCount > 0 {
				resultCh <- "This email or phone number already exists"
				return
			}

			// All checks passed, continue with user creation
			password := Utils.HashPassword(*user.Password)
			user.Password = &password

			// Create extra details for the user object
			user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			user.ID = primitive.NewObjectID()
			user.User_id = user.ID.Hex()

			// Set default user_type if it is empty
			if user.User_type == nil || *user.User_type == "" {
				defaultUserType := "USER"
				user.User_type = &defaultUserType
			}

			// Insert the new user into the user collection
			resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
			if insertErr != nil {
				resultCh <- "User was not created"
				return
			}

			resultCh <- resultInsertionNumber
		}()

		// Wait for the result from the channel
		result := <-resultCh

		// Check the type of the result and respond accordingly
		switch result := result.(type) {
		case string:
			c.JSON(http.StatusInternalServerError, gin.H{"error": result})
		default:
			c.JSON(http.StatusCreated, result)
		}

		defer cancel()
	}
}

func CreateAgent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User

		// convert the JSON data coming from Postman to something that Golang understands
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// validate the data based on the user struct
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		// Create a channel to receive the result of the asynchronous tasks
		resultCh := make(chan interface{})

		// Perform the email and phone number existence checks concurrently
		go func() {
			defer close(resultCh)

			// Check if the email has already been used by another user
			emailCountCh := make(chan int64)
			go func() {
				defer close(emailCountCh)
				count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
				if err != nil {
					log.Panic(err)
					emailCountCh <- -1 // Indicate an error
					return
				}
				emailCountCh <- count
			}()

			// Check if the phone number has already been used by another user
			phoneCountCh := make(chan int64)
			go func() {
				defer close(phoneCountCh)
				count, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
				if err != nil {
					log.Panic(err)
					phoneCountCh <- -1 // Indicate an error
					return
				}
				phoneCountCh <- count
			}()

			// Wait for both emailCountCh and phoneCountCh to receive values
			emailCount, emailCountOk := <-emailCountCh
			phoneCount, phoneCountOk := <-phoneCountCh

			// Check for any errors
			if emailCount < 0 || phoneCount < 0 || !emailCountOk || !phoneCountOk {
				resultCh <- "Error occurred while checking for email or phone number"
				return
			}

			if emailCount > 0 || phoneCount > 0 {
				resultCh <- "This email or phone number already exists"
				return
			}

			// All checks passed, continue with user creation
			password := Utils.HashPassword(*user.Password)
			user.Password = &password

			// Create extra details for the user object
			user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			user.ID = primitive.NewObjectID()
			user.User_id = user.ID.Hex()

			// Set default user_type if it is empty
			if user.User_type == nil || *user.User_type == "" {
				defaultUserType := "AGENT"
				user.User_type = &defaultUserType
			}

			// Insert the new user into the user collection
			resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
			if insertErr != nil {
				resultCh <- "User was not created"
				return
			}

			resultCh <- resultInsertionNumber
		}()

		// Wait for the result from the channel
		result := <-resultCh

		// Check the type of the result and respond accordingly
		switch result := result.(type) {
		case string:
			c.JSON(http.StatusInternalServerError, gin.H{"error": result})
		default:
			c.JSON(http.StatusCreated, result)
		}

		defer cancel()
	}
}


// Read JWT secret key from environment variable
var jwtSecret = []byte(os.Getenv("JWT_SECRET"))


func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password is incorrect"})
			return
		}

		passwordIsValid, msg := Utils.VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()

		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}

		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		}
		

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}


		// Create a new JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"userId": foundUser.User_id,
			"exp":    time.Now().Add(time.Hour * 24).Unix(), // Set token expiration time
		})

		// Sign the token with the secret key
		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
			return
		}


		// Return the token in the response
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
}