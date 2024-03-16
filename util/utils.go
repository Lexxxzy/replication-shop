package util

import (
	"unicode"

	"github.com/BurntSushi/toml"
	"github.com/Lexxxzy/go-echo-template/internal"
	"github.com/labstack/echo/v4"
)

// Boilerplate for json response.
func JsonResponse(c echo.Context, code int, message string) error {
	return c.JSON(code, map[string]string{"message": message})
}

func JsonErrorResponse(c echo.Context, code int, message string) error {
	return c.JSON(code, map[string]string{"error": message})
}

// IsValidPassword checks if a password is valid.
//
// Parameters:
// - password: a string representing the password to be validated.
//
// Returns:
// - bool: a boolean indicating if the password is valid or not.
// - string: a message explaining why the password is invalid if it is not valid.
func IsValidPassword(password string) (bool, string) {
	if len(password) <= 7 {
		return false, "Password must be more than 7 characters."
	}

	hasUpper, hasLower, hasDigit := false, false, false
	for _, r := range password {
		hasUpper = hasUpper || unicode.IsUpper(r)
		hasLower = hasLower || unicode.IsLower(r)
		hasDigit = hasDigit || unicode.IsDigit(r)
		if hasUpper && hasLower && hasDigit {
			return true, "Password is valid."
		}
	}

	message := "Password must contain:"
	if !hasUpper {
		message += " an uppercase letter,"
	}
	if !hasLower {
		message += " a lowercase letter,"
	}
	if !hasDigit {
		message += " a number,"
	}

	return false, message[:len(message)-1] + "."
}

func LoadConfig(path string) (*types.Config, error) {
	var config types.Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
