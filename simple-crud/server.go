package main

import (
	"example/simple-crud/pkg/interfaces/runner"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	runner.Run()
}
