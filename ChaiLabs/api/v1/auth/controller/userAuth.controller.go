package authController

import (
	userAuthInterface "ChaiLabs/ChaiLabs/api/v1/auth/interface"

	userType "ChaiLabs/ChaiLabs/api/v1/user/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthService userAuthInterface.AuthService
}

func NewAuthControllerAuth(authService userAuthInterface.AuthService) AuthController {
	return AuthController{
		AuthService: authService,
	}
}

func (ac *AuthController) GetMessageForSigning(ctx *gin.Context) {
	address := ctx.Param("address")
	userMessage, err := ac.AuthService.SetUserTokenMessage(&address)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Message generated successfully.",
		"data":    userMessage,
	})

}

// desc: fake signature creation
func (ac *AuthController) CreateUserSignatureForTesting(ctx *gin.Context) {
	address := ctx.Param("address")
	userToken, err := ac.AuthService.CreateUserSignatureForTesting(&address)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "Signature generated successfully.",
		"data":    userToken,
	})
}

// desc: login user
func (ac *AuthController) LoginUser(ctx *gin.Context) {

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

	user, err := ac.AuthService.LoginUser(&userToken)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  true,
		"message": "User logged in successfully.",
		"data":    user,
	})
}

//Register Auth Routes

func (ac *AuthController) RegisterUserAuthRoutes(rq *gin.RouterGroup) {
	userAuthRoute := rq.Group("/auth")
	userAuthRoute.GET("/sign-message/:address", ac.GetMessageForSigning)
	userAuthRoute.GET("/fake-signature/:address", ac.CreateUserSignatureForTesting)
	userAuthRoute.POST("/login", ac.LoginUser)
}
