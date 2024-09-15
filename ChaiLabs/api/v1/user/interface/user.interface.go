package userServiceInterface

import (
	communityModel "ChaiLabs/ChaiLabs/api/v1/community/models"
	userModel "ChaiLabs/ChaiLabs/api/v1/user/models"
	userType "ChaiLabs/ChaiLabs/api/v1/user/types"
)

type UserService interface {
	SetUserWalletAddress(*userType.Token, *userModel.User) (*userModel.User, error)
	UploadImage(file interface{}, username *string) (*userModel.User, error)
	SearchByUsername(*string) ([]*userModel.User, error)
	UpdateUserUpdatables(*userType.CreateUserDto) (*userModel.User, error)
	GetAllUsers() ([]*userModel.User, error)
	GetUserDetail(*string) (*userModel.User, error)
	GetUserByUsername(username string, fetchToken bool) (*userModel.User, error)
	GetUserById(*string) (*userModel.User, error)
	GetUserCommunities(*string) ([]*communityModel.Community, error)
	UserLeaveCommunity(*string, *string) (*userModel.User, error)
	GetUsersByUsername([]*string) ([]*userModel.User, error)
}
