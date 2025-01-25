package main

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	app := fiber.New()

	app.Post("/create", createFileHandler)

	log.Info().Msg("Server is starting on port 8080...")
	if err := app.Listen(":8080"); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}

func createFileHandler(c *fiber.Ctx) error {
	fileName := uuid.New().String() + ".txt"

	file, err := os.Create("./fs/" + fileName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create file")
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	err = file.Close()
	if err != nil {
		log.Error().Err(err).Msg("Failed to close file")
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}

	log.Info().Str("file", fileName).Msg("file created successfully")

	return c.SendString("OK")
}
