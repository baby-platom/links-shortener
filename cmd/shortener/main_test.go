package main

import (
	// "fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testingURL = "https://music.yandex.kz/home"
var shortenedURL string
var ts = httptest.NewServer(Router())

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

func testRequest(t *testing.T, ts *httptest.Server, method,
	path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	client := ts.Client()
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	require.NoError(t, err)

	return resp, string(respBody)
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
			res, resBodyString := testRequest(t, ts, http.MethodPost, "/", strings.NewReader(testingURL))

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.NotEmpty(t, resBodyString)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))

			resBodySlice := strings.Split(resBodyString, "/")
			shortenedURL = resBodySlice[len(resBodySlice)-1]
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
						value: testingURL,
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			target := "/" + shortenedURL
			res, _ := testRequest(t, ts, http.MethodGet, target, nil)
			ts.Close()

			assert.Equal(t, test.want.code, res.StatusCode)
			for _, header := range test.want.headers {
				assert.Equal(t, header.value, res.Header.Get(header.name))
			}
		})
	}
}
