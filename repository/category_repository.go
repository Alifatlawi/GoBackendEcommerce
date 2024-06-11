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
	defer rows.Close()

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
	query := `
		INSERT INTO categories (name)
		OUTPUT INSERTED.id
		VALUES (@Name)
	`
	var id int64
	err := db.DB.QueryRow(query, sql.Named("Name", category.Name)).Scan(&id)
	if err != nil {
		log.Println("Failed to create category:", err)
		return 0, err
	}
	return id, nil
}

func UpdateCategory(category models.Category) error {
	_, err := db.DB.Exec("UPDATE categories SET name = @Name WHERE id = @ID", sql.Named("Name", category.Name), sql.Named("ID", category.ID))
	if err != nil {
		log.Println("Failed to update category:", err)
		return err
	}
	return nil
}

func DeleteCategory(id int) error {
	_, err := db.DB.Exec("DELETE FROM categories WHERE id = @ID", sql.Named("ID", id))
	if err != nil {
		log.Println("Failed to delete category:", err)
		return err
	}
	return nil
}
