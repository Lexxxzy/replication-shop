package data

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Lexxxzy/go-echo-template/db"
	"github.com/gocql/gocql"
	"github.com/labstack/gommon/log"
	"net/http"
	"strconv"
	"strings"
)

type Product struct {
	ID           int     `cql:"id" json:"id"`
	Name         string  `cql:"name" json:"name"`
	Price        float64 `cql:"price" json:"price"`
	Manufacturer string  `cql:"manufacturer" json:"manufacturer"`
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

	iter := db.CassandraProxy.GetCurrentSession().Query(`SELECT id, manufacturer, name, price, product_type_id FROM cassandrakeyspace.products`).Iter()

	row := make(map[string]interface{})

	for iter.MapScan(row) {
		productTypeName, err := GetProductTypeName(row["product_type_id"].(int))
		if err != nil {
			return nil, err
		}
		product := Product{
			ID:           row["id"].(int),
			Name:         row["name"].(string),
			Price:        float64(row["price"].(float32)),
			Manufacturer: row["manufacturer"].(string),
			TypeName:     productTypeName,
		}
		products = append(products, product)
		row = make(map[string]interface{})
	}
	if err := iter.Close(); err != nil {
		log.Error("Error fetching all products: ", err)
		return nil, err
	}

	return products, nil
}

func GetProductTypeName(productTypeId int) (string, error) {
	var name string

	query := `SELECT name FROM cassandrakeyspace.product_types WHERE id = ?`
	iter := db.CassandraProxy.GetCurrentSession().Query(query, productTypeId).Iter()

	if iter.Scan(&name) {
		return name, nil
	} else if err := iter.Close(); err != nil {
		return "", err
	}

	return "", gocql.ErrNotFound
}

func FetchProductPrice(productID int) (float64, error) {
	var price float64
	query := `SELECT price FROM cassandrakeyspace.products WHERE id = ?`
	iter := db.CassandraProxy.GetCurrentSession().Query(query, productID).Iter()

	row := make(map[string]interface{})
	if iter.MapScan(row) {
		if rawPrice, ok := row["price"].(float32); ok {
			price = float64(rawPrice)
		} else {
			return 0, fmt.Errorf("failed to fetch or convert price for product ID %d", productID)
		}
	} else {
		return 0, fmt.Errorf("product with ID %d not found", productID)
	}

	if err := iter.Close(); err != nil {
		log.Error("Error fetching product price: ", err)
		return 0, err
	}

	return price, nil
}

func SearchProductByName(name string) ([]Product, error) {
	var products []Product
	query := `SELECT id, name, price, manufacturer, product_type_id FROM cassandrakeyspace.products WHERE name = ?`
	iter := db.CassandraProxy.GetCurrentSession().Query(query, name).Iter()

	row := make(map[string]interface{})

	for iter.MapScan(row) {
		productTypeName, err := GetProductTypeName(row["product_type_id"].(int))
		if err != nil {
			return nil, err
		}
		product := Product{
			ID:           row["id"].(int),
			Name:         row["name"].(string),
			Price:        float64(row["price"].(float32)),
			Manufacturer: row["manufacturer"].(string),
			TypeName:     productTypeName,
		}
		products = append(products, product)
		row = make(map[string]interface{})
	}
	if err := iter.Close(); err != nil {
		log.Error("Error searching product by name: ", err)
		return nil, err
	}
	return products, nil
}

func GetCartItems(userID string) ([]CartItem, error) {
	ctx := context.Background()
	rdb := db.RedisSentinelProxy.GetCurrentClient()

	cartKey := fmt.Sprintf("cart:%s", userID)

	cartItemsData, err := rdb.HGetAll(ctx, cartKey).Result()
	if err != nil {
		log.Error("Error fetching cart items from Redis: ", err)
		return nil, err
	}

	// Early return if no cart items
	if len(cartItemsData) == 0 {
		return []CartItem{}, nil
	}

	var cartItems []CartItem
	db := db.PostgresqlProxy.GetCurrentDB()
	for productIDStr, quantityStr := range cartItemsData {
		productID, _ := strconv.Atoi(productIDStr)
		quantity, _ := strconv.Atoi(quantityStr)

		query := `
            SELECT p.id, p.name, p.price
            FROM products p
            WHERE p.id = ?
        `
		row := db.QueryRow(query, productID)
		var cartItem CartItem
		err := row.Scan(&cartItem.ProductID, &cartItem.Product, &cartItem.Price)
		if err != nil {
			if err != sql.ErrNoRows {
				log.Error("Error fetching product details: ", err)
				return nil, err
			}
			continue
		}
		cartItem.Quantity = quantity
		cartItem.Product = strings.TrimSpace(cartItem.Product)

		cartItems = append(cartItems, cartItem)
	}

	return cartItems, nil
}

