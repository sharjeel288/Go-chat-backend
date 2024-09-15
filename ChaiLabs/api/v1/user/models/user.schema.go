package models

import (
	"ChaiLabs/constants"
	"ChaiLabs/constants/validation"
	"errors"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Token holds data related to user's token.
type Token struct {
	Signature string `json:"signature" bson:"signature,omitempty"`
	Expiry    int64  `json:"expiry" bson:"expiry,omitempty"`
	Message   string `json:"message" bson:"message,omitempty" validate:"required"`
}

// User represents a user in the system.
type User struct {
	ID              primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	WalletAddresses []string             `json:"walletAddresses" bson:"walletAddresses,omitempty" validate:"required"`
	Username        string               `json:"username" bson:"username,omitempty" validate:"required,usernameValid"`
	DisplayName     string               `json:"displayName" bson:"displayName,omitempty" validate:"nameValid"`
	Token           Token                `json:"token" bson:"token,omitempty" validate:"required,dive"`
	Communities     []primitive.ObjectID `json:"communities" bson:"communities,omitempty"`
	JoinedAt        time.Time            `json:"joinedAt" bson:"joinedAt,omitempty"`
	ProfileImage    string               `json:"profileImage" bson:"profileImage,omitempty"`
	BannerImage     string               `json:"bannerImage" bson:"bannerImage,omitempty"`
	Bio             string               `json:"bio" bson:"bio,omitempty" validate:"min=0,max=200"`
}

// ValidateUser performs validations on a User instance.
func (u *User) ValidateUser() error {
	validate := validator.New()
	if err := validate.RegisterValidation("usernameValid", ValidateUsername); err != nil {
		return errors.New("Invalid:username")
	}
	if err := validate.RegisterValidation("nameValid", ValidateDisplayName); err != nil {
		return errors.New("Invalid:displayName")
	}

	return validate.Struct(u)
}

// ValidateUsername checks if the username is valid and not reserved.
func ValidateUsername(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	re := regexp.MustCompile(validation.UserValidationConstantsInstance.Username.Regex)

	if val := re.MatchString(value); val {
		return !IsReservedUsername(value)
	}
	return false
}

// validateName checks if the display name is valid and not reserved.
func ValidateDisplayName(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	re := regexp.MustCompile(validation.UserValidationConstantsInstance.DisplayName.Regex)

	if val := re.MatchString(value); val {
		return !IsReservedUsername(value)
	}
	return false
}

// IsReservedUsername checks if a username is in the reserved usernames list.
func IsReservedUsername(username string) bool {
	for _, reserved := range constants.ReservedUsernames {
		if username == reserved {
			return true
		}
	}
	return false
}
