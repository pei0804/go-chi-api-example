package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	httpdoc "go.mercari.io/go-httpdoc"
)

func TestListController(t *testing.T) {
	document := &httpdoc.Document{
		Name: "List Controller",
	}
	defer func() {
		if err := document.Generate("doc/list.md"); err != nil {
			t.Fatalf("err: %s", err)
		}
	}()

	// mux := http.NewServeMux()
	router := chi.NewRouter()
	c := NewController()
	router.Method("GET", "/api/members", httpdoc.Record(handler(c.List), document, &httpdoc.RecordOption{
		Description: "get user list",
		WithValidate: func(validator *httpdoc.Validator) {
			validator.RequestParams(t, []httpdoc.TestCase{})
			validator.RequestHeaders(t, []httpdoc.TestCase{
				{Target: "Authorization", Expected: "admin", Description: "auth token"},
			})
			validator.ResponseStatusCode(t, http.StatusOK)
			// FIXME slice validation
			// validator.ResponseBody(t, []httpdoc.TestCase{
			// 	{Target: "id", Expected: 1, Description: ""},
			// 	{Target: "name", Expected: "hoge", Description: ""},
			// }, &[]User{})
		},
	}))

	testServer := httptest.NewServer(router)
	defer testServer.Close()
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/api/members", nil)
	req.Header.Add("Authorization", "admin")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	resp, err := http.DefaultClient.Do(req)

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	jsonStr := `[
    {
        "id": 1,
        "name": "hoge"
    },
    {
        "id": 2,
        "name": "foo"
    },
    {
        "id": 3,
        "name": "bar"
    }
]`
	if string(b) != jsonStr {
		t.Fatalf("err: %s", err)
	}
}

func TestShowController(t *testing.T) {
	document := &httpdoc.Document{
		Name: "Show Controller",
	}
	defer func() {
		if err := document.Generate("doc/show.md"); err != nil {
			t.Fatalf("err: %s", err)
		}
	}()
	// mux := mux.NewRouter()
	c := NewController()
	router := chi.NewRouter()
	router.Method("GET", "/api/members/{id}", httpdoc.Record(handler(c.Show), document, &httpdoc.RecordOption{
		Description: "get user show",
		WithValidate: func(validator *httpdoc.Validator) {
			validator.RequestHeaders(t, []httpdoc.TestCase{
				{Target: "Authorization", Expected: "admin", Description: "auth token"},
			})
			validator.ResponseStatusCode(t, http.StatusOK)
			validator.ResponseBody(t, []httpdoc.TestCase{
				{Target: "ID", Expected: 1, Description: "user id"},
				{Target: "Name", Expected: "name_1", Description: "user name"}},
				&User{},
			)
		},
	}))
	testServer := httptest.NewServer(router)
	defer testServer.Close()
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/api/members/1", nil)
	req.Header.Add("Authorization", "admin")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestLoginController(t *testing.T) {
	document := &httpdoc.Document{
		Name: "Login Controller",
	}
	defer func() {
		if err := document.Generate("doc/login.md"); err != nil {
			t.Fatalf("err: %s", err)
		}
	}()
	// mux := mux.NewRouter()
	c := NewController()
	router := chi.NewRouter()
	router.Method("GET", "/api/auth/login", httpdoc.Record(handler(c.Login), document, &httpdoc.RecordOption{
		Description: "auth login",
		WithValidate: func(validator *httpdoc.Validator) {
			validator.RequestParams(t, []httpdoc.TestCase{
				{Target: "token", Expected: "token", Description: "token"},
			})
			validator.ResponseStatusCode(t, http.StatusOK)
			validator.ResponseBody(t, []httpdoc.TestCase{
				{Target: "Authorization", Expected: "admin", Description: "Authorization info"}},
				&AuthInfo{},
			)
		},
	}))
	testServer := httptest.NewServer(router)
	defer testServer.Close()
	req, err := http.NewRequest(http.MethodGet, testServer.URL+"/api/auth/login?token=token", nil)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("err: %s", err)
	}
}
