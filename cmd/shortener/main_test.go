package main

import (
	"github.com/mrsavage322/foryandex/internal/app"
	"github.com/mrsavage322/foryandex/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
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

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var response *http.Response
			var err error

			app.URLMap = storage.NewURLMapStorage()

			if test.method == http.MethodPost {
				request := httptest.NewRequest(test.method, test.request, strings.NewReader(test.body))
				recorder := httptest.NewRecorder()
				app.HandlePost(recorder, request)
				response = recorder.Result()
			} else if test.method == http.MethodGet {
				// Сначала добавим короткий URL
				id := app.GenerateRandomID(5)
				app.URLMap.Set(id, "https://example.com")
				request := httptest.NewRequest(test.method, test.request, nil)
				response = httptest.NewRecorder().Result()
				app.Redirect(httptest.NewRecorder(), request)
			}
			defer response.Body.Close()

			require.NoError(t, err) // Используем err

			assert.Equal(t, test.expectedStatusCode, response.StatusCode)

			if test.expectedLocationHeader != "" {
				locationHeader := response.Header.Get("Location")
				assert.Equal(t, test.expectedLocationHeader, locationHeader)
			}
		})
	}
}
