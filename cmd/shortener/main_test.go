package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/baby-platom/links-shortener/internal/app"
	"github.com/baby-platom/links-shortener/internal/models"
	"github.com/baby-platom/links-shortener/internal/shortid"
)

const defaultContentType = "text/plain"
const testingURL = "https://music.yandex.kz/home"

var ts = httptest.NewServer(app.Router())

type header struct {
	name  string
	value string
}
type want struct {
	code           int
	contentType    string
	headers        []header
	body           string
	bodyIsNotEmpty bool
}
type request struct {
	method      string
	path        string
	body        io.Reader
	contentType string
}
type test struct {
	name    string
	request request
	want    want
}

func testRequest(
	t *testing.T,
	ts *httptest.Server,
	test test,
) {
	requestData := test.request
	wantData := test.want

	req, err := http.NewRequest(requestData.method, ts.URL+requestData.path, requestData.body)
	assert.NoError(t, err, "Error creating request")
	contentType := requestData.contentType
	if contentType == "" {
		contentType = defaultContentType
	}
	req.Header.Set("Content-Type", contentType)

	client := ts.Client()
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	resp, err := client.Do(req)
	assert.NoError(t, err, "Error making HTTP request")

	respBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err, "Error reading response body")
	defer resp.Body.Close()
	respBodyString := string(respBody)

	assert.Equal(t, wantData.code, resp.StatusCode, "Response code didn't match expected")
	if wantData.contentType != "" {
		assert.Equal(t, wantData.contentType, resp.Header.Get("Content-Type"))
	}

	if wantData.bodyIsNotEmpty {
		assert.NotEmpty(t, respBodyString)
	}
	if wantData.body != "" {
		assert.JSONEq(t, wantData.body, respBodyString)
	}

	for _, header := range wantData.headers {
		assert.Equal(t, header.value, resp.Header.Get(header.name))
	}
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
			testRequest(t, ts, test)
		})
	}
}

func TestRestoreURLHandler(t *testing.T) {
	shortenedUrlsByID := shortid.ShortenedUrlsByIDType{
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
			testRequest(t, ts, test)
		})
	}
}

func TestShortenAPIHandler(t *testing.T) {
	contentTypeJSON := "application/json"
	path := "/api/shorten"

	positiveBody, err := json.Marshal(
		models.ShortentRequest{
			URL: testingURL,
		},
	)
	require.NoError(t, err)

	tests := []test{
		{
			name: "positive test #0",
			request: request{
				method:      http.MethodPost,
				path:        path,
				body:        bytes.NewBuffer(positiveBody),
				contentType: contentTypeJSON,
			},
			want: want{
				code:        http.StatusCreated,
				contentType: contentTypeJSON,
			},
		},
		{
			name: "negative test #0",
			request: request{
				method:      http.MethodPost,
				path:        path,
				contentType: contentTypeJSON,
			},
			want: want{
				code: http.StatusInternalServerError,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testRequest(t, ts, test)
		})
	}
	ts.Close()
}
