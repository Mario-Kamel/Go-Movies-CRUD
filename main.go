package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

var movies []Movie

func main() {
	r := mux.NewRouter()
	movies = append(movies, Movie{ID: "1", Title: "Movie1", Director: &Director{FirstName: "Director", LastName: "One"}})
	movies = append(movies, Movie{ID: "2", Title: "Movie2", Director: &Director{FirstName: "Director", LastName: "Two"}})

	r.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(movies)
	}).Methods("GET")

	r.HandleFunc("/movies/{id}", func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		id := params["id"]
		for i, v := range movies {
			if v.ID == id {
				movies = append(movies[:i], movies[i+1:]...)
				break
			}
		}
		w.WriteHeader(http.StatusOK)
	})
	server := http.Server{
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  120 * time.Second,
		Addr:         ":8000",
		Handler:      r,
	}

	fmt.Printf("Starting server at port 8000\n")
	log.Fatal(server.ListenAndServe())

}
