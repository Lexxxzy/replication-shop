package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/Lexxxzy/go-echo-template/db"
	"github.com/Lexxxzy/go-echo-template/handlers"
)

func main() {
	e, err := initializeAppEnvironment()
	if err != nil {
		panic(err)
	}

	e.Logger.Fatal(e.Start(":1323"))
}

func initializeAppEnvironment() (*echo.Echo, error) {
	var isDevelopment bool
	originPath := "http://frontend"

	flag.BoolVar(&isDevelopment, "dev", false, "Use dev.env file as environment")
	flag.Parse()

	if isDevelopment {
		if err := godotenv.Load("dev.env"); err != nil {
			return nil, fmt.Errorf("error reading dev.env: %s", err.Error())
		}
		originPath = "*"
	} else {
		if err := godotenv.Load(); err != nil {
			return nil, fmt.Errorf("error reading .env: %s", err.Error())
		}
	}

	if err := db.Init(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %s", err.Error())
	}

	gob.Register(uuid.UUID{})

	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{originPath},
		AllowMethods: []string{"*"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization", "Set-Cookie"},

		AllowCredentials: true,
	}))

	store := sessions.NewCookieStore([]byte(os.Getenv("SECRET_SESSION")))
	store.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 86400 * 3, // 3 days
		// Secure: true   // HTTPS
		HttpOnly: true,
	}
	e.Use(session.Middleware(store))
	// e.Use(middleware.Secure()) // HTTPS cookies, XSS protection

	initRoutes(e)

	return e, nil
}

// initRoutes initializes the routes for the given Echo instance.
//
// e: The Echo instance to initialize the routes.
// No return values.
func initRoutes(e *echo.Echo) {
	e.POST("/login", handlers.LoginUser)
	e.POST("/register", handlers.Register)
	e.GET("/products", handlers.GetProducts)
	e.POST("/logout", handlers.LogoutUser, handlers.WithAuthentication)

	my := e.Group("/my", handlers.WithAuthentication)
	my.GET("/cart", handlers.GetCart)
	my.PUT("/cart/add", handlers.AddProductToCart)
	my.DELETE("/cart/remove", handlers.RemoveProductFromCart)

	my.GET("/orders", handlers.GetOrders)
	my.POST("/orders/add", handlers.PlaceOrder)
	my.DELETE("/orders/cancel", handlers.CancelOrder)
}
