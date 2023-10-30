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
	}).Methods("DELETE")

	r.HandleFunc("/movies/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		id := params["id"]
		for _, v := range movies {
			if v.ID == id {
				json.NewEncoder(w).Encode(v)
				break
			}
		}
	}).Methods("GET")

	r.HandleFunc("/movies", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var movie Movie
		err := json.NewDecoder(r.Body).Decode(&movie)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Wrong Body Format"))
		}
		movie.ID = strconv.Itoa(rand.Intn(1000))
		movies = append(movies, movie)
		json.NewEncoder(w).Encode(movie)
	}).Methods("POST")

	r.HandleFunc("/movies/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		params := mux.Vars(r)

		for i, v := range movies {
			if v.ID == params["id"] {
				movies = append(movies[:i], movies[i+1:]...)
				var movie Movie
				json.NewDecoder(r.Body).Decode(&movie)
				movie.ID = params["id"]
				movies = append(movies, movie)
				json.NewEncoder(w).Encode(movie)
				break
			}
		}
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
