package config

import (
	"log"
	"os"
	"strconv"
	"sync"
)

var once sync.Once
var mongodbUri string
var privateKey string
var authMessage string
var maxImageUploadSizeInMb int
var cloudinaryCloudName string
var cloudinaryApiKey string
var cloudinaryApiSecret string

func loadEnv() {
	once.Do(func() {
		mongodbUri = os.Getenv("MONGODB_URI")
		privateKey = os.Getenv("PRIVATE_KEY")
		authMessage = os.Getenv("PROFILE_AUTH_MESSAGE")
		cloudinaryCloudName = os.Getenv("CLOUDINARY_CLOUD_NAME")
		cloudinaryApiKey = os.Getenv("CLOUDINARY_API_KEY")
		cloudinaryApiSecret = os.Getenv("CLOUDINARY_API_SECRET")
		var err error
		maxImageUploadSizeInMb, err = strconv.Atoi(os.Getenv("MAX_IMAGE_UPLOAD_SIZE_IN_MB"))
		if err != nil {
			log.Println("Error converting MAX_IMAGE_UPLOAD_SIZE_IN_MB:", err)
			// Handle the error appropriately
		}
	})
}

func MongoDbUri() string {
	loadEnv()
	return mongodbUri
}

func AuthMessage() string {
	loadEnv()
	return authMessage
}

func PrivateKey() string {
	loadEnv()
	return privateKey
}

func CloudinaryCloudName() string {
	loadEnv()
	return cloudinaryCloudName
}

func CloudinaryApiKey() string {
	loadEnv()
	return cloudinaryApiKey
}

func CloudinaryApiSecret() string {
	loadEnv()
	return cloudinaryApiSecret
}

func MaxImageUploadSizeInMb() int {
	loadEnv()
	return maxImageUploadSizeInMb
}
