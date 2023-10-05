package controller

import (
	"context"
	//"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	
     
	//"github.com/cloudinary/cloudinary-go/v2/uploader"
	"github.com/gbubemi22/go-rentals/database"
	"github.com/gbubemi22/go-rentals/models"
	"github.com/gin-gonic/gin"
	//"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/gbubemi22/go-rentals/helpers"
	//"github.com/gbubemi22/go-rentals/models"
	//"github.com/cloudinary/cloudinary-go/v2"
)



var carCollection *mongo.Collection = database.OpenCollection(database.Client, "car")

func CreateCar() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// Convert the JSON data coming from the client to the Car struct
		var car models.Car
		if err := c.ShouldBindJSON(&car); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check if the user exists in the userCollection
		if !userExists(ctx, car.User_ID.Hex()) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user does not exist"})
			return
		}

		// Create additional details for the car object
		car.CreatedAt = time.Now()
		car.UpdatedAt = time.Now()
		car.ID = primitive.NewObjectID()
		car.CarID = car.ID.Hex()

		// Insert the new car into the car collection
		if err := insertCar(ctx, car); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create car"})
			return
		}

		c.JSON(http.StatusCreated, car)
		defer cancel()
	}
}

func userExists(ctx context.Context, userID string) bool {
	count, err := userCollection.CountDocuments(ctx, bson.M{"user_id": userID})
	if err != nil {
		log.Println(err)
		return false
	}
	return count > 0
}

func insertCar(ctx context.Context, car models.Car) error {
	_, err := carCollection.InsertOne(ctx, car)
	return err
}




// GetAllCar retrieves all cars from the car collection
func GetAllCars() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// Find all cars in the car collection
		result, err := carCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cars"})
			return
		}

		// Iterate over the cursor and collect the cars
		var cars []bson.M

		if err := result.All(ctx, &cars); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get cars"})
			return
		}

		if len(cars) == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": " cars  not found"})
			return
		}

		c.JSON(http.StatusOK, cars)
		fmt.Println(cars)
	}
}

// GetOneCar retrieves a single car by its ID
func GetOneCar() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// Get the car ID from the request parameters
		carID := c.Param("car_id")

		// Find the car by its ID in the car collection
		var car models.Car
		err := carCollection.FindOne(ctx, bson.M{"car_id": carID}).Decode(&car)
		defer cancel()
		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get car"})
			return
		}

		

		c.JSON(http.StatusOK, car)
	}
}

// UpdateCar updates a car by its ID
func GetAllCarsByAgent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// Get the user ID from the request parameters
		userID := c.Param("user_id")
		log.Printf("User ID: %s", userID)


		// Retrieve cars from the database based on the user ID
		result, err := carCollection.Find(context.TODO(), bson.M{"user_id": userID})
		if err != nil {
			log.Printf("Error finding cars: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch cars"})
			return
		}
		defer cancel()

		

		// Retrieve cars from the database based on the agent ID
		
		if err != nil {
			log.Printf("Error finding cars: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch cars"})
			return
		}
		defer cancel()

		var cars []models.Car
		if err := result.All(ctx, &cars); err != nil {
			log.Printf("Error decoding cars: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch cars"})
			return
		}

		log.Printf("Retrieved %d cars", len(cars))
		c.JSON(http.StatusOK, cars)

	}
}





// UpdateCar updates a car by its ID
func UpdateCar() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		

		// Get the car ID from the request parameters
		carID := c.Param("car_id")


		var car models.Car
             
		err := carCollection.FindOne(ctx, bson.M{"car_id": carID}).Decode(&car)
if err == mongo.ErrNoDocuments {
    // Document not found, return a 404 response
    c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
    return
} else if err != nil {
    // Handle other errors that might occur during the query
    fmt.Println("Error finding car:", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get car"})
    return
}



		

		var updateObj primitive.D

		if car.CarName != nil {
			updateObj = append(updateObj, bson.E{"carname", car.CarName})
		}

		if car.CarModel != nil {
			updateObj = append(updateObj, bson.E{"carmodel", car.CarModel})
		}

		if car.CarBrand != nil {
			updateObj = append(updateObj, bson.E{"carbrand", car.CarBrand})
		}

		if len(car.CarYear) != 0 {
			updateObj = append(updateObj, bson.E{"carYear", car.CarYear})
		}

		if car.RentalPrice != nil {
			updateObj = append(updateObj, bson.E{"rentalPrice", car.RentalPrice})
		}

		car.UpdatedAt = time.Now()
		updateObj = append(updateObj, bson.E{"updatedAt", car.UpdatedAt})

		upsert := false // Set to false to avoid creating a new car

		filter := bson.M{"car_id": carID}

		fmt.Println("CHECKING", filter)

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := carCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)

		if err != nil {
			msg := fmt.Sprint("car item update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, result)
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


func UpLoadCarImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)

		
		carID := c.Param("car_id")

		fmt.Println("THE CARD ID", carID)
	
	


	var car models.Car
	   
	err := carCollection.FindOne(ctx, bson.M{"car_id": carID}).Decode(&car)
if err == mongo.ErrNoDocuments {
// Document not found, return a 404 response
c.JSON(http.StatusNotFound, gin.H{"error": "Car not found"})
return
} else if err != nil {
// Handle other errors that might occur during the query
fmt.Println("Error finding car:", err)
c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get car"})
return
}


		// Obtain the image file from the request payload
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
			return
		}

		// Read the file content
		imageFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image file"})
			return
		}
		defer imageFile.Close()

		// Create a byte array to hold the file content
		fileContent, err := ioutil.ReadAll(imageFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read image file"})
			return
		}

		// Create a temporary file to hold the image data
