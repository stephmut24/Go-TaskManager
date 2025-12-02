package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// LoadEnv attempts to load a .env file. It first tries the working directory
// and if not found it walks upward through parent directories until it finds a
// .env file or reaches the filesystem root.
func LoadEnv() {
	// try cwd first
	if err := godotenv.Load(); err == nil {
		return
	}

	// if not found in cwd, walk up parents looking for .env
	dir, err := os.Getwd()
	if err != nil {
		log.Println("LoadEnv: unable to get working directory:", err)
		return
	}

	for {
		candidate := filepath.Join(dir, ".env")
		if _, statErr := os.Stat(candidate); statErr == nil {
			if err := godotenv.Load(candidate); err == nil {
				return
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir { // reached root
			break
		}
		dir = parent
	}

	log.Println(".env not found in working directory or any parent")
}
func GetEnv(key string) string {
	return os.Getenv(key)
}
