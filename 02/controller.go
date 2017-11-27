package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

// Controller ハンドラ用
type Controller struct {
}

// NewController コンストラクタ
func NewController() *Controller {
	return &Controller{}
}

// User user
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Show endpoint
func (c *Controller) Show(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	res := User{ID: id, Name: fmt.Sprint("name_", id)}
	return http.StatusOK, res, nil
}

// List endpoint
func (c *Controller) List(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	users := []User{
		{1, "hoge"},
		{2, "foo"},
		{3, "bar"},
	}
	return http.StatusOK, users, nil
}

// AuthInfo 何らかの認証後にトークン発行するようなもの
type AuthInfo struct {
	Authorization string `json:"authorization"`
}

// Login endpoint
func (c *Controller) Login(w http.ResponseWriter, r *http.Request) (int, interface{}, error) {
	token := r.URL.Query().Get("token")
	if token != "token" {
		return http.StatusUnauthorized, nil, fmt.Errorf("有効でないトークンです: %s", token)
	}
	res := AuthInfo{Authorization: "admin"}
	return http.StatusOK, res, nil
}
