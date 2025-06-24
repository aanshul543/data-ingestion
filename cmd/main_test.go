package main_test

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	_ = godotenv.Load() // loads .env file into os.Environ
	os.Exit(m.Run())
}
