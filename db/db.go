package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/microsoft/go-mssqldb"
)

var DB *sql.DB

func InitDB() {
	server := os.Getenv("AZURE_SQL_SERVER_NAME")
	user := os.Getenv("AZURE_SQL_DATABASE_USER")
	password := os.Getenv("AZURE_SQL_DATABASE_PASSWORD")
	port, err := strconv.Atoi(os.Getenv("AZURE_SQL_SERVER_PORT"))
	if err != nil {
		log.Fatalf("Error converting port to integer: %v", err)
	}
	database := os.Getenv("AZURE_SQL_DATABASE_NAME")

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	DB, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	createTables()
}

func createTables() {
	createCategoriesTable := `
    IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='categories' AND xtype='U')
    CREATE TABLE categories (
        id INT PRIMARY KEY IDENTITY(1,1),
        name NVARCHAR(100) NOT NULL UNIQUE
    );`

	createProductsTable := `
	IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='products' AND xtype='U')
	CREATE TABLE products (
		id INT PRIMARY KEY IDENTITY(1,1),
		name NVARCHAR(100) NOT NULL,
		description NVARCHAR(255) NOT NULL,
		img_url NVARCHAR(255) NOT NULL,
		price NVARCHAR(255) NOT NULL,
		category_id INT,
		FOREIGN KEY (category_id) REFERENCES categories(id)
	);`

	createOrdersTable := `
	IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='orders' AND xtype='U')
	CREATE TABLE orders (
		id INT PRIMARY KEY IDENTITY(1,1),
		product_id INT,
		address NVARCHAR(255) NOT NULL,
		phone_number NVARCHAR(20) NOT NULL,
		FOREIGN KEY (product_id) REFERENCES products(id)
	);`

	_, err := DB.Exec(createCategoriesTable)
	if err != nil {
		log.Fatalf("Could not create categories table: %v", err)
	}

	_, err = DB.Exec(createProductsTable)
	if err != nil {
		log.Fatalf("Could not create products table: %v", err)
	}

	_, err = DB.Exec(createOrdersTable)
	if err != nil {
		log.Fatalf("Could not create orders table: %v", err)
	}

	log.Println("Database tables created")
}
