package middleware

import (
	"fmt"
	"net/http"

	helper "github.com/clinton-felix/golang-JWT-auth-project/helpers"
	"github.com/gin-gonic/gin"
)

/*
	the Authenticat function ensures that the protected routes are not publicly
	accessible to just anyone, but instead to those who have a token and are
	indexed in the system. The other routes like the signup and login routes can
	be accessed publicly and do not need a middleware authentication layer
*/
func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// check if the request query has a token
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error":fmt.Sprintf("No Authorization header provided")})
			c.Abort()
			return
		}

		// claims are of type signedDetails
		claims, err := helper.ValidateToken(clientToken)
		// we use empty string here because we are return err as a string msg from ValidateToken()
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"err": err})
			c.Abort()
			return
		}
		c.Set("email", claims.Email)
		c.Set("first_name", claims.First_name)
		c.Set("last_name", claims.Last_name)
		c.Set("uid", claims.Uid)
		c.Set("user_type", claims.User_type)
		c.Next()
	}
}