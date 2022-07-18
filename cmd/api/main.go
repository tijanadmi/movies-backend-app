package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"github/tijanadmi/movies-backend-app/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
	db struct {
		dsn string
	}
	jwt struct {
		secret string
	}
}

type AppStatus struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

type application struct {
	config config
	logger *log.Logger
	models models.Models
}


func main() {
	var cfg config
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
		return
	}

	user := os.Getenv("user")
	if user == "" {
		fmt.Println(err)
		return
	}

	password := os.Getenv("password")
	if password == "" {
		fmt.Println(err)
		return
	}

	host := os.Getenv("host")
	if host == "" {
		fmt.Println(err)
		return
	}
	port := os.Getenv("port")
	if host == "" {
		fmt.Println(err)
		return
	}
	dbname := os.Getenv("dbname")
	if host == "" {
		fmt.Println(err)
		return
	}
	jwt := os.Getenv("jwt-secret")
	if host == "" {
		fmt.Println(err)
		return
	}
	flag.IntVar(&cfg.port, "port", 4000, "Server port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "Application environment (development|production")
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)
	/*flag.StringVar(&cfg.db.dsn, "dsn", "postgres://postgres:postgres@localhost/postgres?sslmode=disable", "Postgres connection string")*/
	flag.StringVar(&cfg.db.dsn, "dsn", psqlInfo, "Postgres connection string")
	flag.StringVar(&cfg.jwt.secret, "jwt-secret", jwt, "secret")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	/*psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)*/

	/*db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		logger.Fatal(err)
		}
	defer db.Close()*/

	db, err := openDB(cfg)
	
	
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()
	
	app := &application{
		config: cfg,
		logger: logger,
		models: models.NewModels(db),
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Println("Starting server on port", cfg.port)

	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}

}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
