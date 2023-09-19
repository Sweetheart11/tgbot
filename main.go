package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("error loading .env file", err)
	}

	host := os.Getenv("POSTGRESDB_HOST")
	port := os.Getenv("POSTGRESDB_PORT")
	password := os.Getenv("POSTGRESDB_PASSWORD")
	dbname := os.Getenv("POSTGRESDB_NAME")
	user := os.Getenv("POSTGRESDB_USER")
	//
	// connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
	// 	user, password, host, port, dbname)
	// connStr := "hostuser=user dbname=mydb password=pass sslmode=disable"
	connStr := fmt.Sprintf(
		"host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		host,
		port,
		user,
		password,
		dbname,
	)
	fmt.Println(connStr)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func main() {
	store, err := NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", store)
}
