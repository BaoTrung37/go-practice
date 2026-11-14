package main

import (
	"fmt"
	"net/http"
	"test/internal/handler"
	"test/internal/middleware"
	"test/internal/repository"
	"test/internal/service"
)

func main() {
	repo := repository.NewMemoryWalletRepository()

	svc := service.NewWalletService(repo)

	h := handler.NewWalletHandler(svc)

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	wrappedMux := middleware.Recovery(middleware.Logging(mux))

	fmt.Println("🚀 Server running on http://localhost:8080")
	if err := http.ListenAndServe("localhost:8080", wrappedMux); err != nil {
		fmt.Printf("Sever error %v\n", err)
	}
}
