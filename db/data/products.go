package data

import (
	"errors"
	"fmt"
	"github.com/Lexxxzy/go-echo-template/db"
	"github.com/gocql/gocql"
	"github.com/labstack/gommon/log"
	"net/http"
	"strings"
)

type Product struct {
	ID           int     `bun:"type:int,pk" json:"id"`
	Name         string  `bun:"type:char(128),notnull" json:"name"`
	Price        float64 `bun:"type:decimal(10,2),notnull" json:"price"`
	Manufacturer string  `bun:"type:char(64)" json:"manufacturer"`
	TypeName     string  `json:"type_name"`
}

type CartItem struct {
	ProductID int     `bun:"type:int,pk" json:"id"`
	Product   string  `bun:"type:char(128),notnull" json:"product"`
	Price     float64 `bun:"type:decimal(10,2),notnull" json:"price"`
	Quantity  int     `bun:"type:int,notnull" json:"quantity"`
}

type Order struct {
	ID              int        `bun:"type:int,pk" json:"id"`
	DeliveryAddress string     `bun:"type:char(256),notnull" json:"delivery_address"`
	OrderDate       string     `bun:"type:timestamp,notnull" json:"order_date"`
	TotalPrice      float64    `bun:"type:decimal(10,2),notnull" json:"total_price"`
	CartItems       []CartItem `json:"cart_items"`
}

func GetAllProducts() ([]Product, error) {
	var products []Product
	iter := db.CassandraProxy.GetCurrentSession().Query("SELECT id, name, price, manufacturer, type_name FROM products").Iter()

	var (
		id           int
		name         string
		price        float64
		manufacturer string
		typeName     string
	)

	for iter.Scan(&id, &name, &price, &manufacturer, &typeName) {
		product := Product{
			ID:           id,
			Name:         name,
			Price:        price,
			Manufacturer: manufacturer,
			TypeName:     typeName,
		}
		products = append(products, product)
	}

	if err := iter.Close(); err != nil {
		log.Error("Error fetching all products: ", err)
		return nil, err
	}

	return products, nil
}

func SearchProductByName(name string) ([]Product, error) {
	var products []Product
	iter := db.CassandraProxy.GetCurrentSession().Query("SELECT id, name, price, manufacturer, type_name FROM products WHERE name LIKE ?", "%"+name+"%").Iter()

	var (
		id           int
		price        float64
		manufacturer string
		typeName     string
	)

	for iter.Scan(&id, &name, &price, &manufacturer, &typeName) {
		product := Product{
			ID:           id,
			Name:         name,
			Price:        price,
			Manufacturer: manufacturer,
			TypeName:     typeName,
		}
		products = append(products, product)
	}

	if err := iter.Close(); err != nil {
		log.Error("Error searching product by name: ", err)
		return nil, err
	}

	return products, nil
}

func GetCartItems(userID string) ([]CartItem, error) {
	var cartItems []CartItem
	iter := db.CassandraProxy.GetCurrentSession().Query("SELECT product_id, product, price, quantity FROM cart_items WHERE user_id = ?", userID).Iter()

	var (
		productID int
		product   string
		price     float64
		quantity  int
	)

	for iter.Scan(&productID, &product, &price, &quantity) {
		cartItem := CartItem{
			ProductID: productID,
			Product:   product,
			Price:     price,
			Quantity:  quantity,
		}
		cartItems = append(cartItems, cartItem)
	}

	if err := iter.Close(); err != nil {
		log.Error("Error fetching cart items: ", err)
		return nil, err
	}

	return cartItems, nil
}

