package userService

import (
	communityModel "ChaiLabs/ChaiLabs/api/v1/community/models"
	userServiceInterface "ChaiLabs/ChaiLabs/api/v1/user/interface"
	"ChaiLabs/ChaiLabs/api/v1/user/models"
	userType "ChaiLabs/ChaiLabs/api/v1/user/types"
	cloudinaryService "ChaiLabs/config/cloudinary"
	"ChaiLabs/utils"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"reflect"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserServiceImpl struct {
	UserRepo *userType.UserRepo
}

// SetUserWalletAddress implements userServiceInterface.UserService.
func (userService *UserServiceImpl) SetUserWalletAddress(userToken *userType.Token, loggedInUser *models.User) (*models.User, error) {
	//user signature validation
	isVerified := utils.VerifySignature(userToken.Signature, userToken.Address, loggedInUser.Token.Message)

	if !isVerified {
		return nil, errors.New("unauthorized access, Invalid signature")
	}

	filter := bson.M{"username": loggedInUser.Username}
	update := bson.M{"$addToSet": bson.M{"walletAddresses": bson.M{"$each": userToken.Address}}}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After) // Return the updated document

	var updatedUser models.User
	err := userService.UserRepo.MongoCollection.FindOneAndUpdate(userService.UserRepo.Ctx, filter, update, opts).Decode(&updatedUser)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil

}

// UpdateUserUpdatables implements userServiceInterface.UserService.
func (userService *UserServiceImpl) UpdateUserUpdatables(createUser *userType.CreateUserDto) (*models.User, error) {

	// find user by username
	filter := bson.M{"username": createUser.Username}
	var user models.User
	err := userService.UserRepo.MongoCollection.FindOne(userService.UserRepo.Ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// update user field
			type UpdatedFields struct {
				username     string
				displayName  string
				bio          string
				profileImage string // Optional field
				bannerImage  string // Optional field
			}
			updatedUserDto := &UpdatedFields{
				username:    createUser.Username,
				displayName: createUser.DisplayName,
				bio:         createUser.Bio,
			}

			val := reflect.ValueOf(*createUser)

			for i := 0; i < val.NumField(); i++ {
				field := val.Field(i)
				fieldName := val.Type().Field(i).Name

				if fieldName != "ProfileImage" && fieldName != "BannerImage" {
					continue
				}

				// Check if the field is not empty (has at least one file)
				if field.Len() > 0 {
					// Get the first file in the slice
					fileHeader := field.Index(0).Interface().(*multipart.FileHeader)

					file, err := fileHeader.Open()
					if err != nil {
						return nil, err
					}

					defer file.Close()

					// Read the file into a byte slice
					buffer, err := io.ReadAll(file)
					if err != nil {
						return nil, err
					}

					imageUrl, err := cloudinaryService.UploadBufferToCloudinary(userService.UserRepo.Ctx, buffer, "")
					if err != nil {
						return nil, err
					}
					if fieldName == "ProfileImage" {
						updatedUserDto.profileImage = imageUrl
					} else if fieldName == "BannerImage" {
						updatedUserDto.bannerImage = imageUrl
					}

				}
			}

			setFields := bson.M{
				"username":    updatedUserDto.username,
				"displayName": updatedUserDto.displayName,
				"bio":         updatedUserDto.bio,
			}

			// Check if profileImage is populated and add it to the query
			if updatedUserDto.profileImage != "" {
				setFields["profileImage"] = updatedUserDto.profileImage
			}

			// Check if bannerImage is populated and add it to the query
			if updatedUserDto.bannerImage != "" {
				setFields["bannerImage"] = updatedUserDto.bannerImage
			}

			// Construct the final query

			query := bson.M{"$set": setFields}
			opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

			// Update the user in the database
			fmt.Printf("user-Id: %+v\n", createUser.UserId)
			err = userService.UserRepo.MongoCollection.FindOneAndUpdate(userService.UserRepo.Ctx, bson.M{"_id": createUser.UserId}, query, opts).Decode(&user)

			if err != nil {
				return nil, err
			}

			return &user, nil
		}
		return nil, err
	}
	return &user, errors.New("user with this username already exists")
}

// GetAllUsers implements userServiceInterface.UserService.
func (u *UserServiceImpl) GetAllUsers() ([]*models.User, error) {
	panic("unimplemented")
}

// GetUserById implements userServiceInterface.UserService.
func (u *UserServiceImpl) GetUserById(*string) (*models.User, error) {
	panic("unimplemented")
}

// GetUserByUsername implements userServiceInterface.UserService.
func (u *UserServiceImpl) GetUserByUsername(username string, fetchToken bool) (*models.User, error) {
	panic("unimplemented")
}

// GetUserCommunities implements userServiceInterface.UserService.
func (u *UserServiceImpl) GetUserCommunities(*string) ([]*communityModel.Community, error) {
	panic("unimplemented")
}

// GetUserDetail implements userServiceInterface.UserService.
func (u *UserServiceImpl) GetUserDetail(*string) (*models.User, error) {
	panic("unimplemented")
}

// GetUsersByUsername implements userServiceInterface.UserService.
func (u *UserServiceImpl) GetUsersByUsername([]*string) ([]*models.User, error) {
	panic("unimplemented")
}

// SearchByUsername implements userServiceInterface.UserService.
func (u *UserServiceImpl) SearchByUsername(*string) ([]*models.User, error) {
	panic("unimplemented")
}

// UploadImage implements userServiceInterface.UserService.
func (u *UserServiceImpl) UploadImage(file interface{}, username *string) (*models.User, error) {
	panic("unimplemented")
}

// UserLeaveCommunity implements userServiceInterface.UserService.
func (u *UserServiceImpl) UserLeaveCommunity(*string, *string) (*models.User, error) {
	panic("unimplemented")
}

func NewUserService(mr *userType.UserRepo) userServiceInterface.UserService {
	return &UserServiceImpl{
		UserRepo: mr,
	}
}
