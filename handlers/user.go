package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/dchest/uniuri"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"eauth/types"
)

type User struct {
	l    *log.Logger
	coll *mongo.Collection
	sesh *types.SafeMap
}

func NewUser(l *log.Logger, coll *mongo.Collection) User {
	return User{l, coll, types.NewSafeMap()}
}

func (uh User) Register(w http.ResponseWriter, r *http.Request) {
	u := types.User{}
	dbu := types.User{}

	err := json.NewDecoder(r.Body).Decode(&u)

	if err != nil {
		uh.l.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = uh.coll.FindOne(context.Background(), bson.D{{"email", u.Email}}).Decode(&dbu)

	if err == mongo.ErrNoDocuments {
		w.WriteHeader(http.StatusOK)
		update, err := uh.coll.InsertOne(context.Background(), u)

		if err != nil {
			uh.l.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		uh.l.Println(update.InsertedID)
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func (uh User) Login(w http.ResponseWriter, r *http.Request) {
	u := types.User{}
	dbu := types.User{}

	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		uh.l.Println(err)
		w.WriteHeader(400)
		return
	}

	err = uh.coll.FindOne(context.Background(), bson.D{{"email", u.Email}}).Decode(&dbu)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if u.Password != dbu.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	go func() {
		otp := uniuri.NewLen(6)
		_, err = uh.coll.UpdateOne(
			context.Background(),
			bson.D{{"email", u.Email}},
			bson.D{{"$set", bson.D{{"token", otp}}}},
		)

		SendMail(u.Email, otp, uh.l)
	}()

	w.WriteHeader(http.StatusOK)
}

func (uh User) GetUser(w http.ResponseWriter, r *http.Request) {
	uh.l.Println()

	auth := strings.Split(r.Header.Get("Authorization"), " ")
	if len(auth) != 2 || auth[0] != "Bearer" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	email, err := DecryptToken(auth[1])

	if err != nil {
		uh.l.Println("Error: ", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	u := types.FEUser{}
	dbu := types.User{}

	err = uh.coll.FindOne(context.Background(), bson.D{{"email", email}}).Decode(&dbu)

	if err != nil {
		uh.l.Println("Error: ", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	u.Name = dbu.Name
	u.Email = dbu.Email

	json.NewEncoder(w).Encode(u)
}