tempFile, err := ioutil.TempFile("", "temp-image-*.png") // You can choose a suitable file extension
if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temporary file"})
    return
}
defer tempFile.Close()



_, err = tempFile.Write(fileContent)
if err != nil {
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write to temporary file"})
    return
}

		
		 // Prepare the Cloudinary upload options and perform the upload
		 uploadedImageURL, err := helpers.UploadToCloudinary(fileContent, "image", carID)
		 if err != nil {
			fmt.Println("Error uploading image to Cloudinary:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image"})
			return
		 }


		 // Delete the temporary file after the upload
_ = os.Remove(tempFile.Name())
		 

		var updateObj primitive.D

		if uploadedImageURL != "" {
			updateObj = append(updateObj, bson.E{"image", uploadedImageURL})
		}

		car.UpdatedAt = time.Now()
		updateObj = append(updateObj, bson.E{"updatedAt", car.UpdatedAt})

		upsert := false

		filter := bson.M{"car_id": carID}

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		result, err := carCollection.UpdateOne(
			ctx,
			filter,
			bson.D{
				{"$set", updateObj},
			},
			&opt,
		)

		if err != nil {
			msg := fmt.Sprint("car item update failed")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, result)
	}
}


// func UpLoadCarImage() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var ctx, cancel = context.WithTimeout(context.Background(), time.Second)

// 		var car models.Car

// 		carID := c.Param("car_id")

// 		if err := c.BindJSON(&car); err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 			return
// 		}

// 		var updateObj primitive.D

// 		if car.Image != nil {
// 			updateObj = append(updateObj, bson.E{"image", car.Image})
// 		}

// 		car.UpdatedAt = time.Now()
// 		updateObj = append(updateObj, bson.E{"updatedAt", car.UpdatedAt})

// 		upsert := false

// 		filter := bson.M{"car_id": carID}

// 		opt := options.UpdateOptions{
// 			Upsert: &upsert,
// 		}

// 		result, err := carCollection.UpdateOne(
// 			ctx,
// 			filter,
// 			bson.D{
// 				{"$set", updateObj},
// 			},
// 			&opt,
// 		)

// 		if err != nil {
// 			msg := fmt.Sprint("car item update failed")
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
// 			return
// 		}
// 		defer cancel()

// 		c.JSON(http.StatusOK, result)

// 	}
// }

// func UploadToCloudinary(file []byte, folder string, carID string) (string, error) {
// 	err := godotenv.Load()
// 	if err != nil {
// 		return "", err
// 	}

// 	cloudName := os.Getenv("CLOUD_NAME")
// 	apiKey := os.Getenv("API_KEY")
// 	apiSecret := os.Getenv("API_SECRET")

// 	cld, err := cloudinary.NewFromConfig(config.New(config.CloudinaryConfig{
// 		CloudName: cloudName,
// 		APIKey:    apiKey,
// 		APISecret: apiSecret,
// 	}))
// 	if err != nil {
// 		return "", err
// 	}

// 	ctx := context.TODO()

// 	uploadResult, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
// 		PublicID: fmt.Sprintf("car_%s", carID), // Set the public ID for the uploaded image
// 		Folder:   folder,                        // Set the folder in Cloudinary to store the images
// 		// You can add more options as per your requirements, such as transformations, tags, etc.
// 	})
// 	if err != nil {
// 		return "", err
// 	}

//		return uploadResult.SecureURL, nil
//	}
// func UploadToCloudinary(file []byte, folder string, carID string) (string, error) {
	
// 	err := godotenv.Load()
// 	if err != nil {
// 		return "", err
// 	}

// 	cloudName := os.Getenv("CLOUD_NAME")
// 	apiKey := os.Getenv("API_KEY")
// 	apiSecret := os.Getenv("API_SECRET")

// 	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
// 	if err != nil {
// 		return "", err
// 	}

// 	ctx := context.TODO()

// 	uploadParams := cloudinary.UploadParams{
// 		publicID: fmt.Sprintf("car_%s", carID), // Set the public ID for the uploaded image
// 		options:  []string{"folder:" + folder}, // Set the folder in Cloudinary to store the images
// 	}

// 	uploadResult, err := cld.Upload.Upload(ctx, file, uploadParams)
// 	if err != nil {
// 		return "", err
// 	}	
	

// 	return uploadResult.SecureURL, nil
// }




