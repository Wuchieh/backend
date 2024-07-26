package bootstrap

import (
	"log"
	"string_backend_0001/internal/database"
	"string_backend_0001/internal/grpc"
	"string_backend_0001/internal/logger"
)

func Close() {
	logger.Info("Server is shutting down...")

	defer func() {
		logger.Info("Server stopped")

		if err := loggerClose(); err != nil {
			log.Println(err)
		}
	}()

	if err := databaseClose(); err != nil {
		logger.Error(err.Error())
	}

	if err := gRPCServerClose(); err != nil {
		logger.Error(err.Error())
	}
}

func databaseClose() error {
	return database.Close()
}

func loggerClose() error {
	return logger.Close()
}

func gRPCServerClose() error {
	return grpc.Close()
}
