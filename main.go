package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/FrancoMusolino/rssagg/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Missing PORT variable")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("Missing dbURL variable")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cannot connect to DB")
	}

	apiCfg := apiConfig{
		DB: database.New(conn),
	}

	go startScraping(apiCfg.DB, 10, time.Minute)

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
	v1Router.Get("/healthz", handlerReadiness)

	// Users
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))

	// Feeds
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Post("/feeds/{feedId}/follow", apiCfg.middlewareAuth(apiCfg.handlerFollowFeed))
	v1Router.Delete("/feeds/{feedId}/unfollow", apiCfg.middlewareAuth(apiCfg.handlerUnfollowFeed))

	v1Router.Get("/feed-follows", apiCfg.middlewareAuth(apiCfg.handlerGetUserFollowedFeels))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	fmt.Printf("Server starting on port %v", port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(router)
}
