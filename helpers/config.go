package helpers


import (
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"os"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"bytes" 	
)


func UploadToCloudinary(fileData []byte, folder string, carID string) (string, error) {
	
	err := godotenv.Load()
	if err != nil {
		return "", err
	}

	cloudName := os.Getenv("CLOUD_NAME")
	apiKey := os.Getenv("API_KEY")
	apiSecret := os.Getenv("API_SECRET")

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return "", err
	}

	ctx := context.TODO()

	//uploadParams := uploader.UploadParams{Folder: "car_folder"}
	uploadParams := uploader.UploadParams{Folder: folder}


	// Create a reader from the byte slice
	fileReader := bytes.NewReader(fileData)

	uploadResult, err := cld.Upload.Upload(ctx, fileReader, uploadParams)
	if err != nil {
		return "", err
	}	
	fmt.Println(uploadResult.SecureURL)

	return uploadResult.SecureURL, nil
	
}