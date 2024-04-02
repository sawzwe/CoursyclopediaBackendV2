package main

import (
	"BackendCoursyclopedia/route"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	if _, exists := os.LookupEnv("RAILWAY_ENVIRONMENT"); exists == false {
		if err := godotenv.Load(); err != nil {
			log.Fatal("error loading .env file:", err)

		}
	}
	app := fiber.New()

	route.Setup(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// log.Fatal(app.Listen("0.0.0.0" + port))
	log.Fatal(app.Listen(":" + port))
}