func AddProductToCart(userID string, productID int, quantity int) error {
	// Start a new session
	session := db.CassandraProxy.GetCurrentSession()

	// Check if the cart exists for the user
	var cartID int
	err := session.Query("SELECT id FROM carts WHERE user_id = ?", userID).Scan(&cartID)
	if err != nil {
		// If cart does not exist, create a new one
		err = session.Query("INSERT INTO carts (user_id) VALUES (?)", userID).Exec()
		if err != nil {
			log.Error("Error creating a new cart: ", err)
			return err
		}
		// Retrieve the newly created cart ID
		err = session.Query("SELECT id FROM carts WHERE user_id = ?", userID).Scan(&cartID)
		if err != nil {
			log.Error("Error retrieving cart ID: ", err)
			return err
		}
	}

	// Add or update the product quantity in the cart
	err = session.Query(`
        INSERT INTO cart_items (cart_id, product_id, quantity)
        VALUES (?, ?, ?)
        IF NOT EXISTS
    `, cartID, productID, quantity).Exec()
	if err != nil {
		log.Error("Error adding product to cart: ", err)
		return err
	}

	return nil
}

func RemoveProductFromCart(userID string, productID int) error {
	// Start a new session
	session := db.CassandraProxy.GetCurrentSession()

	// Fetch the cart ID for the user
	var cartID int
	err := session.Query("SELECT id FROM carts WHERE user_id = ?", userID).Scan(&cartID)
	if err != nil {
		log.Error("Error fetching cart: ", err)
		return err
	}

	// Decrease the quantity of the product in the cart by 1
	updateQuantity := `
        UPDATE cart_items SET quantity = quantity - 1
        WHERE cart_id = ? AND product_id = ?
        IF quantity > 0
    `
	err = session.Query(updateQuantity, cartID, productID).Exec()
	if err != nil {
		log.Error("Error updating product quantity in cart: ", err)
		return err
	}

	// Remove the product from the cart if its quantity becomes 0
	deleteQuery := `
        DELETE FROM cart_items
        WHERE cart_id = ? AND product_id = ? IF EXISTS
    `
	err = session.Query(deleteQuery, cartID, productID).Exec()
	if err != nil {
		log.Error("Error deleting product from cart: ", err)
		return err
	}

	return nil
}

func GetOrders(userID string) ([]Order, error) {
	// Start a new session
	session := db.CassandraProxy.GetCurrentSession()

	// Fetch orders for the user
	var orders []Order
	iter := session.Query(`
        SELECT id, delivery_address, order_date
        FROM orders
        WHERE user_id = ?
    `, userID).Iter()

	for iter.Scan() {
		var order Order
		if err := iter.Scan(&order.ID, &order.DeliveryAddress, &order.OrderDate); err != true {
			log.Error("Error scanning order: ", err)
			return nil, nil
		}

		order.DeliveryAddress = strings.TrimSpace(order.DeliveryAddress)

		// Fetch order items for the order
		order.CartItems, _ = GetOrderItems(order.ID)

		// Calculate total price for the order
		totalPrice := 0.0
		for _, item := range order.CartItems {
			totalPrice += item.Price * float64(item.Quantity)
		}
		order.TotalPrice = totalPrice

		orders = append(orders, order)
	}

	if err := iter.Close(); err != nil {
		log.Error("Error closing iterator: ", err)
		return nil, err
	}

	return orders, nil
}

func GetOrderItems(orderID int) ([]CartItem, error) {
	// Start a new session
	session := db.CassandraProxy.GetCurrentSession()

	// Fetch order items for the given order ID
	var cartItems []CartItem
	iter := session.Query(`
        SELECT product_id, name, price, quantity
        FROM order_items
        WHERE order_id = ?
    `, orderID).Iter()

	for iter.Scan() {
		var cartItem CartItem
		if err := iter.Scan(&cartItem.ProductID, &cartItem.Product, &cartItem.Price, &cartItem.Quantity); err != true {
			log.Error("Error scanning order item: ", err)
			return nil, nil
		}
		cartItems = append(cartItems, cartItem)
	}

	if err := iter.Close(); err != nil {
		log.Error("Error closing iterator: ", err)
		return nil, err
	}

	return cartItems, nil
}

