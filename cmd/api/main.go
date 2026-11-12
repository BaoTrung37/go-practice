package main

import (
	"fmt"
	"test/internal/handler"
	"test/internal/repository"
	"test/internal/service"
)

func main() {
	fmt.Println("Starting Payment API..")

	repo := repository.NewMemoryWalletRepository()

	svc := service.NewWalletService(repo)

	h := handler.NewWalletHandler(svc)
}
