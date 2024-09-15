package types

import (
	userModel "ChaiLabs/ChaiLabs/api/v1/user/models"
	"ChaiLabs/constants/validation"
	"ChaiLabs/utils"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"regexp"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateUserDto struct {
	ProfileImage []*multipart.FileHeader `form:"profileImage" binding:"-" validate:"omitempty,fileType"`
	BannerImage  []*multipart.FileHeader `form:"bannerImage" binding:"-" validate:"omitempty,fileType"`
	Username     string                  `form:"username" binding:"-" validate:"omitempty,usernameValid"`
	DisplayName  string                  `form:"displayName" binding:"-" validate:"omitempty,displaynameValid"`
	Bio          string                  `form:"bio" binding:"-" validate:"omitempty,bioValid"`
	UserId       primitive.ObjectID      `form:"userId" binding:"-"`
}

// UserRepo handles operations on the user collection
type UserRepo struct {
	MongoCollection *mongo.Collection
	Ctx             context.Context
}

/*
Note: omitempty = ? means that the field is optional
*/
type UserImagesDto struct {
	ProfileImage []*multipart.FileHeader `json:"profileImage,omitempty" validate:"fileType"`
	BannerImage  []*multipart.FileHeader `json:"bannerImage,omitempty" validate:"fileType"`
}

// ValidateCreateUser validates the CreateUserDto struct.
// Returns an error.
func (createUser *CreateUserDto) ValidateCreateUser() error {
	validate := validator.New()
	if err := validate.RegisterValidation("usernameValid", userModel.ValidateUsername); err != nil {
		return err
	}

	if err := validate.RegisterValidation("displaynameValid", userModel.ValidateDisplayName); err != nil {
		return err
	}

	if err := validate.RegisterValidation("bioValid", validateBio); err != nil {
		return err
	}

	if err := validate.RegisterValidation("fileType", imageValidation); err != nil {
		return err
	}

	return validate.Struct(createUser)
}

// Token holds data related to user's auth token.
type Token struct {
	Signature string `json:"signature" validate:"required"`
	Address   string `json:"address" validate:"required,eth_addr"`
}

// ValidateToken validates the Token struct.
func (tokenDto *Token) ValidateToken() error {
	validate := validator.New()
	validate.RegisterValidation("eth_addr", ValidateEthAddress)

	return validate.Struct(tokenDto)
}

// ValidateUserImagesDto validates the UserImagesDto struct.
func ValidateUserImage(fileHeader *multipart.FileHeader) error {

	// Validate each uploaded image

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}

	defer file.Close()

	// Read the file into a byte slice
	buffer, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	isValid := utils.ValidateFileType(buffer, validation.AllowedFileTypes)

	if !isValid {
		return errors.New("invalid file type")
	}

	isValidSize := utils.ValidateFileSize(buffer)

	if !isValidSize {
		return errors.New("invalid file size")
	}

	return nil
}

// func ValidateUserImageFileSize(fileHeader *multipart.FileHeader) error {

// 	file, err := fileHeader.Open()
// 	if err != nil {
// 		return err
// 	}

// 	defer file.Close()

// 	// Read the file into a byte slice
// 	buffer, err := io.ReadAll(file)
// 	if err != nil {
// 		return err
// 	}

// 	isValid := utils.ValidateFileSize(buffer)

// 	if !isValid {
// 		return errors.New("invalid file size")
// 	}

// 	return nil

// }

// //

// func RegisterCustomValidators() {
// 	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
// 		v.RegisterValidation("fileType", func(fl validator.FieldLevel) bool {
// 			fileHeaders, ok := fl.Field().Interface().([]*multipart.FileHeader)
// 			if !ok || len(fileHeaders) == 0 {
// 				return true // No file uploaded, so no need to validate
// 			}
// 			for _, fileHeader := range fileHeaders {
// 				if err := ValidateUserImage(fileHeader); err != nil {
// 					return false
// 				}
// 			}
// 			return true
// 		})

// 		v.RegisterValidation("fileSize", func(fl validator.FieldLevel) bool {
// 			fileHeaders, ok := fl.Field().Interface().([]*multipart.FileHeader)
// 			if !ok || len(fileHeaders) == 0 {
// 				return true // No file uploaded, so no need to validate
// 			}
// 			for _, fileHeader := range fileHeaders {
// 				if err := ValidateUserImageFileSize(fileHeader); err != nil {
// 					return false
// 				}
// 			}
// 			return true
// 		})
// 	}
// }

// validateBio checks if the length of the field is less than or equal to the maximum allowed length.
// fl validator.FieldLevel
// bool

func imageValidation(fl validator.FieldLevel) bool {
	fileHeaders, ok := fl.Field().Interface().([]*multipart.FileHeader)

	if !ok || len(fileHeaders) == 0 {
		return true // No file uploaded, so no need to validate
	}
	for _, fileHeader := range fileHeaders {
		if err := ValidateUserImage(fileHeader); err != nil {
			return false
		}
	}
	return true
}

func validateBio(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	max := validation.UserValidationConstantsInstance.Bio.MaxLength

	return len(field) <= max
}

// Validate Address field
func ValidateEthAddress(fl validator.FieldLevel) bool {
	// Ethereum address regex (simplified version, adjust as needed)
	re := regexp.MustCompile("^0x[a-fA-F0-9]{40}$")
	return re.MatchString(fl.Field().String())
}
