package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() *sql.DB {
	var err error
	db, err = sql.Open("sqlite3", "./products.db")
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS products (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,
        description TEXT,
		count INTEGER,
        price REAL
    );`
	if _, err := db.Exec(createTableSQL); err != nil {
		log.Fatal(err)
	}
	return db
}

// func GetAllProducts() (string, error) {
// 	rows, err := db.Query("SELECT name, description, price FROM products")
// 	if err != nil {
// 		return "", err
// 	}
// 	defer rows.Close()

// 	var result strings.Builder
// 	for rows.Next() {
// 		var name, description string
// 		var price float64
// 		var count int
// 		if err := rows.Scan(&name, &description, &price); err != nil {
// 			return "", err
// 		}
// 		result.WriteString(fmt.Sprintf("Название: %s, Описание: %s, Кол-во: %d, Стоимость: %.2f\n", name, description, count, price))
// 	}
// 	if result.Len() == 0 {
// 		return "Нет позиций в базе данных.", nil
// 	}
// 	return result.String(), nil
// }
