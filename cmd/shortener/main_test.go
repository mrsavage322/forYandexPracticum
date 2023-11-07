package main

import (
	"github.com/mrsavage322/foryandex/internal/app"
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
		expectedResponseBody   string
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
		{
			name:                 "JSON POST request with a valid link",
			method:               http.MethodPost,
			request:              "/api/shorten",
			body:                 `{"url": "https://example.com"}`,
			expectedStatusCode:   http.StatusCreated,
			expectedResponseBody: `{"result": "YourShortURLLogicHere"}`,
		},
		{
			name:                 "JSON POST request with an empty link",
			method:               http.MethodPost,
			request:              "/api/shorten",
			body:                 ``,
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var response *http.Response
			var err error

			app.URLMap = app.NewURLMapStorage()

			if test.method == http.MethodPost {
				request := httptest.NewRequest(test.method, test.request, strings.NewReader(test.body))
				recorder := httptest.NewRecorder()
				app.HandlePost(recorder, request)
				app.HandleJSON(recorder, request)
				response = recorder.Result()
			} else if test.method == http.MethodGet {
				id := app.GenerateRandomID(5)
				app.URLMap.Set(id, "https://example.com")
				request := httptest.NewRequest(test.method, test.request, nil)
				response = httptest.NewRecorder().Result()
				app.Redirect(httptest.NewRecorder(), request)
			}
			defer response.Body.Close()

			require.NoError(t, err)

			assert.Equal(t, test.expectedStatusCode, response.StatusCode)

			if test.expectedLocationHeader != "" {
				locationHeader := response.Header.Get("Location")
				assert.Equal(t, test.expectedLocationHeader, locationHeader)
			}
		})
	}
}
