package controller

import (
	"context"
	"fmt"
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

		//convert the JSON data coming from postman to something that golang understands
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//validate the data based on user struct

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return

		}

		//you'll check if the email has already been used by another user

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()

		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})

			return
		}

		//you'll also check if the phone no. has already been used by another user

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the phone number"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exsits"})
			return
		}

		password := Utils.HashPassword(*user.Password)
		user.Password = &password

		//create some extra details for the user object - created_at, updated_at, ID

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()


		// Set default user_type if it is empty
		if user.User_type == nil || *user.User_type == "" {
			defaultUserType := "USER"
			user.User_type = &defaultUserType
		}

		//if all ok, then you insert this new user into the user collection

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User  was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		//return status OK and send the result back

		c.JSON(http.StatusCreated, resultInsertionNumber)

	}
}
// Load environment variables from .env file


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