package repository

import (
	"database/sql"
	"ecommercebackend/db"
	"ecommercebackend/models"
	"errors"
	"fmt"
	"log"
)

func GetAllOrders() ([]models.Order, error) {
	rows, err := db.DB.Query("SELECT id, product_id, address, phone_number FROM orders")
	if err != nil {
		log.Println("Failed to query orders:", err)
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.ProductID, &order.Address, &order.PhoneNumber); err != nil {
			log.Println("Failed to scan order:", err)
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		log.Println("Rows error:", err)
		return nil, err
	}

	return orders, nil
}

// CreateOrder creates a new order
func CreateOrder(order models.Order) (int, error) {
	// Check if the product exists
	var productID int
	err := db.DB.QueryRow("SELECT id FROM products WHERE id = @ID", sql.Named("ID", order.ProductID)).Scan(&productID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("Product does not exist:", err)
			return 0, fmt.Errorf("product with id %d does not exist", order.ProductID)
		}
		log.Println("Failed to query product:", err)
		return 0, err
	}

	query := `
	INSERT INTO orders (product_id, address, phone_number)
	OUTPUT INSERTED.id
	VALUES (@ProductID, @Address, @PhoneNumber)
	`
	var id int
	err = db.DB.QueryRow(query, sql.Named("ProductID", order.ProductID), sql.Named("Address", order.Address), sql.Named("PhoneNumber", order.PhoneNumber)).Scan(&id)
	if err != nil {
		log.Println("Failed to create order:", err)
		return 0, err
	}
	return id, nil
}

// GetOrderById gets an order by ID
func GetOrderById(id int) (models.Order, error) {
	var order models.Order
	query := "SELECT id, product_id, address, phone_number FROM orders WHERE id = @ID"
	err := db.DB.QueryRow(query, sql.Named("ID", id)).Scan(&order.ID, &order.ProductID, &order.Address, &order.PhoneNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return order, nil
		}
		log.Println("Failed to query order:", err)
		return order, err
	}
	return order, nil
}

// UpdateOrder updates an existing order
func UpdateOrder(order models.Order) error {
	_, err := db.DB.Exec("UPDATE orders SET product_id = @ProductID, address = @Address, phone_number = @PhoneNumber WHERE id = @ID",
		sql.Named("ProductID", order.ProductID),
		sql.Named("Address", order.Address),
		sql.Named("PhoneNumber", order.PhoneNumber),
		sql.Named("ID", order.ID))
	if err != nil {
		log.Println("Failed to update order:", err)
		return err
	}
	return nil
}

// DeleteOrder deletes an order
func DeleteOrder(id int) error {
	_, err := db.DB.Exec("DELETE FROM orders WHERE id = @ID", sql.Named("ID", id))
	if err != nil {
		log.Println("Failed to delete order:", err)
		return err
	}
	return nil
}
