package config

import "github.com/dfodeker/clover/database"

type Config struct {
	DB           database.Client
	JwtSecret    string
	Platform     string
	FilepathRoot string
	AssetsRoot   string
	Port         string
}
