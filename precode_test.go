package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainHandlerReturns200AndNotEmptyBody(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=1&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// Проверка статус-кода
	assert.Equal(t, http.StatusOK, responseRecorder.Code, "expected status code 200")

	// Проверка, что тело не пустое
	body := responseRecorder.Body.String()
	assert.NotEmpty(t, body, "response body not be empty")
}

func TestMainHandlerWhenCityUnknown(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=1&city=qwerty123", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// Проверка статус-кода
	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code, "expected status code 400, unknown city")

	// Проверка сообщения об ошибке
	body := responseRecorder.Body.String()
	assert.Equal(t, "wrong city value", body, "expected error message for unknown city")
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	// Переписал totalCount = 4 на случай если значение изменится в будущем
	totalCount := len(cafeList["moscow"])
	req := httptest.NewRequest("GET", "/cafe?count=10&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// Проверка статус-кода
	require.Equal(t, http.StatusOK, responseRecorder.Code, "expected status code 200")

	// Проверка, что возвращаются все кафе из списка
	body := responseRecorder.Body.String()
	list := strings.Split(body, ",")

	assert.Equal(t, totalCount, len(list), "expected cafe count does not match the response")
}
