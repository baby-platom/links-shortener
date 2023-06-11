package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var testingUrl = "https://music.yandex.kz/home"
var shortenedUrl string

type header struct {
	name  string
	value string
}
type want struct {
	code        int
	contentType string
	headers     []header
}
type test struct {
	name string
	want want
}

func TestShortenURLHandler(t *testing.T) {
	tests := []test{
		{
			name: "positive test #1",
			want: want{
				code:        http.StatusCreated,
				contentType: "text/plain",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(testingUrl))
			w := httptest.NewRecorder()
			ShortenURLHandler(w, request)

			res := w.Result()
			assert.Equal(t, res.StatusCode, test.want.code)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.NotEmpty(t, resBody)
			assert.Equal(t, res.Header.Get("Content-Type"), test.want.contentType)

			resBodyString := string(resBody)
			resBodySlice := strings.Split(resBodyString, "/")
			shortenedUrl = resBodySlice[len(resBodySlice)-1]
		})
	}
}

func TestRestoreURLHandler(t *testing.T) {
	tests := []test{
		{
			name: "positive test #1",
			want: want{
				code: http.StatusTemporaryRedirect,
				headers: []header{
					{
						name:  "Location",
						value: testingUrl,
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			target := fmt.Sprintf("/%s", shortenedUrl)
			request := httptest.NewRequest(http.MethodGet, target, strings.NewReader(testingUrl))
			w := httptest.NewRecorder()
			RestoreURLHandler(w, request)

			res := w.Result()
			assert.Equal(t, res.StatusCode, test.want.code)

			for _, header := range test.want.headers {
				assert.Equal(t, res.Header.Get(header.name), header.value)
			}
		})
	}
}
