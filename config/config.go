package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)


var Key string

func SetupEnvVar() error {
	mode := strings.ToUpper(os.Getenv("APP_ENV"))
	if mode == "" || mode == "DEV" {
		log.Print("Dev mode")
		projectRoot := FindProjectRoot()
		if projectRoot == "" {
			return  os.ErrNotExist
		}
		envPath := filepath.Join(projectRoot, "config", ".env")
		err := godotenv.Load(envPath)
		
		if err != nil {
			return  err
		}
		
	} else {
		log.Print("prod mode")
	}
	Key = os.Getenv("SECRETWORDJWT")
	return  nil

}

func FindProjectRoot() string {
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return  dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return  ""
}