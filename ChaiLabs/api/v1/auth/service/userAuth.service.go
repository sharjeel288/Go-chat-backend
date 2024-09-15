package service

import (
	userAuthInterface "ChaiLabs/ChaiLabs/api/v1/auth/interface"
	"ChaiLabs/ChaiLabs/api/v1/user/models"
	userType "ChaiLabs/ChaiLabs/api/v1/user/types"
	"ChaiLabs/utils"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewUserRepo creates a new instance of UserRepo

type AuthServiceImpl struct {
	UserRepo *userType.UserRepo
}

func NewAuthService(mr *userType.UserRepo) userAuthInterface.AuthService {
	return &AuthServiceImpl{
		UserRepo: mr,
	}
}

func (authService *AuthServiceImpl) SetUserTokenMessage(address *string) (*string, error) {
	//Getting message to sign
	message := utils.CreateMessageForSigning(*address)

	token := &models.Token{
		Signature: "",
		Expiry:    time.Now().Add(time.Minute*1).Unix() * 1000, // in milliseconds
		Message:   message,
	}

	filter := bson.M{"walletAddresses": bson.M{"$in": []string{*address}}}

	var user models.User
	err := authService.UserRepo.MongoCollection.FindOne(authService.UserRepo.Ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			//Creating new user
			newUser := &models.User{
				WalletAddresses: []string{*address},
				Token:           *token,
			}

			if _, err := authService.UserRepo.MongoCollection.InsertOne(authService.UserRepo.Ctx, newUser); err != nil {
				return nil, err
			}
			return &token.Message, nil
		}
		return nil, err
	}

	query := bson.M{"$set": bson.M{
		"token": *token,
	}}

	result, updateErr := authService.UserRepo.MongoCollection.UpdateOne(authService.UserRepo.Ctx, filter, query)
	if updateErr != nil {
		return nil, updateErr

	} else if result.MatchedCount == 0 {
		return nil, errors.New("unauthorized access, user or token not found")
	}

	return &token.Message, nil

}

/*
Another way of implementing SetUserTokenMessage in which we make 2 db calls
but response time is poor as compared to the above implementation
in which we making 3 db calls but it has better response time
*/

// func (authService *AuthServiceImpl) SetUserTokenMessage(address *string) (*string, error) {
// 	// Getting message to sign
// 	message := utils.CreateMessageForSigning(*address)

// 	token := &models.Token{
// 		Signature: "",
// 		Expiry:    time.Now().Add(time.Minute * 1).Unix(),
// 		Message:   message,
// 	}

// 	filter := bson.M{"walletAddresses": bson.M{"$in": []string{*address}}}

// 	// Check if the document exists and whether walletAddresses is an array
// 	var existingUser models.User
// 	err := authService.UserRepo.MongoCollection.FindOne(authService.UserRepo.Ctx, filter).Decode(&existingUser)
// 	if err != nil && err != mongo.ErrNoDocuments {
// 		return nil, fmt.Errorf("error checking user: %v", err)
// 	}

// 	var update bson.M
// 	if err == mongo.ErrNoDocuments {
// 		// Document does not exist, create a new user
// 		update = bson.M{
// 			"$set": bson.M{
// 				"walletAddresses": []string{*address},
// 				"token":           *token,
// 			},
// 		}
// 	} else {
// 		// Document exists, update the token and add to walletAddresses if necessary
// 		update = bson.M{
// 			"$set":      bson.M{"token": *token},
// 			"$addToSet": bson.M{"walletAddresses": *address},
// 		}
// 	}

// 	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)
// 	var updatedUser models.User
// 	err = authService.UserRepo.MongoCollection.FindOneAndUpdate(authService.UserRepo.Ctx, filter, update, opts).Decode(&updatedUser)
// 	if err != nil {
// 		return nil, fmt.Errorf("error updating user token: %v", err)
// 	}

// 	return &token.Message, nil
// }

// CheckValidity implements userAuthInterface.AuthService.
func (authService *AuthServiceImpl) CheckValidity(token *userType.Token) (bool, error) {

	filter := bson.M{
		"walletAddresses": bson.M{"$in": []string{token.Address}},
		"token.signature": token.Signature,
		"token.expiry":    bson.M{"$gt": time.Now().UnixNano() / int64(time.Millisecond)}, // in milliseconds
	}

	var result models.User

	err := authService.UserRepo.MongoCollection.FindOne(authService.UserRepo.Ctx, filter).Decode(&result)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}

	return true, nil

}

