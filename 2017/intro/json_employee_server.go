package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// BEGIN types OMIT

type Employee struct {
	Name       string  `json:"name"`
	Department string  `json:"department"`
	Salary     float64 `json:"salary"`
}

type EmployeeRegistry map[int]Employee

// END types OMIT

// BEGIN handler OMIT

func (er EmployeeRegistry) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	i, err := strconv.Atoi(mux.Vars(req)["id"])
	if err != nil {
		http.Error(rw, "Couldn't decode ID", http.StatusBadRequest)
		return
	}

	if e, found := er[i]; !found {
		http.NotFound(rw, req)
	} else {

		rw.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(rw).Encode(e); err != nil {
			log.Printf("error encoding reply: %v", err)
		}

	}

}

// END handler OMIT

// BEGIN demo OMIT

func main() {

	er := EmployeeRegistry{}

	er[60] = Employee{Name: "Christian", Department: "Development", Salary: 1000}
	er[80] = Employee{Name: "John", Department: "Management", Salary: 500}

	r := mux.NewRouter()
	r.Handle("/employee/{id:[0-9]+}", er)

	log.Fatal(http.ListenAndServe(":8123", r))

}

// END demo OMIT
