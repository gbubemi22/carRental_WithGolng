package controller

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gbubemi22/go-rentals/database"
	"github.com/gbubemi22/go-rentals/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var carCollection *mongo.Collection = database.OpenCollection(database.Client, "car")

func CreateCar() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var car models.Car

		// Convert the JSON data coming from the client to the Car struct
		if err := c.ShouldBindJSON(&car); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create a channel to receive the result
		resultCh := make(chan error)

		// Perform the database operation concurrently using a goroutine
		go func() {
			defer close(resultCh)

			// Perform any validation or additional checks on the car data

			// Check if the agentId exists in the userCollection
			agentExistsCh := make(chan bool)
			go func() {
				defer close(agentExistsCh)
				count, err := userCollection.CountDocuments(ctx, bson.M{"_id": car.AgentID})
				if err != nil {
					log.Panic(err)
					agentExistsCh <- false
					return
				}
				agentExistsCh <- count > 0
			}()

			agentExists := <-agentExistsCh

			if !agentExists {
				resultCh <- errors.New("agentId does not exist")
				return
			}

			// Create some extra details for the car object - createdAt, updatedAt, ID
			car.CreatedAt = time.Now()
			car.UpdatedAt = time.Now()
			car.ID = primitive.NewObjectID()

			// Insert the new car into the car collection
			_, err := carCollection.InsertOne(ctx, car)
			resultCh <- err
		}()

		// Wait for the result from the goroutine
		err := <-resultCh
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create car"})
			return
		}

		c.JSON(http.StatusCreated, car)
	}
}

// GetAllCar retrieves all cars from the car collection
func GetAllCars() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// Find all cars in the car collection
		cursor, err := carCollection.Find(ctx, bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cars"})
			return
		}

		// Iterate over the cursor and collect the cars
		var cars []models.Car
		if err := cursor.All(ctx, &cars); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cars"})
			return
		}

		c.JSON(http.StatusOK, cars)
	}
}

// GetOneCar retrieves a single car by its ID
func GetOneCar() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// Get the car ID from the request parameters
		carID := c.Param("id")

		// Find the car by its ID in the car collection
		var car models.Car
		err := carCollection.FindOne(ctx, bson.M{"_id": carID}).Decode(&car)
		defer cancel()
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.JSON(http.StatusNotFound, gin.H{"error": "car not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get car"})
			return
		}

		c.JSON(http.StatusOK, car)
	}
}

// UpdateCar updates a car by its ID
func UpdateCar() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// Get the car ID from the request parameters
		carID := c.Param("id")

		var updateData bson.M

		// Convert the JSON data coming from the client to a map
		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Set the updated timestamp
		updateData["updatedAt"] = time.Now()

		// Create a bson.M for the update operation
		updateFields := bson.M{"$set": updateData}

		// Update the car in the car collection
		updateResult, err := carCollection.UpdateOne(ctx, bson.M{"_id": carID}, updateFields)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update car"})
			return
		}

		if updateResult.ModifiedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "car not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "car updated successfully"})
	}
}

// DeleteCar deletes a car by its ID
func DeleteCar() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// Get the car ID from the request parameters
		carID := c.Param("id")

		// Delete the car from the car collection
		deleteResult, err := carCollection.DeleteOne(ctx, bson.M{"_id": carID})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete car"})
			return
		}

		if deleteResult.DeletedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "car not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "car deleted successfully"})
	}
}

func UpdateCarImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// Get the car ID from the request parameters
		carID := c.Param("id")

		var car models.Car

		// Retrieve the existing Car object from the database based on the car ID
		err := carCollection.FindOne(ctx, bson.M{"_id": carID}).Decode(&car)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "car not found"})
			cancel() // Cancel the context if car is not found
			return
		}

		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			cancel() // Cancel the context if there is an error
			return
		}

		images := form.File["images"]

		// Load Cloudinary credentials from .env file
		err = godotenv.Load(".env")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load .env file"})
			cancel() // Cancel the context if there is an error
			return
		}

		cloudinaryCloudName := os.Getenv("CLOUD_NAME")
		cloudinaryAPIKey := os.Getenv("API_KEY")
		cloudinaryAPISecret := os.Getenv("API_SECRET")

		// Set up Cloudinary SDK
		cld, err := cloudinary.NewFromParams(cloudinaryCloudName, cloudinaryAPIKey, cloudinaryAPISecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to initialize Cloudinary"})
			cancel() // Cancel the context if there is an error
			return
		}

		var imageUrls []string

		// Upload the images to Cloudinary
		for _, image := range images {
			// Open the image file
			file, err := image.Open()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open image file"})
				cancel() // Cancel the context if there is an error
				return
			}
			defer file.Close()

			// Upload file to Cloudinary

			uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{Folder: "upload_folder"})
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload image to Cloudinary"})
				cancel() // Cancel the context if there is an error
				return
			}

			// Append the image URL to the imageUrls slice
			imageUrls = append(imageUrls, uploadResult.SecureURL)
		}

		// Update the Car object with the image URLs
		car.Image = imageUrls

		// Set the updated timestamp
		car.UpdatedAt = time.Now()

		// Update the car in the car collection
		_, err = carCollection.UpdateOne(ctx, bson.M{"_id": carID}, bson.M{"$set": car})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update car"})
			cancel() // Cancel the context if there is an error
			return
		}

		cancel() // Cancel the context at the end of the function
		c.JSON(http.StatusOK, gin.H{"message": "car image updated successfully"})
	}
}

func GetAllCarsByAgent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// Get the agent ID from the request parameters
		agentID := c.Param("agentId")

		// Retrieve cars from the database based on the agent ID
		cursor, err := carCollection.Find(ctx, bson.M{"agentId": agentID})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch cars"})
			return
		}
		defer cancel()

		var cars []models.Car
		if err := cursor.All(ctx, &cars); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch cars"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": cars})
	}
}
