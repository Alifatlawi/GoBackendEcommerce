package db

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/microsoft/go-mssqldb"
	"log"
	"os"
	"strconv"
)

var DB *sql.DB

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	server := os.Getenv("AZURE_SQL_SERVER_NAME")
	user := os.Getenv("AZURE_SQL_DATABASE_USER")
	password := os.Getenv("AZURE_SQL_DATABASE_PASSWORD")
	port, err := strconv.Atoi(os.Getenv("AZURE_SQL_SERVER_PORT"))
	if err != nil {
		log.Fatal("Error converting port to integer:", err)
	}
	database := os.Getenv("AZURE_SQL_DATABASE_NAME")

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	DB, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Could not connect to database:", err)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	createTables()
}

func createTables() {
	// Check if categories table exists and create if it does not
	createCategoriesTable := `
    IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='categories' AND xtype='U')
    CREATE TABLE categories (
        id INT PRIMARY KEY IDENTITY(1,1),
        name NVARCHAR(100) NOT NULL UNIQUE
    );`

	// Check if products table exists and create if it does not
	createProductsTable := `
    IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='products' AND xtype='U')
    CREATE TABLE products (
        id INT PRIMARY KEY IDENTITY(1,1),
        name NVARCHAR(100) NOT NULL,
        description NVARCHAR(255) NOT NULL,
        img_url NVARCHAR(255) NOT NULL,
        price FLOAT NOT NULL,
        category_id INT,
        FOREIGN KEY (category_id) REFERENCES categories(id)
    );`

	_, err := DB.Exec(createCategoriesTable)
	if err != nil {
		log.Fatal("Could not create categories table:", err)
	}

	_, err = DB.Exec(createProductsTable)
	if err != nil {
		log.Fatal("Could not create products table:", err)
	}

	log.Println("Database tables created")
}
