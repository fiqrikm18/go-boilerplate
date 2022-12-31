package main

import (
	"github.com/fiqrikm18/go-boilerplate/internal/config"
)

func main() {
	_, err := config.NewDbConnection()
	if err != nil {
		panic(err)
	}
}
