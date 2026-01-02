package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/kudo07/webScraper/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	fmt.Println("hello world")
	godotenv.Load(".env")
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalln("DB_URL environment variable not set")
	}
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalln("can't connect to database", err)
	}
	apiCnfg := apiConfig{
		DB: database.New(conn),
	}
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	// when something went wrong this endpoint return the consistent return with error
	v1Router.Get("/err", handleErr)
	v1Router.Post("/users", apiCnfg.hadnlerCreateUser)
	v1Router.Get("/users", apiCnfg.middlewareAuth(apiCnfg.handlerGetUser))
	v1Router.Post("/feeds", apiCnfg.middlewareAuth(apiCnfg.handlerCreateFeed))

	router.Mount("/v1", v1Router)
	serv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Printf("Server starting on port %v", portString)
	err = serv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("PORT:", portString)
}
