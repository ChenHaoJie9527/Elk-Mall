package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ChenHaoJie9527/Elk-Mall/internal/common/response"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

func TestRecover_PreventsProcessCrash(t *testing.T) {
	e := echo.New()
	e.HTTPErrorHandler = response.HTTPErrorHandler
	e.Use(middleware.Recover())
	e.GET("/boom", func(c *echo.Context) error {
		panic("boom")
	})

	req := httptest.NewRequest(http.MethodGet, "/boom", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("HTTP = %d, want 500; Recover 没把 panic 收成 error", rec.Code)
	}

	var body response.Response
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal %q: %v", rec.Body.String(), err)
	}
	if body.Code != http.StatusInternalServerError {
		t.Errorf("code = %d, want 500", body.Code)
	}
	if body.Msg != http.StatusText(http.StatusInternalServerError) {
		t.Errorf("msg = %q, want %q（不应把 panic 内容回给客户端）", body.Msg, http.StatusText(http.StatusInternalServerError))
	}
}

func TestWithoutRecover_PanicEscapes(t *testing.T) {
	e := echo.New()
	e.HTTPErrorHandler = response.HTTPErrorHandler
	e.GET("/boom", func(c *echo.Context) error {
		panic("boom")
	})

	req := httptest.NewRequest(http.MethodGet, "/boom", nil)
	rec := httptest.NewRecorder()

	defer func() {
		if r := recover(); r == nil {
			t.Fatal("没有 Recover 时 ServeHTTP 应当 panic")
		}
	}()
	e.ServeHTTP(rec, req)
	t.Fatal("panic 被接住了，说明这条路径不该走到")
}
