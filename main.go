/*
Copyright © 2025 NAME HERE dfodeker
*/
package main

import (
	"log"
	"os"

	"github.com/dfodeker/clover/cmd"
	"github.com/dfodeker/clover/config"
	"github.com/dfodeker/clover/database"
	"github.com/joho/godotenv"
)

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
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	cfg := config.Config{
		DB:         db,
		JwtSecret:  jwtSecret,
		Platform:   platform,
		AssetsRoot: assetsRoot,
	}

	cmd.Execute(cfg)
}
