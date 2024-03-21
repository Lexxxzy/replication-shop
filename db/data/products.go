package data

import (
	"database/sql"
	"fmt"
	"github.com/Lexxxzy/go-echo-template/db"
	"github.com/gocql/gocql"
	"github.com/labstack/gommon/log"
	"net/http"
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

func SearchProductByName(name string) ([]Product, error) {
	var products []Product
	query := `SELECT id, name, price, manufacturer, product_type_id FROM cassandrakeyspace.products WHERE name ILIKE ?`
	iter := db.CassandraProxy.GetCurrentSession().Query(query, "%"+name+"%").Iter()

	row := make(map[string]interface{})

	for iter.MapScan(row) {
		if !iter.MapScan(row) {
			break
		}
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
	var cartItems []CartItem

	query := `
        SELECT p.id, p.name, p.price, ci.quantity
        FROM cart_items ci
        JOIN products p ON ci.product_id = p.id
        JOIN cart c ON ci.cart_id = c.id
        JOIN product_types pt ON p.product_type_id = pt.id
        WHERE c.user_id = ?
    `

	rows, err := db.PostgresqlProxy.GetCurrentDB().Query(query, userID)
	if err != nil {
		log.Error("Error fetching cart items: ", err)
		return nil, err
	}
	defer rows.Close()
	cartItems, err = MapRowsToCartItems(rows)
	if err != nil {
		return nil, err
	}

	return cartItems, nil
}

func AddProductToCart(userID string, productID int, quantity int) error {
	tx, err := db.PostgresqlProxy.GetCurrentDB().Begin()
	if err != nil {
		log.Error("Error starting transaction: ", err)
		return err
	}
	defer tx.Rollback()

	var cartID int
	// Попытка найти существующую корзину для пользователя
	cartQuery := `SELECT id FROM cart WHERE user_id = ?`
	err = tx.QueryRow(cartQuery, userID).Scan(&cartID)
	if err != nil {
		// Если корзина не найдена, создаем новую
		insertCartQuery := `INSERT INTO cart (user_id) VALUES (?) RETURNING id`
		err = tx.QueryRow(insertCartQuery, userID).Scan(&cartID)
		if err != nil {
			log.Error("Error creating a new cart: ", err)
			return err
		}
	}

	// Попытка добавить товар в корзину или обновить его количество, если он уже там есть
	updateQuery := `
        INSERT INTO cart_items (cart_id, product_id, quantity)
        VALUES (?, ?, ?)
        ON CONFLICT (cart_id, product_id)
        DO UPDATE SET quantity = cart_items.quantity + EXCLUDED.quantity
	`

	_, err = tx.Exec(updateQuery, cartID, productID, quantity)
	if err != nil {
		log.Error("Error adding/updating product in cart: ", err)
		return err
	}

	// Если дошли до сюда без ошибок, подтверждаем транзакцию
	if err = tx.Commit(); err != nil {
		log.Error("Error committing transaction: ", err)
		return err
	}

	return nil
}

func RemoveProductFromCart(userID string, productID int) error {
	tx, err := db.PostgresqlProxy.GetCurrentDB().Begin()
	if err != nil {
		log.Error("Error starting transaction: ", err)
		return err
	}
	defer tx.Rollback()

	var cartID int
	// Попытка найти существующую корзину для пользователя
	cartQuery := `SELECT id FROM cart WHERE user_id = ?`
	err = tx.QueryRow(cartQuery, userID).Scan(&cartID)
	if err != nil {
		log.Error("Error fetching cart: ", err)
		return err
	}

	// Уменьшаем количество товара в корзине на 1
	updateQuantity := `UPDATE cart_items SET quantity = quantity - 1 WHERE cart_id = ? AND product_id = ? RETURNING quantity`
	_, err = tx.Exec(updateQuantity, cartID, productID)
	if err != nil {
		log.Error("Error deleting product from cart: ", err)
		return err
	}

	// Удаляем товар из корзины, если его количество стало равно 0
	deleteQuery := `DELETE FROM cart_items WHERE quantity = 0`
	_, err = tx.Exec(deleteQuery, cartID, productID)
	if err != nil {
		log.Error("Error deleting product from cart: ", err)
		return err
	}

	// Если дошли до сюда без ошибок, подтверждаем транзакцию
	if err = tx.Commit(); err != nil {
		log.Error("Error committing transaction: ", err)
		return err
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

func PlaceOrder(userID string, deliveryAddress string) (int, int, error) {
	// Начало транзакции
	tx, err := db.PostgresqlProxy.GetCurrentDB().Begin()
	if err != nil {
		log.Error("Error starting transaction: ", err)
		return http.StatusInternalServerError, 0, fmt.Errorf("error starting transaction")
	}

	// Шаг 1: Создание заказа
	var orderID int
	orderQuery := `
        INSERT INTO orders (user_id, delivery_address, order_date)
        VALUES (?, ?, NOW())
        RETURNING id
    `
	err = tx.QueryRow(orderQuery, userID, deliveryAddress, userID).Scan(&orderID)
	if err != nil {
		tx.Rollback()
		log.Error("Error creating order: ", err)
		return http.StatusInternalServerError, orderID, fmt.Errorf("error creating order, please try again later")
	}

	// Проверка, что корзина не пуста
	emptyCartQuery := `SELECT COUNT(*) FROM cart_items WHERE cart_id IN (SELECT id FROM cart WHERE user_id = ?)`
	var cartItemCount int
	err = tx.QueryRow(emptyCartQuery, userID).Scan(&cartItemCount)
	if err != nil || cartItemCount == 0 {
		tx.Rollback()
		log.Error("Error checking cart items: ", err)
		return http.StatusBadRequest, orderID, fmt.Errorf("cart is empty")
	}

	// Шаг 2: Копирование содержимого корзины в заказ
	copyQuery := `
        INSERT INTO order_items (order_id, product_id, quantity, price_at_order)
        SELECT ?, ci.product_id, ci.quantity, p.price
        FROM cart_items ci
        JOIN products p ON ci.product_id = p.id
        JOIN cart c ON ci.cart_id = c.id
        WHERE c.user_id = ?
    `
	_, err = tx.Exec(copyQuery, orderID, userID)
	if err != nil {
		tx.Rollback()
		log.Error("Error copying cart items to order items: ", err)
		return http.StatusInternalServerError, orderID, fmt.Errorf("error placing order, please try again later")
	}

	// Шаг 3: Очистка корзины
	clearCartQuery := `
        DELETE FROM cart_items
        WHERE cart_id IN (
            SELECT id FROM cart WHERE user_id = ?
        )
    `
	_, err = tx.Exec(clearCartQuery, userID)
	if err != nil {
		tx.Rollback()
		log.Error("Error clearing cart: ", err)
		return http.StatusInternalServerError, orderID, fmt.Errorf("error placing order, please try again later")
	}

	// Завершение транзакции
	err = tx.Commit()
	if err != nil {
		log.Error("Error committing transaction: ", err)
		return http.StatusInternalServerError, orderID, fmt.Errorf("error placing order, please try again later")
	}

	return http.StatusOK, orderID, nil
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
