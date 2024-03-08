package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// WithAuthentication is a middleware function that adds authentication to the request handling chain.
//
// It takes a next echo.HandlerFunc as a parameter and returns an echo.HandlerFunc.
// The next handler function is called after the authentication is performed.
// It retrieves the session from the echo.Context and checks if the user is authenticated.
// If the session is not found or the user is not authenticated, it returns an error response.
// Otherwise, it sets the userID in the context and calls the next handler function.
// The function returns an error value if there is an error during the authentication process.
func WithAuthentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, err := session.Get("session", c)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to retrieve session."})
		}

		userID, ok := sess.Values["userID"].(uuid.UUID)
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "User not authorized or invalid session data"})
		}

		c.Set("userID", userID)

		return next(c)
	}
}
