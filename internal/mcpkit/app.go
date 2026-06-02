package mcpkit

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type App struct {
	Name     string
	Version  string
	Register func(*mcp.Server)
}

func (a App) Run() error {
	log.SetFlags(0)
	log.SetPrefix("[" + a.Name + "] ")

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "add-to-claude":
			return a.registerInClaude()
		case "remove-from-claude":
			return a.unregisterFromClaude()
		case "-h", "--help", "help":
			a.printUsage()
			return nil
		default:
			a.printUsage()
			return fmt.Errorf("unknown command: %s", os.Args[1])
		}
	}
	return a.serve()
}

func (a App) serve() error {
	envPaths := loadEnv(a.Name)

	log.Printf("Starting %s v%s", a.Name, a.Version)
	log.Println("Transport: stdio")
	for _, p := range envPaths {
		if _, err := os.Stat(p); err == nil {
			log.Printf("Env file found: %s", p)
		} else {
			log.Printf("Env file not found: %s", p)
		}
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    a.Name,
		Version: a.Version,
	}, nil)

	a.Register(server)

	return server.Run(context.Background(), &mcp.StdioTransport{})
}

func (a App) printUsage() {
	fmt.Fprintf(os.Stderr, `Usage: %s [command]

Commands:
  (none)    Run as MCP server over stdio (default).
  add-to-claude       Register this server in Claude Code for the current project.
  remove-from-claude  Unregister this server from Claude Code.
  help      Show this help.
`, a.Name)
}

// loadEnv loads .env files in priority order. godotenv.Load never overwrites
// already-set variables, so the first file found for a given key wins:
//  1. $XDG_CONFIG_HOME/mcp-kit/<appName>.env  (canonical; created by `make install`)
//  2. ./.env                                    (CWD fallback for `go run ./cmd/<app>`)
func loadEnv(appName string) []string {
	var paths []string
	if cfgDir, err := os.UserConfigDir(); err == nil {
		p := filepath.Join(cfgDir, "mcp-kit", appName+".env")
		paths = append(paths, p)
		_ = godotenv.Load(p)
	}
	cwdEnv := ".env"
	if abs, err := filepath.Abs(cwdEnv); err == nil {
		cwdEnv = abs
	}
	paths = append(paths, cwdEnv)
	_ = godotenv.Load()
	return paths
}