func PlaceOrder(userID string, deliveryAddress string) (int, int, error) {
	// Start a new session
	session := db.CassandraProxy.GetCurrentSession()

	// Generate a unique order ID (if required, depending on your schema)
	orderID := generateUniqueOrderID()

	// Step 1: Create order
	err := createOrder(session, orderID, userID, deliveryAddress)
	if err != nil {
		return http.StatusInternalServerError, orderID, fmt.Errorf("error creating order, please try again later")
	}

	// Step 2: Copy cart items to order items
	err = copyCartItemsToOrder(session, orderID, userID)
	if err != nil {
		return http.StatusInternalServerError, orderID, fmt.Errorf("error placing order, please try again later")
	}

	// Step 3: Clear cart
	err = clearCart(session, userID)
	if err != nil {
		return http.StatusInternalServerError, orderID, fmt.Errorf("error placing order, please try again later")
	}

	return http.StatusOK, orderID, nil
}

func createOrder(session *gocql.Session, orderID int, userID string, deliveryAddress string) error {
	// Insert order data into orders table
	err := session.Query(`
        INSERT INTO orders (id, user_id, delivery_address, order_date)
        VALUES (?, ?, ?, toTimestamp(now()))
    `, orderID, userID, deliveryAddress).Exec()
	if err != nil {
		return err
	}
	return nil
}

func copyCartItemsToOrder(session *gocql.Session, orderID int, userID string) error {
	// Select cart items for the user
	iter := session.Query(`
        SELECT product_id, quantity
        FROM cart_items
        WHERE user_id = ?
    `, userID).Iter()

	// Iterate over cart items and insert into order_items table
	for iter.Scan() {
		var productID, quantity int
		if err := iter.Scan(&productID, &quantity); err != true {
			return nil
		}
		if err := insertOrderItem(session, orderID, productID, quantity); err != nil {
			return err
		}
	}
	if err := iter.Close(); err != nil {
		return err
	}
	return nil
}

func insertOrderItem(session *gocql.Session, orderID int, productID int, quantity int) error {
	// Insert order item data into order_items table
	err := session.Query(`
        INSERT INTO order_items (order_id, product_id, quantity)
        VALUES (?, ?, ?)
    `, orderID, productID, quantity).Exec()
	if err != nil {
		return err
	}
	return nil
}

func clearCart(session *gocql.Session, userID string) error {
	// Delete cart items for the user
	err := session.Query(`
        DELETE FROM cart_items
        WHERE user_id = ?
    `, userID).Exec()
	if err != nil {
		return err
	}
	return nil
}

// This function generates a unique order ID according to your application's logic
func generateUniqueOrderID() int {
	// Implement your logic to generate a unique order ID
	return 0 // Placeholder implementation
}

func CancelOrder(userID string, orderID int) (int, error) {
	// Step 0: Check if the order exists and belongs to the user
	var ownerID string
	ownerQuery := "SELECT user_id FROM orders WHERE id = ?"
	err := db.CassandraProxy.GetCurrentSession().Query(ownerQuery, orderID).Scan(&ownerID)
	if err != nil {
		if errors.Is(err, gocql.ErrNotFound) {
			return http.StatusNotFound, fmt.Errorf("order not found")
		}
		log.Printf("Error fetching order owner: %s", err)
		return http.StatusInternalServerError, fmt.Errorf("error cancelling order, please try again later")
	}

	if ownerID != userID {
		return http.StatusNotFound, fmt.Errorf("order not found")
	}

	// Step 1: Delete order items associated with the order
	deleteOrderItemsQuery := "DELETE FROM order_items WHERE order_id = ?"
	if err := db.CassandraProxy.GetCurrentSession().Query(deleteOrderItemsQuery, orderID).Exec(); err != nil {
		log.Printf("Error deleting order items: %s", err)
		return http.StatusInternalServerError, fmt.Errorf("error cancelling order, please try again later")
	}

	// Step 2: Delete the order if it belongs to the user
	deleteOrderQuery := "DELETE FROM orders WHERE id = ? AND user_id = ?"
	if err := db.CassandraProxy.GetCurrentSession().Query(deleteOrderQuery, orderID, userID).Exec(); err != nil {
		log.Printf("Error deleting order: %s", err)
		return http.StatusInternalServerError, fmt.Errorf("error cancelling order, please try again later")
	}

	return http.StatusOK, nil
}
