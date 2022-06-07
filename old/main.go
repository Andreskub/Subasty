package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Auction struct {
	ID     string `json:"id,omitempty"`
	Amount int    `json:"amount,omitempty"`
}

type Person struct {
	ID        string `json:"id,omitempty"`
	FirstName string `json:"firstname,omitempty"`
	LastName  string `json:"lastname,omitempty"`
	Amount    int    `json:"amount,omitempty"`
}

var (
	people   []Person
	auctions []Auction
)

/* APIs PERSONS */

func GetPeopleEndpoint(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(people)
}

func GetPersonEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	json.NewEncoder(w).Encode(&Person{})
}

func CreatePersonEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	var person Person
	_ = json.NewDecoder(req.Body).Decode(&person)

	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

func DeletePersonEndpoint(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(people)
}

/* APIs AUCTIONS */

func getAllAuctions(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(auctions)
}

func createAuction(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	var auction Auction
	_ = json.NewDecoder(req.Body).Decode(&auction)

	auction.ID = params["id"]
	auctions = append(auctions, auction)
	json.NewEncoder(w).Encode(auctions)
}

/* APIs AUCTIONS/:id */

func getAuctionById(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	for _, item := range auctions {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	json.NewEncoder(w).Encode(&Auction{})
}

func updateAuctionById(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	for i := 0; i < len(auctions); i++ {
		if auctions[i].ID == params["id"] {
			amount_value := params["amount"]
			amount_value2, _ := strconv.Atoi(amount_value)
			auctions[i].Amount = amount_value2
			json.NewEncoder(w).Encode(auctions[i])
			return
		}
	}

	json.NewEncoder(w).Encode(&Auction{})

}

func deleteAuction(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	for index, item := range auctions {
		if item.ID == params["id"] {
			auctions = append(auctions[:index], auctions[index+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(auctions)
}

/* Main */

func main() {
	router := mux.NewRouter()

	// examples

	people = append(people, Person{ID: "1", FirstName: "Ryan", LastName: "Ray", Address: &Address{City: "Dublin", State: "California"}})
	people = append(people, Person{ID: "2", FirstName: "Joe", LastName: "Doe"})

	// endpoints
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePersonEndpoint).Methods("DELETE")

	auctions = append(auctions, Auction{ID: "1", Amount: 100})

	router.HandleFunc("/auctions", getAllAuctions).Methods("GET")
	router.HandleFunc("/auctions/{id}", getAuctionById).Methods("GET")
	router.HandleFunc("/auctions", createAuction).Methods("POST")
	router.HandleFunc("/auctions/{id}", updateAuctionById).Methods("PUT")
	router.HandleFunc("/auctions/{id}", deleteAuction).Methods("PUT")

	log.Fatal(http.ListenAndServe(":3000", router))
}

// api/subastas -> post: crear nueva subastas
// 				get: devolver las subastas

// api/subastas/:id -> post: crear una puja
// 					get: info de la subasta

// api/personas -> post: crear persona (/registro)
// 				get: devuelve las personas

// api/personas/:id -> get: info de la persona
// 					put: cuando alguien gano una subasta le resta plata
