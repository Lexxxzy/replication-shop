package handlers

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"

	"github.com/Lexxxzy/go-echo-template/db/data"
	"github.com/Lexxxzy/go-echo-template/util"
)

func GetProducts(c echo.Context) error {
	name := c.QueryParam("title")

	if name == "" {
		products, err := data.GetAllProducts()
		if err != nil {
			log.Error("Database query failed: ", err)
			return util.JsonResponse(c, http.StatusInternalServerError, "Error fetching products. Please try again later.")
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"products": products,
		})
	}

	products, err := GetProductByName(name)
	if err != nil {
		return util.JsonResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"products": products,
	})
}

func GetProductByName(name string) ([]data.Product, error) {
	products, err := data.SearchProductByName(name)
	if err != nil {
		log.Error("Database query failed: ", err)
		return nil, fmt.Errorf("error fetching products, please try again later")
	}

	return products, nil
}

func GetCart(c echo.Context) error {
	owner, ok := c.Get("userID").(uuid.UUID)
	if !ok {
		return util.JsonResponse(c, http.StatusUnauthorized, "Unauthorized.")
	}

	cart, err := data.GetCartItems(owner.String())
	if err != nil {
		log.Error("Database query failed: ", err)
		return util.JsonResponse(c, http.StatusInternalServerError, "Error fetching cart. Please try again later.")
	}
	total := 0.0
	for _, item := range cart {
		total += item.Price * float64(item.Quantity)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"total": total,
		"cart":  cart,
	})
}

func AddProductToCart(c echo.Context) error {
	owner, ok := c.Get("userID").(uuid.UUID)
	if !ok {
		return util.JsonResponse(c, http.StatusUnauthorized, "Unauthorized.")
	}

	var cartItem = struct {
		ID       int `json:"item_id"`
		Quantity int `json:"quantity"`
	}{}

	if err := c.Bind(&cartItem); err != nil {
		log.Error("Error binding request data. Cart item was not added.")
		return util.JsonResponse(c, http.StatusBadRequest, "Invalid request.")
	}

	if err := data.AddProductToCart(owner.String(), cartItem.ID, cartItem.Quantity); err != nil {
		log.Error("Database query failed: ", err)
		return util.JsonResponse(c, http.StatusInternalServerError, "Error adding product to cart. Please try again later.")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Product added to cart.",
	})
}

func RemoveProductFromCart(c echo.Context) error {
	owner, ok := c.Get("userID").(uuid.UUID)
	if !ok {
		return util.JsonResponse(c, http.StatusUnauthorized, "Unauthorized.")
	}

	var cartItem = struct {
		ID int `json:"item_id"`
	}{}

	if err := c.Bind(&cartItem); err != nil {
		log.Error("Error binding request data. Cart item was not removed.")
		return util.JsonResponse(c, http.StatusBadRequest, "Invalid request.")
	}

	if err := data.RemoveProductFromCart(owner.String(), cartItem.ID); err != nil {
		log.Error("Database query failed: ", err)
		return util.JsonResponse(c, http.StatusInternalServerError, "Error removing product from cart. Please try again later.")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Product removed from cart.",
	})
}

func GetOrders(c echo.Context) error {
	owner, ok := c.Get("userID").(uuid.UUID)
	if !ok {
		return util.JsonResponse(c, http.StatusUnauthorized, "Unauthorized.")
	}

	orders, err := data.GetOrders(owner.String())
	if err != nil {
		log.Error("Database query failed: ", err)
		return util.JsonResponse(c, http.StatusInternalServerError, "Error fetching orders. Please try again later.")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"orders": orders,
	})
}

func PlaceOrder(c echo.Context) error {
	owner, ok := c.Get("userID").(uuid.UUID)
	if !ok {
		return util.JsonResponse(c, http.StatusUnauthorized, "Unauthorized.")
	}
	deliveryAddress := c.FormValue("delivery_address")

	status, oid, err := data.PlaceOrder(owner.String(), deliveryAddress)

	if err != nil {
		return util.JsonErrorResponse(c, status, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"order_id": oid,
	})
}

func CancelOrder(c echo.Context) error {
	owner, ok := c.Get("userID").(uuid.UUID)
	if !ok {
		return util.JsonResponse(c, http.StatusUnauthorized, "Unauthorized.")
	}

	var orderID = struct {
		ID int `json:"order_id"`
	}{}

	if err := c.Bind(&orderID); err != nil {
		log.Error("Error binding request data. Order was not cancelled.")
		return util.JsonResponse(c, http.StatusBadRequest, "Invalid request.")
	}

	if status, err := data.CancelOrder(owner.String(), orderID.ID); err != nil {
		log.Error("Database query failed: ", err)
		return util.JsonErrorResponse(c, status, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Order cancelled successfully.",
	})
}
