package main

import (
	// "fmt"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	// "reflect"

	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//Book : defines Book Struct
type Book struct {
	ID     int
	Title  string
	Author string
	Year   string
}

var books []Book

// var db *sql.DB

func main() {
	router := mux.NewRouter()

	books = append(books,
		Book{ID: 1, Title: "Book 1", Author: "Author 1", Year: "1993"},
		Book{ID: 2, Title: "Book 2", Author: "Author 2", Year: "2993"},
		Book{ID: 3, Title: "Book 3", Author: "Author 3", Year: "3993"},
		Book{ID: 5, Title: "Book 4", Author: "Author 4", Year: "4993"},
	)

	//router starts
	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/books/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/books", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", removeBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func getBooks(w http.ResponseWriter, r *http.Request) {

	var book Book
	books = []Book{}

	//creating db connection
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/db_books")

	if err != nil {
		log.Fatal(err.Error())
	}

	defer db.Close()

	results, err := db.Query("SELECT id,title,author,year FROM books")

	if err != nil {
		panic(err.Error())
	}
	defer results.Close()
	log.Println(results)

	for results.Next() {
		err = results.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		if err != nil {
			panic(err.Error())
		}
		log.Printf(book.Title)
	}
	// json.NewEncoder(w).Encode(books)
	log.Println("Get all books")
}

func getBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	log.Println(params)
	i, _ := strconv.Atoi(params["id"])
	for _, book := range books {
		if book.ID == i {
			json.NewEncoder(w).Encode(&book)
		}
	}
	log.Println("Get a single book")
}

func addBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)
	books = append(books, book)
	json.NewEncoder(w).Encode(books)

	log.Println("Add a book")
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	json.NewDecoder(r.Body).Decode(&book)
	for i, item := range books {
		if item.ID == book.ID {
			books[i] = book
		}
	}
	json.NewEncoder(w).Encode(books)
}

func removeBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	json.NewDecoder(r.Body).Decode(&book)

	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	for i, item := range books {
		if item.ID == id {
			books = append(books[:i], books[i+1:]...)
		}
	}
	json.NewEncoder(w).Encode(books)
}
