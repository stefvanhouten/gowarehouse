package products

import (
	"database/sql"
	"gopkg.in/go-playground/validator.v9"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name" validate:"required_with_all"`
	Price float64 `json:"price" validate:"required_with_all,min=0"`
}

// Retrieves a single product by ID.
func ProductByID(db *sql.DB, id int64) ([]Product, error) {
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

// Retrieves all products.
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

func CreateProduct(db *sql.DB, p *Product) (int64, error) {
	v := validator.New()

	if err := v.Struct(p); err != nil {
		return 0, err
	}

	// Insert the product into the database.
	result, err := db.Exec(
		"INSERT INTO products (name, price) VALUES (?, ?)",
		p.Name,
		p.Price,
	)

	if err != nil {
		return 0, err
	}

	recordID, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return recordID, nil
}
