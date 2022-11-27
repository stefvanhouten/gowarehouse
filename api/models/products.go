package products

import (
	"database/sql"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func ProductByID(db *sql.DB, id int) ([]Product, error) {
	rows, err := db.Query("SELECT id, name, price FROM products WHERE id = ?", id)

	if err != nil {
		return nil, err
	}

	// Close the rows when we're done.
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}

		products = append(products, p)
	}
	return products, nil
}

func AllProducts(db *sql.DB) ([]Product, error) {
	rows, err := db.Query("SELECT id, name, price FROM products")

	if err != nil {
		return nil, err
	}

	// Close the rows when we're done.
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}

		products = append(products, p)
	}
	return products, nil
}
