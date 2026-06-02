package main

import (
	"log"

	"mcp-kit/internal/github"
	"mcp-kit/internal/mcpkit"
)

func main() {
	app := mcpkit.App{
		Name:     "mcp-github",
		Version:  "0.1.0",
		Register: github.Register,
	}
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
