package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Тест 1: Запрос сформирован корректно, сервис возвращает код ответа 200 и тело ответа не пустое.
func TestMainHandlerCorrectRequest(t *testing.T) {

	req := httptest.NewRequest("GET", "/cafe?city=moscow&count=2", nil)
	responseRecorder := httptest.NewRecorder()

	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusOK, responseRecorder.Code, "Expected status 200 OK for a valid request")

	body, err := io.ReadAll(responseRecorder.Body)
	require.NoError(t, err, "Failed to read response body")

	assert.NotEmpty(t, body, "Expected response body to be non-empty")

	expectedBody := "Мир кофе,Сладкоежка"
	assert.Equal(t, expectedBody, string(body), "Expected specific cafe list for count=2")

	returnedCafes := strings.Split(string(body), ",")
	assert.Len(t, returnedCafes, 2, "Expected 2 cafes in the response")
}

// Тест 2: Город, который передаётся в параметре city, не поддерживается.
// Сервис возвращает код ответа 400 и ошибку wrong city value в теле ответа.
func TestMainHandlerUnsupportedCity(t *testing.T) {

	req := httptest.NewRequest("GET", "/cafe?city=london&count=1", nil)
	responseRecorder := httptest.NewRecorder()

	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusBadRequest, responseRecorder.Code, "Expected status 400 Bad Request for unsupported city")

	body, err := io.ReadAll(responseRecorder.Body)
	require.NoError(t, err, "Failed to read response body")

	expectedErrorMessage := "wrong city value"
	assert.Equal(t, expectedErrorMessage, string(body), "Expected 'wrong city value' error message in body")
}

// Тест 3: Если в параметре count указано больше, чем есть всего,
// должны вернуться все доступные кафе.
func TestMainHandlerCountMoreThanTotal(t *testing.T) {
	totalCount := len(cafeList["moscow"])

	req := httptest.NewRequest("GET", "/cafe?city=moscow&count=100", nil)
	responseRecorder := httptest.NewRecorder()

	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusOK, responseRecorder.Code, "Expected status 200 OK when count exceeds total")

	body, err := io.ReadAll(responseRecorder.Body)
	require.NoError(t, err, "Failed to read response body")

	assert.NotEmpty(t, body, "Expected response body to be non-empty")

	expectedBody := strings.Join(cafeList["moscow"], ",")
	assert.Equal(t, expectedBody, string(body), "Expected all available cafes in body when count exceeds")

	returnedCafes := strings.Split(string(body), ",")
	assert.Len(t, returnedCafes, totalCount, "Expected to return all available cafes")
}
