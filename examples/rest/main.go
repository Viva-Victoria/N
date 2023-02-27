package main

import (
	"gitea.voopsen/OSS/n"
	"gitea.voopsen/OSS/n/log"
	"net/http"
	"os"
	"os/signal"
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
	router := n.NewRouter("/api/v1", log.NewLog())
	router.Handle("/user/{userId:\\d+}/device/{deviceId:\\d+}/pipeline/{pipeId:\\d+}/script/{scriptId:\\d+}", n.HandlerFunc(func(ctx n.Context) error {
		return ctx.WriteJSON("ok")
	}))
	router.Handle("/user/all", n.HandlerFunc(func(ctx n.Context) error {
		return ctx.WriteJSON(_users)
	}))
	router.Handle("/user/{id:\\d+}", n.HandlerFunc(func(ctx n.Context) error {
		var id int
		if err := ctx.Vars().Get("id", &id); err != nil {
			return n.NewBadRequestError(err)
		}
		if isUserIdInvalid(id) {
			ctx.Status(http.StatusNotFound)
			return nil
		}

		return ctx.WriteJSON(_users[id])
	}))
	router.Post("/user", n.HandlerFunc(func(ctx n.Context) error {
		var user User
		if err := ctx.ReadJSON(user); err != nil {
			return n.NewBadRequestError(err)
		}

		user.Id = len(_users)
		_users = append(_users, user)
		return nil
	}))
	router.Put("/user", n.HandlerFunc(func(ctx n.Context) error {
		var user User
		if err := ctx.ReadJSON(&user); err != nil {
			return n.NewBadRequestError(err)
		}

		if isUserIdInvalid(user.Id) {
			ctx.Status(http.StatusNotFound)
			return nil
		}

		_users[user.Id] = user
		return ctx.WriteJSON(user)
	}))
	router.Delete("/user/{id:\\d+}", n.HandlerFunc(func(ctx n.Context) error {
		var id int
		if err := ctx.Vars().Get("id", &id); err != nil {
			return n.NewBadRequestError(err)
		}

		if isUserIdInvalid(id) {
			ctx.Status(http.StatusNotFound)
			return nil
		}

		_users = append(_users[:id], _users[id+1:]...)
		return nil
	}))

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
