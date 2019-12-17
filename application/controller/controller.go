package controller

import (
	"github.com/becosuke/guestbook-api/domain/model"
	"io/ioutil"
	"net/http"
	"strconv"
)

func Register(serveMux *http.ServeMux) {
	serveMux.HandleFunc("/", Index)
	serveMux.HandleFunc("/ping", Ping)
	serveMux.HandleFunc("/add", Add)
	serveMux.HandleFunc("/range", Range)
	serveMux.HandleFunc("/flush", Flush)
	serveMux.HandleFunc("/random", Random)
}

func Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func Ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := model.Ping(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func Add(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "application/json" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := model.Add(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

func Range(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("p"))
	if err != nil || page < 1 {
		page = 1
	}

	res, err := model.Range(page)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(res)
}

func Flush(w http.ResponseWriter, _ *http.Request) {
	model.Flush()
	w.WriteHeader(http.StatusOK)
}

func Random(w http.ResponseWriter, _ *http.Request) {
	model.Random()
	w.WriteHeader(http.StatusOK)
}
