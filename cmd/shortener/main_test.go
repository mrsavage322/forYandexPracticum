package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainPageHandler(t *testing.T) {
	tests := []struct {
		name                   string
		method                 string
		request                string
		body                   string
		expectedStatusCode     int
		expectedLocationHeader string
	}{
		{
			name:               "POST request with a valid link",
			method:             http.MethodPost,
			request:            "/",
			body:               "https://example.com",
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "POST request with an empty link",
			method:             http.MethodPost,
			request:            "/",
			body:               "",
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	urlMap = make(map[string]string)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var response *http.Response
			var err error // Объявляем переменную err

			if test.method == http.MethodPost {
				request := httptest.NewRequest(test.method, test.request, strings.NewReader(test.body))
				// Используйте тот же recorder для записи ответа
				recorder := httptest.NewRecorder()
				mainPage(recorder, request)
				response = recorder.Result()
			} else if test.method == http.MethodGet {
				// Сначала добавим короткий URL
				id := generateRandomID(5)
				urlMap[id] = "https://example.com"

				request := httptest.NewRequest(test.method, test.request, nil)
				response = httptest.NewRecorder().Result()
				redirect(httptest.NewRecorder(), request)
			}

			require.NoError(t, err) // Используем err

			assert.Equal(t, test.expectedStatusCode, response.StatusCode)

			if test.expectedLocationHeader != "" {
				locationHeader := response.Header.Get("Location")
				assert.Equal(t, test.expectedLocationHeader, locationHeader)
			}
		})
	}
}
