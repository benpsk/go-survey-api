package config

import "os"

var (
	PORT         = os.Getenv("PORT")
	DATABASE_URL = os.Getenv("DATABASE_URL")
	JWT_SECRET   = os.Getenv("JWT_SECRET")
	ENV          = os.Getenv("ENV")
)
