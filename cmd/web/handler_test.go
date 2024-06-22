package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingHandler(t *testing.T) {
	t.Parallel() // run this concurrently just an experiment

	request := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	ping(recorder, request)

	response := recorder.Result()
	defer response.Body.Close() // important, otherwise the TCP connection cannot be reused after after body is not used
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "OK", string(body))
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

func TestPingEndToEnd(t *testing.T) {
	app := newTestApplication(t)
	testServer := newTestServer(t, app.generateRoutes())
	defer testServer.Close()
	code, _, body := testServer.get(t, "/ping")

	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "OK", string(body))
}

func TestShowSnippet(t *testing.T) {

	app := newTestApplication(t)
	testServer := newTestServer(t, app.generateRoutes())
	defer testServer.Close()
	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody []byte
	}{
		{"Valid ID", "/snippet/1", http.StatusOK, []byte("Test Snippet")},
		{"Non-existent ID", "/snippet/2", http.StatusNotFound, nil},
		{"Negative ID", "/snippet/-1", http.StatusNotFound, nil},
		{"Decimal ID", "/snippet/1.23", http.StatusNotFound, nil},
		{"String ID", "/snippet/foo", http.StatusNotFound, nil},
		{"Empty ID", "/snippet/", http.StatusNotFound, nil},
		{"Trailing slash", "/snippet/1/", http.StatusNotFound, nil},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			code, _, body := testServer.get(t, test.urlPath)
			assert.Equal(t, test.wantCode, code)

			// assert.Contains can only checks whether a substring is contained by a string
			// so both test body and test.wantBody need to be casted into string
			assert.Contains(t, string(body), string(test.wantBody))
		})
	}
}

func TestSignupUser(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.generateRoutes())
	defer ts.Close()
	_, _, body := ts.get(t, "/user/signup")
	csrfToken := extractCSRFToken(t, body)

	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		csrfToken    string
		wantCode     int
		wantBody     []byte
	}{
		{"Valid submission", "Bob", "bob@example.com", "validPa$$word", csrfToken,
			http.StatusSeeOther, nil},
		{"Empty name", "", "bob@example.com", "validPa$$word", csrfToken, http.StatusOK,
			[]byte("name is required.")},
		{"Empty email", "Bob", "", "validPa$$word", csrfToken, http.StatusOK,
			[]byte("email is required.")},
		{"Empty password", "Bob", "bob@example.com", "", csrfToken, http.StatusOK,
			[]byte("password is required.")},
		{"Invalid email (incomplete domain)", "Bob", "bob@example.", "validPa$$word",
			csrfToken, http.StatusOK, []byte("email is invalid.")},
		{"Invalid email (missing @)", "Bob", "bobexample.com", "validPa$$word", csrfToken,
			http.StatusOK, []byte("email is invalid.")},
		{"Invalid email (missing local part)", "Bob", "@example.com", "validPa$$word",
			csrfToken, http.StatusOK, []byte("email is invalid.")},
		{"Short password", "Bob", "bob@example.com", "pa$$word", csrfToken, http.StatusOK,
			[]byte("password has to be atleast 10 characters long.")},
		{"Duplicate email", "Bob", "dupe@email.com", "validPa$$word", csrfToken, http.StatusOK,
			[]byte("email is already in use.")},
		{"Invalid CSRF Token", "", "", "", "wrongToken", http.StatusBadRequest, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			form := url.Values{}
			form.Add("name", tt.userName)
			form.Add("email", tt.userEmail)
			form.Add("password", tt.userPassword)
			form.Add("csrf_token", tt.csrfToken)

			code, _, body := ts.postForm(t, "/user/signup", form)
			t.Logf("testing %q for want-code %d and want-body %q", tt.name, tt.wantCode,
				tt.wantBody)

			if code != tt.wantCode {
				t.Errorf("want %d; got %d", tt.wantCode, code)
			}

			if !bytes.Contains(body, tt.wantBody) {
				t.Errorf("want body to contain %q, but got %q", tt.wantBody, body)
			}
		})
	}
}

func TestCreateSnippetForm(t *testing.T) {

	t.Run("Unauthorized", func(t *testing.T) {
		app := newTestApplication(t)
		ts := newTestServer(t, app.generateRoutes())
		defer ts.Close()

		code, header, _ := ts.get(t, "/snippet/create")

		assert.Equal(t, http.StatusSeeOther, code)
		assert.Equal(t, "/user/login", header.Get("Location"))

	})

	t.Run("Authorized", func(t *testing.T) {
		app := newTestApplication(t)
		ts := newTestServer(t, app.generateRoutes())
		defer ts.Close()

		// show login form
		_, _, body := ts.get(t, "/user/login")
		csrfToken := extractCSRFToken(t, body)

		// login attempt
		form := url.Values{}
		form.Add("email", "alice@gmail.com")
		form.Add("password", "Cucibaju123")
		form.Add("csrf_token", csrfToken)
		ts.postForm(t, "/user/login", form)

		// show create snippet form
		code, _, body := ts.get(t, "/snippet/create")

		assert.Equal(t, http.StatusOK, code)
		formTag := `<form action="/snippet/create" method="POST">`
		if !bytes.Contains(body, []byte(formTag)) {
			t.Errorf("want body %s to contain %q", body, formTag)
		}
	})
}
