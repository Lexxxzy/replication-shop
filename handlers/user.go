package handlers

import (
	"fmt"
	"net/http"
	"net/mail"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"

	"github.com/Lexxxzy/go-echo-template/db/data"
	"github.com/Lexxxzy/go-echo-template/util"
)

func LoginUser(c echo.Context) error {
	var user data.User

	if err := c.Bind(&user); err != nil {
		log.Error("Error binding request dumps. User was not logged in.")
		return util.JsonResponse(c, http.StatusBadRequest, "Invalid request.")
	}

	addr, err := mail.ParseAddress(user.Email)
	if err != nil {
		log.Error("Invalid email format. - " + err.Error())
		return util.JsonResponse(c, http.StatusBadRequest, "Invalid email.")
	}

	reqPassword := user.Password
	user, err = data.GetUserByEmail(addr.Address)
	if err != nil {
		log.Error("Database query failed: ", err)
		return util.JsonResponse(c, http.StatusUnauthorized, "Invalid credentials.")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqPassword)); err != nil {
		return util.JsonResponse(c, http.StatusUnauthorized, "Invalid credentials.")
	}

	if err, done := SetupUserSession(c, user); done {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"username": strings.TrimSpace(user.Name),
	})
}

func Register(c echo.Context) error {
	log.Info("Registering user.")
	var reqdata data.User

	if err := c.Bind(&reqdata); err != nil {
		log.Error("Error binding request dumps. User was not created. " + err.Error())
		return util.JsonResponse(c, http.StatusBadRequest, "Invalid request.")
	}

	log.Info(fmt.Sprintf("Request dumps: %s %s", reqdata.Name, reqdata.Email))

	addr, err := mail.ParseAddress(reqdata.Email)
	if err != nil {
		log.Error("Invalid email format. - " + err.Error())
		return util.JsonResponse(c, http.StatusBadRequest, "Invalid email.")
	}

	isExists, _ := data.IsUserExists(addr.Address)
	if isExists {
		log.Error("User already exists.")
		return util.JsonResponse(c, http.StatusBadRequest, "User already exists.")
	}

	/* isValid, message := util.IsValidPassword(reqdata.Password)
	if !isValid {
		return util.JsonResponse(c, http.StatusBadRequest, message)
	} */

	password, err := bcrypt.GenerateFromPassword([]byte(reqdata.Password), 14)
	if err != nil {
		log.Error("Error hashing password: " + err.Error())
		return util.JsonResponse(c, http.StatusInternalServerError, "Something went wrong.")
	}

	user := data.User{Name: reqdata.Name, Email: addr.Address, Password: string(password)}
	if err := data.CreateUser(&user, c); err != nil {
		log.Error("Database query failed: " + err.Error())
		return util.JsonResponse(c, http.StatusInternalServerError, "Something went wrong.")
	}

	if err, done := SetupUserSession(c, user); done {
		log.Error("Error setting up user session: ", err)
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"username": strings.TrimSpace(user.Name),
	})
}

func SetupUserSession(c echo.Context, user data.User) (error, bool) {
	sess, err := session.Get("session", c)

	if err != nil {
		log.Error("Session get error: ", err)
		return util.JsonResponse(c, http.StatusInternalServerError, "Error setting session."), true
	}

	sess.Values["authenticated"] = true
	sess.Values["userID"] = user.ID
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		log.Error("Session save error: ", err)
		return util.JsonResponse(c, http.StatusInternalServerError, "Error setting session."), true
	}

	return nil, false
}

func HealthCheck(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		log.Error("Session get error: ", err)
		return util.JsonResponse(c, http.StatusInternalServerError, "Failed to retrieve session")
	}

	isAuthenticated, ok := sess.Values["authenticated"].(bool)
	if !ok || !isAuthenticated {
		return util.JsonResponse(c, http.StatusUnauthorized, "Not Authenticated")
	}

	user, err := data.GetUserById(sess.Values["userID"].(uuid.UUID))
	if err != nil {
		return util.JsonResponse(c, http.StatusInternalServerError, "Failed to retrieve user")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"username": strings.TrimSpace(user.Name),
	})
}

func LogoutUser(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return util.JsonResponse(c, http.StatusInternalServerError, "Failed to retrieve session")
	}

	// Invalidate the session
	sess.Options.MaxAge = -1

	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return util.JsonResponse(c, http.StatusInternalServerError, "Failed to invalidate session")
	}

	return util.JsonResponse(c, http.StatusOK, "Successfully logged out")
}
