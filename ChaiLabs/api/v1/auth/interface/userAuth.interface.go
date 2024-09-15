package userAuthInterface

import (
	"ChaiLabs/ChaiLabs/api/v1/user/models"
	userType "ChaiLabs/ChaiLabs/api/v1/user/types"
)

type AuthService interface {
	SetUserTokenMessage(*string) (*string, error)
	ValidateUser(*userType.Token, ...bool) (*models.User, error)
	CheckValidity(*userType.Token) (bool, error)
	GetUserMessage(*string, ...bool) (*models.User, error)
	LoginUser(*userType.Token) (*models.User, error)
	CreateUserSignatureForTesting(*string) (*userType.Token, error)
	GetUserMessageIfAuthenticated(*string) (*string, error)
}
