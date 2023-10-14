package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// test PostUrl
func Test_postURL(t *testing.T) {
	type want struct {
		statusCode  int
		contentType string
		body        string
	}
	tests := []struct {
		name    string
		want    want
		request string
		body    string
	}{
		{
			name: "test status code 201 (created)",
			want: want{
				statusCode:  201,
				contentType: "text/plain",
			},
			request: "/",
			body:    "asdasda",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			postURL(w, request)

			res := w.Result()

			assert.Equal(t, tt.want.statusCode, res.StatusCode)
			assert.NotNil(t, res.Body)
		})
	}
}

// test getURL
func Test_getURL(t *testing.T) {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name    string
		want    want
		request string
	}{
		{
			name: "return 307 satatus code",
			want: want{
				statusCode: 307,
			},
			request: "/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()
			getURL(w, request)

			res := w.Result()

			assert.Equal(t, tt.want.statusCode, http.StatusTemporaryRedirect)
			assert.NotNil(t, res.Body)
		})
	}
}