package run_processor

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"mdm/libs/4_common/smart_context"
)

// TestParseJSONBody проверяет корректность парсинга JSON-тела.
func TestParseJSONBody(t *testing.T) {
	jsonStr := `{"key": "value"}`
	req := httptest.NewRequest("POST", "/", strings.NewReader(jsonStr))
	data, err := parseJSONBody(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if data["key"] != "value" {
		t.Errorf("Expected value 'value' for key 'key', got %v", data["key"])
	}
}

// TestParseJSONBodyEmpty проверяет поведение при пустом теле запроса.
func TestParseJSONBodyEmpty(t *testing.T) {
	req := httptest.NewRequest("POST", "/", nil)
	data, err := parseJSONBody(req)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(data) != 0 {
		t.Errorf("Expected empty map for empty body, got %v", data)
	}
}

// TestJSONResponseMiddlewareSuccess проверяет успешное выполнение обработчика.
func TestJSONResponseMiddlewareSuccess(t *testing.T) {
	sctx := smart_context.NewSmartContext()
	// Dummy AppHandler: возвращает переданные данные без изменений.
	dummyHandler := func(sctx smart_context.ISmartContext, data map[string]interface{}) (interface{}, error) {
		return data, nil
	}
	handlerFunc := JSONResponseMiddleware(sctx, dummyHandler)

	jsonBody := `{"test": "data"}`
	req := httptest.NewRequest("POST", "/", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handlerFunc(rr, req)
	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", res.StatusCode)
	}
	var respData map[string]interface{}
	err := json.NewDecoder(res.Body).Decode(&respData)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if respData["test"] != "data" {
		t.Errorf("Expected key 'test' with value 'data', got %v", respData["test"])
	}
}

// TestJSONResponseMiddlewareError проверяет обработку ошибки внутри AppHandler.
func TestJSONResponseMiddlewareError(t *testing.T) {
	sctx := smart_context.NewSmartContext()
	// Dummy AppHandler: всегда возвращает ошибку.
	dummyHandler := func(sctx smart_context.ISmartContext, data map[string]interface{}) (interface{}, error) {
		return nil, io.EOF
	}
	handlerFunc := JSONResponseMiddleware(sctx, dummyHandler)

	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"any": "data"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	handlerFunc(rr, req)
	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status code 500, got %d", res.StatusCode)
	}
	var respData map[string]string
	err := json.NewDecoder(res.Body).Decode(&respData)
	if err != nil {
		t.Fatalf("Failed to decode error response: %v", err)
	}
	if !strings.Contains(respData["error"], "EOF") {
		t.Errorf("Expected error message containing 'EOF', got %v", respData["error"])
	}
}