func AddProductToCart(userID string, productID int, quantity int) error {
	ctx := context.Background()
	rdb := db.RedisSentinelProxy.GetCurrentClient()

	cartKey := fmt.Sprintf("cart:%s", userID)

	_, err := rdb.HIncrBy(ctx, cartKey, fmt.Sprintf("%d", productID), int64(quantity)).Result()
	if err != nil {
		log.Error("Error adding/updating product in cart: ", err)
		return err
	}

	return nil
}

func RemoveProductFromCart(userID string, productID int) error {
	ctx := context.Background()
	rdb := db.RedisSentinelProxy.GetCurrentClient()

	cartKey := fmt.Sprintf("cart:%s", userID)

	newQuantity, err := rdb.HIncrBy(ctx, cartKey, fmt.Sprintf("%d", productID), -1).Result()
	if err != nil {
		log.Error("Error decrementing product quantity in cart: ", err)
		return err
	}

	if newQuantity <= 0 {
		_, err = rdb.HDel(ctx, cartKey, fmt.Sprintf("%d", productID)).Result()
		if err != nil {
			log.Error("Error removing product from cart: ", err)
			return err
		}
	}

	return nil
}

func GetOrders(userID string) ([]Order, error) {
	var orders []Order
	query := `
	SELECT o.id, o.delivery_address, o.order_date 
	FROM orders o
	WHERE user_id = ?
    `
	rows, err := db.PostgresqlProxy.GetCurrentDB().Query(query, userID)
	if err != nil {
		log.Error("Error fetching orders: ", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.DeliveryAddress, &order.OrderDate)

		order.DeliveryAddress = strings.TrimSpace(order.DeliveryAddress)

		if err != nil {
			log.Error("Error scanning order: ", err)
			return nil, err
		}

		order.CartItems, err = GetOrderItems(order.ID)
		if err != nil {
			return nil, err
		}

		totalPrice := 0.0
		for _, item := range order.CartItems {
			totalPrice += item.Price * float64(item.Quantity)
		}
		order.TotalPrice = totalPrice

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		log.Error("Error iterating rows: ", err)
		return nil, err
	}

	return orders, nil
}

func GetOrderItems(orderID int) ([]CartItem, error) {
	var cartItems []CartItem
	query := `
	SELECT p.id, p.name, p.price, oi.quantity
	FROM order_items oi
	JOIN products p ON oi.product_id = p.id
	WHERE oi.order_id = ?
    `
	rows, err := db.PostgresqlProxy.GetCurrentDB().Query(query, orderID)
	if err != nil {
		log.Error("Error fetching order items: ", err)
		return nil, err
	}
	defer rows.Close()
	cartItems, err = MapRowsToCartItems(rows)
	if err != nil {
		return nil, err
	}

	return cartItems, nil
}

