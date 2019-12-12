package controller

import (
	"bytes"
	"encoding/json"
	"github.com/becosuke/guestbook-api/config"
	"github.com/becosuke/guestbook-api/domain/model"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var serveMux = http.NewServeMux()

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		log.Fatal(err)
	}
	ret := m.Run()
	teardown()
	os.Exit(ret)
}

func setup() error {
	model.Flush()
	Register(serveMux)
	return nil
}

func teardown() {
	model.Flush()
}

func add(t *testing.T, name string, body string) {
	document := model.Document{Name: name, Body: body}
	request, _ := json.Marshal(document)
	res, err := model.Add(request)
	if err != nil {
		t.Log(err)
	}
	t.Log(string(res))
}

func TestIndex(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	serveMux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("code: %d", w.Code)
	}

	t.Log(w.Body.String())
}

func TestPing(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/ping", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	serveMux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("code: %d", w.Code)
	}

	t.Log(w.Body.String())
}

func TestAdd(t *testing.T) {
	if !config.IsLocal() {
		t.Log("skip test")
		return
	}

	param := model.Document{Name: "add_name", Body: "add_body"}
	body, _ := json.Marshal(param)
	req, err := http.NewRequest(http.MethodPost, "/add", bytes.NewBufferString(string(body)))
	if err != nil {
		t.Error(err)
	}
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	serveMux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("code: %d", w.Code)
	}

	t.Log(w.Body.String())
}

func TestRange(t *testing.T) {
	for i := 0; i < 10; i++ {
		add(t, "name"+strconv.Itoa(i), "body"+strconv.Itoa(i))
	}

	req, err := http.NewRequest(http.MethodGet, "/range?p=1", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	serveMux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("code: %d", w.Code)
	}

	t.Log(w.Body.String())
}

func TestFlush(t *testing.T) {
	if !config.IsLocal() {
		t.Log("skip test")
		return
	}

	req, err := http.NewRequest(http.MethodGet, "/flush", nil)
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()
	serveMux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("code: %d", w.Code)
	}

	t.Log(w.Body.String())
}
