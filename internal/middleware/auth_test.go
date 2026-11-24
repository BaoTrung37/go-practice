package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMiddleware(t *testing.T) {
	// Handler giả lập phía sau middleware để kiểm tra xem UserID có được gắn vào context không
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := GetUserID(r.Context())
		if !ok || userID == "" {
			t.Errorf("expected userID in context, got empty")
		}
		w.WriteHeader(http.StatusOK)
	})

	// Bọc handler giả lập bằng Auth middleware
	authMiddleware := Auth(nextHandler)

	// Kịch bản 1: Không có header Authorization -> Kỳ vọng 401 Unauthorized
	t.Run("Missing Authorization Header", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/protected", nil)
		rec := httptest.NewRecorder()

		authMiddleware.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected 401 Unauthorized, got %d", rec.Code)
		}
	})

	// Kịch bản 2: Token sai -> Kỳ vọng 401 Unauthorized
	t.Run("Invalid Token", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/protected", nil)
		req.Header.Set("Authorization", "Bearer token-fake-123")
		rec := httptest.NewRecorder()

		authMiddleware.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("expected 401 Unauthorized, got %d", rec.Code)
		}
	})

	// Kịch bản 3: Token hợp lệ -> Kỳ vọng 200 OK và lấy được UserID trong context
	t.Run("Valid Token", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/protected", nil)
		req.Header.Set("Authorization", "Bearer token-trung")
		rec := httptest.NewRecorder()

		authMiddleware.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected 200 OK, got %d", rec.Code)
		}
	})
}
