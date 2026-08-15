package httpapi

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestBug3_SKUWithSlashRejected(t *testing.T) {
	clk := &testClock{t: time.Date(2026, 8, 14, 12, 0, 0, 0, time.UTC)}
	api := NewWithClock(clk.now)
	srv := httptest.NewServer(api.Handler())
	defer srv.Close()

	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/api/products",
		strings.NewReader(`{"sku":"A/B","name":"slash","stock":1}`))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("SKU with slash should be rejected with 400, got %d body=%s", resp.StatusCode, body)
	}
}
