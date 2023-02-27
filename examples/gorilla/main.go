package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type User struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

var (
	_users = []User{
		{
			Id:        0,
			FirstName: "Kirill",
			LastName:  "Ispolnov",
		},
		{
			Id:        1,
			FirstName: "Andrei",
			LastName:  "Scheglov",
		},
		{
			Id:        2,
			FirstName: "Andrei",
			LastName:  "Melnikov",
		},
	}
)

func main() {
	router := mux.NewRouter()
	router.Handle("/api/v1/user/{userId:\\d+}/device/{deviceId:\\d+}/pipeline/{pipeId:\\d+}/script/{scriptId:\\d+}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJson(w, "ok")
	}))
	router.Handle("/api/v1/user/all", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJson(w, _users)
	}))
	router.Handle("/api/v1/user/{id:\\d+}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if isUserIdInvalid(id) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		writeJson(w, _users[id])
	}))
	router.Handle("/api/v1/user", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user User
		if err := readJson(r, user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		user.Id = len(_users)
		_users = append(_users, user)
		return
	})).Methods(http.MethodPost)
	router.Handle("/api/v1/user", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var user User
		if err := readJson(r, &user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if isUserIdInvalid(user.Id) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		_users[user.Id] = user
		writeJson(w, user)
	})).Methods(http.MethodPut)
	router.Handle("/api/v1/user/{id:\\d+}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if isUserIdInvalid(id) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		_users = append(_users[:id], _users[id+1:]...)
	})).Methods(http.MethodDelete)

	server := http.Server{
		Addr:    ":80",
		Handler: router,
	}

	go server.ListenAndServe()
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)
	<-done
}

func isUserIdInvalid(id int) bool {
	return id < 0 || id >= len(_users)
}

func readJson(r *http.Request, a any) error {
	return json.NewDecoder(r.Body).Decode(a)
}

func writeJson(w http.ResponseWriter, a any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(a); err != nil {
		log.Println(err)
	}
}
