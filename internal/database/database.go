package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	Health() map[string]string
}

type service struct {
	db *sql.DB
}

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

func New() Service {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	s := &service{db: db}
	return s
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.PingContext(ctx)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

type product struct {
  ID string
  Name string
  Price string
}

func (p *product) getProduct(s *service) error {
  q :="SELECT name, price FROM products WHERE id=$1" 
	return s.db.QueryRow(q, p.ID).Scan(&p.Name, &p.Price)
}

func (p *product) updateProduct(s *service) error {
  q := "UPDATE products SET name=$1, price=$2 WHERE id=$3"
	_, err := s.db.Exec(q, p.Name, p.Price, p.ID)
	return err
}

func (p *product) deleteProduct(s *service) error {
	_, err := s.db.Exec("DELETE FROM products WHERE id=$1", p.ID)
	return err
}

func (p *product) createProduct(s *service) error {
	err := s.db.QueryRow("INSERT INTO products(name, price) VALUES($1, $2) RETURNING id", p.Name, p.Price).Scan(&p.ID)
	if err != nil {
		return err
	}

	return nil
}

func getProducts(s *service, start, count int) ([]product, error) {
  q := "SELECT id, name,  price FROM products LIMIT $1 OFFSET $2"
	rows, err := s.db.Query(q, count, start)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	products := []product{}

	for rows.Next() {
		var p product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

