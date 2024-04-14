package handlers

import (
	"context"
	"eauth/types"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Verify struct {
	l    *log.Logger
	coll *mongo.Collection
	sesh *types.SafeMap
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewVerify(l *log.Logger, coll *mongo.Collection) *Verify {
	return &Verify{
		l,
		coll,
		types.NewSafeMap(),
	}
}

func (uh Verify) VerifyOTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]
	otp := vars["otp"]

	dbu := &types.User{}

	err := uh.coll.FindOne(context.Background(), bson.D{{"email", email}}).Decode(&dbu)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if dbu.Token == otp {
		token := CreateToken(email)

		w.Write([]byte(fmt.Sprintf("{\"token\":\"%s\"}", token)))

	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func (vh Verify) VerifyQR(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	otp := vars["otp"]
	email := vars["email"]

	dbu := &types.User{}

	err := vh.coll.FindOne(context.Background(), bson.D{{"email", email}}).Decode(&dbu)

	_ = err

	w.Header().Set("Content-Type", "text/html")

	if err != nil || dbu.Token != otp {
		w.Write([]byte("<h1>Bad Auth</h1>"))
	} else {
		vh.sesh.Put(email, true)

		w.Write([]byte("<h1>Scan successful</h1>"))
	}
}

func (vh Verify) WaitQR(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	email := vars["email"]

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		vh.l.Println(err)
		return
	}

	vh.l.Println("ws conn successful")
	defer conn.Close()
	t := 0

	for conn != nil {
		time.Sleep(1 * time.Second)
		t += 1

		if v, ok := vh.sesh.Get(email); ok && v {
			vh.sesh.Delete(email)
			token := CreateToken(email)

			conn.WriteMessage(
				websocket.TextMessage,
				[]byte(fmt.Sprintf("{\"token\":\"%s\"}",
					token,
				)))
			break
		} else if t == 120 {
			break
		}
	}
}
