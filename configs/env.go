package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Username() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error Loading .env File!!")
	}

	return os.Getenv("USER")
}

func Password() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error Loading .env File!!")
	}

	return os.Getenv("PASSWORD")
}

func Host() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error Loading .env File!!")
	}

	return os.Getenv(("HOST"))
}

func Port() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error Loading .env File!!")
	}

	return os.Getenv(("PORT"))
}

func DbName() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error Loading .env File!!")
	}

	return os.Getenv(("DBNAME"))
}
