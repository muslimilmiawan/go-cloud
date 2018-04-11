package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Book is the data structure of our object
type Book struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	ISBN        string `json:"isbn"`
	Description string `json:"description,omitempty"`
}

// ToJSON is to convert Book structure to json
func (b Book) ToJSON() []byte {
	ToJSON, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}
	return ToJSON
}

// FromJSON is to convert Book structure From json to Book object
func FromJSON(data []byte) Book {
	book := Book{}
	err := json.Unmarshal(data, &book)
	if err != nil {
		panic(err)
	}
	return book
}

// Books is the collection of book
var books = map[string]Book{
	"0987654321": Book{Title: "Hitchiker's guide to the galaxy", Author: "Arthur Conan Doyle", ISBN: "0987654321"},
	"0123456789": Book{Title: "Cloud Native Go", Author: "M.L Reiner", ISBN: "0123456789"},
}

// BooksHandleFunction is a function to handle books list request
func BooksHandleFunction(w http.ResponseWriter, r *http.Request) {
	switch method := r.Method; method {
	case http.MethodGet:
		books := AllBooks()
		writeJSON(w, books)
	case http.MethodPost:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		book := FromJSON(body)
		isbn, created := CreateBook(book)
		if created {
			w.Header().Add("Location", "/api/books/"+isbn)
			w.WriteHeader(http.StatusCreated)
		} else {
			w.WriteHeader(http.StatusConflict)
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("This is a bad bad request."))
	}
}

// BookHandleFunction is a function to handle books list request
func BookHandleFunction(w http.ResponseWriter, r *http.Request) {
	isbn := r.URL.Path[len("/api/books/"):]

	switch method := r.Method; method {
	case http.MethodGet:
		book, found := GetBook(isbn)
		if found {
			writeJSON(w, book)
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	case http.MethodPut:
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		book := FromJSON(body)
		exists := updateBook(isbn, book)
		if exists {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Book id (" + isbn + ") is updated"))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	case http.MethodDelete:
		DeleteBook(isbn)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Book id (" + isbn + ") is deleted"))
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unsupported Request Method."))
	}
}

// AllBooks is the function to make a book arrays
func AllBooks() []Book {
	values := make([]Book, len(books))
	idx := 0
	for _, book := range books {
		values[idx] = book
		idx++
	}
	return values
}

func writeJSON(w http.ResponseWriter, i interface{}) {
	b, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Write(b)
}

// CreateBook is the function to create new Book
func CreateBook(book Book) (string, bool) {
	_, exists := books[book.ISBN]
	if exists {
		return "", false
	}
	books[book.ISBN] = book
	return book.ISBN, true
}

// GetBook is method to find book using isbn data
func GetBook(isbn string) (Book, bool) {
	book, found := books[isbn]
	return book, found
}

func updateBook(isbn string, book Book) bool {
	_, exists := books[isbn]
	if exists {
		books[isbn] = book
	}
	return exists
}

// DeleteBook is to delete specific book
func DeleteBook(isbn string) {
	delete(books, isbn)
}
