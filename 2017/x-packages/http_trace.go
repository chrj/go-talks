package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/net/trace"
)

func Handler(rw http.ResponseWriter, req *http.Request) {

	tr := trace.New("mypkg.Handler", req.URL.Path)
	defer tr.Finish()

	id := req.URL.Query().Get("id")

	tr.LazyPrintf("looking up object id:%v in database", id)

	obj, err := lookupObj(id)
	if err != nil {
		tr.LazyPrintf("database lookup failed: %v", err)
		tr.SetError()
	}

	tr.LazyPrintf("rendering response for object id:%v", id)

	render(rw, obj)

}

func lookupObj(id string) (string, error) {
	return fmt.Sprintf("object representing id:%v", id), nil
}

func render(rw http.ResponseWriter, obj string) {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	fmt.Fprintf(rw, obj)
}

func main() {

	trace.DebugUseAfterFinish = true

	http.HandleFunc("/handler", Handler)
	log.Fatal(http.ListenAndServe(":8123", nil))

}
