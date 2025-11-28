/*
Copyright © 2025 NAME HERE dfodeker
*/
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/dfodeker/clover/database"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	db           database.Client
	jwtSecret    string
	platform     string
	filepathRoot string
	assetsRoot   string
	port         string
}

func main() {
	godotenv.Load(".env")

	pathToDB := os.Getenv("DB_PATH")
	if pathToDB == "" {
		log.Fatal("DB_URL must be set")
	}
	db, err := database.NewClient(pathToDB)
	if err != nil {
		log.Fatalf("Couldn't connect to database: %v", err)
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM environment variable is not set")
	}
	assetsRoot := os.Getenv("ASSETS_ROOT")
	if assetsRoot == "" {
		log.Fatal("ASSETS_ROOT environment variable is not set")
	}
	filepathRoot := os.Getenv("FILEPATH_ROOT")
	if assetsRoot == "" {
		log.Fatal("ASSETS_ROOT environment variable is not set")
	}
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	cfg := apiConfig{
		db:           db,
		jwtSecret:    jwtSecret,
		platform:     platform,
		assetsRoot:   assetsRoot,
		filepathRoot: filepathRoot,
		port:         port,
	}

	mux := http.NewServeMux()
	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))
	mux.Handle("/app/", appHandler)

	mux.HandleFunc("POST /api/login", cfg.handlerLogin)
	//cmd.Execute(cfg)
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving on: http://localhost:%s/app/\n", port)
	log.Fatal(srv.ListenAndServe())
}
