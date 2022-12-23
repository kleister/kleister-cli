package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kleister/kleister-cli/pkg/command"
)

func main() {
	if env := os.Getenv("KLEISTER_CLI_ENV_FILE"); env != "" {
		godotenv.Load(env)
	}

	if err := command.Run(); err != nil {
		os.Exit(1)
	}
}
