package userController

import (
	userAuthServiceInterface "ChaiLabs/ChaiLabs/api/v1/auth/interface"
	"ChaiLabs/ChaiLabs/api/v1/auth/middleware"
	userServiceInterface "ChaiLabs/ChaiLabs/api/v1/user/interface"
	"ChaiLabs/ChaiLabs/api/v1/user/models"

	userType "ChaiLabs/ChaiLabs/api/v1/user/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService userServiceInterface.UserService
}

func NewUserController(userService userServiceInterface.UserService) UserController {
	return UserController{
		userService: userService,
	}
}

func (uc *UserController) SetUserWalletAddress(ctx *gin.Context) {

	//extract user from headers
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Unable to extract user from headers",
			"data":    nil,
		})
		return
	}

	userObj, ok := user.(*models.User) // Type assertion to convert interface{} to *models.User

	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Internal server error",
			"data":    nil,
		})
		return
	}

	var userToken userType.Token

	if err := ctx.ShouldBindJSON(&userToken); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	//Input validation for the given token
	err := userToken.ValidateToken()

	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	updatedUser, err := uc.userService.SetUserWalletAddress(&userToken, userObj)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "User wallet address updated successfully.",
		"data":    updatedUser,
	})

}

func (uc *UserController) UpdateUser(ctx *gin.Context) {

	// type UserStruct struct {
	// 	ProfileImage *multipart.FileHeader `form:"profileImage" binding:"-"`
	// 	BannerImage  *multipart.FileHeader `form:"bannerImage" binding:"-"`
	// 	Username     string                `form:"username" binding:"-"`
	// 	DisplayName  string                `form:"displayName" binding:"-"`
	// 	Bio          string                `form:"bio" binding:"-"`
	// }

	createUser := userType.CreateUserDto{}

	if err := ctx.ShouldBind(&createUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if err := createUser.ValidateCreateUser(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	//extract user from headers
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": "Unable to extract user from headers",
			"data":    nil,
		})
		return
	}

	userObj, ok := user.(*models.User) // Type assertion to convert interface{} to *models.User

	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": "Internal server error",
			"data":    nil,
		})
		return
	}

	// //set user id into createUser dto
	createUser.UserId = userObj.ID

	// ctx.JSON(http.StatusOK, gin.H{
	// 	"status":  true,
	// 	"message": "User updated successfully.",
	// 	"data":    createUser,
	// })

	updatedUser, err := uc.userService.UpdateUserUpdatables(&createUser)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "User updated successfully.",
		"data":    updatedUser,
	})

}

func (uc *UserController) RegisterUserRoutes(userAuthService userAuthServiceInterface.AuthService, ctx *gin.RouterGroup) {
	userRoute := ctx.Group("/users")
	userRoute.Use(middleware.AuthMiddleware(userAuthService))
	userRoute.PUT("/", uc.UpdateUser)
	userRoute.POST("/wallet-address", uc.SetUserWalletAddress)
}
