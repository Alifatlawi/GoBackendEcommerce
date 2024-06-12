package repository

import (
	"ecommercebackend/db"
	"log"
)

// DeleteAllData deletes all entries in the database
func DeleteAllData() error {
	// Define the queries to delete data from tables
	queries := []string{
		"DELETE FROM products",
		"DELETE FROM categories",
	}

	// Execute the queries to delete data
	for _, query := range queries {
		_, err := db.DB.Exec(query)
		if err != nil {
			log.Println("Failed to delete data:", err)
			return err
		}
	}

	// Reset identity columns
	resetQueries := []string{
		"DBCC CHECKIDENT ('products', RESEED, 0)",
		"DBCC CHECKIDENT ('categories', RESEED, 0)",
	}

	// Execute the queries to reset the identity columns
	for _, query := range resetQueries {
		_, err := db.DB.Exec(query)
		if err != nil {
			log.Println("Failed to reset identity:", err)
			return err
		}
	}

	return nil
}
