package repository

import (
	authController "ChaiLabs/ChaiLabs/api/v1/auth/controller"
	authService "ChaiLabs/ChaiLabs/api/v1/auth/service"
	userController "ChaiLabs/ChaiLabs/api/v1/user/controller"
	userService "ChaiLabs/ChaiLabs/api/v1/user/service"
	userType "ChaiLabs/ChaiLabs/api/v1/user/types"
	"ChaiLabs/config"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoRepository initializes a new MongoDB client and returns a new MongoRepository instance.
// It takes a mongoURI string as a parameter.
// Returns a pointer to a MongoRepository and an error, if any.

func NewMongoRepository() (*MongoRepository, error) {

	// Create a new context TODO
	ctx := context.TODO()

	// Connect to MongoDB using the provided URI
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoDbUri()))
	if err != nil {
		log.Println("Could not connect to MongoDB:", err)
		return nil, err
	}

	// Ping MongoDB to ensure the connection is valid
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Println("Could not ping MongoDB:", err)
		return nil, err
	}

	// Create a new MongoRepository instance with the connected client and context
	repo := &MongoRepository{
		Client: client,
		Ctx:    ctx,
	}

	log.Println("MongoDB: Connection established")
	return repo, nil
}

// TODO: Initialization of the all Controllers here

// Initialization of userAuth controller
func (r *MongoRepository) GetUserAuthController() authController.AuthController {
	userAuthCollection := r.Client.Database("chaiLabs").Collection("users")
	userAuthService := authService.NewAuthService(&userType.UserRepo{MongoCollection: userAuthCollection, Ctx: r.Ctx})
	userAuthController := authController.NewAuthControllerAuth(userAuthService)

	return userAuthController
}

// GetGamerCollection returns the gamer collection from the database
func (r *MongoRepository) GetUserController() userController.UserController {
	userCollection := r.Client.Database("chaiLabs").Collection("users")
	newUserService := userService.NewUserService(&userType.UserRepo{MongoCollection: userCollection, Ctx: r.Ctx})
	userController := userController.NewUserController(newUserService)

	return userController

}

// Disconnect closes the MongoDB connection
func (r *MongoRepository) Disconnect() {
	if err := r.Client.Disconnect(r.Ctx); err != nil {
		log.Println("Could not disconnect from MongoDB:", err)
	}
}
