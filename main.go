package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Movie struct {
	ID       string    `json:"id"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func main() {
	r := mux.NewRouter()
	movies := make(map[string]Movie)

	movies["1"] = Movie{ID: "1", Title: "Movie1", Director: &Director{FirstName: "Director", LastName: "One"}}
	movies["2"] = Movie{ID: "2", Title: "Movie2", Director: &Director{FirstName: "Director", LastName: "Two"}}

	r.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		list := make([]Movie, 0)
		for _, v := range movies {
			list = append(list, v)
		}
		json.NewEncoder(w).Encode(list)
	}).Methods("GET")

	r.HandleFunc("/movies/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		id := params["id"]
		movie, ok := movies[id]
		if ok {
			json.NewEncoder(w).Encode(movie)
			return
		}
		w.WriteHeader(http.StatusNotFound)

	}).Methods("GET")

	r.HandleFunc("/movies/{id}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]
		delete(movies, id)
		w.WriteHeader(http.StatusOK)
	}).Methods("DELETE")

	r.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var movie Movie
		err := json.NewDecoder(r.Body).Decode(&movie)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Wrong Body Format"))
		}
		movie.ID = strconv.Itoa(rand.Intn(1000))
		movies[movie.ID] = movie
		json.NewEncoder(w).Encode(movie)
	}).Methods("POST")

	r.HandleFunc("/movies/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		params := mux.Vars(r)
		id := params["id"]
		var movie Movie
		json.NewDecoder(r.Body).Decode(&movie)
		movie.ID = id
		movies[id] = movie
		json.NewEncoder(w).Encode(movie)
	}).Methods("PUT")

	server := http.Server{
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         ":8000",
		Handler:      r,
	}

	fmt.Printf("Starting server at port 8000\n")
	log.Fatal(server.ListenAndServe())

}
