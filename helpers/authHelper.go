package helper

import (
	"errors"

	"github.com/gin-gonic/gin"
)

// Validating the userType
func CheckUserType(c *gin.Context, role string) (err error) {
	userType := c.GetString("user_type")
	err = nil

	if userType != role {
		err = errors.New("unauthorized to access this recource")
		return err
	}

	return err
}


/* 
	The Logic below ensures that only ADMIN role gets acces to any users data
	or a user can only get access to his own data and not of any other
*/
func MatchUserTypeToUid(c *gin.Context, userId string) (err error) {
	userType := c.GetString("user_type")
	uId := c.GetString("uid")
	err = nil

	if userType	== "USER" && uId != userId {
		err	= errors.New("unauthorized to access this resource")
		return err
	}

	err = CheckUserType(c, userType)
		return err
}