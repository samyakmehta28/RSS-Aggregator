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
	_ "github.com/lib/pq"
	"github.com/samyakmehta28/RSS-Aggregator/internal/database"
)

type apiConfig struct {
	DB *database.Queries
}

func main(){
	fmt.Println("hello world")
	errEnv := godotenv.Load(".env")
	if errEnv != nil {
		log.Fatal("Error loading .env file")
	}
	port:= os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable not set")
	}

	dbURL:= os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable not set")
	}

	conn, errDB := sql.Open("postgres", dbURL)
	if errDB != nil {
		log.Fatal(errDB)
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}
	// fmt.Println("Server is running on port:", port)


	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/health", healthCheck)
	v1Router.Post("/user", apiCfg.createUserHandler)
	v1Router.Get("/user", apiCfg.middlewareAuth(apiCfg.getUserByAPIKey))
	v1Router.Get("/feeds", apiCfg.getFeedsHandler)
	v1Router.Post("/feed", apiCfg.middlewareAuth(apiCfg.createFeedHandler))
	v1Router.NotFound(pathNotFoundError)


	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	fmt.Printf("Starting server on port %s...\n", port)

	errSrv := srv.ListenAndServe()
	if errSrv != nil {
		log.Fatal( errSrv)
	}


}