// CreateUserSignatureForTesting implements userAuthInterface.AuthService.
func (authService *AuthServiceImpl) CreateUserSignatureForTesting(address *string) (*userType.Token, error) {
	user, err := authService.GetUserMessage(address)
	if err != nil {
		return nil, err
	}
	signature, _, err := utils.CreateSignature(user.Token.Message)
	if err != nil {
		return nil, err
	}
	userToken := &userType.Token{
		Signature: signature,
		Address:   *address,
	}
	return userToken, nil
}

// GetUserMessage implements userAuthInterface.AuthService.
func (authService *AuthServiceImpl) GetUserMessage(address *string, opts ...bool) (*models.User, error) {

	isWsConnection := false
	if len(opts) > 0 {
		isWsConnection = opts[0]
	}

	filter := bson.M{
		"walletAddresses": bson.M{"$in": []string{*address}},
		"token.expiry":    bson.M{"$gt": time.Now().UnixNano() / int64(time.Millisecond)}, // convert to milliseconds
	}

	var user models.User

	err := authService.UserRepo.MongoCollection.FindOne(authService.UserRepo.Ctx, filter).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			if isWsConnection {
				return nil, err
			}
			return nil, err
		}
		return nil, err
	}

	if user.Token != (models.Token{}) {
		return &user, nil
	}

	if isWsConnection {
		return nil, errors.New("token not found or token has expired")
	}
	return nil, errors.New("token not found or token has expired")
}

// GetUserMessageIfAuthenticated implements userAuthInterface.AuthService.
func (authService *AuthServiceImpl) GetUserMessageIfAuthenticated(signature *string) (*string, error) {

	var user models.User
	err := authService.UserRepo.MongoCollection.FindOne(authService.UserRepo.Ctx, bson.M{"token.signature": *signature}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("unauthorized access, user or token not found")
		}
		return nil, err
	}
	if user.Token == (models.Token{}) {
		return nil, errors.New("token not found or token has expired")
	}

	if user.Token.Expiry/1000 < time.Now().Unix() {
		return nil, errors.New("user's token has expired")
	}
	return &user.Token.Message, nil
}

// LoginUser implements userAuthInterface.AuthService.
func (authService *AuthServiceImpl) LoginUser(token *userType.Token) (*models.User, error) {

	filter := bson.M{"walletAddresses": bson.M{"$in": []string{token.Address}}}

	var user models.User
	err := authService.UserRepo.MongoCollection.FindOne(authService.UserRepo.Ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("unauthorized access, user or token not found")
		}
		return nil, err
	}
	if user.Token == (models.Token{}) {
		return nil, errors.New("token not found or token has expired")
	}

	isVerified := utils.VerifySignature(token.Signature, token.Address, user.Token.Message)

	if !isVerified {
		return nil, errors.New("unauthorized access, invalid signature")
	}

	if user.Token.Expiry/1000 < time.Now().Unix() {
		return nil, errors.New("user's token has expired")
	}

	// update user's token expiry and signature

	// set user token expiry
	user.Token.Expiry = time.Now().Add(time.Hour*24).Unix() * 1000 // convert to milliseconds

	// set user token signature
	user.Token.Signature = token.Signature

	query := bson.M{"$set": bson.M{
		"token.signature": token.Signature, "token.expiry": user.Token.Expiry}}

	result, updateErr := authService.UserRepo.MongoCollection.UpdateOne(authService.UserRepo.Ctx, filter, query)
	if updateErr != nil {
		return nil, updateErr

	} else if result.MatchedCount == 0 {
		return nil, errors.New("unauthorized access, user or token not found")
	}

	return &user, nil
}

// ValidateUser implements userAuthInterface.AuthService.
func (authService *AuthServiceImpl) ValidateUser(token *userType.Token, opts ...bool) (*models.User, error) {

	isWsConnection := false
	if len(opts) > 0 {
		isWsConnection = opts[0]
	}

	user, err := authService.GetUserMessage(&token.Address, isWsConnection)

	if err != nil {
		return nil, err
	}

	isValid, validityErr := authService.CheckValidity(token)

	if !isValid || validityErr != nil {
		if isWsConnection {
			return nil, errors.New("unauthorized access, token not found or token has expired")
		}
		return nil, validityErr

	}

	isVerified := utils.VerifySignature(token.Signature, token.Address, user.Token.Message)

	if !isVerified {
		if isWsConnection {
			return nil, errors.New("unauthorized access, invalid signature")
		}
		return nil, errors.New("unauthorized access, invalid signature")
	}

	return user, nil
}