func PlaceOrder(userID string, deliveryAddress string) (int, error) {
	ctx := context.Background()
	rdb := db.RedisSentinelProxy.GetCurrentClient()

	cartKey := fmt.Sprintf("cart:%s", userID)

	cartItemCount, err := rdb.HLen(ctx, cartKey).Result()
	if err != nil || cartItemCount == 0 {
		log.Error("Error checking cart items or cart is empty: ", err)
		return 0, fmt.Errorf("cart is empty")
	}

	var orderID int
	orderQuery := `
        INSERT INTO orders (user_id, delivery_address, order_date)
        VALUES (?, ?, NOW())
        RETURNING id
    `
	err = db.PostgresqlProxy.GetCurrentDB().QueryRow(orderQuery, userID, deliveryAddress).Scan(&orderID)
	if err != nil {
		log.Error("Error creating order: ", err)
		return orderID, fmt.Errorf("error creating order, please try again later")
	}

	cartItemsData, err := rdb.HGetAll(ctx, cartKey).Result()
	if err != nil {
		log.Error("Error fetching cart items from Redis: ", err)
		return orderID, fmt.Errorf("error placing order, please try again later")
	}

	for productIDStr, quantityStr := range cartItemsData {
		productID, _ := strconv.Atoi(productIDStr)
		quantity, _ := strconv.Atoi(quantityStr)
		price, err := FetchProductPrice(productID)

		if err != nil {
			log.Error("Error fetching product price")
			return orderID, fmt.Errorf("error fetching product price, pleace try later")
		}

		insertOrderItemQuery := `
            INSERT INTO order_items (order_id, product_id, quantity, price_at_order)
            VALUES (?, ?, ?, ?)
        `
		_, err = db.PostgresqlProxy.GetCurrentDB().Exec(insertOrderItemQuery, orderID, productID, quantity, price)
		if err != nil {
			log.Error("Error inserting order item: ", err)
			return orderID, fmt.Errorf("error placing order, please try again later")
		}
	}

	_, err = rdb.Del(ctx, cartKey).Result()
	if err != nil {
		log.Error("Error clearing cart: ", err)
		return orderID, fmt.Errorf("error placing order, please try again later")
	}

	return orderID, nil
}

func CancelOrder(userID string, orderID int) (int, error) {
	tx, err := db.PostgresqlProxy.GetCurrentDB().Begin()
	if err != nil {
		log.Error("Error starting transaction: ", err)
		return http.StatusInternalServerError, fmt.Errorf("error cancelling order, please try again later")
	}
	// Шаг 0: Проверка, что заказ принадлежит пользователю
	var ownerID string
	ownerQuery := `SELECT user_id FROM orders WHERE id = ?`
	err = tx.QueryRow(ownerQuery, orderID).Scan(&ownerID)
	if ownerID == "" || ownerID != userID {
		tx.Rollback()
		return http.StatusNotFound, fmt.Errorf("order not found")
	}

	if err != nil {
		tx.Rollback()
		log.Error("Error fetching order owner: ", err)
		return http.StatusInternalServerError, fmt.Errorf("error cancelling order, please try again later")
	}

	// Шаг 1: Удаление содержимого заказа
	deleteOrderItemsQuery := `DELETE FROM order_items WHERE order_id = ?`
	_, err = tx.Exec(deleteOrderItemsQuery, orderID)
	if err != nil {
		tx.Rollback()
		log.Error("Error deleting order items: ", err)
		return http.StatusInternalServerError, fmt.Errorf("error cancelling order, please try again later")
	}

	// Шаг 2: Удаление заказа
	deleteOrderQuery := `DELETE FROM orders WHERE id = ? AND user_id = ?`
	_, err = tx.Exec(deleteOrderQuery, orderID, userID)
	if err != nil {
		tx.Rollback()
		log.Error("Error deleting order: ", err)
		return http.StatusInternalServerError, fmt.Errorf("error cancelling order, please try again later")
	}

	// Завершение транзакции
	err = tx.Commit()
	if err != nil {
		log.Error("Error committing transaction: ", err)
		return http.StatusInternalServerError, fmt.Errorf("error cancelling order, please try again later")
	}

	return http.StatusOK, nil
}

func MapRowsToProducts(rows *sql.Rows) ([]Product, error) {
	var products []Product
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Manufacturer, &product.TypeName)

		product.Name = strings.TrimSpace(product.Name)
		product.Manufacturer = strings.TrimSpace(product.Manufacturer)
		product.TypeName = strings.TrimSpace(product.TypeName)

		if err != nil {
			log.Error("Error scanning product: ", err)
			return nil, err
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		log.Error("Error iterating rows: ", err)
		return nil, err
	}

	return products, nil
}

func MapRowsToCartItems(rows *sql.Rows) ([]CartItem, error) {
	var cartItems []CartItem
	for rows.Next() {
		var cartItem CartItem
		err := rows.Scan(&cartItem.ProductID, &cartItem.Product, &cartItem.Price, &cartItem.Quantity)

		cartItem.Product = strings.TrimSpace(cartItem.Product)

		if err != nil {
			log.Error("Error scanning cart item: ", err)
			return nil, err
		}
		cartItems = append(cartItems, cartItem)
	}

	if err := rows.Err(); err != nil {
		log.Error("Error iterating rows: ", err)
		return nil, err
	}

	return cartItems, nil
}
