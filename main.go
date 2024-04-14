package main

import (
	"context"
	"eauth/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	coll := client.Database("test").Collection("eauth")

	l := log.Default()

	uh := handlers.NewUser(l, coll)
	vh := handlers.NewVerify(l, coll)

	sm := mux.NewRouter()

	post_sr := sm.Methods(http.MethodPost).Subrouter()
	get_sr := sm.Methods(http.MethodGet).Subrouter()

	sm.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})

	post_sr.HandleFunc("/login", uh.Login)
	post_sr.HandleFunc("/register", uh.Register)

	get_sr.HandleFunc("/verifyotp/{email}/{otp}", vh.VerifyOTP)
	get_sr.HandleFunc("/verifyqr/{email}/{otp}", vh.VerifyQR)
	sm.HandleFunc("/waitqr/{email}", vh.WaitQR)

	get_sr.HandleFunc("/user", uh.GetUser)

	server := http.Server{
		Handler: cors.AllowAll().Handler(sm),
		Addr:    ":6969",
	}

	server.ListenAndServe()
}
