package http

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateJobHandler(t *testing.T) {
	// 1. Create a dummy request body
	body := []byte(`{"user_id": "test-user", "prompt": "generate an image"}`)
	req, _ := http.NewRequest("POST", "/jobs", bytes.NewBuffer(body))
	_ = req
	// 2. Create a ResponseRecorder (acts as your client browser/sniffer)
	rr := httptest.NewRecorder()

	// 3. Invoke the handler directly (you'd pass your service mock here)
	// handler := NewHandler(mockService)
	// handler.CreateJob(rr, req)

	// 4. Assert against the response code natively
	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}
