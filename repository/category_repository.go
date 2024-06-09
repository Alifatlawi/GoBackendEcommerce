package repository

import (
	"database/sql"
	"ecommercebackend/db"
	"ecommercebackend/models"
	"log"
)

func GetAllCategories() ([]models.Category, error) {
	rows, err := db.DB.Query("SELECT * FROM categories")
	if err != nil {
		log.Println("Failed to query categories:", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("Failed to close rows:", err)
		}
	}(rows)

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		if err := rows.Scan(&category.ID, &category.Name); err != nil {
			log.Println("Failed to scan category:", err)
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

func CreateCategory(category models.Category) (int64, error) {
	result, err := db.DB.Exec("INSERT INTO categories (name) VALUES (?)", category.Name)
	if err != nil {
		log.Println("Failed to create category:", err)
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Println("Failed to get last insert ID:", err)
		return 0, err
	}
	return id, nil
}
