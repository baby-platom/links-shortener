package main

import (
	// "fmt"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/baby-platom/links-shortener/internal/app"
)

var testingURL = "https://music.yandex.kz/home"
var ts = httptest.NewServer(app.Router())

type header struct {
	name  string
	value string
}
type want struct {
	code        int
	contentType string
	headers     []header
}
type request struct {
	method string
	path   string
	body   io.Reader
}
type test struct {
	name    string
	request request
	want    want
}

func testRequest(
	t *testing.T,
	ts *httptest.Server,
	request request,
) (*http.Response, string) {
	req, err := http.NewRequest(request.method, ts.URL+request.path, request.body)
	require.NoError(t, err)

	client := ts.Client()
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err := client.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	return resp, string(respBody)
}

func TestShortenURLHandler(t *testing.T) {
	tests := []test{
		{
			name: "positive test #0",
			request: request{
				method: http.MethodPost,
				path:   "/",
				body:   strings.NewReader(testingURL),
			},
			want: want{
				code:        http.StatusCreated,
				contentType: "text/plain",
			},
		},
		{
			name: "negative test #0",
			request: request{
				method: http.MethodPost,
				path:   "/",
				body:   strings.NewReader("  "),
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, resBodyString := testRequest(t, ts, test.request)
			res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)
			assert.NotEmpty(t, resBodyString)
			if test.want.contentType != "" {
				assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			}
		})
	}
}

func TestRestoreURLHandler(t *testing.T) {
	shortenedUrlsByID := app.ShortenedUrlsByIDType{
		"some_id": testingURL,
	}
	for key, value := range shortenedUrlsByID {
		app.ShortenedUrlsByID[key] = value
	}

	tests := []test{
		{
			name: "negative test #0",
			request: request{
				method: http.MethodGet,
				path:   "/" + "not_existing_id",
			},
			want: want{
				code: http.StatusBadRequest,
			},
		},
	}
	i := 0
	for key, value := range shortenedUrlsByID {
		tests = append(
			tests,
			test{
				name: fmt.Sprintf("positive test %d", i),
				request: request{
					method: http.MethodGet,
					path:   "/" + key,
				},
				want: want{
					code: http.StatusTemporaryRedirect,
					headers: []header{
						{
							name:  "Location",
							value: value,
						},
					},
				},
			},
		)
		i++
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, _ := testRequest(t, ts, test.request)
			res.Body.Close()

			assert.Equal(t, test.want.code, res.StatusCode)
			for _, header := range test.want.headers {
				assert.Equal(t, header.value, res.Header.Get(header.name))
			}
		})
	}
	ts.Close()
}
