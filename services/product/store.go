package product

import (
	"database/sql"

	"github.com/akshtrikha/golang-ecomm/types"
)

// Store struct to hold the database object
// This will be used to handle the database queries
type Store struct {
	db *sql.DB
}

// NewStore function to return a reference to the Store struct
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

// AddProduct function to add the product to db
// /add-product api endpoint
func (s *Store) AddProduct(product types.AddProductPayload) (int, error) {
	// run command to insert product in the products table
	result, err := s.db.Exec("INSERT INTO products (name, description, image, price, quantity) VALUES (?, ?, ?, ?, ?)", product.Name, product.Description, product.Image, product.Price, product.Quantity)
	if err != nil {
		return 0, err
	}

	// get the lastInsertedId
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// return the id of the product created
	return int(id), nil
}

// GetProducts func to get all the products
// response to /get-products api endpoint
func (s *Store) GetProducts() ([]types.Product, error) {
	// run the query to get all the rows from products table
	rows, err := s.db.Query(`SELECT * FROM products`)
	if err != nil {
		return nil, err
	}

	var products []types.Product

	// scan the rows
	for rows.Next() {
		product, err := scanRowIntoProduct(rows)

		if err != nil {
			return nil, err
		}

		products = append(products, *product)
	}

	// return the result
	return products, err
}

func scanRowIntoProduct(rows *sql.Rows) (*types.Product, error) {
	product := new(types.Product)

	err := rows.Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Image,
		&product.Price,
		&product.Quantity,
		&product.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return product, nil
}
