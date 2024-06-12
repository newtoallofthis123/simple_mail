package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DatabaseUrl  string
	Port         string
	Mail         string
	MailPassword string
	AuthSecret   string
}

// Reads the .env file and returns an Env struct.
func ReadEnv() Env {
	godotenv.Load(".env")

	return Env{
		DatabaseUrl:  getEnv("DATABASE_URL"),
		Port:         getEnv("PORT"),
		Mail:         getEnv("MAIL"),
		MailPassword: getEnv("MAIL_PASSWORD"),
		AuthSecret:   getEnv("AUTH_SECRET"),
	}
}

// Returns the value of the given env var name.
func getEnv(name string) string {
	val, ok := os.LookupEnv(name)
	if !ok {
		panic(fmt.Sprintf("Env var %s not found", name))
	}
	return val
}
