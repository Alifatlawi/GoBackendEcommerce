package repository

import (
	"database/sql"
	"ecommercebackend/db"
	"ecommercebackend/models"
	"log"
)

// GetAllProducts retrieves all products from the database
func GetAllProducts() ([]models.Product, error) {
	rows, err := db.DB.Query("SELECT * FROM products")
	if err != nil {
		log.Println("Failed to query products:", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("Failed to close rows:", err)
		}
	}(rows)

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.Id, &product.Name, &product.Description, &product.ImgUrl, &product.Price, &product.CategoryId); err != nil {
			log.Println("Failed to scan product:", err)
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

// CreateProduct inserts a new product into the database
func CreateProduct(newProduct models.Product) (int64, error) {
	query := `
		INSERT INTO products (name, description, img_url, price, category_id)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := db.DB.Exec(query, newProduct.Name, newProduct.Description, newProduct.ImgUrl, newProduct.Price, newProduct.CategoryId)
	if err != nil {
		log.Println("Failed to create product:", err)
		return 0, err
	}
	return result.LastInsertId()
}

// GetProductById retrieves a product by its ID from the database
func GetProductById(id int) (models.Product, error) {
	var product models.Product
	err := db.DB.QueryRow("SELECT * FROM products WHERE id = ?", id).Scan(&product.Id, &product.Name, &product.Description, &product.ImgUrl, &product.Price, &product.CategoryId)
	if err != nil {
		if err == sql.ErrNoRows {
			return product, nil
		}
		log.Println("Failed to query product:", err)
		return product, err
	}
	return product, nil
}

// UpdateProduct updates an existing product in the database
func UpdateProduct(product models.Product) error {
	_, err := db.DB.Exec("UPDATE products SET name = ?, description = ?, img_url = ?, price = ?, category_id = ? WHERE id = ?", product.Name, product.Description, product.ImgUrl, product.Price, product.CategoryId, product.Id)
	if err != nil {
		log.Println("Failed to update product:", err)
		return err
	}
	return nil
}

// DeleteProduct deletes a product from the database by its ID
func DeleteProduct(id int) error {
	_, err := db.DB.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		log.Println("Failed to delete product:", err)
		return err
	}
	return nil
}

// GetProductsByCategoryID retrieves products by category ID from the database
func GetProductsByCategoryID(categoryID int) ([]models.Product, error) {
	rows, err := db.DB.Query("SELECT * FROM products WHERE category_id = ?", categoryID)
	if err != nil {
		log.Println("Failed to query products:", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("Failed to close rows:", err)
		}
	}(rows)

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.Id, &product.Name, &product.Description, &product.ImgUrl, &product.Price, &product.CategoryId); err != nil {
			log.Println("Failed to scan product:", err)
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}
