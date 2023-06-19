package middleware

// import (
// 	"github.com/gin-gonic/gin"
// 	"net/http"
// )

// func ErrorHandlerMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		c.Next() // Execute the next middleware and any remaining handlers

// 		// Check if there are any errors
// 		errors := c.Errors
// 		if len(errors) > 0 {
// 			// Get the last error
// 			err := errors.Last()

// 			// Define the custom error response
// 			var customError struct {
// 				StatusCode int    `json:"statusCode"`
// 				Msg        string `json:"msg"`
// 			}

// 			// Set the default values
// 			customError.StatusCode = http.StatusInternalServerError
// 			customError.Msg = "Something went wrong, please try again later"

// 			// Handle specific error cases
// 			switch err.Err.(type) {
// 			case gin.ErrorType:
// 				errorType := err.Err.(gin.ErrorType)
// 				switch errorType {
// 				case gin.ErrorTypePublic:
// 					customError.Msg = err.Meta.(string)
// 					customError.StatusCode = http.StatusBadRequest // Set the desired status code
// 				case gin.ErrorTypeBind:
// 					customError.Msg = err.Err.Error()
// 					customError.StatusCode = http.StatusBadRequest // Set the desired status code
// 				case gin.ErrorTypeRender:
// 					customError.Msg = err.Err.Error()
// 					customError.StatusCode = http.StatusInternalServerError // Set the desired status code
// 				}
// 			}

// 			// Return the custom error response as JSON
// 			c.JSON(customError.StatusCode, gin.H{"msg": customError.Msg})
// 			c.Abort() // Stop processing remaining handlers
// 		}
// 	}
// }


