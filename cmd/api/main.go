package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"test/internal/handler"
	"test/internal/middleware"
	"test/internal/repository"
	"test/internal/service"
	"test/internal/worker"
)

func main() {
	// 1. Khởi tạo Worker Pool (3 workers, queue 100) và bật chạy
	pool := worker.NewWorkerPool(3, 100)
	pool.Start()

	// 2. Khởi tạo Repo, Service, Handler
	repo := repository.NewMemoryWalletRepository()
	svc := service.NewWalletService(repo)
	h := handler.NewWalletHandler(svc)

	// 3. Đăng ký routes và bọc Middleware
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)
	wrappedMux := middleware.Recovery(middleware.Logging(mux))

	// 4. Cấu hình HTTP Server
	server := &http.Server{
		Addr:    ":8080",
		Handler: wrappedMux,
	}

	// 5. Chạy Server trong một GOROUTINE RIÊNG BIỆT để không chặn main()
	go func() {
		fmt.Println("🚀 Server running on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	// 6. LẮNG NGHE TÍN HIỆU TẮT MÁY (Ctrl+C hoặc SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit // Dừng ở đây chờ tín hiệu...

	fmt.Println("\n🛑 Shutting down server gracefully...")

	// Cho server 5 giây để dọn dẹp các request đang dang dở
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("server forced to shutdown", "error", err)
	}

	// Dừng Worker Pool (chờ các tác vụ ngầm hoàn tất)
	pool.Stop()

	fmt.Println("✅ Server exited cleanly.")
}
