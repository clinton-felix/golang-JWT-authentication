package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/clinton-felix/golang-JWT-auth-project/database"
	helper "github.com/clinton-felix/golang-JWT-auth-project/helpers"
	"github.com/clinton-felix/golang-JWT-auth-project/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// instantiating a new user..
var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()
// Hashing the password
func HashPasword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic()
	}
	return string(bytes)
}

// verify password function..
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("email or password is incorrect")
		check = false
	}
	return check, msg
}

// creating the signup function. This creates a user in the Database
func Signup()gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User 
		
		if err := c.BindJSON(&user) ; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
			return
		}

		/* 	
			Handling validation. Compares the data in user against the instance
			of user struct that we have declared in models to validate its schema
		*/
		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		/* 
			Count of the user in the database is used in Signup function to help us
			validate a user signup. However in the other functions like login, it
			is used to keep track of the count of users in the system..
		*/
		// Checking Email
		count, err := userCollection.CountDocuments(ctx, bson.M{"email" : user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while checking the email."})
		}

		// setting Password
		password := HashPasword(*user.Password)
		user.Password = &password

		// setting phone details
		count, err = userCollection.CountDocuments(ctx, bson.M{"phone" : user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error" : "Error occured while checking the phone no.."})
		}

		if count > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "This Email or Phone number already exists.."})
		}

				/*	creating instances of remaining struct fields */
		user.Created_at, err = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while checking the email."})
		}
		user.Updated_at, err = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while checking the email."})
		}
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, err := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, user.User_id)
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while checking the email."})
		}
		user.Token = &token
		user.Refresh_token = &refreshToken

		// having set all the params required, we need to insert it into the database...
		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("user item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		// send back the insertion number
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

// creating the Login function..
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user) ; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// finding the user being serched by email, and assigning it to foundUser
		err := userCollection.FindOne(ctx, bson.M{"email":user.Email}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":"Email or Password is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if !passwordIsValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		if foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		}
		token, refreshToken, err := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, *foundUser.User_type, foundUser.User_id)
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error":"Error generating login token"})
		}
		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)
		err = userCollection.FindOne(ctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, foundUser)
	}	
}

/*
	creating the GetUsers and GetUser function
	these are functions that can only be callled by the Admin of the system..
*/
func GetUsers() gin.HandlerFunc{
	return func(c *gin.Context) {
		// confirm that the query is from ADMIN role
		if err := helper.CheckUserType(c, "ADMIN") ; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second) 

		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10		// set recordPerPage to 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))
		if err != nil {
			log.Fatal(err)
		}

		// group users based on id and get total count using the $sum.
		// $sum isn't really a pipelining function, but computational
		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{
			{"_id", bson.D{{"_id", "null"}}}, 
			{"total_count", bson.D{{"$sum", 1}}}, 
			{"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},}}}

		// opening the userCollection Aggregation func
		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage,
		})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while listing user items"})
		}
		var allUsers []bson.M
		if err = result.All(ctx, &allUsers) ; err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allUsers[0])
	}
}

func GetUser() gin.HandlerFunc{
	return func(c *gin.Context) {
		userId := c.Param("user_id")
		if err := helper.MatchUserTypeToUid(c, userId) ; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error" : err.Error()})
			return
		}

		/*	Working now with the database to get the user..*/

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		/* 
			NB: finding the user with the provided userID and decoding since Golang 
			does not understand json. Cancel() to close the connection should be defered
		*/
		err := userCollection.FindOne(ctx, bson.M{"user_id" : userId}).Decode(&user)
		defer cancel()		// this will run after the block execution
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}

}