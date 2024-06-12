package repository

import (
	"database/sql"
	"ecommercebackend/db"
	"ecommercebackend/models"
	"log"
)

func GetAllProducts() ([]models.Product, error) {
	rows, err := db.DB.Query("SELECT * FROM products")
	if err != nil {
		log.Println("Failed to query products:", err)
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.ImgUrl, &product.Price, &product.CategoryId); err != nil {
			log.Println("Failed to scan product:", err)
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func CreateProduct(newProduct models.Product) (int64, error) {
	query := `
		INSERT INTO products (name, description, img_url, price, category_id)
		OUTPUT INSERTED.id
		VALUES (@Name, @Description, @ImgUrl, @Price, @CategoryID)
	`
	var id int64
	err := db.DB.QueryRow(query,
		sql.Named("Name", newProduct.Name),
		sql.Named("Description", newProduct.Description),
		sql.Named("ImgUrl", newProduct.ImgUrl),
		sql.Named("Price", newProduct.Price),
		sql.Named("CategoryID", newProduct.CategoryId)).Scan(&id)
	if err != nil {
		log.Println("Failed to create product:", err)
		return 0, err
	}
	return id, nil
}

func GetProductById(id int) (models.Product, error) {
	var product models.Product
	query := "SELECT * FROM products WHERE id = @ID"
	err := db.DB.QueryRow(query, sql.Named("ID", id)).Scan(&product.ID, &product.Name, &product.Description, &product.ImgUrl, &product.Price, &product.CategoryId)
	if err != nil {
		if err == sql.ErrNoRows {
			return product, nil
		}
		log.Println("Failed to query product:", err)
		return product, err
	}
	return product, nil
}

func UpdateProduct(product models.Product) error {
	query := `
		UPDATE products 
		SET name = @Name, description = @Description, img_url = @ImgUrl, price = @Price, category_id = @CategoryID 
		WHERE id = @ID
	`

	_, err := db.DB.Exec(query,
		sql.Named("Name", product.Name),
		sql.Named("Description", product.Description),
		sql.Named("ImgUrl", product.ImgUrl),
		sql.Named("Price", product.Price),
		sql.Named("CategoryID", product.CategoryId),
		sql.Named("ID", product.ID))
	if err != nil {
		log.Println("Failed to update product:", err)
		return err
	}
	return nil
}

func DeleteProduct(id int) error {
	query := "DELETE FROM products WHERE id = @ID"
	_, err := db.DB.Exec(query, sql.Named("ID", id))
	if err != nil {
		log.Println("Failed to delete product:", err)
		return err
	}
	return nil
}

func GetProductsByCategoryID(categoryID int) ([]models.Product, error) {
	query := "SELECT * FROM products WHERE category_id = @CategoryID"
	rows, err := db.DB.Query(query, sql.Named("CategoryID", categoryID))
	if err != nil {
		log.Println("Failed to query products:", err)
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.ImgUrl, &product.Price, &product.CategoryId); err != nil {
			log.Println("Failed to scan product:", err)
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}
