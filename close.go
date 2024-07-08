package main

import (
	"log"
	"string_backend_0001/internal/database"
	"string_backend_0001/internal/logger"
)

func Close() {
	err := database.Close()
	if err != nil {
		log.Println(err)
		return
	}

	err = logger.Close()
	if err != nil {
		log.Println(err)
		return
	}
}
